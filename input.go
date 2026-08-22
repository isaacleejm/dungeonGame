package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputState struct {
	MoveX float64
	MoveY float64
}

type InputManager struct {
	gamepadIDs []ebiten.GamepadID
	deadzone   float64
}

func NewInputManager(deadzone float64) *InputManager {
	return &InputManager{deadzone: deadzone}
}

func (im *InputManager) Poll() InputState {
	var state InputState
	im.gamepadIDs = ebiten.AppendGamepadIDs(im.gamepadIDs[:0])
	if len(im.gamepadIDs) == 0 {
		return state
	}

	id := im.gamepadIDs[0]
	var rawX, rawY float64

	if ebiten.IsStandardGamepadLayoutAvailable(id) {
		rawX = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		rawY = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
	} else if ebiten.GamepadAxisCount(id) >= 2 {
		rawX = ebiten.GamepadAxisValue(id, 0)
		rawY = ebiten.GamepadAxisValue(id, 1)
	}

	// Apply deadzone filtering
	if math.Abs(rawX) > im.deadzone {
		state.MoveX = rawX
	}
	if math.Abs(rawY) > im.deadzone {
		state.MoveY = rawY
	}

	return state
}
