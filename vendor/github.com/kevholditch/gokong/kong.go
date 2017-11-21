package gokong

import (
	"errors"
	"fmt"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"net/http"
	"os"
)

type kong struct {
	Name        string
	pool        *dockertest.Pool
	resource    *dockertest.Resource
	HostAddress string
}

func NewKong(pool *dockertest.Pool, postgres *postgres) *kong {

	envVars := []string{
		"KONG_DATABASE=postgres",
		fmt.Sprintf("KONG_PG_HOST=%s", postgres.Name),
		fmt.Sprintf("KONG_PG_USER=%s", postgres.DatabaseUser),
		fmt.Sprintf("KONG_PG_PASSWORD=%s", postgres.Password),
	}

	options := &dockertest.RunOptions{
		Repository: "kong",
		Tag:        "0.11",
		Env:        envVars,
		Links:      []string{postgres.Name},
		Cmd:        []string{"kong", "migrations", "up"},
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
		Repository: "kong",
		Tag:        "0.11",
		Env:        envVars,
		Links:      []string{fmt.Sprintf("%s:postgres", postgres.Name)},
	}

	resource, err := pool.RunWithOptions(options)

	kongContainerName := getContainerName(resource)

	kongAddress := fmt.Sprintf("http://localhost:%v", resource.GetPort("8001/tcp"))

	if err := pool.Retry(func() error {
		var err error
		curlEndpoint := fmt.Sprintf("%s/apis", kongAddress)
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

	err = os.Setenv(EnvKongAdminHostAddress, kongAddress)
	if err != nil {
		log.Fatalf("Could not set kong host address env variable: %v", err)
	}

	return &kong{
		Name:        kongContainerName,
		pool:        pool,
		resource:    resource,
		HostAddress: kongAddress,
	}
}

func (kong *kong) Stop() error {
	return kong.pool.Purge(kong.resource)
}
