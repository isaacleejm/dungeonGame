package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	Pos    Vector2
	Speed  float64
	Sprite *ebiten.Image
	Health float64
}

type EnemyManager struct {
	Enemies []*Enemy

	MaxEnemies     int
	EnemyQuota     int
	EnemiesSpawned int

	SpawnDelay int
	SpawnTimer int

	EnemySpeed  float64
	EnemyHealth float64
	EnemySprite *ebiten.Image

	ScreenWidth  int
	ScreenHeight int
	Player       *Player
}

func NewEnemyManager(maxEnemies int, enemyQuota int, initialEnemies int, spawnDelay int, screenWidth int, screenHeight int, speed float64, health float64, sprite *ebiten.Image, player *Player, blocks []*Block) *EnemyManager {
	manager := &EnemyManager{
		Enemies:        make([]*Enemy, 0),
		MaxEnemies:     maxEnemies,
		EnemyQuota:     enemyQuota,
		EnemiesSpawned: 0,
		SpawnDelay:     spawnDelay,
		SpawnTimer:     0,
		EnemySpeed:     speed,
		EnemyHealth:    health,
		EnemySprite:    sprite,
		ScreenWidth:    screenWidth,
		ScreenHeight:   screenHeight,
		Player:         player,
	}

	for i := 0; i < initialEnemies; i++ {
		manager.Spawn(blocks)
	}

	return manager
}

func (m *EnemyManager) Spawn(blocks []*Block) {
	if len(m.Enemies) >= m.MaxEnemies {
		return
	}

	if m.EnemiesSpawned >= m.EnemyQuota {
		return
	}

	enemy := NewEnemy(
		m.ScreenWidth,
		m.ScreenHeight,
		m.EnemySpeed,
		m.EnemySprite,
		m.Player,
		m.EnemyHealth,
		blocks,
	)

	m.Enemies = append(m.Enemies, enemy)
	m.EnemiesSpawned++
}

// Update all enemies.
func (m *EnemyManager) Update(blocks []*Block) {
	for _, enemy := range m.Enemies {
		enemy.Update(m.Player, blocks)
	}

	removed := m.RemoveDead()

	if removed {
		m.SpawnTimer = m.SpawnDelay
	}

	if m.EnemiesSpawned >= m.EnemyQuota {
		return
	}

	if len(m.Enemies) >= m.MaxEnemies {
		return
	}

	if m.SpawnTimer > 0 {
		m.SpawnTimer--
		return
	}

	m.Spawn(blocks)
}

func (m *EnemyManager) RemoveDead() bool {
	aliveEnemies := m.Enemies[:0]
	removed := false

	for _, enemy := range m.Enemies {
		if enemy.Health > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		} else {
			removed = true
		}
	}

	m.Enemies = aliveEnemies

	return removed
}

func (m *EnemyManager) Draw(screen *ebiten.Image) {
	for _, enemy := range m.Enemies {
		enemy.Draw(screen)
	}
}

func NewEnemy(screenWidth, screenHeight int, speed float64, sprite *ebiten.Image, player *Player, health float64, blocks []*Block) *Enemy {
	enemy := &Enemy{
		Pos:    Vector2{},
		Speed:  speed,
		Sprite: sprite,
		Health: health,
	}
	enemy.FindSpawnPosition(screenWidth, screenHeight, player, blocks)

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

func (e *Enemy) FindSpawnPosition(screenWidth int, screenHeight int, player *Player, blocks []*Block) {
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
}

func (m *EnemyManager) ResetPositions(blocks []*Block) {
	for i, enemy := range m.Enemies {
		for {
			enemy.FindSpawnPosition(
				m.ScreenWidth,
				m.ScreenHeight,
				m.Player,
				blocks,
			)

			colliding := false

			for j := 0; j < i; j++ {
				if Collides(
					enemy.CollisionBounds(),
					m.Enemies[j].CollisionBounds(),
				) {
					colliding = true
					break
				}
			}

			if !colliding {
				break
			}
		}
	}
}

func (e *Enemy) Update(p *Player, blocks []*Block) {
	if e.Health <= 0 {
		return
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
	if e.Health <= 0 {
		return false
	}

	e.Health -= damage

	if e.Health <= 0 {
		e.Health = 0
		return true
	}
	// returns false if not dead
	return false
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	if e.Health <= 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.Pos.X, e.Pos.Y)

	screen.DrawImage(e.Sprite, op)
}
