package domain

import (
	"net/http"
	"sort"

	"github.com/robfig/cron/v3"
)

type Game struct {
	ID           string
	PlayersChan  chan *Player
	Players      []*Player
	GameRunning  bool
	RoundCounter int
	StopGame     chan bool
	Cron         *cron.Cron
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

func (g *Game) Reset() {
	g.GameRunning = false
	g.RoundCounter = 0

	for _, p := range g.Players {
		p.ResetPoints()
	}
}

type GameService interface {
	ServeWs(w http.ResponseWriter, r *http.Request)
	Run()
	Join(Player) error
}
