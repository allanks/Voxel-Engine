package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/allanks/Voxel-Engine/src/Mob"
	"github.com/allanks/Voxel-Engine/src/Player"
	"github.com/allanks/Voxel-Engine/src/Terrain"
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
	texParam := gl.GetUniformLocation(program, gl.Str("texParam\x00"))
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	bindProjection(program)

	camera := Player.GetCameraMatrix()
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	fmt.Println("Initialising GCubes")

	Terrain.InitGCubes()

	fmt.Println("Generating Player")

	Player.GenPlayer(5, 68, 5)

	fmt.Println("Initialising Buffers")

	var vao, vertexBuffer, typeBuffer, indexBuffer uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vertexBuffer)
	gl.GenBuffers(1, &typeBuffer)
	gl.GenBuffers(1, &indexBuffer)

	bindBuffers(vao, vertexBuffer, typeBuffer, indexBuffer)

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

		gl.Uniform1f(scale, 1)
		gl.Uniform2f(texParam, 0, 24)
		Terrain.BindCubeVertexBuffers(vertexBuffer, indexBuffer)
		Terrain.RenderSkyBox(vao, typeBuffer, offset, x, y, z)
		Player.Render(vao, typeBuffer, offset)
		gl.Uniform3f(offset, 0.0, 0.0, 0.0)
		Mob.BindVertices(vertexBuffer, indexBuffer, 0)
		Mob.Render(vao, typeBuffer, scale, texParam)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func bindBuffers(vao, vertexBuffer, typeBuffer, indexBuffer uint32) {
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, typeBuffer)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(1, 1)
}

func bindProjection(program uint32) {
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	projection := mgl32.Perspective(70.0, float32(WindowWidth)/WindowHeight, 0.1, 100.0)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
}

func bindTextureBuffer(program uint32) {
	fmt.Println("Loading models")

	mobTexture := Mob.InitModels()

	texture, err := Graphics.NewTexture(textureAtlas)
	if err != nil {
		panic(err)
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	textureDataStorageBlock := gl.GetProgramResourceIndex(program, gl.SHADER_STORAGE_BLOCK, gl.Str("texture_data\x00"))
	textureBuffer := Terrain.GetTextureBuffer()
	textureBuffer = append(textureBuffer, mobTexture...)

	gl.GenBuffers(1, &textureDataStorageBlock)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, textureDataStorageBlock)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(textureBuffer)*4, gl.Ptr(textureBuffer), gl.STATIC_DRAW)
}

func main() {
	//Terrain.PackTextures()
	initializeWindow()
}

var vertexShader string = `
#version 450

uniform mat4 projection;
uniform mat4 camera;
uniform float scale;
uniform vec3 offset;
uniform vec2 texParam;

layout(std430,binding=0) buffer texture_data {
	vec2 textureData[];
}texData;

layout(location=0) in vec3 vert; // cube vertex position
layout(location=1) in vec4 cube; // instance data, unique to each object (instance)

in int gl_VertexID;

out vec2 fragTexCoord;

void main() {
    fragTexCoord =  texData.textureData[gl_VertexID+(int(texParam.x+(cube.w*texParam.y)))];
    vec3 vertexData = vert;
    vertexData.x *= scale;
    vertexData.y *= scale;
    vertexData.z *= scale;
    gl_Position = projection * camera *  (vec4( vertexData + vec3(cube.x + offset.x, cube.y + offset.y, cube.z+offset.z), 1));
}
` + "\x00"

var fragmentShader = `
#version 450

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, (fragTexCoord*0.25));
}
` + "\x00"
