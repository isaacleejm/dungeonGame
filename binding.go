package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ActionBinding struct {
	Key             ebiten.Key
	StandardButton  ebiten.StandardGamepadButton
	RawButtonIndex  ebiten.GamepadButton
	isPressedBefore bool
}

func (b *ActionBinding) JustPressedKey() bool {
	pressed := ebiten.IsKeyPressed(b.Key)
	justPressed := pressed && !b.isPressedBefore
	b.isPressedBefore = pressed
	return justPressed
}

func (b *ActionBinding) JustPressedGamepad(id ebiten.GamepadID) bool {
	var pressed bool
	if ebiten.IsStandardGamepadLayoutAvailable(id) {
		pressed = ebiten.IsStandardGamepadButtonPressed(id, b.StandardButton)
	} else {
		pressed = ebiten.IsGamepadButtonPressed(id, b.RawButtonIndex)
	}

	justPressed := pressed && !b.isPressedBefore
	b.isPressedBefore = pressed
	return justPressed
}
