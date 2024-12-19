package global

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"layout/config"
	zapLogger "layout/pkg/zap_logger"
)

func GetLogger() *zap.SugaredLogger {
	return zapLogger.GetLogger()
}

type Option struct {
	LoggerCfg *config.Logger
}

func Init(o Option) error {
	err := zapLogger.Init(
		zapLogger.SetFileName(o.LoggerCfg.FileName),
		zapLogger.SetLogFileDir(o.LoggerCfg.LogDir),
		zapLogger.SetDevelopment(o.LoggerCfg.IsProd),
		zapLogger.SetMaxSize(o.LoggerCfg.MaxSize),
		zapLogger.SetMaxAge(o.LoggerCfg.MaxAge),
		zapLogger.SetMaxBackups(o.LoggerCfg.MaxBackups))
	if err != nil {
		return errors.Wrap(err, "初始化日志失败")
	}
	return nil
}
