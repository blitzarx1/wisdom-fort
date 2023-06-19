package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

type key string

const keyLogger key = "logger"

const (
	logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix

	prefixEnd       = ": "
	separatorLogger = "->"
)

// New returns new logger with the provided prefix. If parentLogger is nil,
// then new logger will be created with os.Stdout as output.
func New(parentLogger *log.Logger, prefix string) *log.Logger {
	if parentLogger == nil {
		return log.New(os.Stdout, prefix+prefixEnd, logFlags)
	}

	trimmedParentPrefix := strings.TrimSuffix(parentLogger.Prefix(), prefixEnd)
	newPrefix := fmt.Sprintf("%s"+separatorLogger+"%s"+prefixEnd, trimmedParentPrefix, prefix)

	return log.New(parentLogger.Writer(), newPrefix, parentLogger.Flags())
}

// WithCtx packs logger to the provided context. Creates new logger if prefix is not empty.
// Returns new context with a new copy of logger packed inside.
func WithCtx(ctx context.Context, l *log.Logger, prefix string) context.Context {
	if prefix != "" {
		l = New(l, prefix)
	}
	return context.WithValue(ctx, keyLogger, l)
}

// MustFromCtx unpacks Logger from the provided context. Panics if Logger is not present
// in the context.
func MustFromCtx(ctx context.Context) *log.Logger {
	l, ok := ctx.Value(keyLogger).(*log.Logger)
	if !ok {
		panic("logger is not present in the context")
	}

	return l
}
