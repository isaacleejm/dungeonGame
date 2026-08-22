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
	deadzone = 0.2 // Ignores minor stick drift
)

type Game struct{
	playerX, playerY float64
	rotation float64 // Stores current rotation angle in radians
	mirrorWizard *ebiten.Image
	gamepadIDs       []ebiten.GamepadID
}

func (g *Game) Update() error {
	// Detect connected controllers
	g.gamepadIDs = ebiten.AppendGamepadIDs(g.gamepadIDs[:0])
	if len(g.gamepadIDs) == 0 {
		return nil // No controller connected; skip input processing
	}
	id := g.gamepadIDs[0]

	var dx, dy float64

	// Left Analog Stick
	var axisX, axisY float64

	if ebiten.IsStandardGamepadLayoutAvailable(id) {
		axisX = ebiten.StandardGamepadAxisValue(
			id,
			ebiten.StandardGamepadAxisLeftStickHorizontal,
		)
		axisY = ebiten.StandardGamepadAxisValue(
			id,
			ebiten.StandardGamepadAxisLeftStickVertical,
		)
	} else {
		// Raw axis fallback
		if ebiten.GamepadAxisCount(id) >= 2 {
	        axisX = ebiten.GamepadAxisValue(id, 0)
	        axisY = ebiten.GamepadAxisValue(id, 1)
	    }
	}

	if math.Abs(axisX) > deadzone {
		dx += axisX
	}
	if math.Abs(axisY) > deadzone {
		dy += axisY
	}

	length := math.Hypot(dx, dy)
	if length > 0 {
		// Calculate target angle based on direction vector
		g.rotation = math.Atan2(dy, dx)

		// Clamp magnitude so analog stick controls variable walking/running speed
		magnitude := math.Min(1.0, length)

		g.playerX += (dx / length) * speed * magnitude
		g.playerY += (dy / length) * speed * magnitude
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 46, 0xff})

	options := &ebiten.DrawImageOptions{}

	width := float64(g.mirrorWizard.Bounds().Dx())
	height := float64(g.mirrorWizard.Bounds().Dy())

	options.GeoM.Translate(-width/2, -height/2)
	options.GeoM.Rotate(g.rotation)
	options.GeoM.Translate(g.playerX+width/2, g.playerY+height/2)

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
