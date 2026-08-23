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
	Health    float64
}

func NewEnemy(
	screenWidth, screenHeight int,
	speed float64,
	sprite *ebiten.Image,
	player *Player,
	alive bool,
	health float64,
	blocks []*Block,
) *Enemy {
	enemyBounds := sprite.Bounds()

	enemyWidth := float64(enemyBounds.Dx())
	enemyHeight := float64(enemyBounds.Dy())

	enemy := &Enemy{
		Pos:    Vector2{},
		Speed:  speed,
		Sprite: sprite,
		Alive:  alive,
		Collision: Vector4{
			X1: 0,
			X2: enemyWidth,
			Y1: 0,
			Y2: enemyHeight,
		},
		Health: health,
	}
	enemy.Respawn(screenWidth, screenHeight, player, blocks)

	return enemy
}

func (e *Enemy) Center() Vector2 {
	bounds := e.Sprite.Bounds()
	return Vector2{
		X: e.Pos.X + float64(bounds.Dx())/2,
		Y: e.Pos.Y + float64(bounds.Dy())/2,
	}
}

func (e *Enemy) CollisionBounds() Vector4 {
	bounds := e.Sprite.Bounds()

	return Vector4{
		X1: e.Pos.X,
		Y1: e.Pos.Y,
		X2: e.Pos.X + float64(bounds.Dx()),
		Y2: e.Pos.Y + float64(bounds.Dy()),
	}
}

func (e *Enemy) Respawn(screenWidth int, screenHeight int, player *Player, blocks []*Block) {
	bounds := e.Sprite.Bounds()

	enemyWidth := float64(bounds.Dx())
	enemyHeight := float64(bounds.Dy())

	playerCenter := player.Center()

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

		// Enemy is too close to player.
		if distance < minDistance {
			continue
		}

		// Collision bounds for the potential spawn position.
		enemyBounds := Vector4{
			X1: startX,
			Y1: startY,
			X2: startX + enemyWidth,
			Y2: startY + enemyHeight,
		}

		// Check if the potential position overlaps a block.
		colliding := false

		for _, block := range blocks {
			if Collides(enemyBounds, block.CollisionBounds()) {
				colliding = true
				break
			}
		}

		// Position is inside a block.
		if colliding {
			continue
		}

		// Valid position found.
		break
	}

	// Set the new position.
	e.Pos.X = startX
	e.Pos.Y = startY

	// Reset enemy state.
	e.Alive = true
	e.Health = 5
}

func (e *Enemy) Update(p *Player, blocks []*Block) {
	if e.Health <= 0 {
		e.Alive = false
	}

	playerPos := p.Pos
	enemyPos := e.Pos

	dx := playerPos.X - enemyPos.X
	dy := playerPos.Y - enemyPos.Y

	distance := math.Hypot(dx, dy)

	if distance == 0 {
		return
	}

	moveX := (dx / distance) * e.Speed
	moveY := (dy / distance) * e.Speed

	e.Pos.X += moveX

	for _, block := range blocks {
		if Collides(e.CollisionBounds(), block.CollisionBounds()) {
			e.Pos.X -= moveX
		}
	}

	e.Pos.Y += moveY

	for _, block := range blocks {
		if Collides(e.CollisionBounds(), block.CollisionBounds()) {
			e.Pos.Y -= moveY
		}
	}
}

func (e *Enemy) TakeDamage(damage float64) bool {
	if !e.Alive {
		return false
	}

	e.Health -= damage

	if e.Health <= 0 {
		e.Health = 0
		e.Alive = false
		return true
	}
	// returns false if not dead
	return false
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	if !e.Alive {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.Pos.X, e.Pos.Y)

	screen.DrawImage(e.Sprite, op)
}
