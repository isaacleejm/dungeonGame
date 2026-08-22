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

type SingleAttack struct {
	pos         Vector2
	startPos 	Vector2
	endPos      Vector2
	rotation    float64
	startRotation float64
	attackTimer float64
	attackState AttackState
	SwordSprite *ebiten.Image
	BeamSprite *ebiten.Image
	spellState SpellState
}

func (sa *SingleAttack) Update(spellState SpellState, input InputState, playerPos Vector2, playerRotation float64) {
	sa.spellState = spellState;

	if spellState == SwordState{
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
	}else if spellState == MirrorState{
		if input.StartAttack && sa.attackState == NotAttacking {
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

			sa.endPos = Vector2{endX, endY}

			sa.attackTimer -= 0.1

			if sa.attackTimer <= 0 {
				sa.attackState = NotAttacking
				return
			}
		}
	}
}

func (sa *SingleAttack) Draw(screen *ebiten.Image){
	if sa.attackState != MiddleAttack{
		return
	}
	
	if sa.spellState == SwordState{
		options := &ebiten.DrawImageOptions{}

		bounds := sa.SwordSprite.Bounds()
		width := float64(bounds.Dx())
		height := float64(bounds.Dy())

		options.GeoM.Translate(-width/2, -height/2)
		options.GeoM.Rotate(sa.rotation)
		options.GeoM.Translate(sa.pos.X, sa.pos.Y)

		screen.DrawImage(sa.SwordSprite, options)
	} else if sa.spellState == MirrorState{
		options := &ebiten.DrawImageOptions{}

		bounds := sa.BeamSprite.Bounds()
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

		screen.DrawImage(sa.BeamSprite, options)
	}
}

