module github.com/kevholditch/terraform-provider-kong

go 1.16

replace github.com/kevholditch/terraform-provider-kong/kong => ./kong

require (
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/containerd/continuity v0.1.0 // indirect
	github.com/docker/cli v20.10.8+incompatible // indirect
	github.com/docker/docker v20.10.8+incompatible // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
	github.com/kong/go-kong v0.28.0
	github.com/lib/pq v1.0.0
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/ory/dockertest/v3 v3.7.0
	github.com/pkg/errors v0.9.1
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
)
