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

	if err := player.Validate(); err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	player = h.gameService.Join(player)

	bs, err = json.Marshal(player)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func (h *handler) getRankingSnapshot(w http.ResponseWriter, r *http.Request) {
	ranking := h.gameService.GetRankingSnapshot()
	bs, err := json.Marshal(ranking)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("ranking:%v\n", string(bs))
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func (h *handler) getGameSnapshot(w http.ResponseWriter, r *http.Request) {
	game := h.gameService.GetGameSnapshot()
	bs, err := json.Marshal(game)
	if err != nil {
		h.log.Error().Err(err).Sendf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("game:%v\n", string(bs))
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
