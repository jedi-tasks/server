package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"pkg/gracefulShutdown"
	fiber2 "pkg/http/fiber"
	"pkg/log/model"
	"pkg/middleware"
	"pkg/stackTrace"
	"server/internal/metrics"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/jedi-tasks/go-server-grpc/proto"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"pkg/database/pgsql"
	grpcServer "pkg/grpc/server"
	"server/internal/utils/grpc/interceptor"

	"pkg/jwtManager"
	"pkg/log"
	"pkg/migrator"
	"pkg/panicRecover"
	"server/internal/config"
	auditLogRepository "server/internal/modules/auditLog/repository"
	auditLogService "server/internal/modules/auditLog/service"
	authEndpointGRPC "server/internal/modules/auth/endpoint/grpc"
	authModel "server/internal/modules/auth/model"
	authService "server/internal/modules/auth/service"
	blockerEndpointGRPC "server/internal/modules/blocker/endpoint/grpc"
	blockerRepository "server/internal/modules/blocker/repository"
	blockerService "server/internal/modules/blocker/service"
	durationOptionEndpointGRPC "server/internal/modules/durationOption/endpoint/grpc"
	durationOptionRepository "server/internal/modules/durationOption/repository"
	durationOptionService "server/internal/modules/durationOption/service"
	listSectionEndpointGRPC "server/internal/modules/listSection/endpoint/grpc"
	listSectionRepository "server/internal/modules/listSection/repository"
	listSectionService "server/internal/modules/listSection/service"
	priorityEndpointGRPC "server/internal/modules/priority/endpoint/grpc"
	priorityRepository "server/internal/modules/priority/repository"
	priorityService "server/internal/modules/priority/service"
	settingsEndpointGRPC "server/internal/modules/settings/endpoint/grpc"
	settingsRepository "server/internal/modules/settings/repository"
	settingsService "server/internal/modules/settings/service"
	syncEndpointGRPC "server/internal/modules/sync/endpoint/grpc"
	syncService "server/internal/modules/sync/service"

	tagEndpointGRPC "server/internal/modules/tag/endpoint/grpc"
	tagRepository "server/internal/modules/tag/repository"
	tagService "server/internal/modules/tag/service"
	taskEndpointGRPC "server/internal/modules/task/endpoint/grpc"
	taskRepository "server/internal/modules/task/repository"
	taskService "server/internal/modules/task/service"
	taskListEndpointGRPC "server/internal/modules/taskList/endpoint/grpc"
	taskListRepository "server/internal/modules/taskList/repository"
	taskListService "server/internal/modules/taskList/service"
	"server/internal/modules/transactor"
	userEndpointGRPC "server/internal/modules/user/endpoint/grpc"
	userRepository "server/internal/modules/user/repository"
	userService "server/internal/modules/user/service"
	"server/internal/utils/errors"
	"server/internal/utils/syncNotifier"
	pgsqlMigrations "server/migrations/pgsql"
)

