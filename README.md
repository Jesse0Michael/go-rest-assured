# GO REST ASSURED

[![Build](https://github.com/Jesse0Michael/go-rest-assured/workflows/Build/badge.svg)](https://github.com/Jesse0Michael/go-rest-assured/actions?query=branch%3Amaster) [![Coverage Status](https://coveralls.io/repos/github/Jesse0Michael/go-rest-assured/badge.svg?branch=master)](https://coveralls.io/github/Jesse0Michael/go-rest-assured?branch=master)

Go Rest Assured is a small service written in GO intended to be used to mock out REST API applications for testing. The concept is based on the [Rest Assured](http://rest-assured.io/) service written in Java and [other languages](https://github.com/artemave/REST-assured)

Go-Rest-Assured keeps track of the Assured Calls you have stubbed out and the Calls that have been made against the service with the following fields:

- Path
- StatusCode
- Method
- Response
- Headers
- Query
- Delay
- Callbacks

Set these fields as a _Given_ call through the client or a HTTP request to the service directly and they will be returned from the Go Rest Assured API when you hit the _When_ endpoint. The Calls you stub out are uniquely mapped with an identity of their Method and Path. If you stub multiple calls to the same Method and Path, the responses will cycle through your stubs based on the order they were created.

If loading callbacks from a JSON file, the call [unmarshaller](pkg/assured/call.go) will attempt to read the resource field as a relative file, or else a quoted string, or else just a byte slice.

To understand how rest assured is working behind the scenes, or to use rest assured as a standalone application you can run from a command line and use anywhere, or how to serve HTTPS traffic from rest assured, read the application [README](cmd/go-assured/README.md)

### Client

```go
import ("github.com/jesse0michael/go-rest-assured/v4/pkg/assured")

// Create and Serve a new Assured Client
client := assured.NewClientServe()
defer client.Close()
```

## Stubbing

```go
call := assured.Call{
  Path: "test/assured",
  StatusCode: 201,
  Method: "GET",
  Delay: 2,
}
// Stub out an assured call
client.Given(call)
```

_If your stubbed endpoint needs to return a different call on a subsequent request, then try stubbing that Method/Path again. The first time you intercept that endpoint the first call will be returned and then moved to the end of the list._

## Intercepting

To use your assured calls hit the following endpoint with the Method/Path that was used to stub the call 

```go
// Get the URL of the client ex: 'http://localhost:11011/when'
testServer := client.URL()
```

Go-Rest-Assured will return `404 NotFound` error response when a matching stub isn't found

As requests come in, the will be stored

## Callbacks
To have the mock server programmatically make a callback to a specified target, use the Callback field

```go
call := assured.Call{
  Path: "test/assured",
  StatusCode: 201,
  Method: "POST",
  Response: []byte(`{"holler_back":true}`),
  Callbacks: []assured.Callback{
    assured.Callback{
      Method: "POST",
      Target: "http://localhost:8080/hit/me",
      Response: []byte(`holla!!`),
    },
  },
}
// Stub out an assured call with callbacks
client.Given(call)
```

_You cannot clear out an individual callback when using the assured.Client, but you can `ClearAll()`_

## Verifying

To verify the calls made against your go-rest-assured service, use the Verify function.

This function returns a list of calls made against the matching Method/Path

```go
// Get a []*assured.Call for a Method and Path
calls := client.Verify("GET", "test/assured")
```

## Clearing

To clear out the stubbed and made calls for a specific Method/Path, use Clear(method, path)

To clear out all stubbed calls on the server, use ClearAll()

```go
// Clears calls for a Method and Path
client.Clear("GET", "test/assured")

// Clears all calls
client.ClearAll()
```
