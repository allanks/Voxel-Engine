package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/allanks/third-game/src/Graphics"
	"github.com/allanks/third-game/src/Player"
	"github.com/allanks/third-game/src/Terrain"
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
	program, err := Graphics.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	Terrain.GenLevel(0, -5, 0)
	Player.GenPlayer()

	projection := mgl32.Perspective(70.0, float32(WindowWidth)/WindowHeight, 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := Player.GetCameraMatrix()
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	rotate := mgl32.Ident4()
	rotateUniform := gl.GetUniformLocation(program, gl.Str("rotate\x00"))
	gl.UniformMatrix4fv(rotateUniform, 1, false, &rotate[0])

	translateUniform := gl.GetUniformLocation(program, gl.Str("translate\x00"))
	gl.Uniform4f(translateUniform, 0.0, 0.0, 0.0, 0.0)

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	Terrain.InitGCubes()
	Graphics.InitSkybox()
	Terrain.PrintCubePos()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)

		camera := Player.GetCameraMatrix()
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		x, y, z := Player.GetPosition()
		gl.Uniform4f(translateUniform, float32(x), float32(y), float32(z), 0.0)
		Graphics.RenderSkybox(vertAttrib, texCoordAttrib, translateUniform)

		Terrain.RenderLevel(vertAttrib, texCoordAttrib, translateUniform)

		Player.MovePlayer(window)
		Terrain.CheckPlayerCollisions()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func startProgram() {

}

func main() {

	go startProgram()
	initializeWindow()
}

var vertexShader string = `
#version 450

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 rotate;
uniform vec4 translate;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * rotate * (vec4(vert, 1) + translate);
}
` + "\x00"

var fragmentShader = `
#version 450

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
