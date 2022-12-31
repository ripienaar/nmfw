# NATS Micro Service Framework

This is a framework powered by NATS Micro that generates a microservice from a Protobuf service definition.

Using NATS as a Microservice transport has significant advantages over HTTP but there has not really
been a good effort made to leverage those features into a gRPC like framework. 

This is a exploration of how such a framework might look, it would leverage NATS to:

 * Provide Load Balancing for horizontal and vertical scale out
 * Provides GSLB style failover and fallback cross region
 * Supports extending a centralized service to the edge using Leafnodes
 * Supports scaling individual functions within the service differently from others. You could run a movie encoder function on expensive GPU equipped machines while other supporting functions can run on smaller instances
 * No service discovery, registries etc needed as NATS handles that in real time
 * (Eventually) re-using handlers between real time RPC based use cases and Job Queue style use cases

## Features

 * Requires standard go types generated using `protoc-gen-go`
 * Creates a service type that hosts the microservice
 * Creates a CLI tool that runs the microservice
 * Microservice can optionally export Prometheus metrics
 * Creates a Client class that can interact with the service
 * Service handlers are pure business logic and transport agnostic

## Status

This is very early days, it was written primarily to see if it's possible to create a [gRPC](https://grpc.io/) like
service framework using the NATS Micro feature.

At present each function in the service runs as an isolated service, this is not ideal and this
feedback is being incorporated in the NATS service design to support multiple handlers per single
service.

Once the above feedback is implemented this plugin will get a major update to be deployed in that manner
which will be more efficient and easier to manage in reality.

## Example

Given the proto file:

```protobuf
syntax = "proto3";

package calc;

option go_package = "github.com/ripienaar/nmfw/example/service";

message AverageRequest {
    repeated float Values = 1;
}

message CalcResponse {
  string Operation =1;
  float Result=2;
}

service Calc {
  // Calculates the average of a series of numbers
  rpc Average(AverageRequest) returns (CalcResponse) {}
}
```

Calling `protoc` will build a number of files that implement the go microservice:

**NOTE**: You can use the `ripienaar/nmfw` docker container that holds the dependencies already

```nohighlight
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest     
$ go install github.com/ripienaar/nmfw/protoc-gen-go-nmfw@latest
$ protoc -I=`pwd`/ \
    --go_out=/Users/rip/go/src \
    --go-nmfw_opt=version=0.0.2,impl=github.com/ripienaar/nmfw/example/impl \
    --go-nmfw_out=/Users/rip/go/src \
    service.proto
```

 * Generates data types using standard `protoc-gen-go` plugin
 * Generates a service called `CalcService` that binds to the Tech Preview [NATS Micro](https://github.com/nats-io/nats-architecture-and-design/blob/main/adr/ADR-32.md) system in `nats.go`
 * Generates a command called `calc` that runs the service
 * Generates a client called `CalcClient` that can interact with the service

Once generated the client can be used to call the service:

```golang
nc, err := natscontext.Connect(contextName)
if err != nil {
    return err
}

c := service.NewCalcClient(nc, 5*time.Second)

res, err := c.Average(context.Background(), service.AverageRequest{Values: []float32{1.1, 2, 3.5}})
if err != nil {
    return err
}

fmt.Println(res.Result)
```

## Limitations

There are some limitations at present given the young age of this project:

 * Services should be implemented in a single `*.proto` file, it may import other files
 * Some effort has been made to support multiple services in the single proto file but there are many edge cases
 * No streaming responses are supported yet

## TODO

 * Support more `.proto` file behaviors
 * Support streaming responses
 * Pass a context to the handlers that include logger and nats connection
 * More observability, possibly propagate tracing headers
 * Generate a Dockerfile to host the service
 * Generate `Makefile` or similar to rebuild the generated code and containers
 * One `micro` per Service
 * Think about timeout, some functions have different timeouts than others, how to handle?
 * Include the proto schema and expose over the `micro` schemas feature

## Contact

R.I.Pienaar / [@ripienaar](https://twitter.com/ripienaar) / [@ripienaar@devco.social](https://devco.social/@ripienaar) / ripienaar @ NATS Slack
