package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const MAX_SWORDS int = 10
const MAX_ENEMIES int = 10

type Game struct {
	input            *InputManager
	player           *Player
	enemies          [MAX_ENEMIES]*Enemy
	mirrorGame       MirrorGame
	SingleAttack     SingleAttack
	multiSwordAttack [MAX_SWORDS]*Sword
}

func (g *Game) Update() error {
	inputState := g.input.Poll(g.player.Center())
	g.player.Update(inputState)

	g.SingleAttack.Update(g.player.SpellState, inputState, g.player.Center(), g.player.Rotation)

	if inputState.StartArrayAttack && g.player.SpellState == MirrorState {
		g.mirrorGame.Cast(g.player.Center(), g.player.Rotation)
	}
	g.mirrorGame.Update()

	if inputState.StartAttack && g.player.SpellState == SwordState {
		for _, sword := range g.multiSwordAttack {
			if sword.Shoot(g.player.Center(), g.player.Rotation) {
				break
			}
		}
	}

	if inputState.StartArrayAttack && g.player.SpellState == SwordState {
		count := 0
		for _, sword := range g.multiSwordAttack {
			if sword.ShootRandom(g.player.Center(), g.player.Rotation) {
				count++
			}
			if count > 4 {
				break
			}
		}
	}
	for _, enemy := range g.enemies {
		enemy.Update(g.player)
	}
	for _, sword := range g.multiSwordAttack {
		sword.Update()
		// Sword Collision with Enemy
		for _, enemy := range g.enemies {
			if enemy.Alive && sword.Active && enemy.CollidesWithPoint(sword.Center()) {
				sword.Active = false
				// enemy.Alive = false
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{30, 30, 46, 0xff})
	g.player.Draw(screen)
	g.mirrorGame.Draw(screen)
	g.SingleAttack.Draw(screen)
	for _, enemy := range g.enemies {
		enemy.Draw(screen)
	}
	for _, sword := range g.multiSwordAttack {
		sword.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenWidthValue, screenHeightValue
}

const (
	screenWidthValue  = 1000
	screenHeightValue = 750
	speed             = 4.0
	deadzone          = 0.2 // Ignores minor stick drift
)

func main() {
	ebiten.SetWindowSize(1000, 750)
	ebiten.SetWindowTitle("Turtlezard")

	mirrorSprite, _, err := ebitenutil.NewImageFromFile("assets/mirror-wizard.png")
	if err != nil {
		log.Fatal(err)
	}

	swordSprite, _, err := ebitenutil.NewImageFromFile("assets/sword-wizard.png")
	if err != nil {
		log.Fatal(err)
	}

	swordAttackSprite, _, err := ebitenutil.NewImageFromFile("assets/sword.png")
	if err != nil {
		log.Fatal(err)
	}

	beamAttackSprite, _, err := ebitenutil.NewImageFromFile("assets/beam.png")
	if err != nil {
		log.Fatal(err)
	}

	enemySprite, _, err := ebitenutil.NewImageFromFile("assets/enemyRed.png")
	if err != nil {
		log.Fatal(err)
	}

	multiSwordAttackSprite, _, err := ebitenutil.NewImageFromFile("assets/sword.png")
	if err != nil {
		log.Fatal(err)
	}

	player := NewPlayer(
		screenWidthValue,
		screenHeightValue,
		speed,
		mirrorSprite,
		swordSprite,
	)

	beamSprite, _, err := ebitenutil.NewImageFromFile("assets/beam.png")
	if err != nil {
		log.Fatal(err)
	}

	mirrorImage := ebiten.NewImage(50, 2)
	mirrorImage.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	multiSwordAttack := [MAX_SWORDS]*Sword{}
	for i := range multiSwordAttack {
		multiSwordAttack[i] = NewSword(5, multiSwordAttackSprite, player, false, 0)
	}

	enemies := [MAX_ENEMIES]*Enemy{}
	enemyHealth := 5
	enemySpeed := 1.0
	for i := range enemies {
		enemies[i] = NewEnemy(screenWidthValue, screenHeightValue, enemySpeed, enemySprite, player, true, enemyHealth)
	}

	game := &Game{
		input:      NewInputManager(deadzone),
		player:     player,
		enemies:    enemies,
		mirrorGame: NewMirrorGame(beamSprite),
		SingleAttack: SingleAttack{
			attackState:   NotAttacking,
			SwordSprite:   swordAttackSprite,
			BeamSprite:    beamAttackSprite,
			pos:           player.Pos,
			startPos:      player.Pos,
			rotation:      player.Rotation,
			startRotation: player.Rotation,
		},
		multiSwordAttack: multiSwordAttack,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
