package domain

import (
	"errors"
	"sort"
)

type Player struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Points      int    `json:"points"`
	Numbers     []int  `json:"numbers"`
	Winners     int    `json:"winners"`
	GamesPlayed int    `json:"gamesPlayed"`
	Observer    bool   `json:"observer"`
}

func (p *Player) Validate() error {
	if len(p.Numbers) != 2 {
		return errors.New("invalid numbers quantity")
	}

	if p.Name == "" {
		return errors.New("invalid name")
	}

	return nil
}

func (p *Player) ComputeScore(n int) int {
	sort.Ints(p.Numbers)
	lower := p.Numbers[0]
	upper := p.Numbers[1]

	if n < lower || n > upper {
		p.Points--
	}

	if n > lower && n < upper {
		p.Points += 5 - (upper - lower)
	}

	if n == lower || n == upper {
		p.Points += 5
	}

	return p.Points
}

func (p *Player) ResetPoints() {
	p.Points = 0
}
