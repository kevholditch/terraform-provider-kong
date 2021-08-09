package containers

import (
	"log"
	"os"

	"github.com/ory/dockertest/v3"
)

type TestContext struct {
	containers      []container
	KongHostAddress string
}

func StartKong(kongRepository string, kongVersion string, kongLicense string) *TestContext {
	log.SetOutput(os.Stdout)

	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	postgres := NewPostgresContainer(pool)
	kong := NewKongContainer(pool, postgres, kongRepository, kongVersion, kongLicense)

	return &TestContext{containers: []container{postgres, kong}, KongHostAddress: kong.HostAddress}
}

func StopKong(testContext *TestContext) {

	for _, container := range testContext.containers {
		err := container.Stop()
		if err != nil {
			log.Printf("Could not stop container: %v \n", err)
		}
	}

}
