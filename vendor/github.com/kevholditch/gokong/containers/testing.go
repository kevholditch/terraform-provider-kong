package containers

import (
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
)

type TestContext struct {
	containers      []container
	KongHostAddress string
}

func StartKong(kongVersion string) *TestContext {
	log.SetOutput(os.Stdout)

	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	postgres := NewPostgresContainer(pool)
	kong := NewKongContainer(pool, postgres, kongVersion)

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
