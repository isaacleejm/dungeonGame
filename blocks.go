package main

import "github.com/hajimehoshi/ebiten/v2"

type Block struct {
	Pos Vector2
	Sprite *ebiten.Image
}

func (b *Block) Update(posX float64, posY float64){
	b.Pos.X = posX
	b.Pos.Y = posY
}

func (b *Block) CollisionBounds() Vector4 {
	bounds := b.Sprite.Bounds()

	return Vector4{
		X1: b.Pos.X,
		X2: b.Pos.X + float64(bounds.Dx()),
		Y1: b.Pos.Y,
		Y2: b.Pos.Y + float64(bounds.Dy()),
	}
}

func (b *Block) Draw(screen *ebiten.Image){
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.Pos.X, b.Pos.Y)

	screen.DrawImage(b.Sprite, op)
}