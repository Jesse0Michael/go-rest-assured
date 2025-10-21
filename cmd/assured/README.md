# ASSURED - CMD

Assured can be used from the command line to spin up a mock rest api. The assured application will take in some arguments to configure configure a server to mock rest calls.

## Running

1. `go get github.com/jesse0michael/go-rest-assured/cmd/assured`

```
Usage of assured:
  -host string
        a host to use in the client's url. (default "localhost")
  -port int
        a port to listen on. default automatically assigns a port.
  -preload string
        a file to parse preloaded calls from.
  -tlsCert string
        location of tls cert for serving https traffic. tlsKey also required, if specified.
  -tlsKey string
        location of tls key for serving https traffic. tlsCert also required, if specified
  -track
        a flag to enable the storing of calls made to the service. (default true)
```

To load in a default set of stubbed endpoints from a file, follow the [Preload API Reference](preload_reference.md) guide.

You can specify a TLS cert/key to mock out HTTPS traffic using [mkcert](https://github.com/FiloSottile/mkcert) self signed certs and mock HTTPS traffic.

## Stubbing

To stub out an assured call hit the following endpoint
`/assured/given`
You must include a JSON body with the following fields to create your stubbed call

```json
{
  "path": "test/assured",
  "status_code": 201,
  "method": "POST",
  "response": "ASSURED ACCEPTED",
  "headers": {
    "Content-Type": "application/json"
  },
  "delay": 2,
  "callbacks": [
    {
      "target": "http://example.com/callback",
      "delay": 1
    }
  ]
}
```
The following fields are available to set on your stubbed call
- Path: The path to match on the assured server
- StatusCode: The HTTP status code to return
- Method: The HTTP method to match
- Response: The response body to return
- Headers: The headers to include in the response
- Delay: The delay before returning the response
- Callbacks: The callbacks to invoke when the stub is hit

_If your stubbed endpoint needs to return a different call on a subsequent request, then try stubbing that Method/Path again. The first time you intercept that endpoint the first call will be returned and then moved to the end of the list._

## Intercepting

To use your assured calls hit the any matched method:path combination previously stubbed out

Assured will return `404 NotFound` error response when a matching stub isn't found

As requests come in, they will be stored

## Verifying

To verify the calls made against your go-rest-assured server, use the endpoint `/assured/verify`
with the request body

```json
{
  "path": "test/assured",
  "method": "GET"
}
```

This endpoint returns a list of assured calls made against the matching Method/Path

``` json
[
  {
    "path": "test/assured",
    "method": "GET",
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

To clear out the stubbed and made calls for a specific Method/Path, use the endpoint POST `assured/clear`
with the request body:

```json
{
  "path": "test/assured",
  "method": "GET"
}
```

To clear out all stubbed calls on the server, use the endpoint `/clearall`
