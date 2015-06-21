package OpenGL45

import (
	m "math"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/go-gl/glow/gl-core/4.5/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type OpenGL45Game struct {
	Control                                                                 Graphics.OpenGLControl
	cubeProgram, mobProgram                                                 uint32
	vao, vertexBuffer, normalBuffer, typeBuffer, uvBuffer                   uint32
	stateBufferStorageBlock, sunBufferStorageBlock, textureDataStorageBlock uint32
	length, offset, normalMat, scale, mobOffset, mobNormalMat               int32
}

func (game *OpenGL45Game) CreateBuffers() {
	gl.GenVertexArrays(1, &game.vao)
	gl.GenBuffers(1, &game.vertexBuffer)
	gl.GenBuffers(1, &game.normalBuffer)
	gl.GenBuffers(1, &game.typeBuffer)
	gl.GenBuffers(1, &game.uvBuffer)
	gl.GenBuffers(1, &game.stateBufferStorageBlock)
	gl.GenBuffers(1, &game.sunBufferStorageBlock)
	gl.GenBuffers(1, &game.textureDataStorageBlock)
	gl.BindVertexArray(game.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, game.vertexBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, game.normalBuffer)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, game.typeBuffer)
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(2, 1)

	gl.BindBuffer(gl.ARRAY_BUFFER, game.uvBuffer)
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBufferBase(gl.UNIFORM_BUFFER, 0, game.stateBufferStorageBlock)
	gl.BufferData(gl.UNIFORM_BUFFER, 32*4, nil, gl.STATIC_DRAW)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 0, game.stateBufferStorageBlock, 0, 32*4)

	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, game.sunBufferStorageBlock)
	gl.BufferData(gl.UNIFORM_BUFFER, 9*4, nil, gl.STATIC_DRAW)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 1, game.sunBufferStorageBlock, 0, 9*4)

	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, game.textureDataStorageBlock)

	var sun [9]float32
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

	gl.BindBuffer(gl.UNIFORM_BUFFER, game.sunBufferStorageBlock)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 4*len(sun), gl.Ptr(&sun[0]))
}

func (game *OpenGL45Game) CreateUniforms() {

	game.length = gl.GetUniformLocation(game.cubeProgram, gl.Str("length\x00"))
	game.offset = gl.GetUniformLocation(game.cubeProgram, gl.Str("offset\x00"))
	game.normalMat = gl.GetUniformLocation(game.cubeProgram, gl.Str("normalMatrix\x00"))

	game.scale = gl.GetUniformLocation(game.mobProgram, gl.Str("scale\x00"))
	game.mobOffset = gl.GetUniformLocation(game.mobProgram, gl.Str("offset\x00"))
	game.mobNormalMat = gl.GetUniformLocation(game.mobProgram, gl.Str("normalMatrix\x00"))

	ident := mgl32.Ident4()

	gl.UniformMatrix4fv(game.normalMat, 1, true, &ident[0])
	gl.UniformMatrix4fv(game.mobNormalMat, 1, true, &ident[0])
	gl.Uniform3f(game.offset, 0.0, 0.0, 0.0)
}

func (game *OpenGL45Game) BindFragData() {

	gl.BindFragDataLocation(game.mobProgram, 0, gl.Str("outputColor\x00"))
	gl.BindFragDataLocation(game.cubeProgram, 0, gl.Str("outputColor\x00"))
}

func (game *OpenGL45Game) BindProjection(WindowWidth float32, WindowHeight float32) {

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	projection := mgl32.Perspective(70.0, WindowWidth/WindowHeight, 0.1, 100.0)
	gl.BindBuffer(gl.UNIFORM_BUFFER, game.stateBufferStorageBlock)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 4*len(projection), gl.Ptr(&projection[0]))
}

func (game *OpenGL45Game) BindBuffers(bufferData ...[]float32) {
	vertices := bufferData[0]
	normals := bufferData[1]
	ssbo := bufferData[2]
	uv := bufferData[3]

	gl.BindBuffer(gl.ARRAY_BUFFER, game.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, game.normalBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(normals)*4, gl.Ptr(normals), gl.STATIC_DRAW)

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, game.textureDataStorageBlock)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(ssbo)*4, gl.Ptr(ssbo), gl.STATIC_DRAW)
	gl.BufferSubData(gl.SHADER_STORAGE_BUFFER, 0, len(ssbo)*4, gl.Ptr(ssbo))

	gl.BindBuffer(gl.ARRAY_BUFFER, game.uvBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(uv)*4, gl.Ptr(uv), gl.STATIC_DRAW)
}

func (game *OpenGL45Game) BindUniforms(parameters ...[]float32) {
	gl.Uniform1f(game.scale, parameters[0][0])
	gl.Uniform1f(game.length, parameters[1][0])
	gl.Uniform3f(game.offset, parameters[2][0], parameters[2][1], parameters[2][2])
}

func (game *OpenGL45Game) BindTexture(texture uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, texture)
}

func (game *OpenGL45Game) RenderInstances(instances []float32, bufferSize int32) {
	gl.UseProgram(game.cubeProgram)

	gl.BindVertexArray(game.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, game.typeBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(instances)*4, gl.Ptr(instances), gl.STATIC_DRAW)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, bufferSize, int32(len(instances)/4))
}

func (game *OpenGL45Game) StartPrograms() {
	// Configure the vertex and fragment shaders
	game.cubeProgram = game.Control.NewProgram("cubeShader.shad", "cubeFrag.frag")
	game.mobProgram = game.Control.NewProgram("mobShader.shad", "mobFragment.frag")
}

func (game *OpenGL45Game) UpdateProjection(states mgl32.Mat4) {
	gl.BindBuffer(gl.UNIFORM_BUFFER, game.stateBufferStorageBlock)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 4*16, 4*len(states), gl.Ptr(&states[0]))
}
