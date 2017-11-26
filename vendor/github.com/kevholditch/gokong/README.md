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
  "github.com/kevholditch/gokong"
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
Create a new API ([for more information on the API fields see the Kong documentation](https://getkong.org/docs/0.11.x/admin-api/#add-api)):
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
api, err := gokong.NewClient(gokong.NewDefaultConfig()).Apis().GetById("cdf5372e-1c10-4ea5-a3dd-1e4c31bb99f5")
```

Get an API by name:
```go
api, err :=  gokong.NewClient(gokong.NewDefaultConfig()).Apis().GetByName("Example")
```

List all APIs:
```go
apis, err := gokong.NewClient(gokong.NewDefaultConfig()).Apis().List()
```

List all APIs with a filter:
```go
apis, err := gokong.NewClient(gokong.NewDefaultConfig()).Apis().ListFiltered(&gokong.ApiFilter{Id:"936ad391-c30d-43db-b624-2f820d6fd38d", Name:"MyApi"})
```

Delete an API by id:
```go
err :=  gokong.NewClient(gokong.NewDefaultConfig()).Apis().DeleteById("f138641a-a15b-43c3-bd76-7157a68eae24")
```

Delete an API by name:
```go
err :=  gokong.NewClient(gokong.NewDefaultConfig()).Apis().DeleteByName("Example")
```

Update an API by id:
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

updatedApi, err :=  gokong.NewClient(gokong.NewDefaultConfig()).Apis().UpdateById("1213a00d-2b12-4d65-92ad-5a02d6c710c2", apiRequest)
```

Update an API by name:
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

updatedApi, err :=  gokong.NewClient(gokong.NewDefaultConfig()).Apis().UpdateByName("Example", apiRequest)
```


## Consumers
Create a new Consumer ([for more information on the Consumer Fields see the Kong documentation](https://getkong.org/docs/0.11.x/admin-api/#create-consumer)):
```go
consumerRequest := &gokong.ConsumerRequest{
  Username: "User1",
  CustomId: "SomeId",
}

consumer, err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().Create(consumerRequest)
```

Get a Consumer by id:
```go
consumer, err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().GetById("e8ccbf13-a662-45be-9b6a-b549cc739c18")
```

Get a Consumer by username:
```go
consumer, err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().GetByUsername("User1")
```

List all Consumers:
```go
consumers, err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().List()
```

List all Consumers with a filter:
```go
consumers, err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().ListFiltered(&gokong.ConsumerFilter{CustomId:"1234", Username: "User1"})
```

Delete a Consumer by id:
```go
err :=  gokong.NewClient(gokong.NewDefaultConfig()).Consumers().DeleteById("7c8741b7-3cf5-4d90-8674-b34153efbcd6")
```

Delete a Consumer by username:
```go
err :=  gokong.NewClient(gokong.NewDefaultConfig()).Consumers().DeleteByUsername("User1")
```

Update a Consumer by id:
```go
consumerRequest := &gokong.ConsumerRequest{
  Username: "User1",
  CustomId: "SomeId",
}

updatedConsumer, err :=  gokong.NewClient(gokong.NewDefaultConfig()).Consumers().UpdateById("44a37c3d-a252-4968-ab55-58c41b0289c2", consumerRequest)
```

Update a Consumer by username:
```go
consumerRequest := &gokong.ConsumerRequest{
  Username: "User2",
  CustomId: "SomeId",
}

updatedConsumer, err :=  gokong.NewClient(gokong.NewDefaultConfig()).Consumers().UpdateByUsername("User2", consumerRequest)
```

# Contributing
I would love to get contributions to the project so please feel free to submit a PR.  To setup your dev station you need go and docker installed.

Once you have cloned the repository the `make` command will build the code and run all of the tests.  If they all pass then you are good to go!

If when you run the make command you get the following error:
```
gofmt needs running on the following files:
```
Then all you need to do is run `make fmt` this will reformat all of the code (I know awesome)!!

