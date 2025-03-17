package container

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type (
	// PostgreSQLContainer хранит информацию о готовом контейнере.
	PostgreSQLContainer struct {
		testcontainers.Container
		Port     string // Port порт по которому можно подключиться к БД.
		Host     string // Host хост на котором будет доступна БД.
		Database string // Database название БД.
		Username string // Username пользователь БД.
		Password string // Password пароль пользователя БД.
	}

	// PostgreSQLConfig хранит конфигурацию контейнера.
	PostgreSQLConfig struct {
		Database string // Database название базы данных для переменной окружения контейнера POSTGRES_DB.
		Username string // Username пользователь для переменной окружения контейнера POSTGRES_USER.
		Password string // Password пароль для переменной окружения контейнера POSTGRES_PASSWORD.
		ImageTag string // ImageTag тег докер образа.
	}
)

// NewPostgreSQLContainer создает контейнера для базы данных PostgreSQL по переданной конфигурации.
func NewPostgreSQLContainer(ctx context.Context, cfg PostgreSQLConfig) (*PostgreSQLContainer, error) {
	const (
		image = "postgres"
		port  = "5432"
	)

	containerPort := fmt.Sprintf("%s/tcp", port)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Env: map[string]string{
				"POSTGRES_USER":     cfg.Username,
				"POSTGRES_PASSWORD": cfg.Password,
				"POSTGRES_DB":       cfg.Database,
			},
			ExposedPorts: []string{
				containerPort,
			},
			Image: fmt.Sprintf("%s:%s", image, cfg.ImageTag),
			WaitingFor: wait.ForExec([]string{"pg_isready"}).
				WithPollInterval(2 * time.Second).
				WithExitCodeMatcher(func(exitCode int) bool {
					return exitCode == 0
				}),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, err
	}

	return &PostgreSQLContainer{
		Container: container,
		Host:      host,
		Port:      mappedPort.Port(),
		Database:  cfg.Database,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}, nil
}
