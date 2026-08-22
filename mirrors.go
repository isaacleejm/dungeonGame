package main

import (
	"math"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type BeamState int

const (
	NoBeam BeamState = iota
	ActiveBeam
)

type Mirror struct {
	Pos      Vector2
	Rotation float64
	Sprite   *ebiten.Image
}

func (m *Mirror) Draw(screen *ebiten.Image) {
	bounds := m.Sprite.Bounds()
	w, h := float64(bounds.Dx()), float64(bounds.Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(m.Rotation)
	op.GeoM.Translate(m.Pos.X, m.Pos.Y)

	screen.DrawImage(m.Sprite, op)
}

const (
	mirrorWidth = 50
	mirrorThickness = 2
	radius        = 80.0  // Distance to the left/right of the player
	startDistance = 25.0  // Distance in front of the player the rows start
	spacing       = 100.0 // Distance between each mirror in the row
	mirrorsPerRow = 3
	staggerOffset = spacing / 2.0
)

type MirrorGame struct {
	Mirrors []Mirror
	BeamNodes  []Vector2
	BeamState BeamState
	TimerTicks int
	Sprite  *ebiten.Image
	BeamSprite *ebiten.Image
}

func NewMirrorGame(beamSprite *ebiten.Image) MirrorGame {
	mirrorImage := ebiten.NewImage(mirrorWidth, mirrorThickness)
	mirrorImage.Fill(color.RGBA{R: 148, G: 226, B: 213, A: 255})

	return MirrorGame{
		Sprite: mirrorImage,
		BeamSprite: beamSprite,
	}
}

func (m *Mirror) CollisionBounds() Vector4 {
	bounds := m.Sprite.Bounds()

	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	return Vector4{
		X1: m.Pos.X - w/2,
		Y1: m.Pos.Y - h/2,
		X2: m.Pos.X + w/2,
		Y2: m.Pos.Y + h/2,
	}
}

func (m *MirrorGame) Update() {
	if m.BeamState == ActiveBeam && m.TimerTicks > 0 {
		m.TimerTicks--

		if m.TimerTicks == 0 {
			m.BeamState = NoBeam
			m.BeamNodes = make([]Vector2, 0)
		}
	}
}

/// Parameters come from the Player
func (m *MirrorGame) Cast(center Vector2, theta float64, blocks []*Block) {
	// Clear previous mirrors (todo: replace old ones gradually)
	m.Mirrors = make([]Mirror, 0)

	// The beam starts at the player's center
	m.BeamNodes = []Vector2{center}

	// Calculate direction vectors
	// Forward vector (parallel to aim)
	fwdX := math.Cos(theta)
	fwdY := math.Sin(theta)

	// Perpendicular vector (90 degrees / pi/2 offset)
	// -Sin(theta), Cos(theta) gets the vector pointing to the player's left
	sideX := -math.Sin(theta)
	sideY := math.Cos(theta)

	// Generate the mirrors and beam
	for i := range mirrorsPerRow {
		leftForwardOffset := startDistance + float64(i)*spacing
		rightForwardOffset := startDistance + staggerOffset + float64(i)*spacing

		// Left Row
		leftPos := Vector2{
			X: center.X + (sideX * radius) + (fwdX * leftForwardOffset),
			Y: center.Y + (sideY * radius) + (fwdY * leftForwardOffset),
		}
		leftMirror := Mirror{
			Pos: leftPos,
			Rotation: theta,
			Sprite:   m.Sprite,
		}

		blockedLeft := false

		for _, block := range blocks {
			if Collides(leftMirror.CollisionBounds(), block.CollisionBounds()){
				blockedLeft = true
				break
			}
		}

		if !blockedLeft {
			m.Mirrors = append(m.Mirrors, leftMirror)
			m.BeamNodes = append(m.BeamNodes, leftPos)
		}

		// Right Row
		rightPos := Vector2{
			X: center.X - (sideX * radius) + (fwdX * rightForwardOffset),
			Y: center.Y - (sideY * radius) + (fwdY * rightForwardOffset),
		}
		rightMirror := Mirror{
			Pos: rightPos,
			Rotation: theta,
			Sprite:   m.Sprite,
		}

		blockedRight := false

		for _, block := range blocks {
			if Collides(rightMirror.CollisionBounds(), block.CollisionBounds()){
				blockedRight = true
				break
			}
		}

		if !blockedRight {
			m.Mirrors = append(m.Mirrors, rightMirror)
			m.BeamNodes = append(m.BeamNodes, rightPos)
		}
	}

	// Add a final node to let the beam exit the mirror array
	// By projecting a "virtual" left mirror position, we get the exact exit angle
	// naturally
	virtualForwardOffset := startDistance + float64(mirrorsPerRow)*spacing
	virtualExitPos := Vector2{
		X: center.X + (sideX * radius) + (fwdX * virtualForwardOffset),
		Y: center.Y + (sideY * radius) + (fwdY * virtualForwardOffset),
	}
	m.BeamNodes = append(m.BeamNodes, virtualExitPos)

	// Start beam timer
	// 60 TPS (ticks per second) * seconds = total ticks to wait
	m.TimerTicks = int(0.5 * 60)
	m.BeamState = ActiveBeam
}

func (m *MirrorGame) Draw(screen *ebiten.Image) {
	// Draw beams first so they appear to strike the surface of the mirrors
	if len(m.BeamNodes) >= 2 {
		beamBounds := m.BeamSprite.Bounds()
		beamWidth := float64(beamBounds.Dx())
		beamHeight := float64(beamBounds.Dy())

		for i := 0; i < len(m.BeamNodes)-1; i++ {
			A := m.BeamNodes[i]
			B := m.BeamNodes[i+1]

			// Calculate distance and angle between the two nodes
			dx := B.X - A.X
			dy := B.Y - A.Y
			dist := math.Hypot(dx, dy)
			angle := math.Atan2(dy, dx)

			// How much we need to stretch the sprite horizontally to reach point B
			scaleX := dist / beamWidth

			op := &ebiten.DrawImageOptions{}

			// Offset Y by half height so the beam rotates exactly on its center line
			op.GeoM.Translate(0, -beamHeight/2)

			// Stretch the beam horizontally to the exact distance
			op.GeoM.Scale(scaleX, 0.2)

			// Rotate to face the next mirror
			op.GeoM.Rotate(angle)

			// Move to the start point (Mirror A)
			op.GeoM.Translate(A.X, A.Y)

			// Additive blending makes energy beams glow nicely when they overlap
			op.Blend = ebiten.BlendLighter

			screen.DrawImage(m.BeamSprite, op)
		}
	}

	for _, mirror := range m.Mirrors {
		mirror.Draw(screen)
	}
}
