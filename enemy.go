package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	Pos       Vector2
	Rotation  float64
	Speed     float64
	Sprite    *ebiten.Image
	Alive     bool
	Collision Vector4
}

func NewEnemy(screenWidth, screenHeight int, speed float64, sprite *ebiten.Image, player *Player, alive bool) *Enemy {
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
		Alive:  alive,
		Collision: Vector4{
			X1: 0,
			X2: enemyWidth,
			Y1: 0,
			Y2: enemyHeight,
		},
	}
}

func (e *Enemy) Center() Vector2 {
	bounds := e.Sprite.Bounds()
	return Vector2{
		X: e.Pos.X + float64(bounds.Dx())/2,
		Y: e.Pos.Y + float64(bounds.Dy())/2,
	}
}

func (e *Enemy) CollisionBounds() Vector4 {
	return Vector4{
		X1: e.Pos.X + e.Collision.X1,
		X2: e.Pos.X + e.Collision.X2,
		Y1: e.Pos.Y + e.Collision.Y1,
		Y2: e.Pos.Y + e.Collision.Y2,
	}
}

func (e *Enemy) CollidesWithPoint(pos Vector2) bool {
	c := e.CollisionBounds()

	return pos.X >= c.X1 &&
		pos.X <= c.X2 &&
		pos.Y >= c.Y1 &&
		pos.Y <= c.Y2
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
	if !e.Alive {
		return
	}
	bounds := e.Sprite.Bounds()
	width, height := float64(bounds.Dx()), float64(bounds.Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.Pos.X+width/2, e.Pos.Y+height/2)

	screen.DrawImage(e.Sprite, op)
}
