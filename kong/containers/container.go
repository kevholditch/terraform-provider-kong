package containers

import (
	"strings"

	"github.com/ory/dockertest"
)

func getContainerName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}

type container interface {
	Stop() error
}
