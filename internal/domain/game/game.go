package game

import "gitlab.com/evzpav/betting-game/pkg/log"

type service struct {
	log log.Logger
}

func NewService(log log.Logger) *service {
	return &service{
		log: log,
	}
}

func (s *service) StartGame() {
	s.log.Info().Send("START GAME")
}
