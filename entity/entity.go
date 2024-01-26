package entity

import (
	"github.com/google/uuid"
)

type Vector struct {
	X int
	Y int
}
func (v Vector) Add(other Vector) Vector{
	v.X += other.X
	v.Y += other.Y
	return v
}
type Rectangle struct {
	Pos    Vector
	Width  int
	Height int
}

func newBall(x int, y int) Rectangle {
	return Rectangle{Pos: Vector{X: x, Y: y}, Width: PlayerWidth, Height: PlayerWidth}
}
func newPaddle(x int, y int) Rectangle {
	return Rectangle{Pos: Vector{X: x, Y: y}, Width: PlayerWidth, Height: PlayerHeight}
}

func (r Rectangle) Overlaps(other Rectangle) bool {
	if r.Pos.Y+r.Height < other.Pos.Y || r.Pos.Y > other.Pos.Y+other.Height {
		return false
	}
	if r.Pos.X+r.Width < other.Pos.X || r.Pos.X > other.Pos.X+other.Width {
		return false
	}
	return true
}

type Player struct {
	Id   uuid.UUID
	Name string
}
type Game struct {
	LeftPaddle  Rectangle
	LeftScore   int
	RightPaddle Rectangle
	RightScore  int
	Ball        Rectangle
	BallDir     Vector
	Time        int
}

const GameWitdh = 100
const GameHeight = 100
const PlayerWidth = 2
const PlayerHeight = 5

type Simulation struct {
	game Game
}

func NewSimulation() Simulation {
	ball := newBall(0, 0)
	ballDir := Vector{X: 1, Y: 1}
	left := newPaddle(5, 0)
	right := newPaddle(95, 0)
	game := Game{Ball: ball, LeftPaddle: left, RightPaddle: right, BallDir: ballDir}
	return Simulation{game: game}
}

func (s *Simulation) reset() {
	s.game.Ball.Pos = Vector{X: 0, Y: 0}
	s.game.BallDir = Vector{X: 1, Y: 1}
	s.game.LeftPaddle.Pos = Vector{X: 5, Y: 0}
	s.game.RightPaddle.Pos = Vector{X: 95, Y: 0}
}

func (s *Simulation) UpdateLeft(y int) {
	s.game.LeftPaddle.Pos.Y = y
}
func (s *Simulation) UpdateRight(y int) {
	s.game.RightPaddle.Pos.Y = y
}

func (s *Simulation) Compute(dTime int) Game {
	s.game.Time += dTime
	newDir := Vector{X: s.game.BallDir.X * dTime, Y: s.game.BallDir.Y * dTime}
	newBall := newBall(s.game.Ball.Pos.X+newDir.X, s.game.Ball.Pos.Y+newDir.Y)
	for {
		if newBall.Pos.Y >= 0 && newBall.Pos.Y < 100 {
			if newBall.Pos.X <= 0 {
				s.game.RightScore += 1
				s.reset()
				return s.game
			}
			if newBall.Pos.X >= 100 {
				s.game.LeftScore += 1
				s.reset()
				return s.game
			}
			if (s.game.LeftPaddle.Pos.Y + s.game.LeftPaddle.Height) - 	
			s.game.Ball = newBall
			s.game.BallDir = newDir
			return s.game
		} else {
			newDir.Y *= -1

		}
	}

}
