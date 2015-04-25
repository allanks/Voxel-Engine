package main

import (
	"fmt"
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
const textureAtlas string = "resource/texture/textureAtlas.png"

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
	program, err := Graphics.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	scale := gl.GetUniformLocation(program, gl.Str("scale\x00"))
	offset := gl.GetUniformLocation(program, gl.Str("offset\x00"))
	texParam := gl.GetUniformLocation(program, gl.Str("length\x00"))
	textureDataStorageBlock := gl.GetProgramResourceIndex(program, gl.SHADER_STORAGE_BLOCK, gl.Str("texture_data\x00"))
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

	fmt.Println("Creating Texture Buffer")

	bindTextureBuffer(program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

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
		Model.Render(typeBuffer, position, Model.Cube)
		gl.DepthMask(true)
		gl.Uniform3f(offset, 0.0, 0.0, 0.0)
		Player.Render(vao, typeBuffer, offset)

		//Model.BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock, 1)
		//Model.Render(vao, typeBuffer)

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

func bindTextureBuffer(program uint32) {

	texture, err := Graphics.NewTexture(textureAtlas)
	if err != nil {
		panic(err)
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, texture)
}

func main() {
	//Terrain.PackTextures()
	initializeWindow()
}

var vertexShader string = `
#version 450

uniform mat4 projection;
uniform mat4 camera;
uniform vec3 offset;
uniform float scale;
uniform float length;

layout(std430,binding=0) buffer texture_data {
	vec2 textureData[];
}texData;

layout(location=0) in vec3 vert; // vertex position
layout(location=1) in vec3 normal; // normal position
layout(location=2) in vec4 object; // instance data, unique to each object (instance)

in int gl_VertexID;

out vec2 fragData;

void main() {
	int ind = gl_VertexID+(int(object.w*length));
	fragData = texData.textureData[ind];
    vec3 vertexData = vert;
    vertexData.x *= scale;
    vertexData.y *= scale;
    vertexData.z *= scale;
    gl_Position = projection * camera *  (vec4( vertexData + vec3(object.x + offset.x, object.y + offset.y, object.z+offset.z), 1));
}
` + "\x00"

var fragmentShader = `
#version 450

uniform sampler2D tex;

in vec2 fragData;

out vec4 outputColor;

void main() {
	outputColor = texture(tex, fragData);
}
` + "\x00"
