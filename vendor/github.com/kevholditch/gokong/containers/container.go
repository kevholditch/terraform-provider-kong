package containers

import (
	"gopkg.in/ory-am/dockertest.v3"
	"strings"
)

func getContainerName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}

type container interface {
	Stop() error
}
