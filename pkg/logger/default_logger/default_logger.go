package defaultlogger

import (
	"os"

	"github.com/phatnt199/go-infra/pkg/application/constants"
	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/logger/config"
	"github.com/phatnt199/go-infra/pkg/logger/models"
	"github.com/phatnt199/go-infra/pkg/logger/zap"
)

var l logger.Logger

func initLogger() {
	logType := os.Getenv("LogConfig_LogType")

	switch logType {
	case "Zap", "":
		l = zap.NewZapLogger(
			&config.LogOptions{LogType: models.Zap, CallerEnabled: false},
			constants.DEV_ENV,
		)
	default:
	}
}

func GetLogger() logger.Logger {
	if l == nil {
		initLogger()
	}

	return l
}
