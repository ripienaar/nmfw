# NATS Micro Service Framework

This is a framework powered by NATS Micro that generates a microservice from a Protobuf service definition. The goal is
to go from `service.proto` to running a service in 5 minutes after supplying transport agnostic business logic requiring
little or no NATS knowledge.

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
 * Timeouts are propagated from Client to Service

## Status

This is very early days, it was written primarily to see if it's possible to create a [gRPC](https://grpc.io/) like
service framework using the NATS Micro feature.

At present each function in the service runs as an isolated service, this is not ideal and this
feedback is being incorporated in the NATS service design to support multiple handlers per single
service.

Once the above feedback is implemented this plugin will get a major update to be deployed in that manner
which will be more efficient and easier to manage in reality.

## Example

### Generating types, services and tools

Given the proto file:

```protobuf
syntax = "proto3";

package calc;

option go_package = "github.com/ripienaar/nmfw/example/service";

message AddRequest {
  repeated float Values = 1;
}

message AverageRequest {
  repeated float Values = 1;
}

message ExpressionRequest {
  string Expression = 1;
}

message CalcResponse {
  string Operation = 1;
  float Result = 2;
}

service Calc {
  // Calculates the average of a series of numbers
  rpc Average(AverageRequest) returns (CalcResponse) {}

  // Calculates the sum of a series of numbers
  rpc Add(AddRequest) returns (CalcResponse) {}

  // Compiles and executes a expr expression, expression must return a number
  rpc Expression(ExpressionRequest) returns (CalcResponse) {}
}
```

I suggest using the `ripienaar/nmfw` docker container that holds the dependencies already to run `protoc`

```
$ docker run --ti --rm -v `pwd`/go/src ripienaar/nmfw:latest
# export VERSION=0.0.2
# export IMPL=github.com/ripienaar/nmfw/example/impl
# export TARGET=service
# export PROTO=service.proto
# protoc -I=`pwd`/ \
    --go_out="/go/src/${TARGET?}" \
    --go_opt=paths=source_relative \
    --go-nmfw_opt="paths=source_relative,version=${VERSION?},impl=${IMPL?}" \
    --go-nmfw_out="/go/src/${TARGET}" \
    "${PROTO?}"
# exit
```

 * Generates data types using standard `protoc-gen-go` plugin info `/go/src/service`
 * Generates a service called `CalcService` that binds to the Tech Preview [NATS Micro](https://github.com/nats-io/nats-architecture-and-design/blob/main/adr/ADR-32.md) system in `nats.go` into `/go/src/service`
 * Generates a command called `calc` that runs the service with version `0.0.2` in `/go/src/service/calc`
 * Generates a client called `CalcClient` that can interact with the service into `/go/src/service`
 * Requires the user to create implementation methods in `github.com/ripienaar/nmfw/example/impl`

**NOTE** This is how the `example` directory in this repository was created

### Implementation

In this case we said our implementation will be in `github.com/ripienaar/nmfw/example/impl` and you will be shown which
functions to create there.

Here's an example for teh `Add()` function:

```golang
func AddHandler(ctx context.Context, req service.AddRequest) (*service.CalcResponse, error) {
	resp := service.CalcResponse{Operation: "add"}
	if len(req.Values) == 0 {
		return &resp, nil
	}

	for _, v := range req.Values {
		resp.Result += v
	}

	return &resp, nil
}
```

The context will have a deadline set which is propagated from the timeout supplied by the client.

The context has a helper value that provides access to prepared loggers, nats connection and theo original `micro` request. 
Use this to log to the service and access things like JetStream without requiring new connections per invocation.

Here's a part of the implementation that defines the interface and then accesses it.  We have to define the interface here to 
avoid cyclic imports. You can define just the helper method you actually need here.

```golang
type request interface {
    Logger() *logrus.Entry
    Conn() *nats.Conn
    Request() *micro.Request
}

func ExpressionHandler(ctx context.Context, req service.ExpressionRequest) (*service.CalcResponse, error) {
    helper := ctx.Value("nmfw").(request)
	log := helper.Logger()
	
    log.Infof("Calculating expression %s", req.Expression)
    
    // ...
}
```

### Running 

The service host uses a NATS Context for connection properties, create it using `nats context` and then run the service after
compiling it: in `service/calc`.

```nohighlight
$ cd service/calc
$ go build
$ ./calc --help
usage: calc [<flags>] <command> [<args> ...]

Micro Service powered by NATS Micro

Commands:
  run  Runs the service

Global Flags:
  --help     Show context-sensitive help
  --version  Show application version.
  --debug    Log at debug level ($DEBUG)
  
$ ./calc run --help
usage: calc run [<flags>]

Runs the service

Flags:
  --context="MICRO"  NATS Context to use for connection ($CONTEXT)
  --port=PORT        Prometheus port for statistics ($PORT)
  --max-recon=60     Maximum reconnection attempts ($MAX_RECON)
```

The command takes some flags for connection properties and logging, future versions will include tools to help administrators
discover and introspect running instances.  We start Prometheus metrics on port `8222`

```nohighlight
$ ./calc run --context AUTH_CALLOUT --port 8222
{"level":"info","msg":"Starting Prometheus listener on :8222/metrics","service":"calc","time":"2022-12-31T14:19:52+01:00","version":"0.0.2"}
{"level":"info","msg":"Connected to nats://127.0.0.1:10222","service":"calc","time":"2022-12-31T14:19:52+01:00","version":"0.0.2"}
{"function":"Expression","level":"info","msg":"Started on subject nmfw.calc.Expression with ID poB9M20P4GBUeK49wkCurR","service":"calc","time":"2022-12-31T14:19:52+01:00","version":"0.0.2"}
{"function":"Average","level":"info","msg":"Started on subject nmfw.calc.Average with ID poB9M20P4GBUeK49wkCuxp","service":"calc","time":"2022-12-31T14:19:52+01:00","version":"0.0.2"}
{"function":"Add","level":"info","msg":"Started on subject nmfw.calc.Add with ID poB9M20P4GBUeK49wkCuud","service":"calc","time":"2022-12-31T14:19:52+01:00","version":"0.0.2"}
```

### Client

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
 * ~~Think about timeout, some functions have different timeouts than others, how to handle?~~ Propagated using the `Nmfw-Deadline` header and passed to handlers as a `context.Context`.
 * Include the proto schema and expose over the `micro` schemas feature

## Contact

R.I.Pienaar / [@ripienaar](https://twitter.com/ripienaar) / [@ripienaar@devco.social](https://devco.social/@ripienaar) / ripienaar @ NATS Slack
