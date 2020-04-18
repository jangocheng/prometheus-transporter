package model

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

/*
 * For Configuration data model
 */
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Config struct {
	HTTP      string    `toml:"http"`
	Transfers []string  `toml:"transfers"`
	Logger    LogConfig `toml:"logger"`
}

/*
 * For Logger data model
 */
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

func NewLogger(base log.Logger) *Logger {
	return &Logger{
		baseLogger:  base,
		debugLogger: level.Debug(base),
		infoLogger:  level.Info(base),
		warnLogger:  level.Warn(base),
		errorLogger: level.Error(base),
	}
}
