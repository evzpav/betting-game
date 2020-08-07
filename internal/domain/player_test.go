package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/evzpav/betting-game/internal/domain"
)

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
}
