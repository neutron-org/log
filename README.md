# Neutron Logger

This repository contains a helper zap.Logger constructor that injects a mandatory "context" field to the logger fields. The context field is expected to carry information about the execution context the logger is used in to ease logs browsing and indexing.

## Configuration

### General configuration

Basically, the logger is configured with the [zap.NewProductionLogger](https://github.com/uber-go/zap/blob/d6ce3b9b283401bc6cf975de6ac3e5ed5aec5341/config.go#L115) with modified logs timestamps as RFC3339-formatted string with nanosecond precision.

For the default logger configuration the logging level is set using an env variable `LOGGER_LEVEL`. According to the [zap.Logger docs](https://github.com/uber-go/zap/blob/d6ce3b9b283401bc6cf975de6ac3e5ed5aec5341/level.go#L28), the expected `LOGGER_LEVEL` values are (the `dpanic`/`DPANIC` level is omitted, higher levels are more important):
- `debug` or `DEBUG` — typically voluminous logs that are usually disabled in production;
- `info` or `INFO` — the default logging priority;
- `warn` or `WARN` — logs that are more important than Info, but don't need individual human review;
- `error` or `ERROR` — high-priority logs. If an application is running smoothly, it shouldn't generate any error-level logs;
- `panic` or `PANIC` — logs a message, then panics;
- `fatal` or `FATAL` — logs a message, then calls os.Exit(1).

Although the default cfg is made for logger instantiation and usage simplicity, a user can still configure the logger the way the [zap.Config](https://github.com/uber-go/zap/blob/d6ce3b9b283401bc6cf975de6ac3e5ed5aec5341/config.go#L45) allows it by assigning a path to a config file to the `LOGGER_CFG_PATH` env variable. Supported extensions for the config file are `.json`, `.yml` and `.yaml`. Configuration made by the config file overwrites the `LOGGER_LEVEL` and should contain comprehensive configuration (this is why the example file below is so detailed although it doesn't contain all the configurable fields).

## Examples

### Lazy logger configuration

```go
package main

import nlogger "github.com/neutron-org/neutron-logger"

func main() {
	l, err := nlogger.NewForContext("my_application")
	if err != nil {
		panic(err)
	}
	l.Debug("debug")
	l.Info("info")
	l.Warn("warm")
	l.Error("error")
}
```

`.env` file:
```
LOGGER_LEVEL=warn
```

results in (stack trace and caller messages are removed for simplicity):
```
{"level":"warn","ts":"2022-09-27T09:00:19.779911+03:00","msg":"warm","context":"my_application"}
{"level":"error","ts":"2022-09-27T09:00:19.78001+03:00","msg":"error","context":"my_application"}
```

### Providing a cfg file

define a `cfg.json` (you can use the following as a boilerplate):
```json
{
    "level": "warn",
    "outputPaths": [
        "stderr"
    ],
    "errorOutputPaths": [
        "stderr"
    ],
    "encoding": "console",
    "sampling": {
        "initial": 100,
        "thereafter": 100
    },
    "encoderConfig": {
        "timeKey": "ts",
        "levelKey": "level",
        "nameKey": "logger",
        "callerKey": "caller",
        "messageKey": "msg",
        "stacktraceKey": "stacktrace",
        "lineEnding": "\n",
        "timeEncoder": "ISO8601"
    }
}
```

make the config file available via `LOGGER_CFG_PATH` env variable
```bash
export LOGGER_CFG_PATH=cfg.json
```

```go
package main

import nlogger "github.com/neutron-org/neutron-logger"

func main() {
	l, err := nlogger.NewForContext("my_application")
	if err != nil {
		panic(err)
	}
	l.Debug("debug")
	l.Info("info")
	l.Warn("warm")
	l.Error("error")
}
```

results in:
```
2022-09-27T08:55:35.379+0300    warm    {"context": "my_application"}
2022-09-27T08:55:35.379+0300    error   {"context": "my_application"}
```

### Loggers registry

If your application contains a lot of different contexts, it may be annoying to create a logger per each context manually with creation error handling. Loggers registry is meant to ease loggers creation and access.

```go
package main

import nlogger "github.com/neutron-org/neutron-logger"

func main() {
    dbContext := "db"
    apiContext := "api"
	registry, err := nlogger.NewRegistry(dbContext, apiContext)
	if err != nil {
		panic(err)
	}
	registry.Get(dbContext).Info("db info message")
	registry.Get(apiContext).Info("api info message")
}
```

results in:
```
{"level":"info","ts":"2022-09-27T12:36:59.919689+03:00","msg":"db info message","context":"db"}
{"level":"info","ts":"2022-09-27T12:36:59.919814+03:00","msg":"api info message","context":"api"}
```