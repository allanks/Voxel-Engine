package Terrain

import (
	"fmt"
	"math"

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
)

type Cube struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	XPos     float64
	YPos     float64
	ZPos     float64
	CubeType float64
}

func (cube *Cube) GetCubeType() float64 {
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

func (cube *Cube) CheckCollision(xPos, yPos, zPos, moveSpeed float64) (float64, float64, float64) {

	easyCompare := func(a, b float64) bool {
		return a-b <= 1+(collisionDistance+moveSpeed) && a-b > 0-(collisionDistance+moveSpeed)
	}
	if easyCompare(xPos, cube.XPos) &&
		easyCompare(yPos, cube.YPos) &&
		easyCompare(zPos, cube.ZPos) {
		fmt.Println("Collision Detected")
		selectEdge := func(a, b float64) float64 {
			if a-b > 0.5 {
				return b + 1 + (collisionDistance + moveSpeed)
			} else {
				return b - (collisionDistance + moveSpeed)
			}
		}
		if (math.Abs(xPos-cube.XPos-0.5) > math.Abs(yPos-cube.YPos-0.5)) && (math.Abs(xPos-cube.XPos-0.5) > math.Abs(zPos-cube.ZPos-0.5)) {
			return selectEdge(xPos, cube.XPos), yPos, zPos
		} else if math.Abs(yPos-cube.YPos-0.5) > math.Abs(zPos-cube.ZPos-0.5) {
			return xPos, selectEdge(yPos, cube.YPos), zPos
		} else {
			return xPos, yPos, selectEdge(zPos, cube.ZPos)
		}
	}
	return xPos, yPos, zPos
}

var gCubes = []GCube{
	GCube{},
	GCube{},
}

func InitGCubes() {
	gCubes[Dirt].initialiseCube(textureDir + "Dirt")
	gCubes[Grass].initialiseCube(textureDir + "Grass")
}

type GCube struct {
	topTexture,
	bottomTexture,
	frontTexture,
	backTexture,
	leftTexture,
	rightTexture,
	vertexArrayObject,
	vertexBufferObject uint32
}

func (cube *GCube) initialiseCube(dir string) {
	gl.GenVertexArrays(1, &cube.vertexArrayObject)
	gl.BindVertexArray(cube.vertexArrayObject)

	gl.GenBuffers(1, &cube.vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

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

func (cube *GCube) RenderCube(vertAttrib, texCoordAttrib uint32, translateUniform int32) {
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.BindVertexArray(cube.vertexArrayObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.vertexBufferObject)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cube.frontTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cube.backTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 4, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cube.leftTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 8, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cube.rightTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 12, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cube.topTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 16, 4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cube.bottomTexture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 20, 4)
}

func Render(cube *Cube, vertAttrib, texCoordAttrib uint32, translateUniform int32) {

	xPos, yPos, zPos := cube.GetPos()
	gl.Uniform4f(translateUniform, float32(xPos), float32(yPos), float32(zPos), 0.0)
	gCube := gCubes[uint(cube.GetCubeType())]
	gCube.RenderCube(vertAttrib, texCoordAttrib, translateUniform)

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
	1.0, 0.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 0.0, 0.0,
	0.0, 0.0, 1.0, 0.0, 1.0,
}
