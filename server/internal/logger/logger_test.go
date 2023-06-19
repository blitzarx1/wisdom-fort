package logger

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		parentLogger *log.Logger
		prefix       string
	}
	tests := []struct {
		name       string
		args       args
		wantPrefix string
	}{
		{
			name: "No Parent Logger",
			args: args{
				parentLogger: nil,
				prefix:       "test",
			},
			wantPrefix: "test: ",
		},
		{
			name: "With Parent Logger",
			args: args{
				parentLogger: log.New(os.Stdout, "parent: ", logFlags),
				prefix:       "test",
			},
			wantPrefix: "parent->test: ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.parentLogger, tt.args.prefix)
			if got.Prefix() != tt.wantPrefix {
				t.Errorf("NewLogger() = %v, want %v", got.Prefix(), tt.wantPrefix)
			}
		})
	}
}

func TestWithCtxAndMustFromCtx(t *testing.T) {
	parentLogger := log.New(os.Stdout, "parent: ", logFlags)
	prefix := "test"

	ctx := WithCtx(context.Background(), parentLogger, prefix)
	// test if logger is added in context
	l := MustFromCtx(ctx)
	if l.Prefix() != "parent->test: " {
		t.Errorf("WithCtx() = %v, want %v", l.Prefix(), "parent->test: ")
	}

	// test panic case
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustFromCtx was supposed to panic")
		}
	}()
	_ = MustFromCtx(context.Background())
}
