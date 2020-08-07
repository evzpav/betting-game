package main

import (
	"fmt"
	"os"

	"net/http"

	"gitlab.com/evzpav/betting-game/internal/domain/game"
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
	go gameService.RunHub()

	// HTTP Server

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gameService.ServeWs(w, r)
	})
	err := http.ListenAndServe(":8787", nil)
	if err != nil {
		log.Fatal().Sendf("ListenAndServe: ", err)
	}
	// handler := http.NewHandler(gameService, log)
	// server := http.New(handler, getProjectHost(), getProjectPort(), log)
	// server.ListenAndServe()

	// Graceful shutdown
	// stopChan := make(chan os.Signal, 1)
	// signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	// <-stopChan
	// server.Shutdown()
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

func serveHome(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	fmt.Println(dir)

	http.ServeFile(w, r, dir+"/frontend/dist/index.html")
}
