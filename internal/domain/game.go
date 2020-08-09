package domain

import (
	"math/rand"
	"net/http"
	"sort"
	"time"
)

const magicNumber int = 21

type Game struct {
	ID           string    `json:"id"`
	Players      []*Player `json:"players"`
	Observers    []*Player `json:"-"`
	GameRunning  bool      `json:"gameRunning"`
	RoundCounter int       `json:"roundCounter"`
	RandomNumber int       `json:"randomNumber"`
	GameCounter  int       `json:"gameCounter"`
	Rules        Rules     `json:"rules"`
}

type Rules struct {
	MinPlayersToStart int `json:"minPlayersToStart"`
	MaxRoundsPerGame  int `json:"maxRoundsPerGame"`
	IntervalSeconds   int `json:"intervalSeconds"`
}

func (g *Game) ResolveWinner() *Player {
	g.SortPlayersByPoints()
	return g.Players[0]
}

func (g *Game) SortPlayersByPoints() {
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
}

func (g *Game) SortPlayersByWinners() {
	sort.SliceStable(g.Players, func(i, j int) bool {
		a := g.Players[i]
		b := g.Players[j]

		if a.Winners == b.Winners {
			return a.Name < b.Name
		}

		return a.Winners > b.Winners
	})
}

func (g *Game) IncrementGamesPlayed() {
	for _, p := range g.Players {
		p.GamesPlayed++
	}
}

func (g *Game) ComputeScores(randomNumber int) {
	for _, p := range g.Players {
		p.ComputeScore(randomNumber)
	}
}

func (g *Game) ResolveWinnerByPoints() *Player {
	var winner *Player
	for _, p := range g.Players {
		if p.Points == magicNumber {
			return winner
		}
	}

	return nil
}

func (g *Game) GenerateRandomNumber() int {
	min := 1
	max := 10
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max+1-min) + min
}

type OverallRanking []Player

func (or OverallRanking) SortPlayersByWinners() {
	sort.SliceStable(or, func(i, j int) bool {
		a := or[i]
		b := or[j]

		if a.Winners == b.Winners {
			return a.Name < b.Name
		}

		return a.Winners > b.Winners
	})
}

type GameService interface {
	ServeWs(w http.ResponseWriter, r *http.Request)
	Run()
	Join(Player) (Player, error)
	GetRankingSnapshot() OverallRanking
	GetGameSnapshot() Game
}
