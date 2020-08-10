package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/evzpav/betting-game/internal/domain"

	"gitlab.com/evzpav/betting-game/pkg/log"
)

type service struct {
	overallRanking map[string]domain.Player
	game           *domain.Game
	hub            *domain.Hub
	log            log.Logger
	playersChan    chan *domain.Player
	cron           *cron.Cron
}

func newHub() *domain.Hub {
	return &domain.Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *domain.Client),
		Unregister: make(chan *domain.Client),
		Clients:    make(map[*domain.Client]bool),
	}
}

func NewService(log log.Logger) *service {
	return &service{
		overallRanking: make(map[string]domain.Player),
		playersChan:    make(chan *domain.Player),
		hub:            newHub(),
		log:            log,
	}
}

func (s *service) SetGameRules(minPlayersToStart, maxRoundsPerGame, intervalSecond, magicNumberMatch int) {
	s.game = domain.NewGame(minPlayersToStart, maxRoundsPerGame, intervalSecond, magicNumberMatch)
	s.log.Info().Sendf("Game rules: %+v", s.game.Rules)
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
			fmt.Printf("UNREGISTER: %+v\n", s.hub.Clients)
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
	for {
		select {
		case p := <-s.playersChan:
			s.game.Observers = append(s.game.Observers, p)

			var once sync.Once
			canStart := func() {
				if len(s.game.Observers) >= s.game.Rules.MinPlayersToStart {
					s.StartGame()
				}
			}

			once.Do(canStart)

		}
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

func (s *service) StartGame() {
	s.log.Debug().Send("Start game")

	s.ResetGame()

	s.game.ID = generateNewID()

	s.Broadcast(domain.StartType, s.game.ID)

	s.startCron()
}

func (s *service) startCron() {
	s.game.GameCounter++
	s.game.GameRunning = true

	s.cron = cron.New(cron.WithSeconds())
	everySecond := "*/1 * * * * *"
	s.cron.AddFunc(everySecond, s.runRound)

	s.cron.Start()
}

func (s *service) runRound() {

	s.game.RoundCounter++
	randomNumber := s.game.GenerateRandomNumber()

	s.game.RandomNumber = randomNumber
	s.log.Debug().Sendf("Round:%v; Number:%v", s.game.RoundCounter, randomNumber)

	s.game.ComputeScores(randomNumber)
	s.game.SortPlayersByPoints()

	if err := s.Broadcast(domain.RoundType, s.game); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
		return
	}

	var winner *domain.Player
	if winner = s.game.ResolveWinnerByPoints(); winner == nil {
		if s.game.RoundCounter >= s.game.Rules.MaxRoundsPerGame {
			winner = s.game.ResolveWinner()
		}
	}

	if winner != nil {
		winner.Winners++
		s.game.Winner = winner
		if err := s.Broadcast(domain.EndType, s.game); err != nil {
			s.log.Error().Err(err).Sendf("%v", err)
		}

		s.stopGame()
	}

}

func (s *service) updateOverallRanking() domain.OverallRanking {
	for _, p := range s.game.Players {
		s.overallRanking[p.ID] = *p
	}

	return s.GetRankingSnapshot()
}

func (s *service) stopGame() {
	s.cron.Stop()
	s.game.GameRunning = false
	s.game.IncrementGamesPlayed()

	ranking := s.updateOverallRanking()

	if err := s.Broadcast(domain.OverallRankingType, ranking); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
	}

	time.Sleep(time.Duration(s.game.Rules.IntervalSeconds) * time.Second)

	s.ResetGame()
	s.startCron()
}

func (s *service) ResetGame() {
	s.game.Winner = nil
	s.game.RoundCounter = 0

	s.game.Players = append(s.game.Players, s.game.Observers...)
	s.game.Observers = make([]*domain.Player, 0)

	for _, p := range s.game.Players {
		p.Observer = false
		p.ResetPoints()
	}

	if err := s.Broadcast(domain.RestartType, true); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
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

	cli := domain.NewClient(s.hub, conn)
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

func (s *service) RegisterNewClient(conn *websocket.Conn) {
	cli := domain.NewClient(s.hub, conn)
	s.Register(cli)

	s.WritePump(cli)
	s.ReadPump(cli)
}

func (s *service) Join(player domain.Player) (domain.Player, error) {
	player.ID = generateNewID()

	if s.game.GameRunning {
		player.Observer = true
	}

	for _, p := range s.game.Players {
		if p.Name == player.Name {
			return domain.Player{}, errors.New("name already in use")
		}
	}

	s.playersChan <- &player

	return player, nil
}

func (s *service) GetRankingSnapshot() domain.OverallRanking {
	var ranking domain.OverallRanking
	for _, p := range s.overallRanking {
		ranking = append(ranking, p)
	}

	ranking.SortPlayersByWinners()
	return ranking
}

func (s *service) GetGameSnapshot() domain.Game {
	return *s.game
}

func generateNewID() string {
	sID := uuid.NewV4()
	return sID.String()
}
