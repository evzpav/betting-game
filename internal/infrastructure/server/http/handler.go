package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	mux.HandleFunc("/", h.staticHandler("frontend/dist"))
	mux.HandleFunc("/api/ws", h.serveWs)
	mux.HandleFunc("/api/game/join", h.postJoin)
	mux.HandleFunc("/api/game/snapshot", h.getGameSnapshot)
	mux.HandleFunc("/api/ranking/snapshot", h.getRankingSnapshot)

	return cors.Default().Handler(mux)

}

func (h *handler) staticHandler(staticPath string) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		path = filepath.Join(staticPath, path)

		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.fileServerHandler(staticPath).ServeHTTP(w, r)
	})

}

func (h *handler) fileServerHandler(staticPath string) http.Handler {
	dir, err := os.Getwd()
	if err != nil {
		h.log.Fatal().Sendf("failed to get working directory: %v", err)
	}

	return http.FileServer(http.Dir(fmt.Sprintf("%s/%s/", dir, staticPath)))
}
