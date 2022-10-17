# Neutron Logger

This repository contains a helper zap.Logger constructor that injects a mandatory "context" field to the logger fields. The context field is expected to carry information about the execution context the logger is used in to ease logs browsing and indexing.

## Configuration

### General configuration

Basically, the logger is configured with the [zap.NewProductionLogger](https://github.com/uber-go/zap/blob/d6ce3b9b283401bc6cf975de6ac3e5ed5aec5341/config.go#L115) with modified logs timestamps as RFC3339-formatted string with nanosecond precision.

For the default logger configuration the logging level is set using an env variable `LOGGER_LEVEL`. According to the [zap.Logger docs](https://github.com/uber-go/zap/blob/d6ce3b9b283401bc6cf975de6ac3e5ed5aec5341/level.go#L28), the expected `LOGGER_LEVEL` values are (higher levels are more important):
- `debug` or `DEBUG` — typically voluminous logs that are usually disabled in production;
- `info` or `INFO` — the default logging priority;
- `warn` or `WARN` — logs that are more important than Info, but don't need individual human review;
- `error` or `ERROR` — high-priority logs. If an application is running smoothly, it shouldn't generate any error-level logs;
- `dpanic` or `DPANIC` — logs a message, then panics. Works only if logger's development field set to true;
- `panic` or `PANIC` — logs a message, then panics;
- `fatal` or `FATAL` — logs a message, then calls os.Exit(1).

Although the default cfg is introduced for logger instantiation and usage simplicity, a user can still configure the logger the way the [zap.Config](https://github.com/uber-go/zap/blob/d6ce3b9b283401bc6cf975de6ac3e5ed5aec5341/config.go#L45) allows it via _setting_ env variables. If no value is set to an env variable, the production logger config value is used by default. See the [env configuration example](#configuration-via-env-variables) for details.

Reminder: a _set_ env variable is a variable of any value stored in the system, even an empty one. So, e.g. `export LOGGER_ENCODERCONFIG_TIMEKEY=` will not apply a default value to the logger, but will result in the timestamp key absence. To set a default value, use `unset LOGGER_ENCODERCONFIG_TIMEKEY`. This is true for all config parameters.

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

### Configuration via env variables

configure all needed logger parameters by exporting corresponding env variables, e.g.:
```bash
export LOGGER_LEVEL=warn
export LOGGER_ENCODING=console
export LOGGER_ENCODERCONFIG_ENCODETIME=ISO8601
```

As you can see, just like LOGGER_LEVEL, all logger-related env variables are prefixed with `LOGGER_` and then writen solidly with `_` as structure level separator.

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