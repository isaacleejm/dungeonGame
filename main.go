package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	input  *InputManager
	player *Player
	enemy  *Enemy
	mirrorGame MirrorGame
	SingleAttack SingleAttack
}

func (g *Game) Update() error {
	inputState := g.input.Poll(g.player.Center())
	g.player.Update(inputState)

	if g.enemy == nil {
		enemySprite, _, err := ebitenutil.NewImageFromFile("assets/enemyRed.png")
		if err != nil {
			log.Fatal(err)
		}
		g.enemy = NewEnemy(
			screenWidthValue,
			screenHeightValue,
			1.0,
			enemySprite,
			g.player,
		)
	}

	g.enemy.Update(g.player)

	g.SingleAttack.Update(g.player.SpellState, inputState, g.player.Center(), g.player.Rotation)

	if inputState.StartArrayAttack && g.player.SpellState == MirrorState {
		g.mirrorGame.Cast(g.player.Center(), g.player.Rotation)
	}
	g.mirrorGame.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{30, 30, 46, 0xff})
	g.player.Draw(screen)
	g.enemy.Draw(screen)
	g.mirrorGame.Draw(screen)
	g.SingleAttack.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenWidthValue, screenHeightValue
}

const (
	screenWidthValue  = 320
	screenHeightValue = 240
	speed             = 4.0
	deadzone          = 0.2 // Ignores minor stick drift
)

func main() {
	ebiten.SetWindowSize(640, 480)
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
	)

	beamSprite, _, err := ebitenutil.NewImageFromFile("assets/beam.png")
	if err != nil {
		log.Fatal(err)
	}

	mirrorImage := ebiten.NewImage(50, 2)
	mirrorImage.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	game := &Game{
		input:  NewInputManager(deadzone),
		player: player,
		enemy:  enemy,
		mirrorGame: NewMirrorGame(beamSprite),
		SingleAttack: SingleAttack{
			attackState: NotAttacking,
			SwordSprite: swordAttackSprite,
			BeamSprite: beamAttackSprite,
			pos: player.Pos,
			startPos: player.Pos,
			rotation: player.Rotation,
			startRotation: player.Rotation,
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
