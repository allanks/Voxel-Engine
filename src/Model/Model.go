package Model

import (
	"github.com/allanks/Voxel-Engine/src/ObjectLoader"
	"github.com/go-gl/glow/gl-core/4.5/gl"
)

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
	vertices, normals, textures []float32
	scale                       float32
}

func InitModels() {
	TestMob.xPos = 3
	TestMob.yPos = 62
	TestMob.zPos = 5
	models[Cube].vertices, models[Cube].normals, models[Cube].textures = ObjectLoader.LoadObjFile("cube/cube.obj")
	models[Cube].scale = 1.00
	models[Cube].textures = getTextureBuffer()
	models[Gopher].vertices, models[Gopher].normals, models[Gopher].textures = ObjectLoader.LoadObjFile("gopher-3d-master/gopher.obj")
	models[Gopher].scale = 0.05
}

func Render(typeBuffer uint32, instances []float32, modelType int) {
	gl.BindBuffer(gl.ARRAY_BUFFER, typeBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(instances)*4, gl.Ptr(instances), gl.STATIC_DRAW)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(len(models[modelType].vertices)/3), int32(len(instances)/4))

}

func BindBuffers(vertexBuffer, normalBuffer, textureDataStorageBlock uint32, scale, texParam int32, modelType int) {
	vertices := models[modelType].vertices
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	normals := models[modelType].normals
	gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(normals)*4, gl.Ptr(normals), gl.STATIC_DRAW)

	textures := models[modelType].textures
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, textureDataStorageBlock)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(textures)*4, gl.Ptr(textures), gl.STATIC_DRAW)

	gl.Uniform1f(scale, models[modelType].scale)
	gl.Uniform1f(texParam, float32(len(vertices)/3))
}
