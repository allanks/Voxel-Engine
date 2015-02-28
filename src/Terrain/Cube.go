package Terrain

type Cube struct {
	xPos, yPos, zPos int32
}

func (cube *Cube) GetPos() (int32, int32, int32) {
	return cube.xPos, cube.yPos, cube.zPos
}

func (cube *Cube) setPos(xPos, yPos, zPos int32) {
	cube.xPos = xPos
	cube.yPos = yPos
	cube.zPos = zPos
}

func (cube *Cube) CheckCubeCollision(xPos, yPos, zPos int32) bool {
	return xPos == cube.xPos && yPos == cube.yPos && zPos == cube.zPos
}

func GenCube(x, y, z int32) *Cube {
	return &Cube{x, y, z}
}
