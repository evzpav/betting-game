package domain_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/evzpav/betting-game/internal/domain"
)

func TestValidatePlayer(t *testing.T) {
	t.Run("invalid numbers quantity", func(t *testing.T) {
		p := &domain.Player{
			Name:    "John",
			Numbers: []int{4},
		}

		assert.EqualError(t, p.Validate(), "invalid numbers quantity")
	})

	t.Run("invalid name", func(t *testing.T) {
		p := &domain.Player{
			Name:    "",
			Numbers: []int{4,2},
		}

		assert.EqualError(t, p.Validate(), "invalid name")
	})
}

func TestComputeScore(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		p := &domain.Player{
			Numbers: []int{4, 6},
		}

		p.ComputeScore(4)
		assert.Equal(t, p.Points, 5)
	})

	t.Run("inside bounds", func(t *testing.T) {
		p := &domain.Player{
			Numbers: []int{3, 8},
		}

		p.ComputeScore(7)
		assert.Equal(t, p.Points, 0)
	})

	t.Run("out of bounds", func(t *testing.T) {
		p := &domain.Player{
			Numbers: []int{4, 6},
		}

		p.ComputeScore(8)
		assert.Equal(t, p.Points, -1)
	})

	t.Run("multiple scores", func(t *testing.T) {
		p1 := &domain.Player{
			Numbers: []int{3, 8},
		}
		p2 := &domain.Player{
			Numbers: []int{5, 7},
		}
		p3 := &domain.Player{
			Numbers: []int{3, 7},
		}

		var players = []*domain.Player{
			p1, p2, p3,
		}

		expectedResults := map[int][]int{
			9:  {-1, -1, -1},
			1:  {-2, -2, -2},
			4:  {-2, -3, -1},
			10: {-3, -4, -2},
			7:  {-3, 1, 3},
			5:  {-3, 6, 4},
			3:  {2, 5, 9},
		}

		randomNumbers := []int{9, 1, 4, 10, 7, 5, 3}

		for _, n := range randomNumbers {

			for i, p := range players {
				p.ComputeScore(n)
				assert.Equal(t, expectedResults[n][i], p.Points, fmt.Sprintf("Random number is:%d", n))
			}
		}

	})

}
