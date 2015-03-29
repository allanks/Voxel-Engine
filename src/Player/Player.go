package Player

import (
	m "math"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	user                                                       player
	movingForward, movingBackward, strafingLeft, strafingRight chan float64
)

const (
	forward, backwards, turnSpeed, moveSpeed, stop, Height float64 = 1, -1, 0.5, 0.1, 0, 1
)

type player struct {
	xPos, yPos, zPos, pitch, turn float64
}

type moveFunc func(float64)

func GenPlayer() {
	user = player{0.0, 1.0 + Height, 0.0, -180.0, 0.0}
	movingForward = make(chan float64)
	movingBackward = make(chan float64)
	strafingLeft = make(chan float64)
	strafingRight = make(chan float64)
	go loopMove(move, movingForward)
	go loopMove(move, movingBackward)
	go loopMove(strafe, strafingLeft)
	go loopMove(strafe, strafingRight)
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

func loopMove(fn moveFunc, input chan float64) {
	xMove := float64(0)

	// Spawn listener for movement
	go func() {
		for {
			xMove = <-input
		}
	}()

	for {
		//fmt.Println("Moving User")
		fn(xMove)
		time.Sleep(16 * time.Millisecond)
	}
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
	user.yPos = user.yPos + (direction * yLook * moveSpeed)
	user.zPos = user.zPos + (direction * zLook * moveSpeed)
}

func strafe(direction float64) {
	xLook := float64(m.Cos(float64(user.turn) * m.Pi / 180))
	zLook := float64(m.Sin(float64(user.turn) * m.Pi / 180))

	user.xPos = user.xPos + (-1 * direction * zLook * moveSpeed)
	user.zPos = user.zPos + (direction * xLook * moveSpeed)
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
	case glfw.KeyW:
		if action == glfw.Press {
			movingForward <- forward
		} else if action == glfw.Release {
			movingForward <- stop
		}
	case glfw.KeyA:
		if action == glfw.Press {
			strafingLeft <- forward
		} else if action == glfw.Release {
			strafingLeft <- stop
		}
	case glfw.KeyS:
		if action == glfw.Press {
			movingBackward <- backwards
		} else if action == glfw.Release {
			movingBackward <- stop
		}
	case glfw.KeyD:
		if action == glfw.Press {
			strafingRight <- backwards
		} else if action == glfw.Release {
			strafingRight <- stop
		}
	}
}
