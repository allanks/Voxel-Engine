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
	vertexBuffer,
	textureBuffer,
	positionBuffer uint32
)

func InitSkybox() {
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)

	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(skyVertices)*4, gl.Ptr(skyVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.GenBuffers(1, &textureBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(skyTexVert)*4, gl.Ptr(skyTexVert), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.GenBuffers(1, &positionBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.VertexAttribDivisor(2, 1)

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

func RenderSkybox(x, y, z float64) {
	gl.DepthMask(false)

	gl.BindVertexArray(vertexArrayObject)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, topTexture)

	position := []float32{float32(x), float32(y), float32(z)}

	gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(position)*4, gl.Ptr(position), gl.STATIC_DRAW)
	//gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 24)
	//gl.DrawElementsInstanced(gl.TRIANGLE_STRIP, 24, gl.UNSIGNED_INT, gl.Ptr(position), 3)

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
	//  X, Y, Z
	// Front face
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,
	// Back face
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,
	// Left face
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	// Right face
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	// Top face
	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	// Bottom face
	1.0, -1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,
}

var skyTexVert = []float32{
	// U, V
	// Front face
	1.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
	// Back face
	1.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
	// Left face
	1.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
	// Right face
	1.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
	// Top face
	0.0, 1.0,
	0.0, 0.0,
	1.0, 1.0,
	1.0, 0.0,
	// Bottom face
	0.0, 1.0,
	0.0, 0.0,
	1.0, 1.0,
	1.0, 0.0,
}
