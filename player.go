package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vector2 struct {
	X, Y float64
}

type Player struct {
	Pos      Vector2
	Rotation float64
	Speed    float64
	Sprite   *ebiten.Image
}

func NewPlayer(screenWidth, screenHeight int, speed float64, sprite *ebiten.Image) *Player {
	bounds := sprite.Bounds()
	startX := float64(screenWidth/2 - bounds.Dx()/2)
	startY := float64(screenHeight/2 - bounds.Dy()/2)

	return &Player{
		Pos:    Vector2{X: startX, Y: startY},
		Speed:  speed,
		Sprite: sprite,
	}
}

func (p *Player) Update(input InputState) {
	dx, dy := input.MoveX, input.MoveY
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}

	p.Rotation = math.Atan2(dy, dx)
	magnitude := math.Min(1.0, length)

	p.Pos.X += (dx / length) * p.Speed * magnitude
	p.Pos.Y += (dy / length) * p.Speed * magnitude
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.Sprite.Bounds()
	width, height := float64(bounds.Dx()), float64(bounds.Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-width/2, -height/2)
	op.GeoM.Rotate(p.Rotation)
	op.GeoM.Translate(p.Pos.X+width/2, p.Pos.Y+height/2)

	screen.DrawImage(p.Sprite, op)
}
