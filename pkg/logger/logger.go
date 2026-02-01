package logger

import "go.uber.org/zap"

func New(level string) *zap.Logger {
	var cfg zap.Config
	switch level {
	case "prod":
		cfg = zap.NewProductionConfig()
	case "dev":
		cfg = zap.NewDevelopmentConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		return zap.NewExample()
	}

	return logger
}
