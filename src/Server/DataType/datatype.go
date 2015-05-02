package DataType

import (
	m "math"

	"gopkg.in/mgo.v2/bson"
)

type Chunk struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	XPos, ZPos, Map int
}

type CubeChunk struct {
	XPos, ZPos int
	Cubes      []float32
}

type Pos struct {
	XPos, YPos, ZPos float32
}

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

func FloorToInt(x float32) int {
	return int(m.Floor(float64(x)))
}
