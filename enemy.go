package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	Pos      Vector2
	Rotation float64
	Speed    float64
	Sprite   *ebiten.Image
}

func NewEnemy(screenWidth, screenHeight int, speed float64, sprite *ebiten.Image, player *Player) *Enemy {
	enemyBounds := sprite.Bounds()

	enemyWidth := float64(enemyBounds.Dx())
	enemyHeight := float64(enemyBounds.Dy())

	playerCenter := player.Center()

	// Distance from player enemy can spawn
	const minDistance = 100.0

	var startX, startY float64

	for {
		// Generate a position that keeps the entire enemy inside the screen.
		startX = rand.Float64() * (float64(screenWidth) - enemyWidth)
		startY = rand.Float64() * (float64(screenHeight) - enemyHeight)

		// Calculate the center of the newly spawned enemy.
		enemyCenterX := startX + enemyWidth/2
		enemyCenterY := startY + enemyHeight/2

		// Distance between enemy center and player center.
		dx := enemyCenterX - playerCenter.X
		dy := enemyCenterY - playerCenter.Y

		distance := math.Hypot(dx, dy)

		// Break if enemy is far enough from player
		if distance >= minDistance {
			break
		}
	}

	return &Enemy{
		Pos:    Vector2{X: startX, Y: startY},
		Speed:  speed,
		Sprite: sprite,
	}
}

func (e *Enemy) Center() Vector2 {
	bounds := e.Sprite.Bounds()
	return Vector2{
		X: e.Pos.X + float64(bounds.Dx())/2,
		Y: e.Pos.Y + float64(bounds.Dy())/2,
	}
}

func (e *Enemy) Update(p *Player) {
	playerPos := p.Pos
	enemyPos := e.Pos

	dx := playerPos.X - enemyPos.X
	dy := playerPos.Y - enemyPos.Y

	distance := math.Hypot(dx, dy)

	if distance == 0 {
		return
	}

	e.Pos.X += (dx / distance) * e.Speed
	e.Pos.Y += (dy / distance) * e.Speed
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	bounds := e.Sprite.Bounds()
	width, height := float64(bounds.Dx()), float64(bounds.Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-width/2, -height/2)
	op.GeoM.Rotate(e.Rotation)
	op.GeoM.Translate(e.Pos.X+width/2, e.Pos.Y+height/2)

	screen.DrawImage(e.Sprite, op)
}
