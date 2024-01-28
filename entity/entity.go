package entity

import (
	"github.com/google/uuid"
)

type Vector struct {
	X int
	Y int
}

func (v Vector) Add(other Vector) Vector {
	v.X += other.X
	v.Y += other.Y
	return v
}
func (v Vector) Dot(other Vector) int {
	return v.X*other.X + v.Y*other.Y
}

type Rectangle struct {
	Pos    Vector
	Width  int
	Height int
}

func (r Rectangle) TopRight() Vector {
	y := r.Pos.Y 
	x := r.Pos.X + r.Width
	return Vector{X: x, Y: y}
}
func (r Rectangle) BottomLeft() Vector{
	y := r.Pos.Y + r.Height
	x := r.Pos.X
	return Vector{X: x, Y: y}
}

func newBall(x int, y int) Rectangle {
	return Rectangle{Pos: Vector{X: x, Y: y}, Width: PlayerWidth, Height: PlayerWidth}
}
func newPaddle(x int, y int) Rectangle {
	return Rectangle{Pos: Vector{X: x, Y: y}, Width: PlayerWidth, Height: PlayerHeight}
}

func (r Rectangle) Overlaps(other Rectangle) bool {
	if r.TopRight().Y < other.BottomLeft().Y || r.BottomLeft().Y > other.TopRight().Y {
		return false
	}
	if r.TopRight().X < other.BottomLeft().X || r.BottomLeft().X > other.TopRight().X {
		return false
	}
	return true
}

type Player struct {
	Id   uuid.UUID
	Name string
}
type GameState struct {
	LeftPaddle  Rectangle
	LeftScore   int
	RightPaddle Rectangle
	RightScore  int
	Ball        Rectangle
	BallDir     Vector
	Time        int
}

const GameWitdh = 200
const GameHeight = 100
const PlayerWidth = 4
const PlayerHeight = 14

type Simulation struct {
	game GameState
}

func NewSimulation() Simulation {
	ball := newBall(0, 0)
	ballDir := Vector{X: 1, Y: 1}
	left := newPaddle(0, 0)
	right := newPaddle(0, 0)
	game := GameState{Ball: ball, LeftPaddle: left, RightPaddle: right, BallDir: ballDir}
	simulation := Simulation{game: game}
	simulation.reset()
	return simulation
}

func (s *Simulation) reset() {
	x := GameWitdh/2 - PlayerWidth/2
	y := GameHeight/2 - PlayerWidth/2
	s.game.Ball.Pos = Vector{X: x, Y: y}
	s.game.BallDir = Vector{X: 1, Y: 1}
	x = PlayerWidth * 2
	y = GameHeight/2 - PlayerHeight/2
	s.game.LeftPaddle.Pos = Vector{X: x, Y: y}
	x = GameWitdh - PlayerWidth*2
	s.game.RightPaddle.Pos = Vector{X: x, Y: y}
}

func (s *Simulation) UpdateLeft(y int) {
	s.game.LeftPaddle.Pos.Y = y
}
func (s *Simulation) UpdateRight(y int) {
	s.game.RightPaddle.Pos.Y = y
}

func (s *Simulation) Compute(dTime int) GameState {
	s.game.Time += dTime
	newDir := Vector{X: s.game.BallDir.X * dTime, Y: s.game.BallDir.Y * dTime}
	newBall := newBall(s.game.Ball.Pos.X+newDir.X, s.game.Ball.Pos.Y+newDir.Y)
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
		if s.game.LeftPaddle.Overlaps(newBall) || s.game.RightPaddle.Overlaps(newBall) {
			s.game.BallDir.X *= -1
			return s.game
		}
		s.game.Ball = newBall
		s.game.BallDir = newDir
		return s.game
	} else {
		s.game.BallDir.Y *= -1
		return s.game
	}

}
