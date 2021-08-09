package containers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ory/dockertest/v3"
)

type kongContainer struct {
	Name        string
	pool        *dockertest.Pool
	resource    *dockertest.Resource
	HostAddress string
}

func NewKongContainer(pool *dockertest.Pool, postgres *postgresContainer, kongRepository string, kongVersion string, kongLicense string) *kongContainer {

	envVars := []string{
		"KONG_DATABASE=postgres",
		"KONG_ADMIN_LISTEN=0.0.0.0:8001",
		fmt.Sprintf("KONG_LICENSE_DATA=%s", kongLicense),
		fmt.Sprintf("KONG_PG_HOST=%s", postgres.Name),
		fmt.Sprintf("KONG_PG_USER=%s", postgres.DatabaseUser),
		fmt.Sprintf("KONG_PG_PASSWORD=%s", postgres.Password),
	}

	options := &dockertest.RunOptions{
		Repository: kongRepository,
		Tag:        kongVersion,
		Env:        envVars,
		Links:      []string{postgres.Name},
		Cmd:        []string{"kong", "migrations", "bootstrap"},
	}

	migrations, err := pool.RunWithOptions(options)

	if err := pool.Retry(func() error {
		migrationsContainer, err := pool.Client.InspectContainer(migrations.Container.ID)
		migrationsContainerName := getContainerName(migrations)
		if err != nil {
			log.Fatalf("Could not get state of migrations container %v", err)
		}

		if migrationsContainer.State.Running {
			log.Printf("Kong Migrations (%v): waiting for migration", migrationsContainerName)
			return errors.New(fmt.Sprintf("Kong Migrations (%v): Error waiting for migration to finish", migrationsContainerName))
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to kong: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not start kong: %s", err)
	}

	options = &dockertest.RunOptions{
		Repository: kongRepository,
		Tag:        kongVersion,
		Env:        envVars,
		Links:      []string{fmt.Sprintf("%s:postgres", postgres.Name)},
	}

	resource, err := pool.RunWithOptions(options)

	kongContainerName := getContainerName(resource)

	kongAddress := fmt.Sprintf("http://localhost:%v", resource.GetPort("8001/tcp"))

	if err := pool.Retry(func() error {
		var err error
		curlEndpoint := fmt.Sprintf("%s/status", kongAddress)
		if err != nil {
			return err
		}

		resp, err := http.Get(curlEndpoint)
		if err != nil {
			return err
		}

		if resp.StatusCode >= 400 {
			return errors.New(fmt.Sprintf("Kong not ready: %+v", resp))
		}

		log.Printf("Kong (%v): up", kongContainerName)

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to kong: %s", err)
	}

	return &kongContainer{
		Name:        kongContainerName,
		pool:        pool,
		resource:    resource,
		HostAddress: kongAddress,
	}
}

func (kong *kongContainer) Stop() error {
	return kong.pool.Purge(kong.resource)
}
