package http

import (
	"fmt"
	"io/ioutil"
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
	// h := &handler{
	// 	gameService: gameService,
	// 	log:         log,
	// }

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	fs := http.FileServer(http.Dir(dir + "/frontend/dist/"))
	mux := http.NewServeMux()

	mux.Handle("/", fs)
	mux.HandleFunc("/api/ws", gameService.ServeWs)
	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	mux.HandleFunc("/api/join", func(w http.ResponseWriter, r *http.Request) {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Sendf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Println(string(bs))

		// gameService.Join()
	})

	handler := cors.Default().Handler(mux)

	return handler

}
