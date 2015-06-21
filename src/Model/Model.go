package Model

import (
	"fmt"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/allanks/Voxel-Engine/src/ObjectLoader"
)

const textureAtlas string = "resource/texture/textureAtlas.png"
const textureGopher string = "resource/texture/Gopher/gopher.png"
const (
	Cube = iota
	Gopher
)

var TestMob Mob

var models = []model{
	model{},
	model{},
}

type Mob struct {
	xPos, yPos, zPos float64
	mType            uint
}

type model struct {
	vertices, normals, uv, ssbo []float32
	scale                       float32
	texture                     uint32
}

var Controller Graphics.OpenGLController
var Control Graphics.OpenGLControl

func InitModels() {
	TestMob.xPos = 3
	TestMob.yPos = 62
	TestMob.zPos = 5
	models[Cube].vertices, models[Cube].normals, models[Cube].uv = ObjectLoader.LoadObjFile("cube/cube.obj")
	models[Cube].scale = 1.00
	models[Cube].ssbo = getTextureBuffer()
	models[Cube].texture = Control.CreateTexture(textureAtlas)

	fmt.Printf("Cube vertices %v\n", len(models[Cube].vertices))
	fmt.Printf("Cube normals %v\n", len(models[Cube].normals))
	fmt.Printf("Cube uvs %v\n", len(models[Cube].uv))

	models[Gopher].vertices, models[Gopher].normals, models[Gopher].uv = ObjectLoader.LoadObjFile("gopher-3d-master/gopher.obj")
	models[Gopher].scale = 1.0
	models[Gopher].ssbo = []float32{0}
	models[Gopher].texture = Control.CreateTexture(textureGopher)

	fmt.Printf("Gopher vertices %v\n", len(models[Gopher].vertices))
	fmt.Printf("Gopher normals %v\n", len(models[Gopher].normals))
	fmt.Printf("Gopher uvs %v\n", len(models[Gopher].uv))
}

func Render(instances []float32, modelType int) {
	Controller.RenderInstances(instances, int32(len(models[modelType].vertices)/3))
}

func BindBuffers(offset []float32, modelType int) {
	Controller.BindBuffers(models[modelType].vertices, models[modelType].normals, models[modelType].ssbo, models[modelType].uv)
	Controller.BindUniforms([]float32{models[modelType].scale}, []float32{float32(len(models[modelType].vertices) / 3)}, offset)
	Controller.BindTexture(models[modelType].texture)
}
