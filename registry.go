package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// NewRegistry creates a logger for each provided context.
func NewRegistry(ctxs ...string) (*Registry, error) {
	registry := &Registry{
		mu:      &sync.Mutex{},
		loggers: make(map[string]*zap.Logger, len(ctxs)),
	}
	for _, ctx := range ctxs {
		logger, err := NewForContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create a logger for ctx %s: %w", ctx, err)
		}
		registry.loggers[ctx] = logger
	}
	return registry, nil
}

// Registry contains a list of loggers mapped to their context.
type Registry struct {
	mu      *sync.Mutex
	loggers map[string]*zap.Logger
}

// Get retrieves a logger for the given context. If no logger set for the context, it panics.
func (r *Registry) Get(ctx string) *zap.Logger {
	r.mu.Lock()
	defer r.mu.Unlock()
	logger, ex := r.loggers[ctx]
	if !ex {
		panic(fmt.Sprintf("no logger set for context %s", ctx))
	}
	return logger
}

// Set sets the logger for the given context.
func (r *Registry) Set(ctx string, logger *zap.Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.loggers[ctx] = logger
}
