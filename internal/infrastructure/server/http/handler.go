package http

import (
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

	mux := http.NewServeMux()

	mux.Handle("/", h.fileServerHandler())
	mux.HandleFunc("/api/ws", h.serveWs)
	mux.HandleFunc("/api/game/join", h.postJoin)
	mux.HandleFunc("/api/game/snapshot", h.getGameSnapshot)
	mux.HandleFunc("/api/ranking/snapshot", h.getRankingSnapshot)

	return cors.Default().Handler(mux)

}

func (h *handler) fileServerHandler() http.Handler {
	dir, err := os.Getwd()
	if err != nil {
		h.log.Fatal().Sendf("failed to get working directory: %v", err)
	}

	return http.FileServer(http.Dir(dir + "/frontend/dist/"))
}
