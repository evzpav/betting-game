package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/evzpav/betting-game/internal/domain/game"
	"gitlab.com/evzpav/betting-game/internal/domain/hub"
	http "gitlab.com/evzpav/betting-game/internal/infrastructure/server/http"
	localstorage "gitlab.com/evzpav/betting-game/internal/infrastructure/storage/localStorage"
	"gitlab.com/evzpav/betting-game/pkg/env"
	"gitlab.com/evzpav/betting-game/pkg/log"
)

const (
	envVarHost        = "HOST"
	envVarPort        = "PORT"
	envVarLoggerLevel = "LOGGER_LEVEL"

	defaultProjectPort = "8787"
	defaultLoggerLevel = "info"
)

var (
	version, build, date string
)

func main() {
	log := log.NewZeroLog("betting-game", version, log.Level(getLoggerLevel()))

	log.Info().Sendf("betting-game - build:%s; date:%s", build, date)

	env.CheckRequired(log)

	// game rules
	const (
		minPlayersToStart int = 2
		maxRoundsPerGame  int = 30
		intervalSeconds   int = 10
		magicNumberMatch  int = 21
	)

	// storage
	gameStorage := localstorage.NewGameStorage(log)

	// services
	hubService := hub.NewService(log)
	gameService := game.NewService(gameStorage, hubService, log)
	gameService.SetGameRules(minPlayersToStart, maxRoundsPerGame, intervalSeconds, magicNumberMatch)
	gameService.Run()

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
