# krakend-sizelimit

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

Add `sizelimit` to the endpoint's `extra_config` and define the maximum request
body size in bytes:

```json
"endpoints": [
  {
    "endpoint": "/test",
    "extra_config": {
      "kivra/sizelimit": {
        "max_bytes": 10
      }
    },
    "backend": [ "..." ]
  }
]
```
