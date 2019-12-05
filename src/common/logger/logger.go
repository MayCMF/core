package logger

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

// Define key name
const (
	TraceIDKey      = "trace_id"
	UserUUIDKey     = "user_uuid"
	UserIDKey       = "user_id"
	SpanTitleKey    = "span_title"
	SpanFunctionKey = "span_function"
	VersionKey      = "version"
)

// TraceIDFunc - Define the function to get the tracking ID
type TraceIDFunc func() string

var (
	version     string
	traceIDFunc TraceIDFunc
)

// Logger - Define a log alias
type Logger = logrus.Logger

// Hook - Define a log hook alias
type Hook = logrus.Hook

// StandardLogger - Get standard logs
func StandardLogger() *Logger {
	return logrus.StandardLogger()
}

// SetLevel - Set the log level
func SetLevel(level int) {
	logrus.SetLevel(logrus.Level(level))
}

// SetFormatter - Set the log output format
func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

// SetOutput - Set log output
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetVersion - Setting version
func SetVersion(v string) {
	version = v
}

// SetTraceIDFunc - Set the processing function of the tracking ID
func SetTraceIDFunc(fn TraceIDFunc) {
	traceIDFunc = fn
}

// AddHook - Increase log hooks
func AddHook(hook Hook) {
	logrus.AddHook(hook)
}

func getTraceID() string {
	if traceIDFunc != nil {
		return traceIDFunc()
	}
	return ""
}

type (
	traceIDContextKey  struct{}
	spanIDContextKey   struct{}
	userUUIDContextKey struct{}
)

// NewTraceIDContext - Create a tracking ID context
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey{}, traceID)
}

// FromTraceIDContext - Get the tracking ID from the context
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return getTraceID()
}

// NewUserUUIDContext - Create a user ID context
func NewUserUUIDContext(ctx context.Context, userUUID string) context.Context {
	return context.WithValue(ctx, userUUIDContextKey{}, userUUID)
}

// FromUserUUIDContext - Get the user ID from the context
func FromUserUUIDContext(ctx context.Context) string {
	v := ctx.Value(userUUIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type spanOptions struct {
	Title    string
	FuncName string
}

// SpanOption - Define the data item of the tracking unit
type SpanOption func(*spanOptions)

// SetSpanTitle - Set the title of the tracking unit
func SetSpanTitle(title string) SpanOption {
	return func(o *spanOptions) {
		o.Title = title
	}
}

// SetSpanFuncName - Set the function name of the tracking unit
func SetSpanFuncName(funcName string) SpanOption {
	return func(o *spanOptions) {
		o.FuncName = funcName
	}
}

// StartSpan - Start a tracking unit
func StartSpan(ctx context.Context, opts ...SpanOption) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	var o spanOptions
	for _, opt := range opts {
		opt(&o)
	}

	fields := map[string]interface{}{
		UserUUIDKey: FromUserUUIDContext(ctx),
		TraceIDKey:  FromTraceIDContext(ctx),
		VersionKey:  version,
	}
	if v := o.Title; v != "" {
		fields[SpanTitleKey] = v
	}
	if v := o.FuncName; v != "" {
		fields[SpanFunctionKey] = v
	}

	return newEntry(logrus.WithFields(fields))
}

// StartSpanWithCall - Start a tracking unit (callback execution)
func StartSpanWithCall(ctx context.Context, opts ...SpanOption) func() *Entry {
	return func() *Entry {
		return StartSpan(ctx, opts...)
	}
}

// Debugf - Write debug log
func Debugf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Debugf(format, args...)
}

// Infof - Write to the message log
func Infof(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Infof(format, args...)
}

// Printf - Write to the message log
func Printf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Printf(format, args...)
}

// Warnf - Write warning log
func Warnf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Warnf(format, args...)
}

// Errorf - Write error log
func Errorf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Errorf(format, args...)
}

// Fatalf Write a major error log
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	StartSpan(ctx).Fatalf(format, args...)
}

func newEntry(entry *logrus.Entry) *Entry {
	return &Entry{entry: entry}
}

// Entry - Define a unified log write mode
type Entry struct {
	entry *logrus.Entry
}

func (e *Entry) checkAndDelete(fields map[string]interface{}, keys ...string) {
	for _, key := range keys {
		if _, ok := fields[key]; ok {
			delete(fields, key)
		}
	}
}

// WithFields - Structured field write
func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	e.checkAndDelete(fields,
		TraceIDKey,
		SpanTitleKey,
		SpanFunctionKey,
		VersionKey)
	return newEntry(e.entry.WithFields(fields))
}

// WithField - Structured field write
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(map[string]interface{}{key: value})
}

// Fatalf - Major error log
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

// Errorf - Error log
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

// Warnf - Warning log
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

// Infof - Message log
func (e *Entry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

// Printf - Message log
func (e *Entry) Printf(format string, args ...interface{}) {
	e.entry.Printf(format, args...)
}

// Debugf - Write debug log
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}
