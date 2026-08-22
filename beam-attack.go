package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type BeamAttack struct {
	pos         Vector2
	startPos 	Vector2
	rotation    float64
	startRotation float64
	attackState AttackState
	endPosition Vector2
	attackTimer float64
	Sprite *ebiten.Image
}

func (sa *BeamAttack) Update(input InputState, playerPos Vector2, playerRotation float64) {
	if input.StartBeamAttack && sa.attackState == NotAttacking {
		sa.attackState = MiddleAttack

		sa.pos = playerPos
		sa.startPos = playerPos

		sa.rotation = playerRotation
		sa.startRotation = playerRotation

		sa.attackTimer = 1
	}

	if sa.attackState == MiddleAttack {
		Range := 25.0

		endX := sa.startPos.X + math.Cos(sa.startRotation)*Range
		endY := sa.startPos.Y + math.Sin(sa.startRotation)*Range

		sa.endPosition = Vector2{endX, endY}

		sa.attackTimer -= 0.1

		if sa.attackTimer <= 0 {
			sa.attackState = NotAttacking
			return
		}
	}
}

func (sa *BeamAttack) Draw(screen *ebiten.Image){
	options := &ebiten.DrawImageOptions{}

	bounds := sa.Sprite.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	sideOffset := 60.0

	offsetX := math.Cos(sa.startRotation) * sideOffset
	offsetY := math.Sin(sa.startRotation) * sideOffset

	options.GeoM.Translate(-width/2, -height/2)
	options.GeoM.Scale(0.5, 1)
	options.GeoM.Rotate(sa.startRotation + math.Pi/2)
	options.GeoM.Translate(
		sa.startPos.X+offsetX,
		sa.startPos.Y+offsetY,
	)

	screen.DrawImage(sa.Sprite, options)
}