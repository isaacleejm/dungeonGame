package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidthValue = 320
	screenHeightValue = 240
	speed = 4.0
)

type Game struct{
	playerX, playerY float64
	rotation float64 // Stores current rotation angle in radians
	mirrorWizard *ebiten.Image
}

func (g *Game) Update() error {
	var dx, dy float64

	// WASD Movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx += speed
	}

	if dx != 0 || dy != 0 {
		// Calculate target angle based on direction vector
		g.rotation = math.Atan2(dy, dx)

		// Normalize speed so diagonal movement isn't faster
		length := math.Hypot(dx, dy)
		g.playerX += (dx / length) * speed
		g.playerY += (dy / length) * speed
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 46, 0xff})

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Rotate(g.rotation)
	options.GeoM.Translate(g.playerX, g.playerY)

	screen.DrawImage(g.mirrorWizard, options)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenWidthValue, screenHeightValue
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Turtlezard")

	mirrorWizard, _, err := ebitenutil.NewImageFromFile("assets/mirror-wizard.png")
	if err != nil {
		log.Fatal(err)
	}
	width := mirrorWizard.Bounds().Dx()
	height := mirrorWizard.Bounds().Dy()

	game := &Game{
		playerX: float64(screenWidthValue/2 - width/2),
		playerY: float64(screenHeightValue/2 - height/2),
		rotation: 0,
		mirrorWizard: mirrorWizard,
	}
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
