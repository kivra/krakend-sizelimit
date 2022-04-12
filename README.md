# sizelimit

The `sizelimit` middleware can be used to limit the request size of incoming
requests.

## Installation

To install `sizelimit` from GitHub:

    go get -u github.com/kivra/krakend-sizelimit@<commit hash>

Then add `sizelimit` to the KrakenD [`handler_factory`](https://github.com/devopsfaith/krakend-ce/blob/master/handler_factory.go)
chain:

```go
handlerFactory = sizelimit.HandlerFactory(handlerFactory)
```

## Usage

The `sizelimit` middleware can be added to an endpoint's `extra_config` and allows
to define the maximum request body size in bytes:

```json
"endpoints": [
  {
    "endpoint": "/test",
    "method": "POST",
    "extra_config": {
      "kivra/sizelimit": {
        "max_bytes": 10
      }
    },
    "backend": [ "..." ]
  }
]
```
