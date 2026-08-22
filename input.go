package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputState struct {
	MoveX        float64
	MoveY        float64
	TargetAngle  float64
	HasAngleLock bool // True if aiming with mouse or moving stick
	ToggleSpell  bool // True for a single frame on button press
}

type InputManager struct {
	gamepadIDs   []ebiten.GamepadID
	deadzone     float64
	lastMousePos Vector2
	toggleBinding ActionBinding
}

func NewInputManager(deadzone float64) *InputManager {
	return &InputManager{
		deadzone: deadzone,
		toggleBinding: ActionBinding{
			Key:            ebiten.KeySpace,
			StandardButton: ebiten.StandardGamepadButtonRightBottom, // A
			RawButtonIndex: 0,
		},
	}
}

func (im *InputManager) Poll(playerCenter Vector2) InputState {
	im.gamepadIDs = ebiten.AppendGamepadIDs(im.gamepadIDs[:0])
	hasGamepad := len(im.gamepadIDs) > 0

	if !hasGamepad {
		return im.pollKBM(playerCenter)
	}

	var state InputState

	id := im.gamepadIDs[0]
	state.ToggleSpell = im.toggleBinding.JustPressedGamepad(id)

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

	// Gamepad rotation faces movement direction
	if math.Hypot(state.MoveX, state.MoveY) > 0 {
		state.TargetAngle = math.Atan2(state.MoveY, state.MoveX)
		state.HasAngleLock = true
	}

	return state
}

// / Keyboard and mouse controls. The player aims towards the cursor,
// independent of WASD movement
func (im *InputManager) pollKBM(playerCenter Vector2) InputState {
	// WASD movement
	var state InputState

	state.ToggleSpell = im.toggleBinding.JustPressedKey()

	var kx, ky float64
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ky -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ky += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		kx -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		kx += 1.0
	}

	if kx != 0 || ky != 0 {
		state.MoveX = kx
		state.MoveY = ky
	}

	// Cursor aiming
	mx, my := ebiten.CursorPosition()
	currentMouse := Vector2{X: float64(mx), Y: float64(my)}

	// Aim with mouse if the mouse has moved, buttons are pressed, or no gamepad input is present
	mouseMoved := currentMouse != im.lastMousePos
	im.lastMousePos = currentMouse

	if mouseMoved || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || (kx != 0 || ky != 0) {
		dx := currentMouse.X - playerCenter.X
		dy := currentMouse.Y - playerCenter.Y
		state.TargetAngle = math.Atan2(dy, dx)
		state.HasAngleLock = true
	}

	return state
}
