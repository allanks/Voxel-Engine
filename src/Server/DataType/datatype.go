package DataType

import "gopkg.in/mgo.v2/bson"

type Chunk struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	XPos, ZPos, Map int
}

type CubeChunk struct {
	XPos, ZPos int
	Cubes      []float32
}
