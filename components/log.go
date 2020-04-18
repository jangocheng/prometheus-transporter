package components

import (
	"os"
	"prometheus-transporter/model"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var logger *model.Logger

func InitLogger(conf model.LogConfig) {
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

	logger = model.NewLogger(baseLogger)
}

func GetLogger() *model.Logger {
	return logger
}
