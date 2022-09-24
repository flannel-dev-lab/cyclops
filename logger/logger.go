// Package logger contains the functions related to logging
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Debug Prints out the debug logs
func Debug(ctx context.Context, msg string) {
	values := GetAll(ctx)
	if values == nil {
		values = make(map[string]string)
	}

	values["DEBUG"] = msg
	manageLog(values)
}

// Info Prints out the info logs
func Info(ctx context.Context, msg string) {
	values := GetAll(ctx)
	if values == nil {
		values = make(map[string]string)
	}

	values["INFO"] = msg
	manageLog(values)
}

// Warn Prints out to warn logs
func Warn(ctx context.Context, msg string, err error) {
	values := GetAll(ctx)
	if values == nil {
		values = make(map[string]string)
	}

	values["WARN"] = msg
	values["CAUSE"] = fmt.Sprintf("%v", err)

	manageLog(values)
}

// Error prints out the error logs
func Error(ctx context.Context, msg string, err error) {
	values := GetAll(ctx)
	if values == nil {
		values = make(map[string]string)
	}

	values["ERROR"] = msg
	values["CAUSE"] = fmt.Sprintf("%v", err)

	manageLog(values)
}

// AddKey adds key and value to existing context
func AddKey(ctx context.Context, key, value string) context.Context {
	keys := copyMap(getStoredKeys(ctx))

	keys = append(keys, key)

	ctx = context.WithValue(ctx, "keys", keys)
	ctx = context.WithValue(ctx, key, value)

	return ctx
}

// GetAll gets all the keys and values stored in the context
func GetAll(ctx context.Context) map[string]string {
	values := make(map[string]string)
	keys := getStoredKeys(ctx)

	for _, k := range keys {
		if v, ok := ctx.Value(k).(string); ok {
			values[k] = v
		}
	}

	return values
}

func copyMap(from []string) []string {
	to := make([]string, len(from))

	for idx, k := range from {
		to[idx] = k
	}

	return to
}

func getStoredKeys(ctx context.Context) []string {
	var keys []string

	if k, ok := ctx.Value("keys").([]string); ok {
		keys = k
	}

	return keys
}

func manageLog(values map[string]string) {
	values["timestamp"] = time.Now().Format(time.RFC3339)

	log, err := json.Marshal(values)
	if err != nil {
		fmt.Println("Could not create JSON from values")
		return
	}

	fmt.Println(string(log))
}
