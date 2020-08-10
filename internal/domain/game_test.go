package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/betting-game/internal/domain"
)

func TestResolveWinner(t *testing.T) {

	t.Run("check order", func(t *testing.T) {
		g := domain.Game{
			Players: []*domain.Player{
				{
					Name:    "harry",
					Points:  13,
					Numbers: []int{1, 5},
				},
				{
					Name:    "stace",
					Points:  14,
					Numbers: []int{1, 5},
				},
				{
					Name:    "evandro",
					Points:  12,
					Numbers: []int{9, 10},
				},
			},
		}

		winner := g.ResolveWinner()
		assert.Equal(t, "stace", winner.Name)
		g.Players[1].Name = "harry"
		g.Players[2].Name = "evandro"
	})

	t.Run("diff name - same points, same upper and lower bounds", func(t *testing.T) {
		g := domain.Game{
			Players: []*domain.Player{
				{
					Name:    "harry",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "stace",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "evandro",
					Points:  10,
					Numbers: []int{1, 5},
				},
			},
		}

		winner := g.ResolveWinner()
		assert.Equal(t, "evandro", winner.Name)
	})

	t.Run("diff lower - same points, same name, same upper bounds", func(t *testing.T) {
		g := domain.Game{
			Players: []*domain.Player{
				{
					Name:    "harry",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "harry",
					Points:  10,
					Numbers: []int{3, 5},
				},
				{
					Name:    "harry",
					Points:  10,
					Numbers: []int{2, 5},
				},
			},
		}

		winner := g.ResolveWinner()
		assert.Equal(t, 3, winner.Numbers[0])
	})

	t.Run("diff points - same name, same bounds", func(t *testing.T) {
		g := domain.Game{
			Players: []*domain.Player{
				{
					Name:    "stace",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "stace",
					Points:  21,
					Numbers: []int{1, 5},
				},
				{
					Name:    "stace",
					Points:  12,
					Numbers: []int{1, 5},
				},
			},
		}

		winner := g.ResolveWinner()
		assert.Equal(t, 21, winner.Points)
	})

	t.Run("diff upper - same name, same points, same lower", func(t *testing.T) {
		g := domain.Game{
			Players: []*domain.Player{
				{
					Name:    "evandro",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "evandro",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "evandro",
					Points:  10,
					Numbers: []int{1, 6},
				},
			},
		}

		winner := g.ResolveWinner()
		assert.Equal(t, 6, winner.Numbers[1])
	})

}

func TestResolveWinnerByPoints(t *testing.T) {

	t.Run("One match 21", func(t *testing.T) {
		g := domain.Game{
			Rules: domain.Rules{MagicNumberMatch: 21},
			Players: []*domain.Player{
				{
					Name:    "evandro",
					Points:  40,
					Numbers: []int{1, 5},
				},
				{
					Name:    "john",
					Points:  21,
					Numbers: []int{1, 5},
				},
				{
					Name:    "rob",
					Points:  22,
					Numbers: []int{1, 6},
				},
			},
		}

		winner := g.ResolveWinnerByPoints()
		assert.Equal(t, "john", winner.Name)
	})

	t.Run("Two match 21", func(t *testing.T) {
		g := domain.Game{
			Rules: domain.Rules{MagicNumberMatch: 21},
			Players: []*domain.Player{
				{
					Name:    "evandro",
					Points:  40,
					Numbers: []int{1, 5},
				},

				{
					Name:    "carlos",
					Points:  21,
					Numbers: []int{1, 5},
				},
				{
					Name:    "vinicius",
					Points:  21,
					Numbers: []int{1, 5},
				},
			},
		}

		winner := g.ResolveWinnerByPoints()
		assert.Equal(t, "carlos", winner.Name)
	})

	t.Run("None match 21", func(t *testing.T) {
		g := domain.Game{
			Rules: domain.Rules{MagicNumberMatch: 21},
			Players: []*domain.Player{
				{
					Name:    "evandro",
					Points:  40,
					Numbers: []int{1, 5},
				},
				{
					Name:    "john",
					Points:  24,
					Numbers: []int{1, 5},
				},
				{
					Name:    "rob",
					Points:  26,
					Numbers: []int{1, 6},
				},
			},
		}

		winner := g.ResolveWinnerByPoints()
		assert.Nil(t, winner)
	})

}
