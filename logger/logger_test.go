package logger

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()

	ctx = AddKey(ctx, "hello", "world")
	ctx = AddKey(ctx, "good", "bye")

	Debug(ctx, "debug")
	Info(ctx, "info")
	Warn(ctx, "warn", errors.New("warn"))
	Error(ctx, "error", errors.New("error"))

	reflect.DeepEqual(GetAll(ctx), map[string]string{"good": "bye", "hello": "world"})
}

func TestLogger_Nil(t *testing.T) {
	ctx := context.Background()

	Debug(ctx, "debug")
	Info(ctx, "info")
	Warn(ctx, "warn", errors.New("warn"))
	Error(ctx, "error", errors.New("error"))

	reflect.DeepEqual(GetAll(ctx), map[string]string{"good": "bye", "hello": "world"})
}
