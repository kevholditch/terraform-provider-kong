package containers

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

type postgresContainer struct {
	Name             string
	ConnectionString string
	Password         string
	DatabaseName     string
	DatabaseUser     string
	pool             *dockertest.Pool
	resource         *dockertest.Resource
}

func NewPostgresContainer(pool *dockertest.Pool) *postgresContainer {

	var db *sql.DB

	password := "kong"
	databaseName := "kong"
	databaseUser := "kong"

	resource, err := pool.Run("postgres", "9.6", []string{
		fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		fmt.Sprintf("POSTGRES_DB=%s", databaseName),
		fmt.Sprintf("POSTGRES_USER=%s", databaseUser),
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	connectionString := fmt.Sprintf("postgres://kong:kong@localhost:%s/kong?sslmode=disable", resource.GetPort("5432/tcp"))
	containerName := getContainerName(resource)

	if err = pool.Retry(func() error {
		var err error

		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	log.Printf("Postgres (%v): up", containerName)

	return &postgresContainer{
		Name:             containerName,
		ConnectionString: connectionString,
		Password:         password,
		DatabaseName:     databaseName,
		DatabaseUser:     databaseUser,
		pool:             pool,
		resource:         resource,
	}
}

func (postgres *postgresContainer) Stop() error {
	return postgres.pool.Purge(postgres.resource)
}
