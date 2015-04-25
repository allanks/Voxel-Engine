package Model

import "github.com/allanks/Voxel-Engine/src/ObjectLoader"

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

func getTextureBuffer() []float32 {

	textureBuffer := []float32{}
	for i := 0; i < 36; i++ {
		textureBuffer = append(textureBuffer, 0, 0)
	}

	textureBuffer = append(textureBuffer, sky.texture...)

	for _, gCube := range GCubes {
		if gCube.Gtype == 0 {
			continue
		}
		textureBuffer = append(textureBuffer, gCube.Texture...)
	}
	return textureBuffer
}

func InitGCubes() {
	_, _, sky.texture = ObjectLoader.LoadObjFile("cube/skybox.obj")
	_, _, GCubes[Dirt].Texture = ObjectLoader.LoadObjFile("cube/dirt.obj")
	_, _, GCubes[Grass].Texture = ObjectLoader.LoadObjFile("cube/grass.obj")
	_, _, GCubes[Stone].Texture = ObjectLoader.LoadObjFile("cube/stone.obj")
	_, _, GCubes[CobbleStone].Texture = ObjectLoader.LoadObjFile("cube/cobblestone.obj")
	_, _, GCubes[Gravel].Texture = ObjectLoader.LoadObjFile("cube/gravel.obj")

	GCubes[Dirt].Gtype = Dirt
	GCubes[Grass].Gtype = Grass
	GCubes[Stone].Gtype = Stone
	GCubes[CobbleStone].Gtype = CobbleStone
	GCubes[Gravel].Gtype = Gravel
}
