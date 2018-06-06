# GO REST ASSURED
[![CircleCI](https://circleci.com/gh/Jesse0Michael/go-rest-assured.svg?style=svg&circle-token=afd5de8a46297d388679dcfc404d4bcc4eceab7a)](https://circleci.com/gh/Jesse0Michael/go-rest-assured) [![Coverage Status](https://coveralls.io/repos/github/Jesse0Michael/go-rest-assured/badge.svg?branch=master)](https://coveralls.io/github/Jesse0Michael/go-rest-assured?branch=master)
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

Set these fields as a *Given* call through the client or a HTTP request to the service  directly and they will be returned from the Go Rest Assured API when you hit the *When* endpoint. The Calls you stub out are uniquely mapped with an identity of their Method and Path. If you stub multiple calls to the same Method and Path, the responses will cycle through your stubs based on the order they were created.

If loading callbacks from a JSON file, the call unmarshaller will attempt to read the resource field as a relative file, or else a quoted string, or else just a byte slice.

## Running

### Standalone 
Please have [GO](https://golang.org/) and [Glide](https://github.com/Masterminds/glide) installed on your machine

1. `go get github.com/jesse0michael/go-rest-assured`
2. `make install-deps`
2. `make build`
3. `bin/go-rest-assured`

```bash
Usage of bin/go-rest-assured:
  -logger string
    	a file to send logs. default logs to STDOUT.
  -port int
    	a port to listen on. default automatically assigns a port.
  -preload string
    	a file to parse preloaded calls from.
  -track
    	a flag to enable the storing of calls made to the service. (default true)
```

### Client
```go
import ("github.com/jesse0michael/go-rest-assured/assured")

// Create and run a new Assured Client
client := assured.NewDefaultClient()
defer client.Close()
```

## Stubbing
To stub out an assured call hit the following endpoint
`/given/{path:.*}`

The HTTP Method you use will be stored in the Assured Call

The Request Body, if present, will be stored in the Assured Call

The stored Status Code will be `200 OK` unless you specify a `"Assured-Status": "[0-9]+"` HTTP Header

You can also set a response delay with the HTTP Header `Assured-Delay` with a number of seconds

Or..

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

*If your stubbed endpoint needs to return a different call on a subsequent request, then try stubbing that Method/Path again. The first time you intercept that endpoint the first call will be returned and then moved to the end of the list.*

## Intercepting
To use your assured calls hit the following endpoint with the Method/Path that was used to stub the call `/when/{path:.*}`

Or..

```go
// Get the URL of the client ex: 'http://localhost:11011/when'
testServer := client.URL()
```

Go-Rest-Assured will return `404 NotFound` error response when a matching stub isn't found

As requests come in, the will be stored

## Callbacks
To include callbacks from Go-Rest-Assured when a stubbed endpoint is hit, create them by hitting the endpoint `/callbacks`
To create a callbacks you must include the HTTP header `Assured-Callback-Target` with the specified endpoint you want your callbacks to be sent to
You must also include the HTTP header `Assured-Callback-Key` with a key with the call to the `/callbacks` endpoint as well as the `/given/{path:.*}` endpoint that for the stubbed call you want the callback to be associated with
You can also set a callback delay with the HTTP Header `Assured-Callback-Delay` with a number of seconds

Or...

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
*You cannot clear out an individual callback when using the assured.Client, but you can `ClearAll()`*

## Verifying
To verify the calls made against your go-rest-assured service, use the endpoint `/verify/{path:.*}`

This endpoint returns a list of assured calls made against the matching Method/Path

Or..

```go
// Get a []*assured.Call for a Method and Path
calls := client.Verify("GET", "test/assured")
```


## Clearing
To clear out the stubbed and made calls for a specific Method/Path, use the endpoint `/clear/{path:.*}`
 *Including the HTTP Header `Assured-Callback-Key` will clear all callbacks associated with that key (independent of path)*

To clear out all stubbed calls on the server, use the endpoint `/clear`

Or..

``` go
// Clears calls for a Method and Path
client.Clear("GET", "test/assured")

// Crears all calls
client.ClearAll()
```

