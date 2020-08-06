package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/evzpav/betting-game/internal/domain/game"
	"gitlab.com/evzpav/betting-game/internal/infrastructure/server/http"
	"gitlab.com/evzpav/betting-game/pkg/env"
	"gitlab.com/evzpav/betting-game/pkg/log"
)

const (
	envVarHost        = "HOST"
	envVarPort        = "PORT"
	envVarLoggerLevel = "LOGGER_LEVEL"

	defaultProjectPort = "5001"
	defaultLoggerLevel = "info"
)

var (
	version, build, date string
)

func main() {
	log := log.NewZeroLog("betting-game", version, log.Level(getLoggerLevel()))

	log.Info().Sendf("betting-game - build:%s; date:%s", build, date)

	env.CheckRequired(log)

	// db, err := mysql.New(getMySQLURL())
	// if err != nil {
	// 	log.Fatal().Err(err).Sendf("failed to connect to mysql: %v", err)
	// }

	// defer func() {
	// 	if err := db.Close(); err != nil {
	// 		log.Error().Err(err).Sendf("error closing database: %v", err)
	// 	}
	// }()

	// // storages
	// userStorage, err := mysql.NewUserStorage(db, log)
	// if err != nil {
	// 	log.Fatal().Err(err).Sendf("error creating storage: %v", err)
	// }

	// services
	gameService := game.NewService(log)

	// HTTP Server
	handler := http.NewHandler(gameService, log)
	server := http.New(handler, getProjectHost(), getProjectPort(), log)
	server.ListenAndServe()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()
}

func getProjectHost() string {
	return env.GetString(envVarHost)
}

func getProjectPort() string {
	return env.GetString(envVarPort, defaultProjectPort)
}

func getLoggerLevel() string {
	return env.GetString(envVarLoggerLevel, defaultLoggerLevel)
}

// func getMySQLURL() string {
// 	return env.GetString(envVarMySQLURL)
// }
