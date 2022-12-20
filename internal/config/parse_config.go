package config

import (
	"github.com/spf13/viper"
	"log"
	"go.uber.org/zap/zapcore"
	"auth/internal/models"
)


var logLevelSeverity = map[string]zapcore.Level{
    "DEBUG" : zapcore.DebugLevel,
    "INFO" : zapcore.InfoLevel,
    "WARNING" : zapcore.WarnLevel,
	"ERROR" : zapcore.ErrorLevel,
    "CRITICAL" : zapcore.DPanicLevel,
    "ALERT" : zapcore.PanicLevel,
    "EMERGENCY" : zapcore.FatalLevel,
}

func GetPort() int {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
    	log.Fatalf("reading config failed: %#v", err)
	}
	
	return viper.GetInt("app.port")
}

func GetMetricsPort() int {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
    	log.Fatalf("reading config failed: %#v", err)
	}
	
	return viper.GetInt("app.metrics_port")
}

func GetUsersFromDb() map[string]*models.User {
	mapAns := make(map[string]*models.User)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
    	log.Fatalf("reading config failed: %#v", err)
	}

	mapUsersParsed := make(map[string]*models.User)

	if err := viper.UnmarshalKey("users", &mapUsersParsed); err != nil {
    	log.Fatalf("unmarshaling config failed: %#v", err)
	}
	
	for _, user := range mapUsersParsed {
		mapAns[user.Login] = user
	}
	return mapAns
}

func GetLoggerLevel() zapcore.Level {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
    	log.Fatalf("reading config failed: %#v", err)
	}
	
	return logLevelSeverity[viper.GetString("logger.level")]
}
