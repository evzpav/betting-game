package domain

import (
	"math/rand"
	"net/http"
	"sort"
	"time"
)

type Game struct {
	ID           string    `json:"id"`
	Players      []*Player `json:"players"`
	Observers    []*Player `json:"observers"`
	GameRunning  bool      `json:"gameRunning"`
	RoundCounter int       `json:"roundCounter"`
}

func (g *Game) ResolveWinner() *Player {
	sort.SliceStable(g.Players, func(i, j int) bool {
		a := g.Players[i]
		b := g.Players[j]

		if a.Points == b.Points {
			if a.Numbers[1] == b.Numbers[1] {
				if a.Numbers[0] == b.Numbers[0] {
					return a.Name < b.Name
				}

				return a.Numbers[0] > b.Numbers[0]
			}

			return a.Numbers[1] > b.Numbers[1]
		}

		return a.Points > b.Points
	})

	return g.Players[0]
}

func (g *Game) GenerateRandomNumber() int {
	min := 1
	max := 10
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max+1-min) + min
}

type GameService interface {
	ServeWs(w http.ResponseWriter, r *http.Request)
	Run()
	Join(Player)
}
