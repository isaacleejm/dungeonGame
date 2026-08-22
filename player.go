package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpellState int

const (
	MirrorState SpellState = iota
	SwordState
)

type Vector2 struct {
	X, Y float64
}

type Player struct {
	Pos          Vector2
	Rotation     float64
	Speed        float64
	MirrorSprite *ebiten.Image
	SwordSprite  *ebiten.Image
	SpellState   SpellState
}

func NewPlayer(
	screenWidth,
	screenHeight int,
	speed float64,
	mirrorSprite *ebiten.Image,
	swordSprite *ebiten.Image,
) *Player {
	bounds := mirrorSprite.Bounds()
	startX := float64(screenWidth/2 - bounds.Dx()/2)
	startY := float64(screenHeight/2 - bounds.Dy()/2)

	return &Player{
		Pos:          Vector2{X: startX, Y: startY},
		Speed:        speed,
		MirrorSprite: mirrorSprite,
		SwordSprite: swordSprite,
	}
}

func (p *Player) Center() Vector2 {
	bounds := p.MirrorSprite.Bounds()
	return Vector2{
		X: p.Pos.X + float64(bounds.Dx())/2,
		Y: p.Pos.Y + float64(bounds.Dy())/2,
	}
}

func (p *Player) Update(input InputState) {
	if input.HasAngleLock {
		p.Rotation = input.TargetAngle
	}

	if input.ToggleSpell {
		switch p.SpellState {
			case MirrorState:
				p.SpellState = SwordState
			case SwordState:
				p.SpellState = MirrorState
		}
	}

	dx, dy := input.MoveX, input.MoveY
	length := math.Hypot(dx, dy)
	if length > 0 {
		magnitude := math.Min(1.0, length)
		p.Pos.X += (dx / length) * p.Speed * magnitude
		p.Pos.Y += (dy / length) * p.Speed * magnitude
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.MirrorSprite.Bounds()
	width, height := float64(bounds.Dx()), float64(bounds.Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-width/2, -height/2)
	op.GeoM.Rotate(p.Rotation)
	op.GeoM.Translate(p.Pos.X+width/2, p.Pos.Y+height/2)

	switch p.SpellState {
		case MirrorState:
			screen.DrawImage(p.MirrorSprite, op)
		case SwordState:
			screen.DrawImage(p.SwordSprite, op)
	}
}
