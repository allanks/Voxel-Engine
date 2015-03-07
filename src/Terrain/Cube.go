package Terrain

import (
	"gopkg.in/mgo.v2/bson"
)

type Cube struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	XPos int32
	YPos int32
	ZPos int32
}

func (cube *Cube) GetPos() (int32, int32, int32) {
	return cube.XPos, cube.YPos, cube.ZPos
}

func (cube *Cube) setPos(xPos, yPos, zPos int32) {
	cube.XPos = xPos
	cube.YPos = yPos
	cube.ZPos = zPos
}
