package Terrain

import (
	"fmt"

	"github.com/allanks/third-game/src/Graphics"
	"github.com/go-gl/glow/gl-core/4.5/gl"
	"gopkg.in/mgo.v2/bson"
)

const (
	collisionDistance float64 = 0.15
	textureDir        string  = "resource/texture/"
)
const (
	// Cube Types
	Dirt = iota
	Grass
	Stone
	CobbleStone
	Gravel
)

type Cube struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ChunkID  bson.ObjectId `bson:"omitempty"`
	XPos     float64
	YPos     float64
	ZPos     float64
	CubeType int
}

func (cube *Cube) GetCubeType() int {
	return cube.CubeType
}

func (cube *Cube) GetPos() (float64, float64, float64) {
	return cube.XPos, cube.YPos, cube.ZPos
}

func (cube *Cube) setPos(xPos, yPos, zPos float64) {
	cube.XPos = xPos
	cube.YPos = yPos
	cube.ZPos = zPos
}

func (cube *Cube) PrintCollision(xPos, yPos, zPos float64) {
	fmt.Printf("%v %v %v\n", xPos-cube.XPos, yPos-cube.YPos, zPos-cube.ZPos)
}

func (cube *Cube) CheckCollision(xPos, yPos, zPos, moveSpeed float64) bool {

	if xPos == cube.XPos &&
		yPos == cube.YPos &&
		zPos == cube.ZPos {
		return true
	}
	return false
}

var gCubes = []GCube{
	GCube{},
	GCube{},
	GCube{},
	GCube{},
	GCube{},
}

var (
	vao, vertexBuffer, textureBuffer uint32
)

var instances int32

func InitGCubes() {
	gCubes[Dirt].backColor = []float32{1.0, 1.0, 0.0}
	gCubes[Dirt].frontColor = []float32{1.0, 1.0, 0.0}
	gCubes[Dirt].leftColor = []float32{1.0, 1.0, 0.0}
	gCubes[Dirt].rightColor = []float32{1.0, 1.0, 0.0}
	gCubes[Dirt].topColor = []float32{1.0, 1.0, 0.0}
	gCubes[Dirt].bottomColor = []float32{1.0, 1.0, 0.0}

	gCubes[Grass].backColor = []float32{0.0, 0.0, 1.0}
	gCubes[Grass].frontColor = []float32{0.0, 0.0, 1.0}
	gCubes[Grass].leftColor = []float32{0.0, 0.0, 1.0}
	gCubes[Grass].rightColor = []float32{0.0, 0.0, 1.0}
	gCubes[Grass].topColor = []float32{0.0, 0.0, 1.0}
	gCubes[Grass].bottomColor = []float32{0.0, 0.0, 1.0}

	gCubes[Stone].backColor = []float32{1.0, 0.0, 1.0}
	gCubes[Stone].frontColor = []float32{1.0, 0.0, 1.0}
	gCubes[Stone].leftColor = []float32{1.0, 0.0, 1.0}
	gCubes[Stone].rightColor = []float32{1.0, 0.0, 1.0}
	gCubes[Stone].topColor = []float32{1.0, 0.0, 1.0}
	gCubes[Stone].bottomColor = []float32{1.0, 0.0, 1.0}

	gCubes[CobbleStone].backColor = []float32{1.0, 1.0, 1.0}
	gCubes[CobbleStone].frontColor = []float32{1.0, 1.0, 1.0}
	gCubes[CobbleStone].leftColor = []float32{1.0, 1.0, 1.0}
	gCubes[CobbleStone].rightColor = []float32{1.0, 1.0, 1.0}
	gCubes[CobbleStone].topColor = []float32{1.0, 1.0, 1.0}
	gCubes[CobbleStone].bottomColor = []float32{1.0, 1.0, 1.0}

	gCubes[Gravel].backColor = []float32{1.0, 0.0, 0.0}
	gCubes[Gravel].frontColor = []float32{1.0, 0.0, 0.0}
	gCubes[Gravel].leftColor = []float32{1.0, 0.0, 0.0}
	gCubes[Gravel].rightColor = []float32{1.0, 0.0, 0.0}
	gCubes[Gravel].topColor = []float32{1.0, 0.0, 0.0}
	gCubes[Gravel].bottomColor = []float32{1.0, 0.0, 0.0}
	/*gCubes[Dirt].initCubeTextures(textureDir + "Dirt")
	gCubes[Grass].initCubeTextures(textureDir + "Grass")
	gCubes[Stone].initCubeTextures(textureDir + "Stone")
	gCubes[CobbleStone].initCubeTextures(textureDir + "CobbleStone")
	gCubes[Gravel].initCubeTextures(textureDir + "Gravel")*/
}

type GCube struct {
	topTexture,
	bottomTexture,
	frontTexture,
	backTexture,
	leftTexture,
	rightTexture uint32
	topColor,
	bottomColor,
	rightColor,
	leftColor,
	frontColor,
	backColor []float32
}

func (cube *GCube) getColors() []float32 {
	colors := []float32{}
	colors = append(colors, cube.frontColor...)
	colors = append(colors, cube.backColor...)
	colors = append(colors, cube.leftColor...)
	colors = append(colors, cube.rightColor...)
	colors = append(colors, cube.topColor...)
	colors = append(colors, cube.bottomColor...)
	return colors
}

func (cube *GCube) initCubeTextures(dir string) {
	fmt.Printf("%v\n", "Loading Textures from: "+dir)

	// Load the textures
	var err error
	cube.topTexture, err = Graphics.NewTexture(dir + "/topFace.png")
	if err != nil {
		panic(err)
	}
	cube.bottomTexture, err = Graphics.NewTexture(dir + "/bottomFace.png")
	if err != nil {
		panic(err)
	}
	cube.frontTexture, err = Graphics.NewTexture(dir + "/frontFace.png")
	if err != nil {
		panic(err)
	}
	cube.backTexture, err = Graphics.NewTexture(dir + "/backFace.png")
	if err != nil {
		panic(err)
	}
	cube.leftTexture, err = Graphics.NewTexture(dir + "/rightFace.png")
	if err != nil {
		panic(err)
	}
	cube.rightTexture, err = Graphics.NewTexture(dir + "/leftFace.png")
	if err != nil {
		panic(err)
	}
}

func InitialiseGCubeBuffers() {
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.GenBuffers(1, &textureBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeTexCoords)*4, gl.Ptr(cubeTexCoords), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	positions, colors := []float32{}, []float32{}

	for _, c := range gameMap.chunks {
		positions = append(positions, c.getPositions()...)
		colors = append(colors, c.getColors()...)
	}

	fmt.Printf("%v%v\n", "Position Buffer has entries ", len(positions))
	fmt.Printf("%v%v\n", "Color Buffer has entries ", len(colors))

	var positionBuffer, colorBuffer uint32
	gl.GenBuffers(1, &positionBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
	gl.EnableVertexAttribArray(2)
	gl.BufferData(gl.ARRAY_BUFFER, len(positions)*4, gl.Ptr(positions), gl.STATIC_DRAW)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(2, 6)

	gl.GenBuffers(1, &colorBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
	gl.EnableVertexAttribArray(3)
	gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.STATIC_DRAW)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(3, 1)

	instances = int32(len(colors) / 3)
}

func RenderLevel() {

	for i, gc := range gCubes {
		gc.RenderCubes(i)
	}
}

func (cube *GCube) RenderCubes(i int) {

	gl.BindVertexArray(vao)
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 24, int32(instances))
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

var cubeTexCoords = []float32{
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
	1.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
	// Bottom face
	1.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
}
