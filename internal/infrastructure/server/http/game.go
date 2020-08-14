package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"gitlab.com/evzpav/betting-game/internal/domain"
)

type respMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (h *handler) responseBadRequest(w http.ResponseWriter, err error) {
	var msg = respMessage{
		Message: err.Error(),
		Status:  http.StatusBadRequest,
	}

	bs, err := json.Marshal(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write(bs); err != nil {
		h.log.Error().Sendf("failed to write response: %v", err)
	}
}

func (h *handler) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		return
	}

	h.gameService.AddNewWebsocketClient(conn)

}

func (h *handler) postJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var player domain.Player
	if err := json.Unmarshal(bs, &player); err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := player.Validate(); err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		h.responseBadRequest(w, err)
		return
	}

	player, err = h.gameService.Join(player)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		h.responseBadRequest(w, err)
		return
	}

	bs, err = json.Marshal(player)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bs); err != nil {
		h.log.Error().Sendf("failed to write response: %v", err)
	}
}

func (h *handler) getRankingSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ranking := h.gameService.GetRankingSnapshot()
	bs, err := json.Marshal(ranking)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bs); err != nil {
		h.log.Error().Sendf("failed to write response: %v", err)
	}
}

func (h *handler) getGameSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	game := h.gameService.GetGameSnapshot()

	bs, err := json.Marshal(game)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bs); err != nil {
		h.log.Error().Sendf("failed to write response: %v", err)
	}
}