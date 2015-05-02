package Player

import (
	"fmt"
	m "math"
	"time"

	"github.com/allanks/Voxel-Engine/src/Server/DataType"
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
	gameMap                             *Terrain.Level
}

type moveFunc func(float64)

func GenPlayer(xPos, yPos, zPos float64) {
	lastFrameTime = glfw.GetTime()
	user = player{xPos, yPos, zPos, -180.0, 0.0, 0.0, true, true, &Terrain.Level{}}

	Terrain.StartConnection()
	user.gameMap.InitChunk(int(m.Floor(float64(xPos))), int(m.Floor(float64(zPos))))
	go user.loopChunkLoader()
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
		user.moveY(user.yPos + (moveSpeed))
	case window.GetKey(glfw.KeyLeftShift) == glfw.Press && user.freeMovement:
		user.moveY(user.yPos + (-1 * moveSpeed))
	case !user.freeMovement:
		user.moveY(user.yPos + (user.fall * moveSpeed) - 1)
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

func checkCubes(cubes []DataType.Pos, current, final DataType.Pos) bool {
	for _, cube := range cubes {
		if (current.XPos < cube.XPos && final.XPos > cube.XPos) || (current.XPos > cube.XPos && final.XPos < cube.XPos) ||
			(current.YPos < cube.YPos && final.YPos > cube.YPos) || (current.YPos > cube.YPos && final.YPos < cube.YPos) ||
			(current.ZPos < cube.ZPos && final.ZPos > cube.ZPos) || (current.ZPos > cube.ZPos && final.ZPos < cube.ZPos) {
			return false
		}
	}
	return true
}

func (user *player) moveY(newY float64) {
	upperCube := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos + 1), ZPos: float32(user.zPos)}
	lowerCube := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos - Height - 1), ZPos: float32(user.zPos)}

	cubes := user.gameMap.GetCubes([]DataType.Pos{upperCube, lowerCube})

	topCurrent := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos), ZPos: float32(user.zPos)}
	topFinal := DataType.Pos{XPos: float32(user.xPos), YPos: float32(newY), ZPos: float32(user.zPos)}
	bottomCurrent := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos)}
	bottomFinal := DataType.Pos{XPos: float32(user.xPos), YPos: float32(newY - Height), ZPos: float32(user.zPos)}

	if checkCubes(cubes, topCurrent, topFinal) || checkCubes(cubes, bottomCurrent, bottomFinal) {
		user.yPos = newY
	} else {
		user.fall = 0.0
		user.yPos = m.Floor(user.yPos)
	}
}

func moveXZ(newX, newZ float64) {
	upperFrontCube := DataType.Pos{XPos: float32(user.xPos + 1), YPos: float32(user.yPos), ZPos: float32(user.zPos)}
	lowerFrontCube := DataType.Pos{XPos: float32(user.xPos + 1), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos)}
	upperBackCube := DataType.Pos{XPos: float32(user.xPos - 1), YPos: float32(user.yPos), ZPos: float32(user.zPos)}
	lowerBackCube := DataType.Pos{XPos: float32(user.xPos - 1), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos)}
	upperLeftCube := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos), ZPos: float32(user.zPos - 1)}
	lowerLeftCube := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos - 1)}
	upperRightCube := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos), ZPos: float32(user.zPos + 1)}
	lowerRightCube := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos + 1)}

	cubes := user.gameMap.GetCubes([]DataType.Pos{upperFrontCube, lowerFrontCube, upperBackCube, lowerBackCube, upperLeftCube, lowerLeftCube, upperRightCube, lowerRightCube})

	topCurrent := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos), ZPos: float32(user.zPos)}
	topFinal := DataType.Pos{XPos: float32(newX), YPos: float32(user.yPos), ZPos: float32(user.zPos)}
	bottomCurrent := DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos)}
	bottomFinal := DataType.Pos{XPos: float32(newX), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos)}
	if checkCubes(cubes, topCurrent, topFinal) || checkCubes(cubes, bottomCurrent, bottomFinal) {
		user.xPos = newX
	}

	topCurrent = DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos), ZPos: float32(user.zPos)}
	topFinal = DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos), ZPos: float32(newZ)}
	bottomCurrent = DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos - Height), ZPos: float32(user.zPos)}
	bottomFinal = DataType.Pos{XPos: float32(user.xPos), YPos: float32(user.yPos - Height), ZPos: float32(newZ)}
	if checkCubes(cubes, topCurrent, topFinal) || checkCubes(cubes, bottomCurrent, bottomFinal) {
		user.zPos = newZ
	}
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
	if user.freeMovement {
		user.moveY(newY)
	}
	moveXZ(newX, newZ)
}

func strafe(direction float64) {
	xLook := float64(m.Cos(float64(user.turn) * m.Pi / 180))
	zLook := float64(m.Sin(float64(user.turn) * m.Pi / 180))
	newX := user.xPos + (-1 * direction * zLook * moveSpeed)
	newZ := user.zPos + (direction * xLook * moveSpeed)
	moveXZ(newX, newZ)
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
			user.fall = 0.0
		}
	case glfw.KeySpace:
		if action == glfw.Press && !user.freeMovement && user.yPos == m.Floor(user.yPos) {
			user.fall = jumpSpeed
		}
	case glfw.KeyP:
		fmt.Printf("Player X %v, Y %v, Z %v Free %v\n", int(m.Floor(user.xPos)), int(m.Floor(user.yPos)), int(m.Floor(user.zPos)), user.freeMovement)
	case glfw.KeyC:
		fmt.Printf("Camera %v\n", GetCameraMatrix())
	case glfw.KeyG:
		fmt.Printf("Near Y Cubes %v\n", user.gameMap.GetYCubes(user.xPos, user.yPos, user.zPos, Height))
	}
}

func (p *player) loopChunkLoader() {
	for {
		p.gameMap.LoopChunkLoader(p.xPos, p.zPos)
		time.Sleep(1 * time.Second)
	}
}

func Render(vao, typeBuffer uint32, offset int32) {
	user.gameMap.RenderLevel(vao, typeBuffer, offset)
}
