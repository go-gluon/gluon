package log

import (
	"fmt"
)

var (
	Log Logger = SimpleLogger{}
)

type Fields map[string]interface{}

func (f Fields) Add(k string, v interface{}) Fields {
	f[k] = v
	return f
}

func (f Fields) Err(err error) Fields {
	f["error"] = err
	return f
}

func Add(k string, v interface{}) Fields {
	return Fields{k: v}
}

func Err(err error) Fields {
	return Fields{}.Err(err)
}

func ErrorE(msg string, err error) {
	Log.Error(msg, Err(err))
}

func Error(msg string, fields ...map[string]interface{}) {
	Log.Error(msg, fields...)
}

func Trace(msg string, fields ...map[string]interface{}) {
	Log.Trace(msg, fields...)
}

func Debug(msg string, fields ...map[string]interface{}) {
	Log.Debug(msg, fields...)
}

func Info(msg string, fields ...map[string]interface{}) {
	Log.Info(msg, fields...)
}

func Warn(msg string, fields ...map[string]interface{}) {
	Log.Warn(msg, fields...)
}

type Logger interface {
	Trace(msg string, fields ...map[string]interface{})
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
}

type SimpleLogger struct {
}

func (d SimpleLogger) Trace(msg string, fields ...map[string]interface{}) {
	fmt.Printf("[TRACE] %s %v\n", msg, fmt.Sprint(fields))
}

func (d SimpleLogger) Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		fmt.Printf("[DEBUG] %s %v\n", msg, fmt.Sprint(fields))
	}
	fmt.Printf("[DEBUG] %s\n", msg)
}

func (d SimpleLogger) Info(msg string, fields ...map[string]interface{}) {
	fmt.Printf(" [INFO] %s %v\n", msg, fmt.Sprint(fields))
}

func (d SimpleLogger) Warn(msg string, fields ...map[string]interface{}) {
	fmt.Printf(" [WARN] %s %v\n", msg, fmt.Sprint(fields))
}

func (d SimpleLogger) Error(msg string, fields ...map[string]interface{}) {
	fmt.Printf("[ERROR] %s %v\n", msg, fmt.Sprint(fields))
}
