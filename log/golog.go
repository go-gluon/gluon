package log

import (
	"fmt"
	"log"
	"os"
)

type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
)

type SimpleLogger struct {
	level Level
	warn  *log.Logger
	info  *log.Logger
	erro  *log.Logger
	debu  *log.Logger
	trac  *log.Logger
}

func NewSimpleLogger() *SimpleLogger {
	d := SimpleLogger{level: InfoLevel}
	d.debu = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile|log.Lmsgprefix|log.Lmicroseconds)
	d.trac = log.New(os.Stderr, "TRACE: ", log.Ldate|log.Ltime|log.Llongfile|log.Lmsgprefix)
	d.info = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Llongfile|log.Lmsgprefix)
	d.warn = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile|log.Lmsgprefix)
	d.erro = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile|log.Lmsgprefix)
	return &d
}

func (d *SimpleLogger) Level(lvl Level) *SimpleLogger {
	d.level = lvl
	return d
}

func (d *SimpleLogger) Trace(msg string, fields ...map[string]interface{}) {
	d.log(d.debu, TraceLevel, msg, fields...)
}

func (d *SimpleLogger) Debug(msg string, fields ...map[string]interface{}) {
	d.log(d.debu, DebugLevel, msg, fields...)
}

func (d *SimpleLogger) Info(msg string, fields ...map[string]interface{}) {
	d.log(d.debu, InfoLevel, msg, fields...)
}

func (d *SimpleLogger) Warn(msg string, fields ...map[string]interface{}) {
	d.log(d.debu, WarnLevel, msg, fields...)
}

func (d *SimpleLogger) Error(msg string, fields ...map[string]interface{}) {
	d.log(d.debu, ErrorLevel, msg, fields...)
}

func (d *SimpleLogger) log(log *log.Logger, lvl Level, msg string, fields ...map[string]interface{}) {
	if lvl < d.level {
		return
	}
	if len(fields) > 0 {
		log.Printf("%s %+v\n", msg, fmt.Sprint(fields))
		return
	}
	log.Printf("%s\n", msg)
}
