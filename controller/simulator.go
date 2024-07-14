package controller

import (
	"log"
	"pong/entity"
	"time"
)

type ServerMessage struct {
	Game entity.GameState
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
		var waiting Observer 
		isWaiting := false
		for player := range s.register {
			if !isWaiting {
				waiting = player
				isWaiting = true
				log.Println(waiting.Player.Id, " waits for a match")
			} else {
				startGame(waiting, player, s.ticks)
				isWaiting = false
			}
		}
	}()
}

func startGame(left Observer, right Observer, ticks int) {
	log.Println("Start match of left: ", *left.Player, " right: ", *right.Player)
	go func() {
		ticker := time.NewTicker(time.Duration(ticks) * time.Millisecond)
		lastTick := time.Now()
		leftInput := 0
		rightInput := 0
		simulation := entity.NewSimulation()
	EndGame:
		for {
			select {
			case t := <-ticker.C:
				dTime := t.Sub(lastTick)
				lastTick = t
				simulation.UpdateLeft(leftInput)
				simulation.UpdateRight(rightInput)
				state := simulation.Compute(int(dTime))
				update := ServerMessage{Game: state}
				left.Client <- update
				tmp := state.LeftPaddle
				state.LeftPaddle = state.RightPaddle
				state.RightPaddle = tmp
				update = ServerMessage{Game: state}
				right.Client <- update

				leftInput = 0
				rightInput = 0
			case input, more := <-left.Server:
				if !more {
					break EndGame
				}
				leftInput = input.Input
				log.Printf("Left paddle moved %v\n", input)
			case input, more := <-right.Server:
				if !more {
					break EndGame
				}
				rightInput = input.Input
				log.Printf("Right paddle moved %v\n", input)
			}
		}
		log.Println("Game closed left: ", *left.Player, " right: ", *right.Player)
	}()
}
