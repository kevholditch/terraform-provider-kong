[![Build Status](https://travis-ci.org/kevholditch/gokong.svg?branch=master)](https://travis-ci.org/kevholditch/gokong)

GoKong
======
A kong go client fully tested with no mocks!!

## GoKong
GoKong is a easy to use api client for [kong](https://getkong.org/).  The difference with the gokong library is all of its tests are written against a real running kong running inside a docker container, yep that's right you won't see a horrible mock anywhere!!

## Run build
Ensure docker is installed then run:
`make`

## Importing

To add gokong via `go get`:
```
go get github.com/kevholditch/gokong
```

To add gokong via `govendor`:
```
govendor fetch github.com/kevholditch/gokong
```

## Usage

Import gokong
```go
import (
  gokong "github.com/kevholditch/gokong"
)
```

To create a default config for use with the client:
```go
config := gokong.NewDefaultConfig()
```

`NewDefaultConfig` creates a config with the host address set to the value of the env variable `KONG_ADMIN_ADDR`.
If the env variable is not set then the address is defaulted to `http://localhost:8001`.

You can of course create your own config with the address set to whatever you want:
```go
config := gokong.Config{HostAddress:"http://localhost:1234"}
```


Getting the status of the kong server:
```go
kongClient := gokong.NewClient(gokong.NewDefaultConfig())
status, err := kongClient.Status().Get()
```

Gokong is fluent so we can combine the above two lines into one:
```go
status, err := gokong.NewClient(gokong.NewDefaultConfig()).Status().Get()
```

## APIs
Create a new API:
```go
apiRequest := &gokong.ApiRequest{
	Name:                   "Example",
	Hosts:                  []string{"example.com"},
	Uris:                   []string{"/example"},
	Methods:                []string{"GET", "POST"},
	UpstreamUrl:            "http://localhost:4140/testservice",
	StripUri:               true,
	PreserveHost:           true,
	Retries:                3,
	UpstreamConnectTimeout: 1000,
	UpstreamSendTimeout:    2000,
	UpstreamReadTimeout:    3000,
	HttpsOnly:              true,
	HttpIfTerminated:       true,
}

api, err := gokong.NewClient(gokong.NewDefaultConfig()).Apis().Create(apiRequest)
```

Get an API by id:
```go
api, err := gokong.NewClient(gokong.NewDefaultConfig()).Apis().GetById("ExampleApi")
```

Get all apis:
```go
apis, err := gokong.NewClient(gokong.NewDefaultConfig()).Apis().GetAll()
```

Get all apis with a filter:
```go
filtered, err := gokong.NewClient(gokong.NewDefaultConfig()).Apis().GetAllFiltered(&gokong.GetAllFilter{Id:"936ad391-c30d-43db-b624-2f820d6fd38d", Name:"MyApi"})
```


