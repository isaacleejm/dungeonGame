package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Level struct {
	Value      int
	Pos        Vector2
	FontSource *text.GoTextFaceSource
	GameFace   *text.GoTextFace
}

func NewLevel(position Vector2, fontSource *text.GoTextFaceSource, gameFace *text.GoTextFace) *Level {
	return &Level{
		Value:      1,
		Pos:        position,
		FontSource: fontSource,
		GameFace:   gameFace,
	}
}

func (s *Level) Add(points int) {
	s.Value += points
}

func (s *Level) Reset() {
	s.Value = 0
}

func (s *Level) Draw(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("Level: %d", s.Value)

	background := ebiten.NewImage(120, 40)
	background.Fill(color.RGBA{0, 0, 0, 180})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		s.Pos.X-10,
		s.Pos.Y-5,
	)

	screen.DrawImage(background, op)

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(
		s.Pos.X,
		s.Pos.Y,
	)

	text.Draw(screen, scoreText, s.GameFace, textOp)
}
