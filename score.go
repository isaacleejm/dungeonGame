package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Score struct {
	Value      int
	Pos        Vector2
	fontSource *text.GoTextFaceSource
	gameFace   *text.GoTextFace
}

func NewScore(position Vector2, fontSource *text.GoTextFaceSource, gameFace *text.GoTextFace) *Score {
	return &Score{
		Value:      0,
		Pos:        position,
		fontSource: fontSource,
		gameFace:   gameFace,
	}
}

func (s *Score) Add(points int) {
	s.Value += points
}

func (s *Score) Reset() {
	s.Value = 0
}

func (s *Score) Draw(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("Score: %d", s.Value)

	op := &text.DrawOptions{}
	op.GeoM.Translate(
		s.Pos.X,
		s.Pos.Y,
	)

	text.Draw(screen, scoreText, s.gameFace, op)
}
