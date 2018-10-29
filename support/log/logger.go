package log

import (
	"os"
	"strings"
)

type Level int
type Format int

const (
	EnvKeyLogCtx         = "FLOGO_LOG_CTX"
	EnvKeyLogDateFormat  = "FLOGO_LOG_DTFORMAT"
	DefaultLogDateFormat = "2006-01-02 15:04:05.000"
	EnvKeyLogLevel       = "FLOGO_LOG_LEVEL"
	DefaultLogLevel      = LevelInfo
	EnvKeyLogFormat      = "FLOGO_LOG_FORMAT"
	DefaultLogFormat     = FormatConsole

	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError

	FormatConsole Format = iota
	FormatJson
)

type Logger interface {
	DebugEnabled() bool
	TraceEnabled() bool

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Tracef(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})

	Structured() StructuredLogger
}

type StructuredLogger interface {
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field = interface{}

var (
	rootLogger Logger
	ctxLogging bool
)

func init() {
	configureLogging()
}

func CtxLoggingEnabled() bool {
	return ctxLogging
}

func RootLogger() Logger {
	return rootLogger
}

func SetLogLevel(logger Logger, level Level) {
	setZapLogLevel(logger, level)
}

func ChildLogger(logger Logger, name string) Logger {

	childLogger, err := newZapChildLogger(logger, name)
	if err != nil {
		rootLogger.Warnf("unable to create child logger named: %s - %s", name, err.Error())
		childLogger = logger
	}

	return childLogger
}

func ChildLoggerWithFields(logger Logger, fields ...Field) Logger {
	childLogger, err := newZapChildLoggerWithFields(logger, fields...)
	if err != nil {
		rootLogger.Warnf("unable to create child logger with fields: %s", err.Error())
		childLogger = logger
	}

	return childLogger
}

func Sync() {
	zapSync(rootLogger)
}

func configureLogging() {
	envLogCtx := os.Getenv(EnvKeyLogCtx)
	if strings.ToLower(envLogCtx) == "true" {
		ctxLogging = true
	}

	rootLogLevel := DefaultLogLevel

	envLogLevel := strings.ToUpper(os.Getenv(EnvKeyLogLevel))
	switch envLogLevel {
	case "DEBUG":
		rootLogLevel = LevelDebug
	case "INFO":
		rootLogLevel = LevelInfo
	case "WARN":
		rootLogLevel = LevelWarn
	case "ERROR":
		rootLogLevel = LevelError
	default:
		rootLogLevel = DefaultLogLevel
	}

	logFormat := DefaultLogFormat
	envLogFormat := strings.ToUpper(os.Getenv(EnvKeyLogFormat))
	if envLogFormat == "JSON" {
		logFormat = FormatJson
	}

	rootLogger = newZapRootLogger("flogo", logFormat)
	SetLogLevel(rootLogger, rootLogLevel)
}