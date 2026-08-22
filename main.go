package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	input                  *InputManager
	player                 *Player
	enemy                  *Enemy
	mirrorGame             MirrorGame
	SingleAttack           SingleAttack
	multiSwordAttack       [10]*Sword
	multiSwordAttackSprite *ebiten.Image
}

func (g *Game) Update() error {
	inputState := g.input.Poll(g.player.Center())
	g.player.Update(inputState)
	g.enemy.Update(g.player)

	g.SingleAttack.Update(g.player.SpellState, inputState, g.player.Center(), g.player.Rotation)

	if inputState.StartArrayAttack && g.player.SpellState == MirrorState {
		g.mirrorGame.Cast(g.player.Center(), g.player.Rotation)
	}
	g.mirrorGame.Update()
	if inputState.StartArrayAttack && g.player.SpellState == SwordState {
		count := 0
		for _, sword := range g.multiSwordAttack {
			if sword.Shoot(g.player.Center(), g.player.Rotation) {
				count++
			}
			if count > 4 {
				break
			}
		}
	}
	for _, sword := range g.multiSwordAttack {
		sword.Update()
		if g.enemy.Alive && sword.Active && g.enemy.CollidesWithPoint(sword.Center()) {
			sword.Active = false
			g.enemy.Alive = false
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{30, 30, 46, 0xff})
	g.player.Draw(screen)
	g.enemy.Draw(screen)
	g.mirrorGame.Draw(screen)
	g.SingleAttack.Draw(screen)
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

	enemy := NewEnemy(
		screenWidthValue,
		screenHeightValue,
		1.0,
		enemySprite,
		player,
		true,
	)

	beamSprite, _, err := ebitenutil.NewImageFromFile("assets/beam.png")
	if err != nil {
		log.Fatal(err)
	}

	mirrorImage := ebiten.NewImage(50, 2)
	mirrorImage.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	multiSwordAttack := [10]*Sword{
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
		NewSword(5, multiSwordAttackSprite, player, false, 0),
	}

	game := &Game{
		input:      NewInputManager(deadzone),
		player:     player,
		enemy:      enemy,
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
		multiSwordAttackSprite: multiSwordAttackSprite,
		multiSwordAttack:       multiSwordAttack,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
