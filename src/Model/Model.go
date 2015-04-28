package Model

import (
	"fmt"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/allanks/Voxel-Engine/src/ObjectLoader"
	"github.com/go-gl/glow/gl-core/4.5/gl"
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

func InitModels() {
	var err error
	TestMob.xPos = 3
	TestMob.yPos = 62
	TestMob.zPos = 5
	models[Cube].vertices, models[Cube].normals, models[Cube].uv = ObjectLoader.LoadObjFile("cube/cube.obj")
	models[Cube].scale = 1.00
	models[Cube].ssbo = getTextureBuffer()
	models[Cube].texture, err = Graphics.NewTexture(textureAtlas)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Cube vertices %v\n", len(models[Cube].vertices))
	fmt.Printf("Cube normals %v\n", len(models[Cube].normals))
	fmt.Printf("Cube uvs %v\n", len(models[Cube].uv))

	models[Gopher].vertices, models[Gopher].normals, models[Gopher].uv = ObjectLoader.LoadObjFile("gopher-3d-master/gopher.obj")
	models[Gopher].scale = 1.0
	models[Gopher].ssbo = []float32{0}
	models[Gopher].texture, err = Graphics.NewTexture(textureGopher)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Gopher vertices %v\n", len(models[Gopher].vertices))
	fmt.Printf("Gopher normals %v\n", len(models[Gopher].normals))
	fmt.Printf("Gopher uvs %v\n", len(models[Gopher].uv))
}

func Render(instanceBuffer uint32, instances []float32, modelType int) {

	gl.BindBuffer(gl.ARRAY_BUFFER, instanceBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(instances)*4, gl.Ptr(instances), gl.STATIC_DRAW)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(len(models[modelType].vertices)/3), int32(len(instances)/4))

}

func BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock, uvBuffer uint32, scale, length int32, modelType int) {
	vertices := models[modelType].vertices
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	normals := models[modelType].normals
	gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(normals)*4, gl.Ptr(normals), gl.STATIC_DRAW)

	ssbo := models[modelType].ssbo
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, textureDataStorageBlock)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(ssbo)*4, gl.Ptr(ssbo), gl.STATIC_DRAW)
	gl.BufferSubData(gl.SHADER_STORAGE_BUFFER, 0, len(ssbo)*4, gl.Ptr(ssbo))

	uv := models[modelType].uv
	gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(uv)*4, gl.Ptr(uv), gl.STATIC_DRAW)

	gl.Uniform1f(scale, models[modelType].scale)
	gl.Uniform1f(length, float32(len(vertices)/3))

	gl.ActiveTexture(gl.TEXTURE0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, models[modelType].texture)
}
