package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/evzpav/betting-game/internal/domain"
)

func (h *handler) serveWs(w http.ResponseWriter, r *http.Request) {
	h.gameService.ServeWs(w, r)
}
func (h *handler) postJoin(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var player domain.Player
	if err := json.Unmarshal(bs, &player); err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(bs))
	fmt.Printf("Player: %+v\n", player)

	if err := player.Validate(); err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playerID := h.gameService.Join(player)

	player.ID = playerID
	
	bs, err = json.Marshal(player)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
