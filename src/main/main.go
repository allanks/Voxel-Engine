package main

import (
	"fmt"
	m "math"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/allanks/Voxel-Engine/src/Model"
	"github.com/allanks/Voxel-Engine/src/Player"
	"github.com/allanks/Voxel-Engine/src/glText45"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/glow/gl-core/4.5/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	eyeVec, viewVec mgl32.Vec3
)

const WindowWidth = 800
const WindowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	// Have all the cpus
	runtime.GOMAXPROCS(runtime.NumCPU())
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
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Query the extensions to determine if we can enable the debug callback
	var numExtensions int32
	gl.GetIntegerv(gl.NUM_EXTENSIONS, &numExtensions)

	extensions := make(map[string]bool)
	for i := int32(0); i < numExtensions; i++ {
		extension := gl.GoStr(gl.GetStringi(gl.EXTENSIONS, uint32(i)))
		extensions[extension] = true
	}

	if _, ok := extensions["GL_ARB_debug_output"]; ok {
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS_ARB)
		gl.DebugMessageCallbackARB(gl.DebugProc(glDebugCallback), gl.Ptr(nil))
		// Trigger an error to demonstrate debug output
		gl.Enable(gl.CONTEXT_FLAGS)
	}
	window.SetKeyCallback(Player.OnKey)
	window.SetCursorPosCallback(Player.OnCursor)
	initOpenGLProgram(window)
}

func initOpenGLProgram(window *glfw.Window) {
	drawTitle := false
	drawGame := true
	for !window.ShouldClose() {
		if drawTitle {
			titleScreenLoop(window, &drawTitle, &drawGame)
		} else if drawGame {
			mainGameLoop(window, &drawGame, &drawTitle)
		}
	}
}

func titleScreenLoop(window *glfw.Window, drawTitle, drawGame *bool) {
	program, err := Graphics.NewProgram("titleShader.shad", "titleFrag.frag")
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
	}

}

func setOrthoganolProjection(window *glfw.Window, stateBufferStorageBlock uint32) {
	width, height := window.GetSize()
	projection := mgl32.Ortho2D(0.0, float32(width), float32(height), 0.0)
	gl.BindBuffer(gl.UNIFORM_BUFFER, stateBufferStorageBlock)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 4*len(projection), gl.Ptr(&projection[0]))
}

