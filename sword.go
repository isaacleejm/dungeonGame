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
	if s.Active || s.Lifetime > 0 {
		return false
	}

	s.Pos = position
	s.Rotation = direction
	s.Lifetime = 150
	s.Active = true
	return true
}

func (s *Sword) ShootRandom(position Vector2, direction float64) bool {
	if s.Active || s.Lifetime > 0 {
		return false
	}

	const spread = math.Pi / 6 // 30 degrees

	s.Pos = position
	s.Rotation = direction + (rand.Float64()*2-1)*spread
	s.Lifetime = 120
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

func (s *Sword) CollisionBounds() Vector4 {
	bounds := s.Sprite.Bounds()

	return Vector4{
		X1: s.Pos.X,
		Y1: s.Pos.Y,
		X2: s.Pos.X + float64(bounds.Dx()),
		Y2: s.Pos.Y + float64(bounds.Dy()),
	}
}

func (s *Sword) Update(blocks []*Block) {
	if !s.Active {
		if s.Lifetime > 0 {
			s.Lifetime--
		}
		return
	}
	moveX := math.Cos(s.Rotation) * s.Speed
	moveY := math.Sin(s.Rotation) * s.Speed

	s.Pos.X += moveX
	s.Pos.Y += moveY

	for _, block := range blocks {
		if Collides(s.CollisionBounds(), block.CollisionBounds()) {
			s.Active = false
		}
	}

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
