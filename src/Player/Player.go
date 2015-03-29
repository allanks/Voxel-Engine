package Player

import (
	m "math"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	user player
)

const (
	forward, backwards, turnSpeed, moveSpeed, stop, Height float64 = 1, -1, 0.5, 0.1, 0, 1
)

type player struct {
	xPos, yPos, zPos, pitch, turn float64
	freeMovement                  bool
}

type moveFunc func(float64)

func GenPlayer() {
	user = player{0.0, 1.0 + Height, 0.0, -180.0, 0.0, false}
}

func MovePlayer(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		move(1)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		move(-1)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		strafe(1)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		strafe(-1)
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		jump(1)
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		jump(-1)
	}
}

func SetPosistion(xPos, yPos, zPos float64) {
	user.xPos = xPos
	user.yPos = yPos
	user.zPos = zPos
}

func GetPosition() (float64, float64, float64) {
	return user.xPos, user.yPos, user.zPos
}

func GetPlayerSpeed() float64 {
	return moveSpeed
}

func GetCameraMatrix() mgl32.Mat4 {
	xLook := float64(m.Sin(float64(user.pitch)*m.Pi/180) * m.Cos(float64(user.turn)*m.Pi/180))
	zLook := float64(m.Sin(float64(user.pitch)*m.Pi/180) * m.Sin(float64(user.turn)*m.Pi/180))
	yLook := -1 * float64(m.Cos(float64(-1*user.pitch)*m.Pi/180))
	return mgl32.LookAtV(
		mgl32.Vec3{float32(user.xPos), float32(user.yPos), float32(user.zPos)},
		mgl32.Vec3{float32(user.xPos + xLook), float32(user.yPos + yLook), float32(user.zPos + zLook)},
		mgl32.Vec3{0, 1, 0})
}

func move(direction float64) {
	xLook := float64(m.Sin(float64(user.pitch)*m.Pi/180) * m.Cos(float64(user.turn)*m.Pi/180))
	zLook := float64(m.Sin(float64(user.pitch)*m.Pi/180) * m.Sin(float64(user.turn)*m.Pi/180))
	yLook := -1 * float64(m.Cos(float64(-1*user.pitch)*m.Pi/180))
	user.xPos = user.xPos + (direction * xLook * moveSpeed)
	if user.freeMovement {
		user.yPos = user.yPos + (direction * yLook * moveSpeed)
	}
	user.zPos = user.zPos + (direction * zLook * moveSpeed)
}

func strafe(direction float64) {
	xLook := float64(m.Cos(float64(user.turn) * m.Pi / 180))
	zLook := float64(m.Sin(float64(user.turn) * m.Pi / 180))

	user.xPos = user.xPos + (-1 * direction * zLook * moveSpeed)
	user.zPos = user.zPos + (direction * xLook * moveSpeed)
}

func jump(direction float64) {
	user.yPos = user.yPos + (direction * moveSpeed)
}

func OnCursor(window *glfw.Window, xPos, yPos float64) {
	if yPos > -5 {
		window.SetCursorPos(xPos, -5)
		yPos = -5
	}
	if yPos < -359 {
		window.SetCursorPos(xPos, -359)
		yPos = -359
	}
	user.turn = float64(int32(float64(xPos)*turnSpeed) % 360)
	user.pitch = float64(int32(float64(yPos)*turnSpeed) % 360)
}

func OnKey(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyRightShift:
		if action == glfw.Press {
			user.freeMovement = !user.freeMovement
		}
	}

}
