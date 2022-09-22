# Neutron Logger

This repository contains a helper zap.Logger constructor that injects a mandatory "context" field to the logger fields. The context field is expected to carry information about the execution context the logger is used in to ease logs browsing and indexing.

The logging level is set using an env variable `LOG_LEVEL`. According to the zap.Logger docs, the expected `LOG_LEVEL` values are (the `dpanic`/`DPANIC` level is omitted, higher levels are more important):
- `debug` or `DEBUG` — typically voluminous logs that are usually disabled in production;
- `info` or `INFO` — the default logging priority;
- `warn` or `WARN` — logs that are more important than Info, but don't need individual human review;
- `error` or `ERROR` — high-priority logs. If an application is running smoothly, it shouldn't generate any error-level logs;
- `panic` or `PANIC` — logs a message, then panics;
- `fatal` or `FATAL` — logs a message, then calls os.Exit(1).

### Example

```go
package main

func main() {
	l, err := NewForContext("my_application")
	if err != nil {
		panic(err)
	}
	l.Debug("debug")
	l.Info("info")
	l.Warn("warm")
	l.Error("error")
	l.Fatal("fatal")
}
```

`.env` file:
```
LOG_LEVEL=error
```

results in (stack trace message is omitted):
```
{"level":"error","ts":1663859412.242798,"caller":"playground/main.go:11","msg":"error","context":"my_application"}
{"level":"fatal","ts":1663859412.242861,"caller":"playground/main.go:12","msg":"fatal","context":"my_application"}
```