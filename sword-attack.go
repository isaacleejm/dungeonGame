package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type AttackState int

const (
	NotAttacking AttackState = iota
	MiddleAttack
)

type SwordAttack struct {
	pos         Vector2
	startPos 	Vector2
	rotation    float64
	startRotation float64
	attackState AttackState
	Sprite *ebiten.Image
}

func (sa *SwordAttack) Update(input InputState, playerPos Vector2, playerRotation float64) {
	if input.StartAttack && sa.attackState == NotAttacking {
		sa.attackState = MiddleAttack

		sa.pos = playerPos
		sa.startPos = playerPos

		sa.rotation = playerRotation
		sa.startRotation = playerRotation
	}

	if sa.attackState == MiddleAttack {
		speed := 1.0

		sa.pos.X += math.Cos(sa.startRotation) * speed
		sa.pos.Y += math.Sin(sa.startRotation) * speed

		dx := sa.pos.X - sa.startPos.X
		dy := sa.pos.Y - sa.startPos.Y

		if math.Hypot(dx, dy) >= 100.0 {
			sa.attackState = NotAttacking
		}
	}
}

func (sa *SwordAttack) Draw(screen *ebiten.Image){
	options := &ebiten.DrawImageOptions{}

	bounds := sa.Sprite.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	options.GeoM.Translate(-width/2, -height/2)
	options.GeoM.Rotate(sa.rotation)
	options.GeoM.Translate(sa.pos.X, sa.pos.Y)

	screen.DrawImage(sa.Sprite, options)
}