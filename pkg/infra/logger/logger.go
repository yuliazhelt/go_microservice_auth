package logger

import (
	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"auth/internal/config"
	
)


const (
	DSN = "https://6815a4754d814fe580391b2d6ef5b599@o4503908268507136.ingest.sentry.io/4503908778180608"
)

func initSentry(log *zap.Logger, sentryAddress, environment string) *zap.Logger {
	if sentryAddress == "" {
		return log
	}

	cfg := zapsentry.Configuration{
		Level: config.GetLoggerLevel(),
		Tags: map[string]string{
			"environment": environment,
			"app":         "demoApp",
		},
	}

	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(sentryAddress))
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}

	return zapsentry.AttachCoreToLogger(core, log)
}

func GetLogger(debug bool, dsn string, env string) (*zap.Logger, error) {
	var err error
	var l *zap.Logger

	if debug {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}

	l = initSentry(l, dsn, env)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = l.Sync()
	}()
	l.Debug("Logger initialized in debug level")

	return l, err
}