package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"user-crud/internal/config"
	"user-crud/internal/repository/postgres"
	"user-crud/internal/services/user"
	notifierMocks "user-crud/internal/tests/mocks/notifier"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/mock"
)

var userService user.Service

func TestMain(m *testing.M) {
	ctx := context.Background()
	conf := testConfig()
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4-alpine",
		Env: []string{
			"POSTGRES_USER=" + conf.Postgres.User,
			"POSTGRES_PASSWORD=" + conf.Postgres.Password,
			"POSTGRES_DB=" + conf.Postgres.DBName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: conf.Postgres.Port},
			},
		},
	}

	resource, err := p.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("could not run docker resource: %s", err)
	}

	resource.Expire(30)

	var postgresClient *postgres.Client

	connect := func() error {
		postgresClient, err = postgres.NewClient(ctx, conf)
		if err != nil {
			return err
		}
		RunTestGooseMigrations(postgresClient.Database())

		return nil
	}

	err = p.Retry(connect)
	if err != nil {
		log.Fatalf("could not connect to the db: %v", err)
	}

	notifier := &notifierMocks.Service{}
	notifier.On("Publish", mock.Anything, mock.Anything).Return(nil)

	userService = user.NewService(postgresClient, notifier)
	code := m.Run()

	if err = p.Purge(resource); err != nil {
		log.Fatalf("could not purge resource %s", err)
	}

	os.Exit(code)
}

func testConfig() *config.Config {
	return &config.Config{
		Postgres: struct {
			DBName   string `yaml:"db_name"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Port     string `yaml:"port"`
			DSN      string `yaml:"dsn"`
		}{
			User:     "postgres",
			Password: "secret",
			DBName:   "test",
			Port:     "5435",
			DSN:      "host=localhost user=%s password=%s database=%s port=%s",
		},
	}
}

func RunTestGooseMigrations(db *sqlx.DB) {
	goose.SetBaseFS(postgres.Migrations)
	if err := goose.Up(db.DB, "migrations"); err != nil {
		panic(err)
	}
}
