package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/evzpav/betting-game/internal/domain"
	"gitlab.com/evzpav/betting-game/pkg/log"
)

type handler struct {
	gameService domain.GameService
	log         log.Logger
}

func NewHandler(gameService domain.GameService, log log.Logger) http.Handler {
	handler := &handler{
		gameService: gameService,
		log:         log,
	}

	r := mux.NewRouter()
	r.Use(handler.logger())

	r.HandleFunc("/", redirectToLogin).Methods("GET")

	return r
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}
