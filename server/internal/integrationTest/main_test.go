// Package integrationTest — интеграционные тесты сервисов на настоящем Postgres (testcontainers-go)
// с миграциями goose из migrations/pgsql. Нужен запущенный Docker
package integrationTest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"pkg/database/pgsql"
	"pkg/migrator"
	"pkg/sql"
	pgsqlMigrations "server/migrations/pgsql"
)

// database — подключение к тестовой базе, общее для всех тестов пакета
var database *sql.DB

// TestMain поднимает контейнер Postgres, накатывает миграции и запускает тесты
func TestMain(m *testing.M) {
	os.Exit(run(m))
}

// run выполняет подготовку окружения и тесты, возвращая код завершения
func run(m *testing.M) int {
	ctx := context.Background()

	// Поднимаем Postgres 17
	container, err := postgres.Run(ctx, "postgres:17",
		postgres.WithDatabase("tasks_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(time.Minute)),
	)
	if err != nil {
		fmt.Println("не удалось поднять Postgres:", err)
		return 1
	}
	defer func() { _ = testcontainers.TerminateContainer(container) }()

	// Подключаемся к базе
	host, err := container.Host(ctx)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	database, err = pgsql.NewClientPgsql(ctx, pgsql.PgsqlConfigEnv{
		Host:     host + ":" + port.Port(),
		User:     "postgres",
		Password: "postgres",
		Database: "tasks_test",
	})
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer database.Close()

	// Накатываем миграции так же, как при старте сервера
	postgreSQLMigrator, err := migrator.NewMigrator(migrator.MigratorConfig{
		Migrations:      nil,
		Conn:            database.DB.DB,
		EmbedMigrations: pgsqlMigrations.EmbedMigrationsPgsql,
		Dir:             "pgsql",
		Dialect:         goose.DialectPostgres,
	})
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if err = postgreSQLMigrator.Up(ctx); err != nil {
		fmt.Println(err)
		return 1
	}

	return m.Run()
}
