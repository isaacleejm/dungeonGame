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

	player := &Player{
		Pos:          Vector2{X: startX, Y: startY},
		Speed:        speed,
		MirrorSprite: mirrorSprite,
		SwordSprite:  swordSprite,
	}

	return player
}

func (p *Player) Center() Vector2 {
	bounds := p.MirrorSprite.Bounds()
	return Vector2{
		X: p.Pos.X + float64(bounds.Dx())/2,
		Y: p.Pos.Y + float64(bounds.Dy())/2,
	}
}

func (p *Player) CollisionBoundsMirror() Vector4{
	bounds := p.MirrorSprite.Bounds()

	return Vector4{
		X1: p.Pos.X,
		Y1: p.Pos.Y,
		X2: p.Pos.X + float64(bounds.Dx()),
		Y2: p.Pos.Y + float64(bounds.Dy()),
	}
}

func (p *Player) CollisionBoundsSword() Vector4{
	bounds := p.SwordSprite.Bounds()

	return Vector4{
		X1: p.Pos.X,
		Y1: p.Pos.Y,
		X2: p.Pos.X + float64(bounds.Dx()),
		Y2: p.Pos.Y + float64(bounds.Dy()),
	}
}

func (p *Player) Update(input InputState, blocks []*Block) {
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
	moveX := float64(0)
	moveY := float64(0)
	if length > 0 {
		magnitude := math.Min(1.0, length)
		moveX = (dx / length) * p.Speed * magnitude
		moveY = (dy / length) * p.Speed * magnitude
	}

	p.Pos.X += moveX

	for _, block := range blocks {
		if Collides(p.CollisionBoundsMirror(), block.CollisionBounds()) || Collides(p.CollisionBoundsSword(), block.CollisionBounds()) {
			p.Pos.X -= moveX
		}
	}

	p.Pos.Y += moveY

	for _, block := range blocks {
		if Collides(p.CollisionBoundsMirror(), block.CollisionBounds()) || Collides(p.CollisionBoundsSword(), block.CollisionBounds()) {
			p.Pos.Y -= moveY
		}
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
