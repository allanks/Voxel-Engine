package Terrain

import (
	"math/rand"
)

var (
	cubes []*Cube
)

func GenLevel(xPos, yPos, zPos int32) {

	cubes = append(cubes, GenCube(xPos, yPos-1, zPos))
	genPaths(cubes[0], 10)
}

func genPaths(cube *Cube, pathLength int32) {

	if pathLength == 0 {
		return
	}

	genPath(cube, int32(rand.Intn(6)), pathLength)
	genPath(cube, int32(rand.Intn(6)), pathLength)

}

func genPath(cube *Cube, decision, pathLength int32) {

	xPos, yPos, zPos := cube.GetPos()
	switch decision / 2 {
	case 0:
		xPos += (-1 + (decision % 2)) + (decision % 2)
	case 1:
		yPos -= (decision % 2)
	case 2:
		zPos += (-1 + (decision % 2)) + (decision % 2)
	}
	if !checkCubeCollisions(xPos, yPos, zPos) {
		newCube := GenCube(xPos, yPos, zPos)
		cubes = append(cubes, newCube)
		genPaths(newCube, pathLength-1)
	}
}

func checkCubeCollisions(xPos, yPos, zPos int32) bool {

	for _, cube := range cubes {
		if cube.CheckCubeCollision(xPos, yPos, zPos) {
			return true
		}
	}
	return false
}

func RenderLevel(vertAttrib, texCoordAttrib uint32, translateUniform int32) {

	for _, cube := range cubes {
		Render(cube, vertAttrib, texCoordAttrib, translateUniform)
	}
}
