# Sensu Go Alerta Handler

[![Bonsai Asset Badge](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/alerta/sensu-alerta-handler) [![build-test](https://github.com/alerta/sensu-alerta-handler/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/alerta/sensu-alerta-handler/actions/workflows/test.yml)

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

**alerta-handler.json**

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "alerta"
    },
    "spec": {
        "type": "pipe",
        "command": "sensu-alerta-handler --endpoint-url https://alerta.example.com/api",
        "env_vars": [
            "ALERTA_API_KEY=G25k9JR2yoZIcHROQGS477nk_Riw4CIghFC6j9NE"
        ],

        "timeout": 30,
        "filters": [
            "is_incident"
        ]
    }
}
```

Create the handler resource:

    $ sensuctl create -f alerta-handler.json

Example Sensu Go check definition:

```
{
    "api_version": "core/v2",
    "type": "CheckConfig",
    "metadata": {
        "namespace": "default",
        "name": "dummy-app-healthz"
    },
    "spec": {
        "command": "check-http -u http://localhost:8080/healthz",
        "subscriptions":[
            "dummy"
        ],
        "publish": true,
        "interval": 10,
        "handlers": [
            "alerta"
        ]
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
  -K, --api-key string        API key for authenticated access
      --endpoint-url string   API endpoint URL (default "http://localhost:8080")
  -E, --environment string    Environment eg. Production, Development (default "Entity Namespace")
  -h, --help                  help for sensu-alerta-handler
```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/alerta/sensu-alerta-handler/releases
