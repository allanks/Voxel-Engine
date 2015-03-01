package Player

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	m "math"
)

var (
	user player
)

type player struct {
	xPos, yPos, zPos, pitch, turn, turnSpeed float32
}

func GenPlayer() {
	user = player{0.0, 0.0, 0.0, 0.0, 0.0, 1.0}
}

func GetCameraMatrix() mgl32.Mat4 {
	zLook := float32(m.Sin(float64(user.pitch)*m.Pi/180) * m.Sin(float64(user.turn)*m.Pi/180))
	xLook := float32(m.Sin(float64(user.pitch)*m.Pi/180) * m.Cos(float64(user.turn)*m.Pi/180))
	yLook := -1 * float32(m.Cos(float64(-1*user.pitch)*m.Pi/180))
	return mgl32.LookAtV(
		mgl32.Vec3{user.xPos, user.yPos, user.zPos},
		mgl32.Vec3{user.xPos + xLook, user.yPos + yLook, user.zPos + zLook},
		mgl32.Vec3{0, 1, 0})
}

func OnCursor(window *glfw.Window, xPos, yPos float64) {
	if yPos > -1 {
		window.SetCursorPos(xPos, 0)
		yPos = -1
	}
	if yPos < -180 {
		window.SetCursorPos(xPos, -180)
		yPos = -180
	}
	user.turn = float32(int32(float32(xPos)*user.turnSpeed) % 360)
	user.pitch = float32(int32(float32(yPos)*user.turnSpeed) % 360)
}

func OnKey(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyW:
		user.yPos += 1
	case glfw.KeyA:
		user.xPos -= 1
	case glfw.KeyS:
		user.yPos -= 1
	case glfw.KeyD:
		user.xPos += 1
	}
}
