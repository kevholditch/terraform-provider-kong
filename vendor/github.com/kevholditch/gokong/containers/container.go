package containers

import (
	"strings"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

func getContainerName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}

type container interface {
	Stop() error
}
