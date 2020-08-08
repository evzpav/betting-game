package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/evzpav/betting-game/internal/domain"

	"gitlab.com/evzpav/betting-game/pkg/log"
)

const (
	minPlayersToStart int = 2
	maxRoundsPerGame  int = 10
	intervalSeconds   int = 3
)

type service struct {
	overallRanking map[string]domain.Player
	game           *domain.Game
	hub            *domain.Hub
	log            log.Logger
	PlayersChan    chan *domain.Player
	cron           *cron.Cron
}

func newHub() *domain.Hub {
	return &domain.Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *domain.Client, 2),
		Unregister: make(chan *domain.Client),
		Clients:    make(map[*domain.Client]bool),
	}
}

func newGame() *domain.Game {
	return &domain.Game{}
}

func NewService(log log.Logger) *service {
	return &service{
		overallRanking: make(map[string]domain.Player),
		PlayersChan:    make(chan *domain.Player, 2),
		game:           newGame(),
		hub:            newHub(),
		log:            log,
	}
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
	for {
		select {
		case p := <-s.PlayersChan:
			if s.game.GameRunning {
				s.game.Observers = append(s.game.Observers, p)
			} else {
				s.game.Players = append(s.game.Players, p)

				if len(s.game.Players) >= minPlayersToStart {
					s.StartGame()
				}
			}

		}
	}
}

func (s *service) Register(c *domain.Client) {
	s.hub.Register <- c
}

func (s *service) Unregister(c *domain.Client) {
	s.hub.Unregister <- c
}

func (s *service) Broadcast(messageType string, data interface{}) error {
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

	s.game.ID = generateNewID()

	s.Broadcast("start", "start")

	s.startCron()
}

func (s *service) startCron() {
	s.game.GameRunning = true
	s.cron = cron.New(cron.WithSeconds())
	everySecond := "*/1 * * * * *"
	s.cron.AddFunc(everySecond, s.runRound)

	s.cron.Start()
}

func (s *service) runRound() {

	s.game.RoundCounter++

	randomNumber := s.game.GenerateRandomNumber()

	msg := fmt.Sprintf("Round %v: %v:\n", s.game.RoundCounter, randomNumber)
	s.log.Info().Send(msg)

	s.game.ComputeScores(randomNumber)
	s.game.SortPlayersByPoints()

	if err := s.Broadcast("round", s.game); err != nil {
		s.log.Error().Err(err).Sendf("%v", err)
		return
	}

	if winner := s.game.ResolveWinnerByPoints(); winner != nil {
		winner.Winners++
		s.stopGame()
		s.Broadcast("end", winner)
		return
	}

	if s.game.RoundCounter >= maxRoundsPerGame {
		winner := s.game.ResolveWinner()
		winner.Winners++
		fmt.Printf("winner is: %v\n", winner.Name)

		s.stopGame()
		s.Broadcast("end", winner)
	}

}

func (s *service) stopGame() {
	s.game.IncrementGamesPlayed()

	for _, p := range s.game.Players {
		s.overallRanking[p.ID] = *p
	}

	var ranking domain.OverallRanking
	for _, p := range s.overallRanking {
		ranking = append(ranking, p)
	}
	ranking.SortPlayersByWinners()
	_ = s.Broadcast("overallranking", ranking)

	fmt.Printf("Ranking: %v\n", s.overallRanking)
	// s.overallRanking = append(s.overallRanking, s.game.Players...)

	s.game.GameRunning = false
	s.cron.Stop()

	s.ResetGame()

	_ = s.Broadcast("finished", "finished")

	time.Sleep(time.Duration(intervalSeconds) * time.Second)
	s.startCron()
}

func (s *service) ResetGame() {
	s.game.RoundCounter = 0

	s.game.Players = append(s.game.Players, s.game.Observers...)
	s.game.Observers = make([]*domain.Player, 0)

	for _, p := range s.game.Players {
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

func (s *service) CloseSend(cli *domain.Client) {
	close(cli.Send)
}

func (s *service) RegisterNewClient(conn *websocket.Conn) {
	cli := domain.NewClient(s.hub, conn)
	s.Register(cli)

	s.WritePump(cli)
	s.ReadPump(cli)
}

func (s *service) Join(player domain.Player) {
	player.ID = generateNewID()
	s.PlayersChan <- &player
}

func generateNewID() string {
	sID := uuid.NewV4()
	return sID.String()
}
