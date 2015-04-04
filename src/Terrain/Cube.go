package Terrain

import (
	"image"
	"image/draw"
	"image/png"
	"os"

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
	texture []float32
}

type GCube struct {
	Texture []float32
	Gtype   uint8
}

type Cube struct {
	ID                         bson.ObjectId `bson:"_id,omitempty"`
	ChunkID                    bson.ObjectId
	XPos, YPos, ZPos, CubeType uint8
}

func (cube *Cube) GetCubeType() uint8 {
	return cube.CubeType
}

func (cube *Cube) CheckCollision(xPos, yPos, zPos uint8) bool {

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

func RenderSkyBox(vao, positionBuffer, textureBuffer uint32, x, y, z float64) {
	position := []float32{float32(x) - 0.5, float32(y) - 0.5, float32(z) - 0.5}
	gl.DepthMask(false)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, textureBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(sky.texture)*4, gl.Ptr(sky.texture), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(position)*4, gl.Ptr(position), gl.STATIC_DRAW)

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
	1.0, 0.0, 0.0,
	0.0, 0.0, 0.0,
	1.0, 1.0, 0.0,
	0.0, 1.0, 0.0,
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
	sky.texture = []float32{
		// Front and back are swapped as we are view from inside the cube

		// front
		2, 2,
		2, 3,
		1, 2,
		1, 3,

		// back
		1, 2,
		1, 3,
		0, 2,
		0, 3,

		// left
		1, 1,
		1, 2,
		0, 1,
		0, 2,

		//right
		2, 1,
		2, 2,
		1, 1,
		1, 2,

		// top
		1, 0,
		1, 1,
		0, 0,
		0, 1,

		// bottom
		2, 0,
		2, 1,
		1, 0,
		1, 1,
	}

	GCubes[Dirt].Gtype = Dirt
	GCubes[Grass].Gtype = Grass
	GCubes[Stone].Gtype = Stone
	GCubes[CobbleStone].Gtype = CobbleStone
	GCubes[Gravel].Gtype = Gravel

	GCubes[Dirt].Texture = []float32{
		3, 1,
		3, 2,
		2, 1,
		2, 2,

		3, 1,
		3, 2,
		2, 1,
		2, 2,

		3, 1,
		3, 2,
		2, 1,
		2, 2,

		3, 1,
		3, 2,
		2, 1,
		2, 2,

		3, 1,
		3, 2,
		2, 1,
		2, 2,

		3, 1,
		3, 2,
		2, 1,
		2, 2,
	}
	GCubes[Grass].Texture = []float32{
		3, 2,
		3, 3,
		2, 2,
		2, 3,

		3, 2,
		3, 3,
		2, 2,
		2, 3,

		3, 2,
		3, 3,
		2, 2,
		2, 3,

		3, 2,
		3, 3,
		2, 2,
		2, 3,

		3, 2,
		3, 3,
		2, 2,
		2, 3,

		3, 2,
		3, 3,
		2, 2,
		2, 3,
	}
	GCubes[Stone].Texture = []float32{
		2, 3,
		2, 4,
		1, 3,
		1, 4,

		2, 3,
		2, 4,
		1, 3,
		1, 4,

		2, 3,
		2, 4,
		1, 3,
		1, 4,

		2, 3,
		2, 4,
		1, 3,
		1, 4,

		2, 3,
		2, 4,
		1, 3,
		1, 4,

		2, 3,
		2, 4,
		1, 3,
		1, 4,
	}
	GCubes[CobbleStone].Texture = []float32{
		3, 0,
		3, 1,
		2, 0,
		2, 1,

		3, 0,
		3, 1,
		2, 0,
		2, 1,

		3, 0,
		3, 1,
		2, 0,
		2, 1,

		3, 0,
		3, 1,
		2, 0,
		2, 1,

		3, 0,
		3, 1,
		2, 0,
		2, 1,

		3, 0,
		3, 1,
		2, 0,
		2, 1,
	}
	GCubes[Gravel].Texture = []float32{
		1, 3,
		1, 4,
		0, 3,
		0, 4,

		1, 3,
		1, 4,
		0, 3,
		0, 4,

		1, 3,
		1, 4,
		0, 3,
		0, 4,

		1, 3,
		1, 4,
		0, 3,
		0, 4,

		1, 3,
		1, 4,
		0, 3,
		0, 4,

		1, 3,
		1, 4,
		0, 3,
		0, 4,
	}
}

func PackTextures() {
	textureAtlas := image.NewRGBA(image.Rect(0, 0, 2048, 2048))

	mr := image.Rectangle{image.Point{0, 0}, image.Point{512, 512}}
	r := image.Rectangle{image.Point{0, 0}, image.Point{512, 512}}
	img := loadPNG("SkyBox/top.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 0}, image.Point{1024, 512}}
	img = loadPNG("SkyBox/bottom.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{0, 512}, image.Point{512, 1024}}
	img = loadPNG("SkyBox/left.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 512}, image.Point{1024, 1024}}
	img = loadPNG("SkyBox/right.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{0, 1024}, image.Point{512, 1536}}
	img = loadPNG("SkyBox/front.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 1024}, image.Point{1024, 1536}}
	img = loadPNG("SkyBox/back.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{1024, 0}, image.Point{1536, 512}}
	img = loadPNG("CobbleStone/cobblestone.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{1024, 512}, image.Point{1536, 1024}}
	img = loadPNG("Dirt/dirt.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{1024, 1024}, image.Point{1536, 1536}}
	img = loadPNG("Grass/grass.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{0, 1536}, image.Point{512, 2048}}
	img = loadPNG("Gravel/gravel.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 1536}, image.Point{1024, 2048}}
	img = loadPNG("Stone/stone.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	pngFile, err := os.Create("resource/texture/textureAtlas.png")
	if err != nil {
		panic(err)
	}
	png.Encode(pngFile, textureAtlas)
}

func loadPNG(filePath string) image.Image {

	var pngFile *os.File
	var err error
	var pngImg image.Image
	pngFile, err = os.Open("resource/texture/" + filePath)
	if err != nil {
		panic(err)
	}
	defer pngFile.Close()
	pngImg, err = png.Decode(pngFile)
	if err != nil {
		panic(err)
	}
	return pngImg
}
