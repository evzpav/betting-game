package game

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"gitlab.com/evzpav/betting-game/internal/domain"

	"gitlab.com/evzpav/betting-game/pkg/log"
)

const (
	minPlayersToStart = 2
	maxRoundsPerGame  = 10
	intervalSeconds   = 3
	magicNumber       = 21
)

type service struct {
	game *domain.Game
	hub  *domain.Hub
	log  log.Logger
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
	return &domain.Game{
		StopGame:    make(chan bool),
		PlayersChan: make(chan *domain.Player, 2),
	}
}

func NewService(log log.Logger) *service {
	return &service{
		game: newGame(),
		hub:  newHub(),
		log:  log,
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
				fmt.Printf("player out\n")
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
		case p := <-s.game.PlayersChan:
			s.game.Players = append(s.game.Players, p)

			if len(s.game.Players) >= minPlayersToStart && !s.game.GameRunning {
				s.StartGame()
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

func (s *service) Broadcast(msg []byte) {
	s.hub.Broadcast <- msg
}

func (s *service) StartGame() {
	s.log.Debug().Send("Start game")
	s.startCron()
}

func (s *service) startCron() {
	s.game.Cron = cron.New(cron.WithSeconds())

	s.game.Cron.AddFunc("*/1 * * * * *", s.runRound)

	s.game.GameRunning = true

	s.game.Cron.Start()
}

func (s *service) runRound() {

	s.game.RoundCounter++

	randomNumber := randomNumber()

	msg := fmt.Sprintf("Round %v: %v:\n", s.game.RoundCounter, randomNumber)
	fmt.Println(msg)

	s.hub.Broadcast <- []byte(msg)

	var winner *domain.Player
	for _, p := range s.game.Players {
		fmt.Printf("Players: %+v\n", p)
		score := p.ComputeScore(randomNumber)
		if score == magicNumber {
			winner = p
			break
		}
	}

	if winner != nil {
		winner.Winners++
		s.StopGame()
		return
	}

	if s.game.RoundCounter == maxRoundsPerGame {
		player := s.game.ResolveWinner()
		player.Winners++
		fmt.Printf("winner is: %v\n", player.Name)
		s.StopGame()
	}

}

func (s *service) StopGame() {
	s.game.Cron.Stop()
	s.game.Reset()

	s.hub.Broadcast <- []byte("finished")

	time.Sleep(intervalSeconds * time.Second)
	s.startCron()
}

// func (s *service) computeScores(n int) {
// 	for _, p := range s.game.Players {
// 		score := p.ComputeScore(n)
// 		if score == 21 {

// 		}
// 	}
// }

func randomNumber() int {
	min := 1
	max := 10
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max+1-min) + min
}

func (s *service) ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve ws")

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

func (s *service) Join(player domain.Player) error {
	fmt.Println(player)
	s.game.PlayersChan <- &player
	return nil
}
