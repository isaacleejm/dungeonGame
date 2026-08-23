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

type BeamAttack struct {
	pos           Vector2
	startPos      Vector2
	endPos        Vector2
	rotation      float64
	startRotation float64
	attackTimer   float64
	attackState   AttackState
	BeamSprite    *ebiten.Image
	spellState    SpellState
	BeamNodes     []Vector2
}

func (sa *BeamAttack) Update(
	spellState SpellState,
	input InputState,
	playerPos Vector2,
	playerRotation float64,
	mirrors []Mirror,
) {
	sa.spellState = spellState
	if spellState == MirrorState {
		if input.StartAttack && sa.attackState == NotAttacking {
			sa.attackState = MiddleAttack

			sa.pos = playerPos
			sa.startPos = playerPos

			sa.rotation = playerRotation
			sa.startRotation = playerRotation

			sa.attackTimer = 1
		}

		if sa.attackState != MiddleAttack {
			return
		}

		Range := 25.0

		endX := sa.startPos.X + math.Cos(sa.startRotation)*Range
		endY := sa.startPos.Y + math.Sin(sa.startRotation)*Range

		sa.endPos = Vector2{endX, endY}

		sa.attackTimer -= 0.1

		if sa.attackTimer <= 0 {
			sa.attackState = NotAttacking
			return
		}

		// Calculate Beam Bounces
		sa.BeamNodes = []Vector2{sa.pos}
		
		currentO := sa.pos
		currentD := Vector2{X: math.Cos(sa.startRotation), Y: math.Sin(sa.startRotation)}
		maxBounces := 5

		for range maxBounces {
			var closestMirror *Mirror
			minT := math.MaxFloat64
			var closestHit Vector2

			// Find the closest mirror intersection
			for i := range mirrors {
				m := &mirrors[i]
				
				const mw float64 = mirrorWidth 
				dx := math.Cos(m.Rotation) * (mw / 2)
				dy := math.Sin(m.Rotation) * (mw / 2)
				
				A := Vector2{X: m.Pos.X - dx, Y: m.Pos.Y - dy}
				B := Vector2{X: m.Pos.X + dx, Y: m.Pos.Y + dy}

				hit, t, hitPos := rayIntersectsSegment(currentO, currentD, A, B)
				if hit && t < minT {
					minT = t
					closestHit = hitPos
					closestMirror = m
				}
			}

			// Handle the hit
			if closestMirror != nil {
				sa.BeamNodes = append(sa.BeamNodes, closestHit)

				// Calculate mirror normal vector (perpendicular to its rotation)
				nx := -math.Sin(closestMirror.Rotation)
				ny := math.Cos(closestMirror.Rotation)

				// Ensure the normal faces the incoming beam
				dot := currentD.X*nx + currentD.Y*ny
				if dot > 0 {
					nx = -nx
					ny = -ny
					dot = -dot
				}

				// Law of Reflection: R = D - 2(D·N)N
				rx := currentD.X - 2*dot*nx
				ry := currentD.Y - 2*dot*ny

				// Set up the next ray
				currentD = Vector2{X: rx, Y: ry}
				currentO = closestHit 
			} else {
				// If no hit, extend out to max range (1000 pixels) and stop bouncing
				finalPos := Vector2{
					X: currentO.X + currentD.X*1000,
					Y: currentO.Y + currentD.Y*1000,
				}
				sa.BeamNodes = append(sa.BeamNodes, finalPos)
				break
			}
		}
	}
}

func (sa *BeamAttack) Draw(screen *ebiten.Image) {
	if sa.attackState != MiddleAttack {
		return
	}

	if sa.spellState == MirrorState {
		DrawBeamPath(screen, sa.BeamSprite, sa.BeamNodes, 0.2)
	}
}

func DrawBeamPath(
	screen *ebiten.Image,
	sprite *ebiten.Image,
	nodes []Vector2,
	thickness float64,
) {
	if len(nodes) < 2 || sprite == nil {
		return
	}

	bounds := sprite.Bounds()
	beamWidth := float64(bounds.Dx())
	beamHeight := float64(bounds.Dy())

	for i := 0; i < len(nodes)-1; i++ {
		A := nodes[i]
		B := nodes[i+1]

		dx := B.X - A.X
		dy := B.Y - A.Y
		dist := math.Hypot(dx, dy)
		angle := math.Atan2(dy, dx)

		scaleX := dist / beamWidth

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, -beamHeight/2)
		op.GeoM.Scale(scaleX, thickness)
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(A.X, A.Y)
		op.Blend = ebiten.BlendLighter

		screen.DrawImage(sprite, op)
	}
}

func rayIntersectsSegment(O, D, A, B Vector2) (bool, float64, Vector2) {
	E := Vector2{X: B.X - A.X, Y: B.Y - A.Y}
	det := D.Y*E.X - D.X*E.Y

	// If det is 0, the ray and mirror are parallel
	if math.Abs(det) < 0.0001 {
		return false, 0, Vector2{}
	}

	dx := A.X - O.X
	dy := A.Y - O.Y

	t := (dy*E.X - dx*E.Y) / det
	u := (D.X*dy - D.Y*dx) / det

	// t > 0.01 ensures we don't instantly hit the mirror we just bounced off
	// 0 <= u <= 1 ensures we actually hit the segment, not empty space next to it
	if t > 0.01 && u >= 0 && u <= 1 {
		hit := Vector2{X: O.X + t*D.X, Y: O.Y + t*D.Y}
		return true, t, hit
	}
	
	return false, 0, Vector2{}
}
