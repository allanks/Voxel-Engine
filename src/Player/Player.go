package Player

import (
	"fmt"
	m "math"

	"github.com/allanks/Voxel-Engine/src/Terrain"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	user          player
	lastFrameTime float64
)

const (
	forward,
	backwards,
	turnSpeed,
	moveSpeed,
	stop,
	Height,
	terminalVelocity,
	jumpSpeed,
	gravity,
	collisionDistance float64 = 1, -1, 0.5, 0.1, 0, 1, -10, 2, 9.8, 0.15
)

type player struct {
	xPos, yPos, zPos, pitch, turn, fall float64
	freeMovement, isFalling             bool
}

type moveFunc func(float64)

func GenPlayer(xPos, yPos, zPos float64) {
	lastFrameTime = glfw.GetTime()
	user = player{xPos, yPos, zPos, -180.0, 0.0, 0.0, false, true}
}

func MovePlayer(window *glfw.Window) {
	frameTime := glfw.GetTime()
	frameRate := frameTime - lastFrameTime
	lastFrameTime = frameTime
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
	switch {
	case window.GetKey(glfw.KeySpace) == glfw.Press && user.freeMovement:
		if !Terrain.IsInCube(user.xPos, user.yPos+(moveSpeed)-1, user.zPos, collisionDistance) {
			user.yPos = user.yPos + (moveSpeed)
		}
	case window.GetKey(glfw.KeyLeftShift) == glfw.Press && user.freeMovement:
		if !Terrain.IsInCube(user.xPos, user.yPos+(-1*moveSpeed)-1, user.zPos, collisionDistance) {
			user.yPos = user.yPos + (-1 * moveSpeed)
		}
	case !user.freeMovement:
		if Terrain.IsInCube(user.xPos, user.yPos+(user.fall*moveSpeed)-1, user.zPos, collisionDistance) {
			user.isFalling = false
			user.fall = 0.0
			user.yPos = m.Floor(user.yPos)
		} else {
			user.isFalling = true
		}
		if user.isFalling {
			user.yPos = user.yPos + (user.fall * moveSpeed)
			if user.fall > terminalVelocity {
				user.fall = user.fall - (gravity * frameRate)
			}
		}
	}
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
	var xLook, zLook float64
	if user.freeMovement {
		xLook = -1 * float64(m.Sin(float64(user.pitch)*m.Pi/180)*m.Cos(float64(user.turn)*m.Pi/180))
		zLook = -1 * float64(m.Sin(float64(user.pitch)*m.Pi/180)*m.Sin(float64(user.turn)*m.Pi/180))
	} else {
		xLook = float64(m.Cos(float64(user.turn) * m.Pi / 180))
		zLook = float64(m.Sin(float64(user.turn) * m.Pi / 180))
	}
	yLook := -1 * float64(m.Cos(float64(-1*user.pitch)*m.Pi/180))
	newX := user.xPos - (direction * xLook * moveSpeed)
	newY := user.yPos + (direction * yLook * moveSpeed)
	newZ := user.zPos - (direction * zLook * moveSpeed)
	if !Terrain.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!Terrain.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.xPos = newX
	}
	if user.freeMovement &&
		!Terrain.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!Terrain.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.yPos = newY
	}
	if !Terrain.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!Terrain.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.zPos = newZ
	}
}

func strafe(direction float64) {
	xLook := float64(m.Cos(float64(user.turn) * m.Pi / 180))
	zLook := float64(m.Sin(float64(user.turn) * m.Pi / 180))
	newX := user.xPos + (-1 * direction * zLook * moveSpeed)
	newZ := user.zPos + (direction * xLook * moveSpeed)

	if !Terrain.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!Terrain.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.xPos = newX
	}
	if !Terrain.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!Terrain.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.zPos = newZ
	}
}

func CheckPlayerCollisions(x, y, z float64, cubes []Terrain.Cube) bool {
	for _, cube := range cubes {
		if cube.CheckCollision(int(m.Floor(x)), int(m.Floor(y-1)), int(m.Floor(z))) {
			return false
		}
	}
	return true
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
	case glfw.KeySpace:
		if action == glfw.Press && !user.freeMovement && user.yPos == m.Floor(user.yPos) {
			user.fall = jumpSpeed
		}
	case glfw.KeyP:
		fmt.Printf("%v%v\n", "Player ", user)
	}

}
