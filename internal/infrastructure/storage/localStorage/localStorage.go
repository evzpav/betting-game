package localstorage

import (
	"gitlab.com/evzpav/betting-game/internal/domain"
	"gitlab.com/evzpav/betting-game/pkg/log"
)

type gameStorage struct {
	game           *domain.Game
	overallRanking domain.OverallRanking
	log            log.Logger
}

func NewGameStorage(log log.Logger) *gameStorage {
	return &gameStorage{
		game:           &domain.Game{},
		overallRanking: domain.OverallRanking{},
		log:            log,
	}
}

func (gs *gameStorage) GetGame() *domain.Game {
	return gs.game
}

func (gs *gameStorage) SetGame(game *domain.Game) {
	gs.game = game
}

func (gs *gameStorage) GetOverallRanking() domain.OverallRanking {
	return gs.overallRanking
}

func (gs *gameStorage) SetOverallRanking(or domain.OverallRanking) {
	gs.overallRanking = or
}