const version = "@{version}"
const build = "@{build}"
const buildDate = "@{buildDate}"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Основной контекст приложения
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Перехватываем возможную панику
	defer func() {
		panicRecover.PanicRecover(func(err error) {
			log.Fatal(err)
		})
	}()

	// Получаем конфиг
	conf := config.Load()

	// Получаем имя хоста
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Инициализируем логгер
	if err := log.InitDefaultLogger(
		model.SystemInfo{
			Version:     version,
			Build:       build,
			ServiceName: conf.ServiceName,
			Env:         conf.Environment,
			Hostname:    hostname,
			BuildDate:   buildDate,
		},
		conf.Logger,
	); err != nil {
		return err
	}

	// Инициализируем все синглтоны
	log.Info("Инициализируем синглтоны")
	if err = initSingletons(conf); err != nil {
		return err
	}

	// Подключаемся к базе данных
	log.Info("Подключаемся к БД")
	pgsql, err := pgsql.NewClientPgsql(ctx, conf.Pgsql)
	if err != nil {
		return err
	}
	defer pgsql.Close()

	// Запускаем миграции в базе данных
	// TODO: Подумать, как откатывать миграции при ошибках
	log.Info("Запускаем миграции")
	postgreSQLMigrator, err := migrator.NewMigrator(
		migrator.MigratorConfig{
			Migrations:      nil,
			Conn:            pgsql.DB.DB,
			EmbedMigrations: pgsqlMigrations.EmbedMigrationsPgsql,
			Dir:             "pgsql",
			Dialect:         goose.DialectPostgres,
		},
	)
	if err != nil {
		return err
	}
	if err = postgreSQLMigrator.Up(ctx); err != nil {
		return err
	}

	// Регистрируем репозитории
	transactor := transactor.NewTransactor(pgsql)
	userRepository := userRepository.NewUserRepository(pgsql)
	auditLogRepository := auditLogRepository.NewAuditLogRepository(pgsql)
	settingsRepository := settingsRepository.NewSettingsRepository(pgsql)
	taskListRepository := taskListRepository.NewTaskListRepository(pgsql)
	listSectionRepository := listSectionRepository.NewListSectionRepository(pgsql)
	priorityRepository := priorityRepository.NewPriorityRepository(pgsql)
	durationOptionRepository := durationOptionRepository.NewDurationOptionRepository(pgsql)
	taskRepository := taskRepository.NewTaskRepository(pgsql)
	tagRepository := tagRepository.NewTagRepository(pgsql)
	blockerRepository := blockerRepository.NewBlockerRepository(pgsql)

	// Будильник активных SubscribeToSync-стримов — одна на весь процесс инстанс, используется и
	// auditLogService (пишет уведомления при TrackMutation), и syncService (подписывает стримы).
	syncNotifier := syncNotifier.New()

	// Регистрируем сервисы
	auditLogService := auditLogService.NewAuditLogService(
		auditLogRepository,
		syncNotifier,
	)

	userService, err := userService.NewUserService(
		userRepository,
		transactor,
		[]byte(conf.Auth.GeneralSalt),
		auditLogService,
	)
	if err != nil {
		return err
	}

	settingsService := settingsService.NewSettingsService(
		settingsRepository,
		transactor,
		auditLogService,
		userService,
		settingsService.Version{
			Version: version,
			Build:   build,
		},
	)

	taskListService := taskListService.NewTaskListService(taskListRepository, transactor, auditLogService, userService)
	listSectionService := listSectionService.NewListSectionService(listSectionRepository, transactor, auditLogService, userService, taskListService)
	priorityService := priorityService.NewPriorityService(priorityRepository, transactor, auditLogService, userService)
	durationOptionService := durationOptionService.NewDurationOptionService(durationOptionRepository, transactor, auditLogService, userService)
	taskService := taskService.NewTaskService(taskRepository, transactor, auditLogService, userService, taskListService, listSectionService, priorityService)
	tagService := tagService.NewTagService(tagRepository, transactor, auditLogService, userService, taskService)
	blockerService := blockerService.NewBlockerService(blockerRepository, transactor, auditLogService, userService, taskService)

	// Сервисы, зависящие друг от друга по кругу, связываем после создания
	// Сервер ничего не удаляет каскадом: удаление отклоняется, пока на запись ссылаются живые строки
	taskListService.SetListContentsCheckers(taskService, listSectionService)
	listSectionService.SetLiveTasksChecker(taskService)
	priorityService.SetLiveTasksChecker(taskService)
	taskService.SetTaskLinksCheckers(tagService, blockerService)

	authService := authService.NewAuthService(
		userRepository,
		transactor,
		[]byte(conf.Auth.GeneralSalt),
		conf.Auth.SignUpEnabled,
		auditLogService,
		[]authService.DefaultsCreator{taskListService, priorityService, durationOptionService},
	)

	// Создаем начального пользователя из окружения, если его еще нет
	if conf.Auth.InitialUser.Email != "" {
		created, err := authService.EnsureUser(ctx, authModel.EnsureUserReq{
			Name:     conf.Auth.InitialUser.Name,
			Email:    conf.Auth.InitialUser.Email,
			Password: conf.Auth.InitialUser.Password,
		})
		if err != nil {
			log.Fatal(err)
		}
		if created {
			log.Info("Создан начальный пользователь")
		}
	}

	syncService := syncService.NewSyncService(
		transactor,
		userService,
		auditLogService,
		taskListService,
		listSectionService,
		priorityService,
		durationOptionService,
		taskService,
		tagService,
		blockerService,
		tagService,
		settingsService,
		syncNotifier,
	)

	// Настраиваем и получаем GRPC-сервер
	grpcServer := grpcServer.NewGRPCServer(&grpcServer.ServerOptions{
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			interceptor.NewErrorHandlerInterceptor().Unary(),
			interceptor.NewTimeInterceptor().Unary(),
			interceptor.NewAuthInterceptor(config.GetOpenForEveryoneMethods()).Unary(),
			interceptor.NewLoggerInterceptor().Unary(),
		},
		StreamInterceptors: []grpc.StreamServerInterceptor{
			interceptor.NewErrorHandlerInterceptor().Stream(),
			interceptor.NewTimeInterceptor().Stream(),
			interceptor.NewAuthInterceptor(config.GetOpenForEveryoneMethods()).Stream(),
			interceptor.NewLoggerInterceptor().Stream(),
		},
	})
	proto.RegisterAuthEndpointServer(grpcServer, authEndpointGRPC.NewAuthEndpoint(authService))
	proto.RegisterUserEndpointServer(grpcServer, userEndpointGRPC.NewUserEndpoint(userService))
	proto.RegisterSettingsEndpointServer(grpcServer, settingsEndpointGRPC.NewSettingsEndpoint(settingsService))
	proto.RegisterTaskListEndpointServer(grpcServer, taskListEndpointGRPC.NewTaskListEndpoint(taskListService))
	proto.RegisterListSectionEndpointServer(grpcServer, listSectionEndpointGRPC.NewListSectionEndpoint(listSectionService))
	proto.RegisterPriorityEndpointServer(grpcServer, priorityEndpointGRPC.NewPriorityEndpoint(priorityService))
	proto.RegisterDurationOptionEndpointServer(grpcServer, durationOptionEndpointGRPC.NewDurationOptionEndpoint(durationOptionService))
	proto.RegisterTaskEndpointServer(grpcServer, taskEndpointGRPC.NewTaskEndpoint(taskService))
	proto.RegisterTagEndpointServer(grpcServer, tagEndpointGRPC.NewTagEndpoint(tagService))
	proto.RegisterBlockerEndpointServer(grpcServer, blockerEndpointGRPC.NewBlockerEndpoint(blockerService))
	proto.RegisterSyncEndpointServer(grpcServer, syncEndpointGRPC.NewSyncEndpoint(syncService))

	// Создаем слушателя порта для gRPC-сервера
	grpcLn, err := net.Listen("tcp", conf.Listen.GRPC)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	defer grpcLn.Close()

	// Настраиваем и получаем HTTP-сервер
	httpServer, err := fiber2.GetDefaultServer(conf.ServiceName, nil)
	if err != nil {
		return err
	}

	// Настраиваем и получаем служебный HTTP-сервер
	serviceHTTPServer, err := fiber2.GetDefaultServer(conf.ServiceName, nil)
	if err != nil {
		return err
	}

	// Настраиваем служебный HTTP-сервер для pprof
	serviceHTTPServer.Use(pprof.New())

	// Настраиваем healthcheck
	serviceHTTPServer.Use(healthcheck.New(healthcheck.Config{
		Next:              nil,
		LivenessProbe:     nil,
		LivenessEndpoint:  "",
		ReadinessProbe:    fiber2.NewReadyHandler(nil),
		ReadinessEndpoint: "",
	}))

	// Настраиваем метрики
	serviceHTTPServer.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Инициализируем эндпоинты
	httpServer.Get("/version", middleware.NewVersionHandler(version, build, buildDate, hostname)) // GET /version

	// Создаем wait группу
	eg, ctx := errgroup.WithContext(ctx)

	// Создаем горутину на запуск HTTP-сервера
	eg.Go(func() error {
		log.Info(fmt.Sprintf("HTTP server is listening :%v", conf.Listen.HTTP))
		if err := httpServer.Listen(conf.Listen.HTTP); err != nil {
			return errors.InternalServer.Wrap(err)
		}
		return nil
	})

	// Создаем горутину на запуск служебного HTTP-сервера
	eg.Go(func() error {
		log.Info(fmt.Sprintf("service HTTP server is listening :%v", conf.Listen.ServiceHTTP))
		if err := serviceHTTPServer.Listen(conf.Listen.ServiceHTTP); err != nil {
			return errors.InternalServer.Wrap(err)
		}
		return nil
	})

	// Создаем горутину на запуск gRPC-сервера
	eg.Go(func() error {
		log.Info(fmt.Sprintf("gRPC server is listening: %v", conf.Listen.GRPC))
		if err := grpcServer.Serve(grpcLn); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			return errors.InternalServer.Wrap(err)
		}
		return nil
	})

	// Добавляем в группу горутину для graceful shutdown
	gracefulShutdown.AddGracefulShutdownErrGroup(eg, ctx, []*fiber.App{httpServer, serviceHTTPServer}, []*grpc.Server{grpcServer})

	// Ждем завершения контекста или ошибок в горутинах
	return eg.Wait()
}

func initSingletons(conf config.Config) error {

	stackTrace.Init("", true, 1)

	// Инициализируем JWT-менеджер
	accessTokenTTL, err := time.ParseDuration(conf.Auth.AccessTokenTTL)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	refreshTokenTTL, err := time.ParseDuration(conf.Auth.RefreshTokenTTL)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	jwtManager.Init([]byte(conf.Auth.SigningKey), accessTokenTTL, refreshTokenTTL)

	return metrics.Init(conf.ServiceName)
}
