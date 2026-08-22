package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Health struct {
	Value      int
	Pos        Vector2
	FontSource *text.GoTextFaceSource
	GameFace   *text.GoTextFace
}

func NewHealth(value int, position Vector2, fontSource *text.GoTextFaceSource, gameFace *text.GoTextFace) *Health {
	return &Health{
		Value:      value,
		Pos:        position,
		FontSource: fontSource,
		GameFace:   gameFace,
	}
}

func (s *Health) Update(value int) {
	s.Value = value
}

func (s *Health) Hit(damage int) {
	s.Value -= damage
}

func (s *Health) Reset() {
	s.Value = 0
}

func (s *Health) Draw(screen *ebiten.Image) {
	healthText := fmt.Sprintf("Health: %d", s.Value)

	op := &text.DrawOptions{}
	op.GeoM.Translate(
		s.Pos.X,
		s.Pos.Y,
	)

	text.Draw(screen, healthText, s.GameFace, op)
}
