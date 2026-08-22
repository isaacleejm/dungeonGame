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
}

func (g *Game) Update() error {
	inputState := g.input.Poll(g.player.Center())
	g.player.Update(inputState)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 46, 0xff})
	g.player.Draw(screen)
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

	mirrorWizard, _, err := ebitenutil.NewImageFromFile("assets/mirror-wizard.png")
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		input: NewInputManager(deadzone),
		player: NewPlayer(
			screenWidthValue,
			screenHeightValue,
			speed,
			mirrorWizard,
		),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
