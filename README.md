# ASSURED

[![Build](https://github.com/Jesse0Michael/go-rest-assured/workflows/Build/badge.svg)](https://github.com/Jesse0Michael/go-rest-assured/actions?query=branch%3Amain) [![Coverage Status](https://coveralls.io/repos/github/Jesse0Michael/go-rest-assured/badge.svg?branch=main)](https://coveralls.io/github/Jesse0Michael/go-rest-assured?branch=main)

Assured is a GO service used to mock out REST API applications for testing. It keeps track of the calls you have stubbed out and the request that have been made against the service. REST responses can be stubbed with the following fields:

- Path
- StatusCode
- Method
- Response
- Headers
- Delay
- Callbacks

Set these fields as a _Given_ call through the client or a HTTP request to the service directly and they will be returned from the Assured Server when you hit the matching stubbed call. The Calls you stub out are uniquely mapped with an identity of their Method and Path. If you stub multiple calls to the same Method and Path, the responses will cycle through your stubs based on the order they were created.

If loading calls from a JSON file, the call [unmarshaller](pkg/assured/call.go) will attempt to read the resource field as a relative file, or else a quoted string, or else just a byte slice.

To understand how assured is working behind the scenes, or to use assured as a standalone application you can run from a command line and use anywhere, or how to serve HTTPS traffic from assured, read the application [README](cmd/assured/README.md)

### Client & Server

```go
import ("github.com/jesse0michael/go-rest-assured/v5/pkg/assured")

// Create an Assured Client and Server
a, err := assured.ServeAssured(ctx)
defer a.Close()
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
a.Given(ctx, call)
```

_If your stubbed endpoint needs to return a different call on a subsequent request, then try stubbing that Method/Path again. The first time you intercept that endpoint the first call will be returned and then moved to the end of the list._

## Intercepting

To use your assured calls hit the following endpoint with the Method/Path that was used to stub the call 

```go
// Get the URL of the server ex: 'http://localhost:11011/when'
testServer := a.URL()
```

Assured will return `404 NotFound` error response when a matching stub isn't found

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
a.Given(ctx, call)
```
  

## Verifying

To verify the calls made against your assured service, use the Verify function.

This function returns a list of calls made against the matching Method/Path

```go
// Get a []*assured.Call for a Method and Path
calls := a.Verify(ctx, "GET", "test/assured")
```

## Clearing

To clear out the stubbed and made calls for a specific Method/Path, use Clear(method, path)

To clear out all stubbed calls on the server, use ClearAll()

```go
// Clears calls for a Method and Path
a.Clear(ctx, "GET", "test/assured")

// Clears all calls
a.ClearAll(ctx)
```
