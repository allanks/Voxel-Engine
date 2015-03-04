package Terrain

import (
	"github.com/allanks/third-game/src/Graphics"
	"github.com/go-gl/glow/gl-core/4.5/gl"
)

var (
	texture, vertexArrayObject, vertexBufferObject uint32
)

func InitCube() {
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)

	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
	// Load the texture
	var err error
	texture, err = Graphics.NewTexture("resource/texture/square.png")
	if err != nil {
		panic(err)
	}
}

func Render(cube *Cube, vertAttrib, texCoordAttrib uint32, translateUniform int32) {

	xPos, yPos, zPos := cube.GetPos()
	gl.Uniform4f(translateUniform, float32(xPos), float32(yPos), float32(zPos), 0.0)

	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.BindVertexArray(vertexArrayObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4*6)

}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Front face
	1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 0.0, 1.0, 0.0, 0.0,
	0.0, 1.0, 1.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0, 0.0,
	// Back face
	0.0, 1.0, 0.0, 0.0, 1.0,
	0.0, 0.0, 0.0, 0.0, 0.0,
	1.0, 1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 0.0, 1.0, 0.0,
	// Left face
	0.0, 1.0, 1.0, 0.0, 1.0,
	0.0, 0.0, 1.0, 0.0, 0.0,
	0.0, 1.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 0.0, 1.0, 0.0,
	// Right face
	1.0, 1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0, 0.0,
	// Top face
	0.0, 1.0, 0.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 0.0, 0.0,
	0.0, 1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 0.0,
	// Bottom face
	1.0, 0.0, 0.0, 0.0, 1.0,
	0.0, 0.0, 0.0, 0.0, 0.0,
	1.0, 0.0, 1.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0, 0.0,
}
