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
	h := &handler{
		gameService: gameService,
		log:         log,
	}

	r := mux.NewRouter()
	r.Use(h.logger())

	// spa := spaHandler{staticPath: "frontend/dist", indexPath: "index.html"}
	// r.PathPrefix("/").Handler(spa).Methods("GET")
	// r.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		h.gameService.ServeWs(w, r)
	})

	// r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	hj, ok := w.(http.Hijacker)
	// 	if !ok {
	// 		fmt.Println("not ok")
	// 		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	conn, bufrw, err := hj.Hijack()
	// 	if err != nil {
	// 		fmt.Println("not huj")

	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	// Don't forget to close the connection:
	// 	defer conn.Close()
	// 	bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
	// 	bufrw.Flush()
	// 	s, err := bufrw.ReadString('\n')
	// 	if err != nil {
	// 		h.log.Info().Sendf("error reading string: %v", err)
	// 		return
	// 	}
	// 	fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
	// 	bufrw.Flush()

	// 	// var upgrader = websocket.Upgrader{
	// 	// 	ReadBufferSize:  1024,
	// 	// 	WriteBufferSize: 1024,
	// 	// }

	// 	// conn, err := upgrader.Upgrade(hj, r, nil)
	// 	// if err != nil {
	// 	// 	h.log.Error().Err(err).Sendf("%v", err)
	// 	// 	return
	// 	// }

	// 	// h.gameService.RegisterNewClient(conn)

	// }).Methods("GET")

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	return r
}
