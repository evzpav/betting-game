package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"
	"gitlab.com/evzpav/betting-game/internal/domain"
	"gitlab.com/evzpav/betting-game/pkg/log"
)

type handler struct {
	gameService domain.GameService
	log         log.Logger
}

func NewHandler(gameService domain.GameService, log log.Logger) http.Handler {
	h := &handler{
		gameService: gameService,
		log:         log,
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	fs := http.FileServer(http.Dir(dir + "/frontend/dist/"))
	mux := http.NewServeMux()

	mux.Handle("/", fs)
	mux.HandleFunc("/api/ws", h.serveWs)
	mux.HandleFunc("/api/join", h.postJoin)

	return cors.Default().Handler(mux)

}
