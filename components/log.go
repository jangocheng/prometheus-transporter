package components

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Logger struct {
	baseLogger  log.Logger
	debugLogger log.Logger
	infoLogger  log.Logger
	warnLogger  log.Logger
	errorLogger log.Logger
}

func (l Logger) Debug(keyvalues ...interface{}) {
	l.debugLogger.Log(keyvalues...)
}

func (l Logger) Info(keyvalues ...interface{}) {
	l.infoLogger.Log(keyvalues...)
}

func (l Logger) Warn(keyvalues ...interface{}) {
	l.warnLogger.Log(keyvalues...)
}

func (l Logger) Error(keyvalues ...interface{}) {
	l.errorLogger.Log(keyvalues...)
}

var logger *Logger

func InitLogger(conf LogConfig) {
	var baseLogger log.Logger

	{
		var lvl level.Option
		switch conf.Level {
		case "error":
			lvl = level.AllowError()
		case "warn":
			lvl = level.AllowWarn()
		case "info":
			lvl = level.AllowInfo()
		case "debug":
			lvl = level.AllowDebug()
		default:
			panic("unexpected log level")
		}

		baseLogger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		if conf.Format == "json" {
			baseLogger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		}
		baseLogger = level.NewFilter(baseLogger, lvl)

		baseLogger = log.With(baseLogger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

	logger = &Logger{
		baseLogger:  baseLogger,
		debugLogger: level.Debug(baseLogger),
		infoLogger:  level.Info(baseLogger),
		warnLogger:  level.Warn(baseLogger),
		errorLogger: level.Error(baseLogger),
	}
}

func GetLogger() *Logger {
	return logger
}
