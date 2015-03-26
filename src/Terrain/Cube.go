package Terrain

import (
	"fmt"
	"math"

	"gopkg.in/mgo.v2/bson"
)

const (
	collisionDistance float64 = 0.15
)

type Cube struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	XPos float64
	YPos float64
	ZPos float64
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
