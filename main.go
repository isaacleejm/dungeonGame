package main

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const MAX_SWORDS int = 10
const MAX_ENEMIES int = 50

type Game struct {
	input                  *InputManager
	player                 *Player
	enemies                [MAX_ENEMIES]*Enemy
	mirrorGame             MirrorGame
	beamAttack             BeamAttack
	multiSwordAttackSprite *ebiten.Image
	blocks                 []*Block
	background             color.RGBA

	blockIndex  int
	layoutIndex int
	themeIndex  int

	level int

	blockSpriteList  []*ebiten.Image
	multiSwordAttack [MAX_SWORDS]*Sword
	score            *Score
	health           *Health
}

var backgroundList = []color.RGBA{
	Backgrounds.Gray,
	Backgrounds.Blue,
	Backgrounds.Red,
	Backgrounds.Green,
	Backgrounds.Purple,
}

var layoutList = [][]string{
	Layouts.Layout1,
	Layouts.Layout2,
	Layouts.Layout3,
	Layouts.Layout4,
	Layouts.Layout5,
	Layouts.Layout6,
	Layouts.Layout7,
	Layouts.Layout8,
	Layouts.Layout9,
	Layouts.Layout10,
}

var levelScoreThresholds = []int{
	10,
	25,
	45,
	70,
	100,
}

var customisablePermissions = false

func (g *Game) UpdateLevel() {
	currentLevel := BlocksFromLayout(
		layoutList[g.layoutIndex],
		g.blockSpriteList[g.blockIndex],
		backgroundList[g.themeIndex],
	)

	g.blocks = currentLevel.Blocks
	g.background = currentLevel.Background

	g.player.Pos = Vector2{80, 80}
}

