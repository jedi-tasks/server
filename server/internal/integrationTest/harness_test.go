package integrationTest

import (
	"context"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	pkgErrors "pkg/errors"
	auditLogRepository "server/internal/modules/auditLog/repository"
	auditLogService "server/internal/modules/auditLog/service"
	blockerRepository "server/internal/modules/blocker/repository"
	blockerService "server/internal/modules/blocker/service"
	listSectionRepository "server/internal/modules/listSection/repository"
	listSectionService "server/internal/modules/listSection/service"
	priorityRepository "server/internal/modules/priority/repository"
	priorityService "server/internal/modules/priority/service"
	tagRepository "server/internal/modules/tag/repository"
	tagService "server/internal/modules/tag/service"
	taskRepository "server/internal/modules/task/repository"
	taskService "server/internal/modules/task/service"
	taskListRepository "server/internal/modules/taskList/repository"
	taskListService "server/internal/modules/taskList/service"
	"server/internal/modules/transactor"
	userRepository "server/internal/modules/user/repository"
	userService "server/internal/modules/user/service"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
	"server/internal/utils/syncNotifier"
)

// services — сервисы, связанные так же, как в main.go
type services struct {
	transactor  *transactor.Transactor                 // Транзакции
	auditLog    *auditLogService.AuditLogService       // Аудит-лог
	taskList    *taskListService.TaskListService       // Списки
	listSection *listSectionService.ListSectionService // Разделы
	priority    *priorityService.PriorityService       // Приоритеты
	task        *taskService.TaskService               // Задачи
	tag         *tagService.TagService                 // Теги
	blocker     *blockerService.BlockerService         // Блокеры
}

// fixture — окружение одного теста: сервисы, отдельный пользователь с устройством и часы
type fixture struct {
	services
	t         *testing.T                         // Текущий тест
	ctx       context.Context                    // Контекст запросов
	necessary necessary.NecessaryUserInformation // Пользователь и устройство
	inboxID   uuid.UUID                          // Inbox пользователя
	clock     time.Time                          // Последнее выданное время изменения
}

// newServices собирает сервисы поверх тестовой базы
func newServices(t *testing.T) services {
	transactor := transactor.NewTransactor(database)
	auditLogService := auditLogService.NewAuditLogService(auditLogRepository.NewAuditLogRepository(database), syncNotifier.New())
	userService, err := userService.NewUserService(userRepository.NewUserRepository(database), transactor, []byte("salt"), auditLogService)
	if err != nil {
		t.Fatal(err)
	}

	taskListService := taskListService.NewTaskListService(taskListRepository.NewTaskListRepository(database), transactor, auditLogService, userService)
	listSectionService := listSectionService.NewListSectionService(listSectionRepository.NewListSectionRepository(database), transactor, auditLogService, userService, taskListService)
	priorityService := priorityService.NewPriorityService(priorityRepository.NewPriorityRepository(database), transactor, auditLogService, userService)
	taskService := taskService.NewTaskService(taskRepository.NewTaskRepository(database), transactor, auditLogService, userService, taskListService, listSectionService, priorityService)
	tagService := tagService.NewTagService(tagRepository.NewTagRepository(database), transactor, auditLogService, userService, taskService)
	blockerService := blockerService.NewBlockerService(blockerRepository.NewBlockerRepository(database), transactor, auditLogService, userService, taskService)

	// Связываем так же, как в main.go
	taskListService.SetListContentsCheckers(taskService, listSectionService)
	listSectionService.SetLiveTasksChecker(taskService)
	priorityService.SetLiveTasksChecker(taskService)
	taskService.SetTaskLinksCheckers(tagService, blockerService)

	return services{
		transactor:  transactor,
		auditLog:    auditLogService,
		taskList:    taskListService,
		listSection: listSectionService,
		priority:    priorityService,
		task:        taskService,
		tag:         tagService,
		blocker:     blockerService,
	}
}

// newFixture заводит нового пользователя с устройством и Inbox
func newFixture(t *testing.T) *fixture {
	ctx := context.Background()
	f := &fixture{
		services:  newServices(t),
		t:         t,
		ctx:       ctx,
		necessary: necessary.NecessaryUserInformation{UserID: uuid.New(), DeviceID: "device-" + uuid.NewString()},
		clock:     time.Now().UTC().Truncate(time.Millisecond),
	}

	// Создаем пользователя и устройство напрямую в базе
	f.exec(sq.Insert("tasks.users").SetMap(map[string]any{
		"id": f.necessary.UserID, "name": "Test", "email": f.necessary.UserID.String() + "@test.local",
		"password_hash": []byte("hash"), "password_salt": []byte("salt"), "time_create": f.clock,
	}))
	f.exec(sq.Insert("tasks.devices").SetMap(map[string]any{
		"user_id": f.necessary.UserID, "device_id": f.necessary.DeviceID, "refresh_token": "token",
		"application_bundle_id": "test", "application_version": "1", "application_build": "1",
		"device_ip_address": "127.0.0.1", "device_user_agent": "test", "device_os_name": "test",
		"device_os_version": "1", "device_name": "test", "device_model_name": "test",
	}))

	// Создаем Inbox и поднимаем чекпоинт устройства до его записи в аудит-логе
	if err := f.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {
		return f.taskList.CreateDefaults(ctxTx, f.necessary.UserID, f.necessary.DeviceID, f.clock)
	}); err != nil {
		t.Fatal(err)
	}
	f.exec(sq.Update("tasks.devices").
		Set("last_synced_audit_log_id", sq.Expr("(SELECT MAX(id) FROM tasks.audit_log WHERE user_id = ?)", f.necessary.UserID)).
		Where(sq.Eq{"user_id": f.necessary.UserID}))

	inboxID, err := f.taskList.GetInboxID(ctx, f.necessary.UserID)
	if err != nil {
		t.Fatal(err)
	}
	f.inboxID = inboxID

	return f
}

// exec выполняет служебный запрос подготовки данных
func (f *fixture) exec(query sq.Sqlizer) {
	f.t.Helper()
	if err := database.Exec(f.ctx, query); err != nil {
		f.t.Fatal(err)
	}
}

// next возвращает время изменения позже всех выданных ранее, чтобы последняя запись побеждала
func (f *fixture) next() time.Time {
	f.clock = f.clock.Add(time.Second)
	return f.clock
}

// ok падает, если мутация вернула ошибку
func (f *fixture) ok(err error) {
	f.t.Helper()
	if err != nil {
		f.t.Fatalf("неожиданная ошибка: %v", err)
	}
}

// rejected проверяет, что мутация отклонена BadRequest с ожидаемым текстом для пользователя
func (f *fixture) rejected(err error, humanText string) {
	f.t.Helper()
	var customError pkgErrors.Error
	if !pkgErrors.As(err, &customError) {
		f.t.Fatalf("ожидалась ошибка %q, получено %v", humanText, err)
	}
	if customError.ErrorType != errors.BadRequest.ErrorType {
		f.t.Fatalf("ожидался BadRequest, получено %v (%v)", customError.ErrorType.HTTPCode, err)
	}
	if customError.HumanText != humanText {
		f.t.Fatalf("ожидался текст %q, получено %q", humanText, customError.HumanText)
	}
}

// count возвращает число строк по запросу COUNT(*)
func (f *fixture) count(query sq.SelectBuilder) int {
	f.t.Helper()
	var result int
	if err := database.Get(f.ctx, &result, query); err != nil {
		f.t.Fatal(err)
	}
	return result
}
