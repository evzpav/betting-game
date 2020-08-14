package game

import (
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"gitlab.com/evzpav/betting-game/internal/domain"

	"gitlab.com/evzpav/betting-game/pkg/log"
)

type service struct {
	gameStorage domain.GameStorage
	hubService  domain.HubService
	log         log.Logger
	playersChan chan *domain.Player
	cron        *cron.Cron
}

func NewService(gameStorage domain.GameStorage, hubService domain.HubService, log log.Logger) *service {
	return &service{
		gameStorage: gameStorage,
		playersChan: make(chan *domain.Player),
		hubService: hubService,
		log:        log,
	}
}

func (s *service) SetGameRules(minPlayersToStart, maxRoundsPerGame, intervalSeconds, magicNumberMatch int) {
	game := domain.NewGame(minPlayersToStart, maxRoundsPerGame, intervalSeconds, magicNumberMatch)
	s.gameStorage.SetGame(game)
	s.log.Info().Sendf("Game rules: %+v", game.Rules)
}

func (s *service) Run() {
	go s.hubService.RunHub()
	go s.WaitForPlayers()
}

func (s *service) WaitForPlayers() {
	for p := range s.playersChan {
		game := s.gameStorage.GetGame()
		game.Observers = append(game.Observers, p)

		var once sync.Once
		canStart := func() {
			if len(game.Observers) >= game.Rules.MinPlayersToStart {
				s.StartGame(game)
			}
		}

		once.Do(canStart)
	}
}

func (s *service) StartGame(game *domain.Game) {
	s.log.Debug().Send("Start game")

	s.ResetGame(game)

	if err := s.hubService.Broadcast(domain.StartType, game); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
		return
	}

	s.startCron(game)

	s.gameStorage.SetGame(game)
}

func (s *service) startCron(game *domain.Game) {
	game.GameCounter++
	game.GameRunning = true

	s.cron = cron.New(cron.WithSeconds())
	everySecond := "*/1 * * * * *"
	_, err := s.cron.AddFunc(everySecond, func() {
		s.runRound(game)
	})
	if err != nil {
		s.log.Error().Err(err).Sendf("failed to add cron func %v", err)
	}

	s.cron.Start()
}

func (s *service) runRound(game *domain.Game) {

	game.RoundCounter++
	randomNumber := game.GenerateRandomNumber()
	game.RandomNumber = randomNumber
	s.log.Debug().Sendf("Round:%v; Number:%v", game.RoundCounter, randomNumber)

	game.ComputeScores(randomNumber)
	game.SortPlayersByPoints()

	if err := s.hubService.Broadcast(domain.RoundType, game); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
		return
	}

	var winner *domain.Player
	if winner = game.ResolveWinnerByPoints(); winner == nil {
		if game.RoundCounter >= game.Rules.MaxRoundsPerGame {
			winner = game.ResolveWinner()
		}
	}

	if winner != nil {
		winner.Winners++
		game.Winner = winner
		if err := s.hubService.Broadcast(domain.EndType, game); err != nil {
			s.log.Error().Err(err).Sendf("%v", err)
		}

		s.stopGame(game)
	}

}

func (s *service) stopGame(game *domain.Game) {
	s.cron.Stop()
	game.GameRunning = false
	game.IncrementGamesPlayed()

	ranking := s.updateOverallRanking(game)

	if err := s.hubService.Broadcast(domain.OverallRankingType, ranking); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
	}

	s.intervalTicker(game)

	s.ResetGame(game)
	s.startCron(game)
}

func (s *service) intervalTicker(game *domain.Game) {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	counter := game.Rules.IntervalSeconds + 1
	go func(counter int) {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				counter--

				s.log.Debug().Sendf("Interval: %d", counter)

				if err := s.hubService.Broadcast(domain.IntervalTickerType, counter); err != nil {
					s.log.Error().Err(err).Sendf("%v", err)
				}
			}
		}
	}(counter)

	time.Sleep(time.Duration(game.Rules.IntervalSeconds) * time.Second)
	ticker.Stop()
	done <- true
}

func (s *service) updateOverallRanking(game *domain.Game) domain.OverallRanking {

	var overallRanking domain.OverallRanking
	for _, p := range game.Players {
		overallRanking = append(overallRanking, *p)
	}

	overallRanking.SortPlayersByWinners()

	s.gameStorage.SetOverallRanking(overallRanking)

	return overallRanking
}

func (s *service) ResetGame(game *domain.Game) {
	game.Reset()

	for _, p := range game.Players {
		p.Observer = false
		p.ResetPoints()
	}
}

func (s *service) Join(player domain.Player) (domain.Player, error) {
	player.ID = domain.GenerateNewID()
	player.Name = strings.ToLower(player.Name)
	player.Observer = true

	game := s.gameStorage.GetGame()
	if err := game.IsNameInUse(player.Name); err != nil {
		return domain.Player{}, err
	}

	s.playersChan <- &player

	return player, nil
}

func (s *service) GetRankingSnapshot() domain.OverallRanking {
	return s.gameStorage.GetOverallRanking()
}

func (s *service) GetGameSnapshot() domain.Game {
	return *s.gameStorage.GetGame()
}

func (s *service) AddNewWebsocketClient(conn *websocket.Conn) {
	s.hubService.AddNewWebsocketClient(conn)
}
