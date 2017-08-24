# GO REST ASSURED
Go Rest Assured is a small service written in GO intended to be used to mock out REST API applications for testing. The concept is based on the [Rest Assured](http://rest-assured.io/) service written in Java and [other languages](https://github.com/artemave/REST-assured)

Go-Rest-Assured keeps track of the Assured Calls you have stubbed out and the Calls that have been made against the service with the following fields:

- Path
- StatusCode
- Method
- Response

## Running

### Standalone 

1. `go get github.com/jesse0michael/go-rest-assured`
2. `make build`
3. `bin/go-rest-assured`

### Client
```go
import ("github.com/jesse0michael/go-rest-assured/assured")

// Create and run a new Assured Client
client := assured.NewClient()
client.Run()
defer client.Close()
```

## Stubbing
To stub out an assured call hit the following endpoint
`/given/{path:.*}`

The HTTP Method you use will be stored in the Assured Call

The Request Body, if present, will be stored in the Assured Call

The stored Status Code will be `200 OK` unless you specify a `"Assured-Status": "[0-9]+"` HTTP Header

Or..

```go
call := assured.Call{
  Path: "test/assured",
  StatusCode: 201,
  Method: "GET",
}
// Stub out an assured call
client.Given(call)
```

*If your stubbed endpoint needs to return a different call on a subsequent request, then try stubbing that Method/Path again. The first time you intercept that endpoint the first call will be returned and then moved to the end of the list.*

## Intercepting
To use your assured calls hit the following endpoint with the Method/Path that was used to stub the call `/when/{path:.*}`

Or..

```go
// Get the URL of the client ex: 'localhost:11011/when'
testServer := client.URL()
```

Go-Rest-Assured will return `404 NotFound` error response when a matching stub isn't found

As requests come in, the will be stored


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

To clear out all stubbed calls on the server, use the endpoint `/clear`

Or..

``` go
// Clears calls for a Method and Path
client.Clear("GET", "test/assured")

// Crears all calls
client.ClearAll()
```

