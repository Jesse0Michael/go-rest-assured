# GO REST ASSURED - CMD

Go Rest Assured can be used from the command line to spin up a mock rest api. The rest assured application will take in some arguments to configure configure a server to mock rest calls.

## Running

1. `go get github.com/jesse0michael/go-rest-assured/cmd/go-assured`

```bash
Usage of go-assured:
  -logger string
    	a file to send logs. default logs to STDOUT.
  -port int
    	a port to listen on. default automatically assigns a port.
  -preload string
    	a file to parse preloaded calls from.
  -track
    	a flag to enable the storing of calls made to the service. (default true)
```

To load in a default set of stubbed endpoints from a file, follow the [Preload API Reference](preload_reference.md) guide.

## Stubbing

To stub out an assured call hit the following endpoint
`/given/{path:.*}`

The HTTP Method you use will be stored in the Assured Call

The Request Body, if present, will be stored in the Assured Call

The stored Status Code will be `200 OK` unless you specify a `"Assured-Status": "[0-9]+"` HTTP Header

You can also set a response delay with the HTTP Header `Assured-Delay` with a number of seconds


_If your stubbed endpoint needs to return a different call on a subsequent request, then try stubbing that Method/Path again. The first time you intercept that endpoint the first call will be returned and then moved to the end of the list._

## Intercepting

To use your assured calls hit the following endpoint with the Method/Path that was used to stub the call `/when/{path:.*}`

Go-Rest-Assured will return `404 NotFound` error response when a matching stub isn't found

As requests come in, the will be stored

## Callbacks

To include callbacks from Go-Rest-Assured when a stubbed endpoint is hit, create them by hitting the endpoint `/callbacks`
To create a callbacks you must include the HTTP header `Assured-Callback-Target` with the specified endpoint you want your callbacks to be sent to
You must also include the HTTP header `Assured-Callback-Key` with a key with the call to the `/callbacks` endpoint as well as the `/given/{path:.*}` endpoint that for the stubbed call you want the callback to be associated with
You can also set a callback delay with the HTTP Header `Assured-Callback-Delay` with a number of seconds

## Verifying

To verify the calls made against your go-rest-assured service, use the endpoint `/verify/{path:.*}`

This endpoint returns a list of assured calls made against the matching Method/Path

```
[
  {
    "path": "test/assured",
    "method": "GET",
    "status_code": 200,
    "delay": 0,
    "response": "eyJhc3N1cmVkIjogdHJ1ZX0=",
    "headers": {
      "Content-Length": "17",
      "User-Agent": "Go-http-client/1.1",
    },
    "query": {
      "assured": "max"
    }
  },
  {
    "path": "test/assured",
    "method": "GET",
    "status_code": 409,
    "delay": 0,
    "response": "ZXJyb3I=",
    "headers": {
      "Content-Length": "5",
      "User-Agent": "Go-http-client/1.1",
    }
  }
]

```

## Clearing

To clear out the stubbed and made calls for a specific Method/Path, use the endpoint `/clear/{path:.*}`
_Including the HTTP Header `Assured-Callback-Key` will clear all callbacks associated with that key (independent of path)_

To clear out all stubbed calls on the server, use the endpoint `/clear`
