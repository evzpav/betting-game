package game

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	hub         *domain.Hub
	log         log.Logger
	playersChan chan *domain.Player
	cron        *cron.Cron
}

func newHub() *domain.Hub {
	return &domain.Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *domain.Client),
		Unregister: make(chan *domain.Client),
		Clients:    make(map[*domain.Client]bool),
	}
}

func NewService(gameStorage domain.GameStorage, log log.Logger) *service {
	return &service{
		gameStorage: gameStorage,
		playersChan: make(chan *domain.Player),
		hub:         newHub(),
		log:         log,
	}
}

func (s *service) SetGameRules(minPlayersToStart, maxRoundsPerGame, intervalSeconds, magicNumberMatch int) {
	game := domain.NewGame(minPlayersToStart, maxRoundsPerGame, intervalSeconds, magicNumberMatch)
	s.gameStorage.Set(game)
	s.log.Info().Sendf("Game rules: %+v", game.Rules)
}

func (s *service) Run() {
	go s.RunHub()
	go s.WaitForPlayers()
}

func (s *service) RunHub() {
	for {
		select {
		case client := <-s.hub.Register:
			s.hub.Clients[client] = true
		case client := <-s.hub.Unregister:
			if _, ok := s.hub.Clients[client]; ok {
				delete(s.hub.Clients, client)
				close(client.Send)
			}
		case message := <-s.hub.Broadcast:
			for client := range s.hub.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(s.hub.Clients, client)
				}
			}
		}
	}
}

func (s *service) WaitForPlayers() {
	for p := range s.playersChan {
		game := s.gameStorage.Get()
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

func (s *service) Register(c *domain.Client) {
	s.hub.Register <- c
}

func (s *service) Unregister(c *domain.Client) {
	s.hub.Unregister <- c
}

func (s *service) Broadcast(messageType domain.MessageType, data interface{}) error {
	var wsMessage = domain.Message{
		MessageType: messageType,
		Data:        data,
	}

	bs, err := json.Marshal(wsMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	s.hub.Broadcast <- bs
	return nil
}

func (s *service) StartGame(game *domain.Game) {
	s.log.Debug().Send("Start game")

	s.ResetGame(game)

	game.ID = domain.GenerateNewID()
	game.GameRunning = true

	if err := s.Broadcast(domain.StartType, game); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
		return
	}

	s.startCron(game)

	s.gameStorage.Set(game)
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

	if err := s.Broadcast(domain.RoundType, game); err != nil {
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
		if err := s.Broadcast(domain.EndType, game); err != nil {
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

	if err := s.Broadcast(domain.OverallRankingType, ranking); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
	}

	time.Sleep(time.Duration(game.Rules.IntervalSeconds) * time.Second)

	s.ResetGame(game)
	s.startCron(game)
}

func (s *service) updateOverallRanking(game *domain.Game) domain.OverallRanking {
	ranking := make(map[string]domain.Player)

	for _, p := range game.Players {
		ranking[p.ID] = *p
	}

	var overallRanking domain.OverallRanking
	for _, p := range ranking {
		overallRanking = append(overallRanking, p)
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

func (s *service) ServeWs(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
		return
	}

	cli := domain.NewClient(s.hub, conn, s.log)
	s.Register(cli)

	s.WritePump(cli)
	s.ReadPump(cli)
}

func (s *service) WritePump(cli *domain.Client) {
	go cli.WritePump()
}

func (s *service) ReadPump(cli *domain.Client) {
	go cli.ReadPump(s.hub)
}

func (s *service) Join(player domain.Player) (domain.Player, error) {
	player.ID = domain.GenerateNewID()
	player.Name = strings.ToLower(player.Name)
	player.Observer = true

	game := s.gameStorage.Get()
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
	return *s.gameStorage.Get()
}
