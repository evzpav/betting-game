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
					Name:    "Harry",
					Points:  13,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Stace",
					Points:  14,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Evandro",
					Points:  12,
					Numbers: []int{9, 10},
				},
			},
		}

		winner := g.ResolveWinner()
		assert.Equal(t, "Stace", winner.Name)
		g.Players[1].Name = "Harry"
		g.Players[2].Name = "Evandro"
	})

	t.Run("diff name - same points, same upper and lower bounds", func(t *testing.T) {
		g := domain.Game{
			Players: []*domain.Player{
				{
					Name:    "Harry",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Stace",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Evandro",
					Points:  10,
					Numbers: []int{1, 5},
				},
			},
		}

		winner := g.ResolveWinner()
		assert.Equal(t, "Evandro", winner.Name)
	})

	t.Run("diff lower - same points, same name, same upper bounds", func(t *testing.T) {
		g := domain.Game{
			Players: []*domain.Player{
				{
					Name:    "Harry",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Harry",
					Points:  10,
					Numbers: []int{3, 5},
				},
				{
					Name:    "Harry",
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
					Name:    "Stace",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Stace",
					Points:  21,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Stace",
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
					Name:    "Evandro",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Evandro",
					Points:  10,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Evandro",
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
					Name:    "Evandro",
					Points:  40,
					Numbers: []int{1, 5},
				},
				{
					Name:    "John",
					Points:  21,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Rob",
					Points:  22,
					Numbers: []int{1, 6},
				},
			},
		}

		winner := g.ResolveWinnerByPoints()
		assert.Equal(t, "John", winner.Name)
	})

	t.Run("Two match 21", func(t *testing.T) {
		g := domain.Game{
			Rules: domain.Rules{MagicNumberMatch: 21},
			Players: []*domain.Player{
				{
					Name:    "Evandro",
					Points:  40,
					Numbers: []int{1, 5},
				},
				{
					Name:    "John",
					Points:  21,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Rob",
					Points:  21,
					Numbers: []int{1, 6},
				},
			},
		}

		winner := g.ResolveWinnerByPoints()
		assert.Equal(t, "Rob", winner.Name)
	})


	t.Run("None match 21", func(t *testing.T) {
		g := domain.Game{
			Rules: domain.Rules{MagicNumberMatch: 21},
			Players: []*domain.Player{
				{
					Name:    "Evandro",
					Points:  40,
					Numbers: []int{1, 5},
				},
				{
					Name:    "John",
					Points:  24,
					Numbers: []int{1, 5},
				},
				{
					Name:    "Rob",
					Points:  26,
					Numbers: []int{1, 6},
				},
			},
		}

		winner := g.ResolveWinnerByPoints()
		assert.Nil(t, winner)
	})

}