func (g *Game) Update() error {
	inputState := g.input.Poll(g.player.Center())
	g.player.Update(inputState, g.blocks)

	g.beamAttack.Update(
		g.player.SpellState,
		inputState,
		g.player.Center(),
		g.player.Rotation,
		g.mirrorGame.Mirrors,
	)

	if inputState.StartArrayAttack {
		switch g.player.SpellState {
		case MirrorState:
			g.mirrorGame.Cast(
				g.player.Center(),
				g.player.Rotation,
				g.blocks,
			)

		case SwordState:
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
	}

	g.mirrorGame.Update()

	if inputState.customisabilityChange {
		customisablePermissions = !customisablePermissions
	}

	if inputState.StartAttack && g.player.SpellState == SwordState {
		for _, sword := range g.multiSwordAttack {
			if sword.Shoot(g.player.Center(), g.player.Rotation) {
				break
			}
		}
	}

	if inputState.blockChange && customisablePermissions {
		g.blockIndex++

		if g.blockIndex >= len(g.blockSpriteList) {
			g.blockIndex = 0
		}

		g.UpdateLevel()
	}

	if inputState.LayoutChange && customisablePermissions {
		g.layoutIndex++

		if g.layoutIndex >= len(layoutList) {
			g.layoutIndex = 0
		}

		g.UpdateLevel()
	}

	if inputState.ThemeChange && customisablePermissions {
		g.themeIndex++

		if g.themeIndex >= len(backgroundList) {
			g.themeIndex = 0
		}

		g.UpdateLevel()
	}

	for _, enemy := range g.enemies {
		enemy.Update(g.player, g.blocks)
	}
	for _, sword := range g.multiSwordAttack {
		sword.Update(g.blocks)
		if sword.Active {
			if g.damageEnemiesInBounds(sword.CollisionBounds(), false) {
				sword.Active = false
				g.score.Add(1)

				if g.score.Value >= levelScoreThresholds[g.layoutIndex] {
					g.layoutIndex++
					g.level++

					if g.layoutIndex >= len(layoutList) {
						g.layoutIndex = 0
					}

					g.UpdateLevel()
				}
			}
		}
	}

	beamThickness := float64(g.mirrorGame.BeamSprite.Bounds().Dx()) * 0.2

	g.checkBeamCollisions(g.beamAttack.BeamNodes, beamThickness)
	g.checkBeamCollisions(g.mirrorGame.BeamNodes, beamThickness)

	g.health.Update(g.player.Health)

	return nil
}

func (g *Game) damageEnemiesInBounds(bounds Vector4, pierce bool) (hitSomething bool) {
	for _, enemy := range g.enemies {
		if enemy.Alive && Collides(bounds, enemy.CollisionBounds()) {
			hitSomething = true

			if enemy.TakeDamage(1) {
				g.score.Add(1)

				if g.score.Value >= levelScoreThresholds[g.layoutIndex] {
					g.layoutIndex++
					g.level++

					if g.layoutIndex >= len(layoutList) {
						g.layoutIndex = 0
					}

					g.UpdateLevel()
				}
			}

			// If the attack doesn't pierce (sword), stop checking other enemies
			if !pierce {
				break
			}
		}
	}
	return hitSomething
}

func (g *Game) checkBeamCollisions(nodes []Vector2, thickness float64) {
	for i := 0; i < len(nodes)-1; i++ {
		A := nodes[i]
		B := nodes[i+1]

		bounds := BeamSegmentBounds(A, B, thickness)
		g.damageEnemiesInBounds(bounds, true)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background color
	// screen.Fill(color.RGBA{30, 30, 46, 0xff})
	screen.Fill(g.background)
	g.player.Draw(screen)
	g.mirrorGame.Draw(screen)
	g.beamAttack.Draw(screen)

	for _, block := range g.blocks {
		block.Draw(screen)
	}

	for _, enemy := range g.enemies {
		enemy.Draw(screen)
	}
	for _, sword := range g.multiSwordAttack {
		sword.Draw(screen)
	}
	g.score.Draw(screen)
	g.health.Draw(screen)
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

//go:embed assets/gamecontrollerdb.txt
var customGamepadDB string

func main() {
	ebiten.SetWindowSize(1000, 750)
	ebiten.SetWindowTitle("Turtlezard")

	mappingsApplied, err := ebiten.UpdateStandardGamepadLayoutMappings(customGamepadDB)
	if !mappingsApplied || err != nil {
		log.Printf("Warning: failed to load embedded gamepad mappings: %v", err)
	}

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

	dungeonBlockSprite, _, err := ebitenutil.NewImageFromFile("assets/dungeon-block.png")
	if err != nil {
		log.Fatal(err)
	}

	logBlockSprite, _, err := ebitenutil.NewImageFromFile("assets/log-block.png")
	if err != nil {
		log.Fatal(err)
	}

	brickBlockSprite, _, err := ebitenutil.NewImageFromFile("assets/brick-block.png")
	if err != nil {
		log.Fatal(err)
	}

	blockSpriteList := []*ebiten.Image{
		dungeonBlockSprite,
		brickBlockSprite,
		logBlockSprite,
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
		multiSwordAttack[i] = NewSword(5, swordAttackSprite, player, false, 0)
	}

	currentLevel := BlocksFromLayout(
		Layouts.Layout1,
		blockSpriteList[0],
		Backgrounds.Gray,
	)

	enemies := [MAX_ENEMIES]*Enemy{}
	enemyHealth := 5
	enemySpeed := 1.0
	for i := range enemies {
		enemies[i] = NewEnemy(screenWidthValue, screenHeightValue, enemySpeed, enemySprite, player, true, enemyHealth, currentLevel.Blocks)
	}

	currentLevel := BlocksFromLayout(
		Layouts.Layout1,
		blockSpriteList[0],
		Backgrounds.Gray,
	)
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	gameFace := &text.GoTextFace{
		Source: fontSource,
		Size:   24,
	}

	game := &Game{
		input:      NewInputManager(deadzone),
		player:     player,
		enemies:    enemies,
		mirrorGame: NewMirrorGame(beamSprite),
		beamAttack: BeamAttack{
			attackState:   NotAttacking,
			BeamSprite:    beamAttackSprite,
			pos:           player.Pos,
			startPos:      player.Pos,
			rotation:      player.Rotation,
			startRotation: player.Rotation,
		},
		multiSwordAttackSprite: swordAttackSprite,
		multiSwordAttack:       multiSwordAttack,
		blocks:                 currentLevel.Blocks,
		background:             currentLevel.Background,

		blockIndex:  0,
		layoutIndex: 0,
		themeIndex:  0,

		level: 1,

		blockSpriteList: blockSpriteList,
		score: NewScore(
			Vector2{X: screenWidthValue - 125, Y: 0},
			fontSource,
			gameFace,
		),
		health: NewHealth(
			player.Health,
			Vector2{X: 0, Y: 0},
			fontSource,
			gameFace,
		),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
