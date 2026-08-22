package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sword struct {
	Pos      Vector2
	Rotation float64
	Speed    float64
	Sprite   *ebiten.Image
	Active   bool
	Lifetime int
}

func NewSword(speed float64, sprite *ebiten.Image, player *Player, active bool, lifetime int) *Sword {

	return &Sword{
		Pos:      Vector2{X: player.Center().X, Y: player.Center().Y},
		Rotation: player.Rotation,
		Speed:    speed,
		Sprite:   sprite,
		Active:   active,
		Lifetime: lifetime,
	}
}

func (s *Sword) Shoot(position Vector2, direction float64) bool {
	if s.Active {
		return false
	}

	const spread = math.Pi / 6 // 30 degrees

	s.Pos = position
	s.Rotation = direction + (rand.Float64()*2-1)*spread
	s.Lifetime = 60
	s.Active = true
	return true
}

func (s *Sword) Center() Vector2 {
	bounds := s.Sprite.Bounds()
	return Vector2{
		X: s.Pos.X + float64(bounds.Dx())/2,
		Y: s.Pos.Y + float64(bounds.Dy())/2,
	}
}

func (s *Sword) Update() {
	if !s.Active {
		return
	}
	s.Pos.X += math.Cos(s.Rotation) * s.Speed
	s.Pos.Y += math.Sin(s.Rotation) * s.Speed

	s.Lifetime--

	if s.Lifetime <= 0 {
		s.Active = false
	}
}

func (s *Sword) Draw(screen *ebiten.Image) {
	if !s.Active {
		return
	}

	bounds := s.Sprite.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-width/2, -height/2)
	op.GeoM.Rotate(s.Rotation + math.Pi/4)
	op.GeoM.Translate(s.Pos.X, s.Pos.Y)

	screen.DrawImage(s.Sprite, op)
}
