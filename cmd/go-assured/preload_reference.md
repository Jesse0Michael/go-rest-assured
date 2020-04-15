## Preload Rest Assured Endpoints

To stub rest assured endpoints with a json file, pass a JSON file to the `-preload` argument that follows this specification:

## Example

```json
{
  "calls": [
    {
      "path": "test/assured",
      "method": "GET",
      "status_code": 201,
      "delay": 0,
      "response": "testdata/assured.json",
      "headers": {
        "Content-Length": "17",
        "User-Agent": "Go-http-client/1.1",
        "Accept-Encoding": "gzip"
      }
    },
    {
      "path": "test/assured",
      "method": "GET",
      "status_code": 200,
      "delay": 0,
      "response": "testdata/image.jpg",
      "headers": {
        "Content-Length": "56000",
        "Content-Type": "image/jpeg",
        "User-Agent": "Go-http-client/1.1"
      }
    },
    {
      "path": "teapot/assured",
      "method": "POST",
      "status_code": 418,
      "delay": 2,
      "headers": {
        "Content-Length": "0",
        "User-Agent": "Go-http-client/1.1",
        "Accept-Encoding": "gzip"
      },
      "callbacks": [
        {
          "target": "http://localhost:9000/when/capture",
          "method": "PUT",
          "delay": 5,
          "response": "Therefore do not worry about tomorrow, for tomorrow will worry about itself. Each day has enough trouble of its own."
        }
      ]
    }
  ]
}
```

### calls

**[object array]** The rest assured calls loaded into the go rest assured application

```json
{
    "calls": [
        ...
    ]
}
```

### calls[x].path
**[string]** The http path to the endpoints. 

```json
{
    "path": "test/assured",
    ...
}
```
*When call this path, to receive the stubbed response you need to include the `/when/` path prefix. e.g. `http://localhost:8888/when/test/assured`*

### calls[x].method
**[string]** The http method to the endpoints. Defaults to "GET".

```json
{
    ...
    "method": "POST",
    ...
}
```

### calls[x].status_code
**[int]** The http status code to respond with. Defaults to 200 OK.

```json
{
    ...
    "status_code": 201,
    ...
}
```

### calls[x].response
**[string]** The http response body to respond with using a custom and complex JSON unmarshall function. Unmarshalling will first check if the data is a local file path that can be read. Else it will check if the data is stringified JSON and un-stringify the data to use. Else it will just use the []byte. Optional.

```json
{
    ...
    "response": "responses/success.json",
    ...
},
{
    ...
    "response": "{\"happy\": true}",
    ...
},
{
    ...
    "response": "string cheese",
    ...
}
```

### calls[x].headers
**[object]** The http headers to include with the response. Keys and values must be strings. 

```json
{
    ...
    "headers": {
      "Content-Length": "17",
      "User-Agent": "Go-http-client/1.1"
    },
    ...
}
```

### calls[x].callbacks
**[object array]** Specified callbacks to be made by the go rest assured application when an endpoint is hit with specified parameters. Optional.

```json
{
    ...
    "callbacks": [
      ...
    ]
}
```

### calls[x].callbacks[x].target
**[string]** The http target too hit with the callback. Required

```json
    {
        "target": "http://localhost:9000/when/capture",
        ...
    }
```

### calls[x].callbacks[x].method
**[string]** The http method to use with the callback. Defaults to "GET".

```json
    {
        "method": "PUT",
        ...
    }
```

### calls[x].callbacks[x].response
**[string]** The http response body to respond with in the callback. uses the same custom and complex JSON unmarshall function as the endpoint's response. Unmarshalling will first check if the data is a local file path that can be read. Else it will check if the data is stringified JSON and un-stringify the data to use. Else it will just use the []byte. Optional.

```json
    {
        ...
        "response": "responses/success.json",
        ...
    },
    {
        ...
        "response": "{\"happy\": true}",
        ...
    },
    {
        ...
        "response": "string cheese",
        ...
    }       
```

### calls[x].callbacks[x].headers
**[object]** The http headers to include with the callback. Keys and values must be strings. 

```json
    {
        ...
        "headers": {
        "Content-Length": "17",
        "User-Agent": "Go-http-client/1.1"
        },
        ...
    }
```

### calls[x].callbacks[x].delay
**[int]** A synthetic delay, in seconds, to delay the callback from triggering. Optional. 

```json
    {
        ...
        "delay": 2
    }
```


---

Follow the go rest assured application [README.md](README.md) for instructions on how to interact with your stub server
