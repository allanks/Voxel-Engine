package main

import (
	"fmt"
	m "math"
	"runtime"
	"time"
	"unsafe"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/allanks/Voxel-Engine/src/Model"
	"github.com/allanks/Voxel-Engine/src/Player"
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

	// Configure the vertex and fragment shaders
	program, err := Graphics.NewProgram("vertexShader.shad", "fragmentShader.frag")
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	scale := gl.GetUniformLocation(program, gl.Str("scale\x00"))
	offset := gl.GetUniformLocation(program, gl.Str("offset\x00"))
	texParam := gl.GetUniformLocation(program, gl.Str("length\x00"))
	textureDataStorageBlock := gl.GetProgramResourceIndex(program, gl.SHADER_STORAGE_BLOCK, gl.Str("texture_data\x00"))
	sunColor := gl.GetUniformLocation(program, gl.Str("sun.vColor\x00"))
	sunDirection := gl.GetUniformLocation(program, gl.Str("sun.vDirection\x00"))
	sunIntensity := gl.GetUniformLocation(program, gl.Str("sun.intensity\x00"))
	normalMat := gl.GetUniformLocation(program, gl.Str("normalMatrix\x00"))
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	bindProjection(program)

	camera := Player.GetCameraMatrix()
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	fmt.Println("Initialising GCubes")

	Model.InitGCubes()

	fmt.Println("Generating Player")

	Player.GenPlayer(5, 68, 5)
	xPos, _ := window.GetCursorPos()
	window.SetCursorPos(xPos, -180)

	fmt.Println("Initialising Buffers")

	var vao, vertexBuffer, normalBuffer, typeBuffer, indexBuffer uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vertexBuffer)
	gl.GenBuffers(1, &normalBuffer)
	gl.GenBuffers(1, &typeBuffer)
	gl.GenBuffers(1, &indexBuffer)
	gl.GenBuffers(1, &textureDataStorageBlock)

	bindBuffers(vao, vertexBuffer, normalBuffer, typeBuffer, indexBuffer)

	fmt.Println("Loading models")

	Model.InitModels()

	gopher := []float32{0, 63, 5, 0}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	fmt.Println("Creating light")

	sunX := float32(m.Cos(250.0 * m.Pi / 180))
	sunY := float32(m.Sin(250.0 * m.Pi / 180))
	sunZ := float32(m.Cos(60 * m.Pi / 180))

	gl.Uniform3f(sunColor, 1.0, 1.0, 1.0)

	ident := mgl32.Ident4()

	gl.UniformMatrix4fv(normalMat, 1, true, &ident[0])

	fmt.Println("Starting Draw Loop")

	go func() {
		ticker := time.Tick(16 * time.Millisecond)
		for {
			Player.MovePlayer(window)
			<-ticker
		}
	}()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		camera = Player.GetCameraMatrix()
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		x, y, z := Player.GetPosition()
		position := []float32{float32(x), float32(y), float32(z), 1}

		gl.BindVertexArray(vao)
		Model.BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock, scale, texParam, Model.Cube)
		gl.DepthMask(false)
		gl.BindVertexArray(vao)
		gl.Uniform3f(offset, -0.5, -0.5, -0.5)
		gl.Uniform3f(sunDirection, 0.0, 0.0, 0.0)
		gl.Uniform1f(sunIntensity, 1.0)
		Model.Render(typeBuffer, position, Model.Cube)
		gl.DepthMask(true)
		gl.Uniform3f(offset, 0.0, 0.0, 0.0)
		gl.Uniform3f(sunDirection, sunX, sunY, sunZ)
		gl.Uniform1f(sunIntensity, 0.5)
		Player.Render(vao, typeBuffer, offset)

		gl.BindVertexArray(vao)
		Model.BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock, scale, texParam, Model.Gopher)
		Model.Render(vao, gopher, Model.Gopher)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func bindBuffers(vao, vertexBuffer, normalBuffer, typeBuffer, indexBuffer uint32) {
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
}

func bindProjection(program uint32) {
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	projection := mgl32.Perspective(70.0, float32(WindowWidth)/WindowHeight, 0.1, 100.0)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
}

func main() {
	//Terrain.PackTextures()
	initializeWindow()
}
