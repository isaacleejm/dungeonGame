package main

func Collides(a, b Vector4) bool {
	return a.X1 < b.X2 &&
		a.X2 > b.X1 &&
		a.Y1 < b.Y2 &&
		a.Y2 > b.Y1
}
