# Sensu Go Alerta Handler
[![Bonsai Asset Badge](https://img.shields.io/badge/CHANGEME-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/CHANGEME/CHANGEME) [![TravisCI Build Status](https://travis-ci.org/alerta/sensu-alerta-handler.svg?branch=master)](https://travis-ci.org/alerta/sensu-alerta-handler)

Forward Sensu events to Alerta.

## Installation

Download the latest version of the sensu-alerta-handler from [releases][1],
or create an executable script from this source.

From the local path of the sensu-alerta-handler repository:

```
go build -o /usr/local/bin/sensu-alerta-handler main.go
```

## Configuration

Example Sensu Go handler definition:

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "alerta"
    },
    "spec": {
        "...": "..."
    }
}
```

## Usage Examples

Help:

```
The Sensu Go Alerta handler for event forwarding

Usage:
  sensu-alerta-handler [flags]

Flags:
  -K, --api-key string        API key for authenticated access.
      --endpoint-url string   API endpoint URL.
  -h, --help                  help for sensu-alerta-handler
```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/alerta/sensu-alerta-handler/releases
