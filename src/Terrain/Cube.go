package Terrain

import (
	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/go-gl/glow/gl-core/4.5/gl"
	"gopkg.in/mgo.v2/bson"
)

const (
	collisionDistance float64 = 0.15
	textureAtlas      string  = "resource/texture/textureAtlas.png"
)
const (
	// Cube Types
	Dirt = iota
	Grass
	Stone
	CobbleStone
	Gravel
)

type skyBox struct {
	texture  []float32
	position []float32
}

type GCube struct {
	Texture []float32
	Gtype   int
}

type Cube struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ChunkID  bson.ObjectId
	XPos     int
	YPos     int
	ZPos     int
	CubeType int
}

func (cube *Cube) GetCubeType() int {
	return cube.CubeType
}

func (cube *Cube) CheckCollision(xPos, yPos, zPos int) bool {

	if xPos == cube.XPos &&
		yPos == cube.YPos &&
		zPos == cube.ZPos {
		return true
	}
	return false
}

var sky skyBox

var GCubes = []GCube{
	GCube{},
	GCube{},
	GCube{},
	GCube{},
	GCube{},
}

var instances int32

func InitialiseGCubeBuffers() (uint32, uint32, uint32) {
	var vao, positionBuffer uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vertexBuffer, textureBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.GenBuffers(1, &textureBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureBuffer)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.GenBuffers(1, &positionBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(2, 1)

	texture, err := Graphics.NewTexture(textureAtlas)
	if err != nil {
		panic(err)
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	return vao, positionBuffer, textureBuffer
}

func RenderSkyBox(vao, positionBuffer, textureBuffer uint32) {
	gl.DepthMask(false)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, textureBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(sky.texture)*4, gl.Ptr(sky.texture), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(sky.position)*4, gl.Ptr(sky.position), gl.STATIC_DRAW)

	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 24, int32(1))

	gl.DepthMask(true)
}

var cubeVertices = []float32{
	//  X, Y, Z,
	// Front face
	1.0, 1.0, 1.0,
	1.0, 0.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 0.0, 1.0,
	// Back face
	0.0, 1.0, 0.0,
	0.0, 0.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 0.0, 0.0,
	// Left face
	0.0, 1.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 0.0,
	// Right face
	1.0, 1.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 1.0, 1.0,
	1.0, 0.0, 1.0,
	// Top face
	0.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	0.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	// Bottom face
	1.0, 0.0, 0.0,
	0.0, 0.0, 0.0,
	1.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
}

func InitGCubes() {
	GCubes[Dirt].Gtype = Dirt
	GCubes[Grass].Gtype = Grass
	GCubes[Stone].Gtype = Stone
	GCubes[CobbleStone].Gtype = CobbleStone
	GCubes[Gravel].Gtype = Gravel

	GCubes[Dirt].Texture = []float32{
		2, 0,
		2, 1,
		1, 0,
		1, 1,

		2, 0,
		2, 1,
		1, 0,
		1, 1,

		2, 0,
		2, 1,
		1, 0,
		1, 1,

		2, 0,
		2, 1,
		1, 0,
		1, 1,

		2, 0,
		2, 1,
		1, 0,
		1, 1,

		2, 0,
		2, 1,
		1, 0,
		1, 1,
	}
	GCubes[Grass].Texture = []float32{
		1, 1,
		1, 2,
		0, 1,
		0, 2,

		1, 1,
		1, 2,
		0, 1,
		0, 2,

		1, 1,
		1, 2,
		0, 1,
		0, 2,

		1, 1,
		1, 2,
		0, 1,
		0, 2,

		1, 1,
		1, 2,
		0, 1,
		0, 2,

		1, 1,
		1, 2,
		0, 1,
		0, 2,
	}
	GCubes[Stone].Texture = []float32{
		2, 2,
		2, 3,
		1, 2,
		1, 3,

		2, 2,
		2, 3,
		1, 2,
		1, 3,

		2, 2,
		2, 3,
		1, 2,
		1, 3,

		2, 2,
		2, 3,
		1, 2,
		1, 3,

		2, 2,
		2, 3,
		1, 2,
		1, 3,

		2, 2,
		2, 3,
		1, 2,
		1, 3,
	}
	GCubes[CobbleStone].Texture = []float32{
		1, 0,
		1, 1,
		0, 0,
		0, 1,

		1, 0,
		1, 1,
		0, 0,
		0, 1,

		1, 0,
		1, 1,
		0, 0,
		0, 1,

		1, 0,
		1, 1,
		0, 0,
		0, 1,

		1, 0,
		1, 1,
		0, 0,
		0, 1,

		1, 0,
		1, 1,
		0, 0,
		0, 1,
	}
	GCubes[Gravel].Texture = []float32{
		2, 1,
		2, 2,
		1, 1,
		1, 2,

		2, 1,
		2, 2,
		1, 1,
		1, 2,

		2, 1,
		2, 2,
		1, 1,
		1, 2,

		2, 1,
		2, 2,
		1, 1,
		1, 2,

		2, 1,
		2, 2,
		1, 1,
		1, 2,

		2, 1,
		2, 2,
		1, 1,
		1, 2,
	}
	sky.position = []float32{
		-0.5, -0.5, -0.5,
	}
	sky.texture = []float32{
		// front
		4, 0,
		4, 1,
		3, 0,
		3, 1,

		// back
		3, 0,
		3, 1,
		2, 0,
		2, 1,

		// left
		4, 1,
		4, 2,
		3, 1,
		3, 2,

		//right
		1, 2,
		1, 3,
		0, 2,
		0, 3,

		// top
		1, 3,
		1, 4,
		0, 3,
		0, 4,

		// bottom
		3, 1,
		3, 2,
		2, 1,
		2, 2,
	}
}
