package domain

import (
	"errors"
	"sort"

	"github.com/gorilla/websocket"
)

type Game struct {
	ID           string    `json:"id"`
	Players      []*Player `json:"players"`
	Observers    []*Player `json:"-"`
	Winner       *Player   `json:"winner,omitempty"`
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
	MagicNumberMatch  int `json:"magicNumberMatch"`
}

func NewGame(minPlayersToStart, maxRoundsPerGame, intervalSeconds, magicNumberMatch int) *Game {
	rules := Rules{
		MinPlayersToStart: minPlayersToStart,
		MaxRoundsPerGame:  maxRoundsPerGame,
		IntervalSeconds:   intervalSeconds,
		MagicNumberMatch:  magicNumberMatch,
	}

	return &Game{
		Rules: rules,
	}
}

func (g *Game) IsNameInUse(name string) error {
	for _, p := range g.Observers {
		if p.Name == name {
			return errors.New("name already in use")
		}
	}

	for _, p := range g.Players {
		if p.Name == name {
			return errors.New("name already in use")
		}
	}

	return nil
}

func (g *Game) ResolveWinner() *Player {
	g.SortPlayersByPoints()
	return g.Players[0]
}

func (g *Game) SortPlayersByPoints() {
	sortPlayersByPoints(g.Players)
}

func (g *Game) SortPlayersByWinners() {
	sortPlayersByWinners(g.Players)
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
	var matchedMagic []*Player
	for _, p := range g.Players {
		if p.Points == g.Rules.MagicNumberMatch {
			matchedMagic = append(matchedMagic, p)
		}
	}

	if len(matchedMagic) == 1 {
		return matchedMagic[0]
	}

	if len(matchedMagic) > 1 {
		sortPlayersByPoints(matchedMagic)
		return matchedMagic[0]
	}

	return nil
}

func (g *Game) Reset() {
	g.ID = GenerateNewID()
	g.Winner = nil
	g.RoundCounter = 0

	g.Players = append(g.Players, g.Observers...)
	g.Observers = make([]*Player, 0)
}

func (g *Game) GenerateRandomNumber() int {
	min := 1
	max := 10

	return GenerateRandomNumber(min, max+1)
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

func sortPlayersByWinners(players []*Player) {
	sort.SliceStable(players, func(i, j int) bool {
		a := players[i]
		b := players[j]

		if a.Winners == b.Winners {
			return a.Name < b.Name
		}

		return a.Winners > b.Winners
	})
}

func sortPlayersByPoints(players []*Player) {
	sort.SliceStable(players, func(i, j int) bool {
		a := players[i]
		b := players[j]

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

type GameService interface {
	AddNewWebsocketClient(conn *websocket.Conn)
	Run()
	Join(Player) (Player, error)
	GetGameSnapshot() Game
	GetRankingSnapshot() OverallRanking
}
type GameStorage interface {
	GetGame() *Game
	SetGame(*Game)
	GetOverallRanking() OverallRanking
	SetOverallRanking(OverallRanking)
}
