package Terrain

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/go-gl/glow/gl-core/4.5/gl"
	"gopkg.in/mgo.v2/bson"
)

const (
	collisionDistance float64 = 0.15
)
const (
	// Cube Types
	Empty = iota
	SkyBox
	Dirt
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

var sky skyBox

var GCubes = []GCube{
	GCube{},
	GCube{},
	GCube{},
	GCube{},
	GCube{},
	GCube{},
	GCube{},
}

var instances int32

func BindCubeVertexBuffers(vertexBuffer, indexBuffer uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(CubeElements)*4, gl.Ptr(CubeElements), gl.STATIC_DRAW)

}

func GetTextureBuffer() []float32 {

	buffer := []float32{}
	for i := 0; i < 48; i++ {
		buffer = append(buffer, 0)
	}

	buffer = append(buffer, sky.texture...)

	for _, gCube := range GCubes {
		if gCube.Gtype == 0 {
			continue
		}
		buffer = append(buffer, gCube.Texture...)
	}
	return buffer
}

func RenderSkyBox(vao, typeBuffer uint32, offset int32, x, y, z float64) {
	gl.DepthMask(false)

	gl.BindVertexArray(vao)

	position := []float32{float32(x), float32(y), float32(z), 1}
	gl.BindBuffer(gl.ARRAY_BUFFER, typeBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 16, gl.Ptr(position), gl.STATIC_DRAW)

	gl.Uniform3f(offset, -0.5, -0.5, -0.5)

	gl.DrawElementsInstanced(gl.TRIANGLES, 36, gl.UNSIGNED_INT, gl.Ptr(nil), int32(1))

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
	1.0, 1.0, 0.0,
	1.0, 0.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 0.0,
	// Left face
	0.0, 1.0, 1.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 0.0,
	// Right face
	1.0, 1.0, 1.0,
	1.0, 1.0, 0.0,
	1.0, 0.0, 1.0,
	1.0, 0.0, 0.0,
	// Top face
	1.0, 1.0, 1.0,
	1.0, 1.0, 0.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 0.0,
	// Bottom face
	1.0, 0.0, 1.0,
	1.0, 0.0, 0.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 0.0,
}

var CubeElements = []uint32{
	// front
	0, 1, 2,
	1, 3, 2,
	// back
	4, 5, 6,
	5, 7, 6,
	// Left
	8, 9, 10,
	9, 11, 10,
	// Right
	12, 13, 14,
	13, 15, 14,
	// Top
	16, 17, 18,
	17, 19, 18,
	// Bottom
	20, 21, 22,
	21, 23, 22,
}

func InitGCubes() {
	sky.texture = []float32{
		// Front and back are swapped as we are view from inside the cube

		// front
		1, 1,
		1, 2,
		2, 1,
		2, 2,

		//back
		1, 1,
		1, 2,
		0, 1,
		0, 2,

		// left
		0, 2,
		1, 2,
		0, 3,
		1, 3,

		// right
		2, 2,
		1, 2,
		2, 3,
		1, 3,

		// top
		0, 1,
		0, 0,
		1, 1,
		1, 0,

		// bottom
		1, 0,
		1, 1,
		2, 0,
		2, 1,
	}

	GCubes[Dirt].Gtype = Dirt
	GCubes[Grass].Gtype = Grass
	GCubes[Stone].Gtype = Stone
	GCubes[CobbleStone].Gtype = CobbleStone
	GCubes[Gravel].Gtype = Gravel

	GCubes[Dirt].Texture = loadTexCoords(2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1)
	GCubes[Grass].Texture = loadTexCoords(2, 1, 2, 1, 2, 1, 2, 1, 2, 2, 2, 1)
	GCubes[Stone].Texture = loadTexCoords(1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3)
	GCubes[CobbleStone].Texture = loadTexCoords(2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0)
	GCubes[Gravel].Texture = loadTexCoords(0, 3, 0, 3, 0, 3, 0, 3, 0, 3, 0, 3)
}

func loadTexCoords(u1, v1, u2, v2, u3, v3, u4, v4, u5, v5, u6, v6 float32) []float32 {
	tex := []float32{
		u1, v1,
		u1, v1 + 1,
		u1 + 1, v1,
		u1 + 1, v1 + 1,
		u2, v2,
		u2, v2 + 1,
		u2 + 1, v2,
		u2 + 1, v2 + 1,
		u3, v3,
		u3, v3 + 1,
		u3 + 1, v3,
		u3 + 1, v3 + 1,
		u4, v4,
		u4, v4 + 1,
		u4 + 1, v4,
		u4 + 1, v4 + 1,
		u5, v5,
		u5, v5 + 1,
		u5 + 1, v5,
		u5 + 1, v5 + 1,
		u6, v6,
		u6, v6 + 1,
		u6 + 1, v6,
		u6 + 1, v6 + 1,
	}
	return tex
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
