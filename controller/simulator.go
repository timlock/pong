package controller

import (
	"log"
	"pong/entity"
	"time"
)

type ServerMessage struct {
	Game entity.Game
}

type ClientMessage struct {
	Input int
}

type Lobby struct {
	register chan Observer
	ticks    int
}
type Observer struct {
	Player *entity.Player
	Client chan ServerMessage
	Server chan ClientMessage
}

func NewSimulator(ticks int) Lobby {
	return Lobby{register: make(chan Observer), ticks: ticks}
}
func (s Lobby) Add(observer Observer) {
	s.register <- observer
}
func (s Lobby) Run() {
	go func() {
		var waiting *Observer = nil
		for player := range s.register {
			if waiting == nil {
				waiting = &player
				log.Println(waiting.Player.Id, " waits for a match")
			} else {
				startGame(*waiting, player, s.ticks)
				waiting = nil

			}
		}
	}()
}

func startGame(left Observer, right Observer, ticks int) {
	log.Println("Start match of left: ", *left.Player, " right: ", *right.Player)
	go func() {
		ticker := time.NewTicker(time.Duration(ticks) * time.Millisecond)
		ball := entity.Vector{X: 0, Y: 0}
		ballDir := entity.Vector{X: 50, Y: 1}
		leftPos := entity.Vector{X: 5, Y: 0}
		rightPos := entity.Vector{X: 95, Y: 0}
		leftInput := 0
		rightInput := 0
		lastTick := time.Now()
	EndGame:
		for {
			select {
			case t := <-ticker.C:
				dTime := t.Sub(lastTick)
				lastTick = t
				leftPos.Y += leftInput * int(dTime)
				rightPos.Y += rightInput * int(dTime)
				ball.X += ballDir.X * int(dTime)
				ball.Y += ballDir.Y * int(dTime)
				leftInput = 0
				rightInput = 0

				game := entity.Game{LeftPaddle: leftPos, RightPaddle: rightPos, Ball: ball}
				update := ServerMessage{Game: game}
				left.Client <- update
				right.Client <- update
			case input, more := <-left.Server:
				if !more {
					break EndGame
				}
				leftInput = input.Input
			case input, more := <-right.Server:
				if !more {
					break EndGame
				}
				rightInput = input.Input
			}
		}
		log.Println("Game closed left: ", *left.Player, " right: ", *right.Player)
	}()
}