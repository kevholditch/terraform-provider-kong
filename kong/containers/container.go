package containers

import (
	"strings"

	"github.com/ory/dockertest/v3"
)

func getContainerName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}

type container interface {
	Stop() error
}
