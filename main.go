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
}

func (g *Game) Update() error {
	inputState := g.input.Poll(g.player.Center())
	g.player.Update(inputState)
	g.enemy.Update(g.player)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{30, 30, 46, 0xff})
	g.player.Draw(screen)
	g.enemy.Draw(screen)
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
		&SwordAttack{
			attackState: NotAttacking,
			Sprite: swordAttackSprite,
		},
	)

	enemy := NewEnemy(
		screenWidthValue,
		screenHeightValue,
		1.0,
		enemySprite,
		player,
	)

	game := &Game{
		input:  NewInputManager(deadzone),
		player: player,
		enemy:  enemy,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
