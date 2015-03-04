package Graphics

import (
	"github.com/go-gl/glow/gl-core/4.5/gl"
)

var (
	topTexture,
	bottomTexture,
	frontTexture,
	backTexture,
	leftTexture,
	rightTexture,
	vertexArrayObject,
	vertexBufferObject uint32
)

func InitSkybox() {
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)

	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(skyVertices)*4, gl.Ptr(skyVertices), gl.STATIC_DRAW)

	// Load the textures
	var err error
	topTexture, err = NewTexture("resource/texture/SkyBox/jajlands1_up.png")
	if err != nil {
		panic(err)
	}
	bottomTexture, err = NewTexture("resource/texture/SkyBox/jajlands1_dn.png")
	if err != nil {
		panic(err)
	}
	frontTexture, err = NewTexture("resource/texture/SkyBox/jajlands1_ft.png")
	if err != nil {
		panic(err)
	}
	backTexture, err = NewTexture("resource/texture/SkyBox/jajlands1_bk.png")
	if err != nil {
		panic(err)
	}
	leftTexture, err = NewTexture("resource/texture/SkyBox/jajlands1_rt.png")
	if err != nil {
		panic(err)
	}
	rightTexture, err = NewTexture("resource/texture/SkyBox/jajlands1_lf.png")
	if err != nil {
		panic(err)
	}
}

func RenderSkybox(vertAttrib, texCoordAttrib uint32, translateUniform int32) {
	gl.DepthMask(false)

	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.BindVertexArray(vertexArrayObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, frontTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, backTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 4, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, leftTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 8, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, rightTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 12, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, topTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 16, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, bottomTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 20, 4)

	gl.DepthMask(true)
}

var skyVertices = []float32{
	//  X, Y, Z, U, V
	// Front face
	1.0, 1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	// Back face
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 0.0, 1.0,
	// Left face
	-1.0, 1.0, 1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 1.0,
	// Right face
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 1.0,
	1.0, 1.0, 1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 1.0,
	// Top face
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 0.0,
	// Bottom face
	1.0, -1.0, -1.0, 0.0, 1.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 1.0, 0.0,
}
