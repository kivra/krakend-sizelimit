# krakend-sizelimit

The `sizelimit` middleware can be used to limit the request size of incoming
requests.

## Installation

To install `sizelimit` from GitHub:

    go get -u github.com/kivra/krakend-sizelimit@<commit hash>

Then add `sizelimit` to the KrakenD [`handler_factory`](https://github.com/krakendio/krakend-ce/blob/master/handler_factory.go)
chain:

```go
handlerFactory = sizelimit.HandlerFactory(handlerFactory)
```

## Usage

Add `sizelimit` to the endpoint's `extra_config` and define the maximum request
body size. If no unit is specified, `max_size` is assumed to have unit `bytes`.
Other units can be specified explictly: `B` (bytes, same as no unit), `kB`,
`MB`, `GB`, `TB`.

```json
"endpoints": [
  {
    "endpoint": "/test",
    "extra_config": {
      "kivra/sizelimit": {
        "max_size": "10MB"
      }
    },
    "backend": [ "..." ]
  }
]
```
