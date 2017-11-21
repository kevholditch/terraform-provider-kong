package gokong

import (
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
)

type TestContext struct {
	containers []container
	code       int
}

func StartTestContainers() *TestContext {
	log.SetOutput(os.Stdout)

	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	postgres := NewPostgres(pool)
	kong := NewKong(pool, postgres)

	return &TestContext{containers: []container{postgres, kong}}
}

func StopTestContainers(testContext *TestContext) {

	for _, container := range testContext.containers {
		err := container.Stop()
		if err != nil {
			log.Printf("Could not stop container: %v \n", err)
		}
	}

}
