package logger

import "go.uber.org/zap"

const (
	envLocal = "Development"
)

// Logger methods interface
func SetupLogger(env string) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error

	switch env {
	case envLocal:
		logger, err = zap.NewDevelopment()
	default:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	sugarLogger := logger.Sugar()

	return sugarLogger, nil
}