func mainGameLoop(window *glfw.Window, drawGame, drawTitle *bool) {
	// Configure the vertex and fragment shaders
	cubeProgram, err := Graphics.NewProgram("cubeShader.shad", "cubeFrag.frag")
	if err != nil {
		panic(err)
	}
	mobProgram, err := Graphics.NewProgram("mobShader.shad", "mobFragment.frag")
	if err != nil {
		panic(err)
	}

	length := gl.GetUniformLocation(cubeProgram, gl.Str("length\x00"))
	offset := gl.GetUniformLocation(cubeProgram, gl.Str("offset\x00"))
	normalMat := gl.GetUniformLocation(cubeProgram, gl.Str("normalMatrix\x00"))
	gl.BindFragDataLocation(cubeProgram, 0, gl.Str("outputColor\x00"))

	scale := gl.GetUniformLocation(mobProgram, gl.Str("scale\x00"))
	mobOffset := gl.GetUniformLocation(mobProgram, gl.Str("offset\x00"))
	mobNormalMat := gl.GetUniformLocation(mobProgram, gl.Str("normalMatrix\x00"))
	gl.BindFragDataLocation(mobProgram, 0, gl.Str("outputColor\x00"))

	camera := Player.GetCameraMatrix()

	fmt.Println("Initialising GCubes")

	Model.InitGCubes()

	fmt.Println("Generating Player")

	Player.GenPlayer(5, 68, 5)
	xPos, _ := window.GetCursorPos()
	window.SetCursorPos(xPos, -180)

	fmt.Println("Initialising Buffers")

	var vao, vertexBuffer, normalBuffer, typeBuffer, uvBuffer, stateBufferStorageBlock, sunBufferStorageBlock, textureDataStorageBlock uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vertexBuffer)
	gl.GenBuffers(1, &normalBuffer)
	gl.GenBuffers(1, &typeBuffer)
	gl.GenBuffers(1, &uvBuffer)
	gl.GenBuffers(1, &stateBufferStorageBlock)
	gl.GenBuffers(1, &sunBufferStorageBlock)
	gl.GenBuffers(1, &textureDataStorageBlock)

	bindBuffers(vao, vertexBuffer, normalBuffer, typeBuffer, uvBuffer, textureDataStorageBlock, stateBufferStorageBlock, sunBufferStorageBlock)

	fmt.Println("Loading models")

	Model.InitModels()

	gopher := []float32{0, 63, 5, 0}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	fmt.Println("Creating light")

	var sun [9]float32

	bindProjection(stateBufferStorageBlock)

	// Sun color
	sun[0] = 1.0
	sun[1] = 1.0
	sun[2] = 1.0
	sun[3] = 0.0
	// Sun position
	sun[4] = float32(m.Cos(250.0 * m.Pi / 180))
	sun[5] = float32(m.Sin(250.0 * m.Pi / 180))
	sun[6] = float32(m.Cos(60 * m.Pi / 180))
	sun[7] = 0.0
	// Sun intensity
	sun[8] = 0.8

	gl.BindBuffer(gl.UNIFORM_BUFFER, sunBufferStorageBlock)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 4*len(sun), gl.Ptr(&sun[0]))

	ident := mgl32.Ident4()

	gl.UniformMatrix4fv(normalMat, 1, true, &ident[0])
	gl.UniformMatrix4fv(mobNormalMat, 1, true, &ident[0])
	gl.Uniform3f(mobOffset, 0.0, 0.0, 0.0)

	fmt.Println("Starting Draw Loop")

	go func() {
		ticker := time.Tick(16 * time.Millisecond)
		for *drawGame {
			Player.MovePlayer(window)
			<-ticker
		}
	}()

	for !window.ShouldClose() && *drawGame {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		camera = Player.GetCameraMatrix()

		gl.BindBuffer(gl.UNIFORM_BUFFER, stateBufferStorageBlock)
		gl.BufferSubData(gl.UNIFORM_BUFFER, 4*16, 4*len(camera), gl.Ptr(&camera[0]))

		x, y, z := Player.GetPosition()
		position := []float32{float32(x), float32(y), float32(z), 1}

		gl.UseProgram(cubeProgram)
		gl.BindVertexArray(vao)
		Model.BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock, uvBuffer, scale, length, Model.Cube)
		gl.DepthMask(false)
		gl.Uniform3f(offset, -0.5, -0.5, -0.5)
		Model.Render(typeBuffer, position, Model.Cube)
		gl.DepthMask(true)
		gl.Uniform3f(offset, 0.0, 0.0, 0.0)
		Player.Render(vao, typeBuffer, offset)

		gl.UseProgram(mobProgram)
		gl.BindVertexArray(vao)
		Model.BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock, uvBuffer, scale, length, Model.Gopher)
		Model.Render(vao, gopher, Model.Gopher)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func bindProjection(stateBufferStorageBlock uint32) {
	projection := mgl32.Perspective(70.0, float32(WindowWidth)/WindowHeight, 0.1, 100.0)
	gl.BindBuffer(gl.UNIFORM_BUFFER, stateBufferStorageBlock)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 4*len(projection), gl.Ptr(&projection[0]))
}

func bindBuffers(vao, vertexBuffer, normalBuffer, typeBuffer, uvBuffer, textureDataStorageBlock, stateBufferStorageBlock, sunBufferStorageBlock uint32) {
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, typeBuffer)
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(2, 1)

	gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBufferBase(gl.UNIFORM_BUFFER, 0, stateBufferStorageBlock)
	gl.BufferData(gl.UNIFORM_BUFFER, 32*4, nil, gl.STATIC_DRAW)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 0, stateBufferStorageBlock, 0, 32*4)

	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, sunBufferStorageBlock)
	gl.BufferData(gl.UNIFORM_BUFFER, 9*4, nil, gl.STATIC_DRAW)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 1, sunBufferStorageBlock, 0, 9*4)

	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, textureDataStorageBlock)
}

func main() {
	//Terrain.PackTextures()
	initializeWindow()
}
