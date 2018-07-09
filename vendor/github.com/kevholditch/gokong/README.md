[![Build Status](https://travis-ci.org/kevholditch/gokong.svg?branch=master)](https://travis-ci.org/kevholditch/gokong)

GoKong
======
A kong go client fully tested with no mocks!!

## Notice
As per version v1.0.0 all values are now pointer to type, this is to allow detection between setting the zero value of the type versus not setting it.
For example if you want to set a string to "" this will be omitted when serializing to json if you use `string` so to get round this we can use *string.
This is as per the way the aws go sdk does it.


## GoKong
GoKong is a easy to use api client for [kong](https://getkong.org/).  The difference with the gokong library is all of its tests are written against a real running kong running inside a docker container, yep that's right you won't see a horrible mock anywhere!!

## Supported Kong Versions
As per [travis build](https://travis-ci.org/kevholditch/gokong):
```
KONG_VERSION=0.13.0
```

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

There are a number of options you can set via config either by explicitly setting them when creating a config instance or
 by simply using the `NewDefaultConfig` method and using env variables.  Below is a table of the fields, the env variables that can be used
 to set them and their default values if you do not provide one via an env variable:

| Config property       | Env variable         | Default if not set    | Use                                                                             |
|:----------------------|:---------------------|:----------------------|:--------------------------------------------------------------------------------|
| HostAddress           | KONG_ADMIN_ADDR      | http://localhost:8001 | The url of the kong admin api                                                   |
| Username              | KONG_ADMIN_USERNAME  | not set               | Username for the kong admin api                                                 |
| Password              | KONG_ADMIN_PASSWORD  | not set               | Password for the kong admin api                                                 |
| InsecureSkipVerify    | TLS_SKIP_VERIFY      | false                 | Whether to skip tls certificate verification for the kong api when using https  |
| ApiKey                | KONG_API_KEY         | not set               | The api key you have used to lock down the kong admin api (via key-auth plugin) |



You can of course create your own config with the address set to whatever you want:
```go
config := gokong.Config{HostAddress:"http://localhost:1234"}
```

Also you can apply Username and Password for admin-api Basic Auth:
```go
config := gokong.Config{HostAddress:"http://localhost:1234",Username:"adminuser",Password:"yoursecret"}
```

If you need to ignore TLS verification, you can set InsecureSkipVerify:
```go
config := gokong.Config{InsecureSkipVerify: true}
```
This might be needed if your Kong installation is using a self-signed certificate, or if you are proxying to the Kong admin port.

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
Create a new API ([for more information on the API fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#api-object):
```go
apiRequest := &gokong.ApiRequest{
	Name:                   "Example",
	Hosts:                  gokong.StringSlice([]string{"example.com"}),
  Uris:                   gokong.StringSlice([]string{"/example"}),
  Methods:                gokong.StringSlice([]string{"GET", "POST"}),
  UpstreamUrl:            gokong.String("http://localhost:4140/testservice"),
  StripUri:               gokong.Bool(true),
  PreserveHost:           gokong.Bool(true),
  Retries:                gokong.Int(3),
  UpstreamConnectTimeout: gokong.Int(1000),
  UpstreamSendTimeout:    gokong.Int(2000),
  UpstreamReadTimeout:    gokong.Int(3000),
  HttpsOnly:              gokong.Bool(true),
  HttpIfTerminated:       gokong.Bool(true),
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
  	Hosts:                  gokong.StringSlice([]string{"example.com"}),
    Uris:                   gokong.StringSlice([]string{"/example"}),
    Methods:                gokong.StringSlice([]string{"GET", "POST"}),
    UpstreamUrl:            gokong.String("http://localhost:4140/testservice"),
    StripUri:               gokong.Bool(true),
    PreserveHost:           gokong.Bool(true),
    Retries:                gokong.Int(3),
    UpstreamConnectTimeout: gokong.Int(1000),
    UpstreamSendTimeout:    gokong.Int(2000),
    UpstreamReadTimeout:    gokong.Int(3000),
    HttpsOnly:              gokong.Bool(true),
    HttpIfTerminated:       gokong.Bool(true),
}

updatedApi, err :=  gokong.NewClient(gokong.NewDefaultConfig()).Apis().UpdateById("1213a00d-2b12-4d65-92ad-5a02d6c710c2", apiRequest)
```

Update an API by name:
```go
apiRequest := &gokong.ApiRequest{
 	Name:                   "Example",
 	Hosts:                  gokong.StringSlice([]string{"example.com"}),
  Uris:                   gokong.StringSlice([]string{"/example"}),
  Methods:                gokong.StringSlice([]string{"GET", "POST"}),
  UpstreamUrl:            gokong.String("http://localhost:4140/testservice"),
  StripUri:               gokong.Bool(true),
  PreserveHost:           gokong.Bool(true),
  Retries:                gokong.Int(3),
  UpstreamConnectTimeout: gokong.Int(1000),
  UpstreamSendTimeout:    gokong.Int(2000),
  UpstreamReadTimeout:    gokong.Int(3000),
  HttpsOnly:              gokong.Bool(true),
  HttpIfTerminated:       gokong.Bool(true),

updatedApi, err :=  gokong.NewClient(gokong.NewDefaultConfig()).Apis().UpdateByName("Example", apiRequest)
```


## Consumers
Create a new Consumer ([for more information on the Consumer Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#consumer-object)):
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

## Plugins
Create a new Plugin to be applied to all APIs and consumers do not set `ApiId` or `ConsumerId`.  Not all plugins can be configured in this way
 ([for more information on the Plugin Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#add-plugin)):

```go
pluginRequest := &gokong.PluginRequest{
  Name: "response-ratelimiting",
  Config: map[string]interface{}{
    "limits.sms.minute": 20,
  },
}

createdPlugin, err := gokong.NewClient(gokong.NewDefaultConfig()).Plugins().Create(pluginRequest)
```

Create a new Plugin for a single API (only set `ApiId`), not all plugins can be configured in this way ([for more information on the Plugin Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#plugin-object)):
```go
client := gokong.NewClient(gokong.NewDefaultConfig())

apiRequest := &gokong.ApiRequest{
  Name:                   "test-api",
  Hosts:                  []string{"example.com"},
  Uris:                   []string{"/example"},
  Methods:                []string{"GET", "POST"},
  UpstreamUrl:            "http://localhost:4140/testservice",
  StripUri:               true,
  PreserveHost:           true,
  Retries:                "3",
  UpstreamConnectTimeout: 1000,
  UpstreamSendTimeout:    2000,
  UpstreamReadTimeout:    3000,
  HttpsOnly:              true,
  HttpIfTerminated:       true,
}

createdApi, err := client.Apis().Create(apiRequest)

pluginRequest := &gokong.PluginRequest{
  Name: "response-ratelimiting",
  ApiId: createdApi.Id,
  Config: map[string]interface{}{
    "limits.sms.minute": 20,
  },
}

createdPlugin, err :=  client.Plugins().Create(pluginRequest)
```

Create a new Plugin for a single Consumer (only set `ConsumerId`), Not all plugins can be configured in this way ([for more information on the Plugin Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#plugin-object)):
```go
client := gokong.NewClient(gokong.NewDefaultConfig())

consumerRequest := &gokong.ConsumerRequest{
  Username: "User1",
  CustomId: "test",
}

createdConsumer, err := client.Consumers().Create(consumerRequest)

pluginRequest := &gokong.PluginRequest{
  Name: "response-ratelimiting",
  ConsumerId: createdConsumer.Id,
  Config: map[string]interface{}{
    "limits.sms.minute": 20,
  },
}

createdPlugin, err :=  client.Plugins().Create(pluginRequest)
```

Create a new Plugin for a single Consumer and Api (set `ConsumerId` and `ApiId`), Not all plugins can be configured in this way ([for more information on the Plugin Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#plugin-object)):
```go
client := gokong.NewClient(gokong.NewDefaultConfig())

consumerRequest := &gokong.ConsumerRequest{
  Username: "User1",
  CustomId: "test",
}

createdConsumer, err := client.Consumers().Create(consumerRequest)

apiRequest := &gokong.ApiRequest{
  Name:                   "test-api",
  Hosts:                  []string{"example.com"},
  Uris:                   []string{"/example"},
  Methods:                []string{"GET", "POST"},
  UpstreamUrl:            "http://localhost:4140/testservice",
  StripUri:               true,
  PreserveHost:           true,
  Retries:                "3",
  UpstreamConnectTimeout: 1000,
  UpstreamSendTimeout:    2000,
  UpstreamReadTimeout:    3000,
  HttpsOnly:              true,
  HttpIfTerminated:       true,
}

createdApi, err := client.Apis().Create(apiRequest)

pluginRequest := &gokong.PluginRequest{
  Name:       "response-ratelimiting",
  ConsumerId: createdConsumer.Id,
  ApiId:      createdApi.Id,
  Config: map[string]interface{}{
    "limits.sms.minute": 20,
  },
}

createdPlugin, err :=  client.Plugins().Create(pluginRequest)
```

Get a plugin by id:
```go
plugin, err := gokong.NewClient(gokong.NewDefaultConfig()).Plugins().GetById("04bda233-d035-4b8a-8cf2-a53f3dd990f3")
```

List all plugins:
```go
plugins, err := gokong.NewClient(gokong.NewDefaultConfig()).Plugins().List()
```

List all plugins with a filter:
```go
plugins, err := gokong.NewClient(gokong.NewDefaultConfig()).Plugins().ListFiltered(&gokong.PluginFilter{Name: "response-ratelimiting", ConsumerId: "7009a608-b40c-4a21-9a90-9219d5fd1ac7"})
```

Delete a plugin by id:
```go
err := gokong.NewClient(gokong.NewDefaultConfig()).Plugins().DeleteById("f2bbbab8-3e6f-4d9d-bada-d486600b3b4c")
```

Update a plugin by id:
```go
updatePluginRequest := &gokong.PluginRequest{
  Name:       "response-ratelimiting",
  ConsumerId: createdConsumer.Id,
  ApiId:      createdApi.Id,
  Config: map[string]interface{}{
    "limits.sms.minute": 20,
  },
}

updatedPlugin, err := gokong.NewClient(gokong.NewDefaultConfig()).Plugins().UpdateById("70692eed-2293-486d-b992-db44a6459360", updatePluginRequest)
```
## Configure a plugin for a Consumer
To configure a plugin for a consumer you can use the `CreatePluginConfig`, `GetPluginConfig` and `DeletePluginConfig` methods on the `Consumers` endpoint.
  Some plugins require configuration for a consumer for example the [jwt plugin[(https://getkong.org/plugins/jwt/#create-a-jwt-credential).

Create a plugin config for a consumer:
```go
createdPluginConfig, err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().CreatePluginConfig("f6539872-d8c5-4d6c-a2f2-923760329e4e", "jwt", "{\"key\": \"a36c3049b36249a3c9f8891cb127243c\"}")
```

Get a plugin config for a consumer by plugin config id:
```
pluginConfig, err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().GetPluginConfig("58c5229-dc92-4632-91c1-f34d9b84db0b", "jwt", "22700b52-ba59-428e-b03b-ba429b1e775e")
```

Delete a plugin config for a consumer by plugin config id:
```
err := gokong.NewClient(gokong.NewDefaultConfig()).Consumers().DeletePluginConfig("3958a860-ceac-4a6c-9bbb-ff8d69a585d2", "jwt", "bde04c3a-46bb-45c9-9006-e8af20d04342")
```

## Certificates
Create a Certificate ([for more information on the Certificate Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#certificate-object)):

```go
certificateRequest := &gokong.CertificateRequest{
  Cert: gokong.String("public key --- 123"),
  Key:  gokong.String("private key --- 456"),
}

createdCertificate, err := gokong.NewClient(gokong.NewDefaultConfig()).Certificates().Create(certificateRequest)
```

Get a Certificate by id:
```go
certificate, err := gokong.NewClient(gokong.NewDefaultConfig()).Certificates().GetById("0408cbd4-e856-4565-bc11-066326de9231")
```

List all certificates:
```go
certificates, err := gokong.NewClient(gokong.NewDefaultConfig()).Certificates().List()
```

Delete a Certificate:
```go
err := gokong.NewClient(gokong.NewDefaultConfig()).Certificates().DeleteById("db884cf2-9dd7-4e33-9ef5-628165076a42")
```

Update a Certificate:
```go
updateCertificateRequest := &gokong.CertificateRequest{
  Cert: gokong.String("public key --- 789"),
  Key:  gokong.String("private key --- 111"),
}

updatedCertificate, err := gokong.NewClient(gokong.NewDefaultConfig()).Certificates().UpdateById("1dc11281-30a6-4fb9-aec2-c6ff33445375", updateCertificateRequest)
```

# Routes

Create a Route ([for more information on the Sni Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#route-object)):
```go
serviceRequest := &gokong.ServiceRequest{
  Name:     gokong.String("service-name" + uuid.NewV4().String()),
  Protocol: gokong.String("http"),
  Host:     gokong.String("foo.com"),
}

client := gokong.NewClient(NewDefaultConfig())

createdService, err := client.Services().AddService(serviceRequest)

routeRequest := &RouteRequest{
  Protocols:    gokong.StringSlice([]string{"http"}),
  Methods:      gokong.StringSlice([]string{"GET"}),
  Hosts:        gokong.StringSlice([]string{"foo.com"}),
  StripPath:    gokong.Bool(true),
  PreserveHost: gokong.Bool(true),
  Service:      &RouteServiceObject{Id: *createdService.Id},
  Paths:        gokong.StringSlice([]string{"/bar"})
}

createdRoute, err := client.Routes().AddRoute(routeRequest)
```

Get a route by ID:
```go
result, err := gokong.NewClient(gokong.NewDefaultConfig()).Routes().GetRoute(createdRoute.Id)
```

Get all routes:
```go
result, err := gokong.NewClient(gokong.NewDefaultConfig()).Routes().GetRoutes(&RouteQueryString{})
```

Get routes from service ID or Name:
```go
result, err := gokong.NewClient(gokong.NewDefaultConfig()).Routes().GetRoutesFromServiceId(createdService.Id)
```

Update a route:
```go
routeRequest := &RouteRequest{
  Protocols:    gokong.StringSlice([]string{"http"}),
  Methods:      gokong.StringSlice([]string{"GET"}),
  Hosts:        gokong.StringSlice([]string{"foo.com"}),
  Paths:        gokong.StringSlice([]string{"/bar"}),
  StripPath:    gokong.Bool(true),
  PreserveHost: gokong.Bool(true),
  Service:      &RouteServiceObject{Id: *createdService.Id},
}

createdRoute, err := gokong.NewClient(gokong.NewDefaultConfig()).Routes().AddRoute(routeRequest)

routeRequest.Paths = gokong.StringSlice([]string{"/qux"})
updatedRoute, err := gokong.NewClient(gokong.NewDefaultConfig()).Routes().UpdateRoute(*createdRoute.Id, routeRequest)
```

Delete a route:
```go
client.Routes().DeleteRoute(createdRoute.Id)
```

# Services

Create an Service ([for more information on the Sni Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#service-object)):
```go
serviceRequest := &ServiceRequest{
		Name:     gokong.String("service-name-0"),
		Protocol: gokong.String("http"),
		Host:     gokong.String("foo.com"),
	}

	client := gokong.NewClient(gokong.NewDefaultConfig())

	createdService, err := client.Services().AddService(serviceRequest)
```

Get information about a service with the service ID or Name
```go
serviceRequest := &ServiceRequest{
		Name:     gokong.String("service-name-0"),
    Protocol: gokong.String("http"),
    Host:     gokong.String("foo.com")
	}

client := gokong.NewClient(gokong.NewDefaultConfig())

createdService, err := client.Services().AddService(serviceRequest)

resultFromId, err := client.Services().GetServiceById(createdService.Id)

resultFromName, err := client.Services().GetServiceByName(createdService.Id)
```

Get information about a service with the route ID
```go
result, err := gokong.NewClient(gokong.NewDefaultConfig()).Services().GetServiceRouteId(routeInformation.Id)
```

Get many services information
```go
result, err := gokong.NewClient(gokong.NewDefaultConfig()).Services().GetServices(&ServiceQueryString{
	Size: 500
	Offset: 300
})
```

Update a service with the service ID or Name
```go
serviceRequest := &ServiceRequest{
  Name:     gokong.String("service-name-0"),
  Protocol: gokong.String("http"),
  Host:     gokong.String("foo.com"),
}

client := gokong.NewClient(gokong.NewDefaultConfig())

createdService, err := client.Services().AddService(serviceRequest)

serviceRequest.Host = gokong.String("bar.io")
updatedService, err := client.Services().UpdateServiceById(createdService.Id, serviceRequest)
result, err := client.Services().GetServiceById(createdService.Id)
```

Update a service by the route ID
```go
serviceRequest := &ServiceRequest{
  Name:     gokong.String("service-name-0"),
  Protocol: gokong.String("http"),
  Host:     gokong.String("foo.com"),
}

client := gokong.NewClient(gokong.NewDefaultConfig())

createdService, err := client.Services().AddService(serviceRequest)

serviceRequest.Host = "bar.io"
updatedService, err := client.Services().UpdateServiceById(createdService.Id, serviceRequest)
result, err := client.Services().UpdateServicebyRouteId(routeInformation.Id)
```

Delete a service
```go
err = client.Services().DeleteServiceById(createdService.Id)
```

## SNIs
Create an SNI ([for more information on the Sni Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#sni-objects)):
```go
client := gokong.NewClient(gokong.NewDefaultConfig())

certificateRequest := &gokong.CertificateRequest{
  Cert: "public key --- 123",
  Key:  "private key --- 111",
}

certificate, err := client.Certificates().Create(certificateRequest)

snisRequest := &gokong.SnisRequest{
  Name:             "example.com",
  SslCertificateId: certificate.Id,
}

sni, err := client.Snis().Create(snisRequest)
```

Get an SNI by name:
```go
sni, err := client.Snis().GetByName("example.com")
```

List all SNIs:
```
snis, err := client.Snis().List()
```

Delete an SNI by name:
```go
err := client.Snis().DeleteByName("example.com")
```

Update an SNI by name:
```go
updateSniRequest := &gokong.SnisRequest{
  Name:             "example.com",
  SslCertificateId: "a9797703-3ae6-44a9-9f0a-4ebb5d7f301f",
}

updatedSni, err := client.Snis().UpdateByName("example.com", updateSniRequest)
```

## Upstreams
Create an Upstream ([for more information on the Upstream Fields see the Kong documentation](https://getkong.org/docs/0.13.x/admin-api/#upstream-objects)):
```go
upstreamRequest := &gokong.UpstreamRequest{
  Name: "test-upstream",
  Slots: 10,
}

createdUpstream, err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().Create(upstreamRequest)
```

Get an Upstream by id:
```go
upstream, err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().GetById("3705d962-caa8-4d0b-b291-4f0e85fe227a")
```

Get an Upstream by name:
```go
upstream, err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().GetByName("test-upstream")
```

List all Upstreams:
```go
upstreams, err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().List()
```

List all Upstreams with a filter:
```go
upstreams, err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().ListFiltered(&gokong.UpstreamFilter{Name:"test-upstream", Slots:10})
```

Delete an Upstream by id:
```go
err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().DeleteById("3a46b122-47ee-4c5d-b2de-49be84a672e6")
```

Delete an Upstream by name:
```go
err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().DeleteById("3a46b122-47ee-4c5d-b2de-49be84a672e6")
```

Delete an Upstream by id:
```go
err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().DeleteByName("test-upstream")
```

Update an Upstream by id:
```
updateUpstreamRequest := &gokong.UpstreamRequest{
  Name: "test-upstream",
  Slots: 10,
}

updatedUpstream, err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().UpdateById("3a46b122-47ee-4c5d-b2de-49be84a672e6", updateUpstreamRequest)
```

Update an Upstream by name:
```go
updateUpstreamRequest := &gokong.UpstreamRequest{
  Name: "test-upstream",
  Slots: 10,
}

updatedUpstream, err := gokong.NewClient(gokong.NewDefaultConfig()).Upstreams().UpdateByName("test-upstream", updateUpstreamRequest)
```

# Contributing
I would love to get contributions to the project so please feel free to submit a PR.  To setup your dev station you need go and docker installed.

Once you have cloned the repository the `make` command will build the code and run all of the tests.  If they all pass then you are good to go!

If when you run the make command you get the following error:
```
gofmt needs running on the following files:
```
Then all you need to do is run `make goimports` this will reformat all of the code (I know awesome)!!

