package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	gamegl "github.com/allanks/Voxel-Engine/src/Graphics/Game/OpenGL45"
	controlgl "github.com/allanks/Voxel-Engine/src/Graphics/OpenGL45"
	"github.com/allanks/Voxel-Engine/src/Model"
	"github.com/allanks/Voxel-Engine/src/Player"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	eyeVec, viewVec                 mgl32.Vec3
	gameController, titleController Graphics.OpenGLController
	openGLControl                   Graphics.OpenGLControl
)

const WindowWidth = 800
const WindowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	// Have all the cpus
	runtime.GOMAXPROCS(runtime.NumCPU())
	openGLControl = &controlgl.OpenGLControl{}
	gameController = &gamegl.OpenGL45Game{Control: openGLControl}
	Model.Controller = gameController
	Model.Control = openGLControl
}

func glDebugCallback(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {
	fmt.Printf("Debug source=%d type=%d severity=%d: %s\n", source, gltype, severity, message)
}

func initializeWindow() {

	// Initialize GLFW for window management
	if glfw.Init() != nil {
		panic("failed to initialize glfw")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetCursorPos(0, 0)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Initialize Glow
	openGLControl.Init()

	window.SetKeyCallback(Player.OnKey)
	window.SetCursorPosCallback(Player.OnCursor)
	initOpenGLProgram(window)
}

func initOpenGLProgram(window *glfw.Window) {
	drawTitle := false
	drawGame := true
	for !window.ShouldClose() {
		if drawTitle {
			//titleScreenLoop(window, &drawTitle, &drawGame)
			drawTitle = false
			drawGame = true
		} else if drawGame {
			mainGameLoop(window, &drawGame, &drawTitle)
		}
	}
}

func titleScreenLoop(window *glfw.Window, drawTitle, drawGame *bool) {
	/*program, err := Graphics.NewProgram("titleShader.shad", "titleFrag.frag")
	if err != nil {
		panic(err)
	}
	color := gl.GetUniformLocation(program, gl.Str("vColor\x00"))
	gl.Uniform4f(color, 0, 0, 0, 1)

	var vao, vertexBuffer, textureDataStorageBlock, objectBuffer, stateBufferStorageBlock uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &objectBuffer)
	gl.GenBuffers(1, &vertexBuffer)
	gl.GenBuffers(1, &textureDataStorageBlock)
	gl.GenBuffers(1, &stateBufferStorageBlock)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, textureDataStorageBlock)

	setOrthoganolProjection(window, stateBufferStorageBlock)

	gl.BindBuffer(gl.ARRAY_BUFFER, objectBuffer)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(1, 1)

	fmt.Println("Loading ttf")
	ttf, err := os.Open("resource/fonts/Arial/arial.ttf")
	if err != nil {
		panic(err)
	}
	fmt.Println("Generating Fonts")
	font, err := glText45.LoadTruetype(ttf, 12, 32, 127, glText45.LeftToRight, vertexBuffer, textureDataStorageBlock)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting Draw loop Draw Title %v\n", *drawTitle)
	for !window.ShouldClose() && *drawTitle {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		font.DisplayString(0, 0, objectBuffer, "%v", "Printing String")

		window.SwapBuffers()
		glfw.PollEvents()
	}*/

}

func setOrthoganolProjection(window *glfw.Window, stateBufferStorageBlock uint32) {
	/*width, height := window.GetSize()
	projection := mgl32.Ortho2D(0.0, float32(width), float32(height), 0.0)
	gl.BindBuffer(gl.UNIFORM_BUFFER, stateBufferStorageBlock)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 4*len(projection), gl.Ptr(&projection[0]))*/
}

func mainGameLoop(window *glfw.Window, drawGame, drawTitle *bool) {
	gameController.StartPrograms()
	gameController.CreateUniforms()
	gameController.CreateBuffers()

	camera := Player.GetCameraMatrix()
	Model.InitGCubes()
	Model.InitModels()

	Player.GenPlayer(5, 68, 5)
	xPos, _ := window.GetCursorPos()
	window.SetCursorPos(xPos, -180)

	//gopher := []float32{0, 63, 5, 0}

	gameController.BindProjection(float32(WindowWidth), float32(WindowHeight))

	fmt.Println("Starting Draw Loop")

	go func() {
		ticker := time.Tick(16 * time.Millisecond)
		for *drawGame {
			Player.MovePlayer(window)
			<-ticker
		}
	}()

	for !window.ShouldClose() && *drawGame {
		openGLControl.Clear()

		camera = Player.GetCameraMatrix()

		gameController.UpdateProjection(camera)

		x, y, z := Player.GetPosition()
		position := []float32{float32(x), float32(y), float32(z), Model.SkyBox}

		openGLControl.DepthToggle(false)

		Model.BindBuffers([]float32{-0.5, -0.5, -0.5}, Model.Cube)
		Model.Render(position, Model.Cube)

		openGLControl.DepthToggle(true)

		Model.BindBuffers([]float32{0.0, 0.0, 0.0}, Model.Cube)
		Player.Render()

		//gl.UseProgram(mobProgram)
		//gl.BindVertexArray(vao)
		//Model.BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock, uvBuffer, scale, length, Model.Gopher)
		//Model.Render(vao, gopher, Model.Gopher)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func main() {
	//Terrain.PackTextures()
	initializeWindow()
}
