package Player

import (
	"fmt"
	"log"
	m "math"
	"os"
	"sync"
	"time"

	"github.com/allanks/Voxel-Engine/src/Terrain"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/glow/gl-core/4.5/gl"
	"github.com/go-gl/mathgl/mgl32"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	user          player
	lastFrameTime float64
	mongoSession  *mgo.Session
)

const (
	forward,
	backwards,
	turnSpeed,
	moveSpeed,
	stop,
	Height,
	terminalVelocity,
	jumpSpeed,
	gravity,
	collisionDistance float64 = 1, -1, 0.5, 0.1, 0, 1, -10, 2, 9.8, 0.15
	mongodb string = "localhost:27017"
)

type player struct {
	xPos, yPos, zPos, pitch, turn, fall float64
	freeMovement, isFalling             bool
	gameMap                             level
	logFile                             *os.File
}

type level struct {
	chunks []*Chunk
}

type moveFunc func(float64)

func GenPlayer(xPos, yPos, zPos float64) {
	lastFrameTime = glfw.GetTime()
	createDatabaseLink()
	user = player{xPos, yPos, zPos, -180.0, 0.0, 0.0, false, true, level{}, &os.File{}}
	user.loadGameMap(mongoSession)
}

func MovePlayer(window *glfw.Window) {
	frameTime := glfw.GetTime()
	frameRate := frameTime - lastFrameTime
	lastFrameTime = frameTime
	if window.GetKey(glfw.KeyW) == glfw.Press {
		move(1)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		move(-1)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		strafe(1)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		strafe(-1)
	}
	switch {
	case window.GetKey(glfw.KeySpace) == glfw.Press && user.freeMovement:
		if !user.gameMap.IsInCube(user.xPos, user.yPos+(moveSpeed)-1, user.zPos, collisionDistance) {
			user.yPos = user.yPos + (moveSpeed)
		}
	case window.GetKey(glfw.KeyLeftShift) == glfw.Press && user.freeMovement:
		if !user.gameMap.IsInCube(user.xPos, user.yPos+(-1*moveSpeed)-1, user.zPos, collisionDistance) {
			user.yPos = user.yPos + (-1 * moveSpeed)
		}
	case !user.freeMovement:
		if user.gameMap.IsInCube(user.xPos, user.yPos+(user.fall*moveSpeed)-1, user.zPos, collisionDistance) {
			user.isFalling = false
			user.fall = 0.0
			user.yPos = m.Floor(user.yPos)
		} else {
			user.isFalling = true
		}
		if user.isFalling {
			user.yPos = user.yPos + (user.fall * moveSpeed)
			if user.fall > terminalVelocity {
				user.fall = user.fall - (gravity * frameRate)
			}
		}
	}
}

func GetPosition() (float64, float64, float64) {
	return user.xPos, user.yPos, user.zPos
}

func GetPlayerSpeed() float64 {
	return moveSpeed
}

func GetCameraMatrix() mgl32.Mat4 {
	xLook := float64(m.Sin(float64(user.pitch)*m.Pi/180) * m.Cos(float64(user.turn)*m.Pi/180))
	zLook := float64(m.Sin(float64(user.pitch)*m.Pi/180) * m.Sin(float64(user.turn)*m.Pi/180))
	yLook := -1 * float64(m.Cos(float64(-1*user.pitch)*m.Pi/180))
	return mgl32.LookAtV(
		mgl32.Vec3{float32(user.xPos), float32(user.yPos), float32(user.zPos)},
		mgl32.Vec3{float32(user.xPos + xLook), float32(user.yPos + yLook), float32(user.zPos + zLook)},
		mgl32.Vec3{0, 1, 0})
}

func move(direction float64) {
	var xLook, zLook float64
	if user.freeMovement {
		xLook = -1 * float64(m.Sin(float64(user.pitch)*m.Pi/180)*m.Cos(float64(user.turn)*m.Pi/180))
		zLook = -1 * float64(m.Sin(float64(user.pitch)*m.Pi/180)*m.Sin(float64(user.turn)*m.Pi/180))
	} else {
		xLook = float64(m.Cos(float64(user.turn) * m.Pi / 180))
		zLook = float64(m.Sin(float64(user.turn) * m.Pi / 180))
	}
	yLook := -1 * float64(m.Cos(float64(-1*user.pitch)*m.Pi/180))
	newX := user.xPos - (direction * xLook * moveSpeed)
	newY := user.yPos + (direction * yLook * moveSpeed)
	newZ := user.zPos - (direction * zLook * moveSpeed)
	if !user.gameMap.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!user.gameMap.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.xPos = newX
	}
	if user.freeMovement &&
		!user.gameMap.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!user.gameMap.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.yPos = newY
	}
	if !user.gameMap.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!user.gameMap.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.zPos = newZ
	}
}

func strafe(direction float64) {
	xLook := float64(m.Cos(float64(user.turn) * m.Pi / 180))
	zLook := float64(m.Sin(float64(user.turn) * m.Pi / 180))
	newX := user.xPos + (-1 * direction * zLook * moveSpeed)
	newZ := user.zPos + (direction * xLook * moveSpeed)

	if !user.gameMap.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!user.gameMap.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.xPos = newX
	}
	if !user.gameMap.IsInCube(user.xPos, user.yPos, user.zPos, collisionDistance) &&
		!user.gameMap.IsInCube(user.xPos, user.yPos+Height, user.zPos, collisionDistance) {
		user.zPos = newZ
	}
}

func OnCursor(window *glfw.Window, xPos, yPos float64) {
	if yPos > -5 {
		window.SetCursorPos(xPos, -5)
		yPos = -5
	}
	if yPos < -359 {
		window.SetCursorPos(xPos, -359)
		yPos = -359
	}
	user.turn = float64(int32(float64(xPos)*turnSpeed) % 360)
	user.pitch = float64(int32(float64(yPos)*turnSpeed) % 360)
}

func OnKey(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyRightShift:
		if action == glfw.Press {
			user.freeMovement = !user.freeMovement
		}
	case glfw.KeySpace:
		if action == glfw.Press && !user.freeMovement && user.yPos == m.Floor(user.yPos) {
			user.fall = jumpSpeed
		}
	case glfw.KeyP:
		fmt.Printf("%v%v\n", "Player ", user)
	}
}

func (gameMap *level) IsInCube(xPos, yPos, zPos, collisionDistance float64) bool {
	pX := int32(m.Floor(xPos))
	pY := int32(m.Floor(yPos))
	pZ := int32(m.Floor(zPos))
	pXpC := int32(m.Floor(xPos + collisionDistance))
	pXmC := int32(m.Floor(xPos - collisionDistance))
	pZpC := int32(m.Floor(zPos + collisionDistance))
	pZmC := int32(m.Floor(zPos - collisionDistance))
	for _, c := range gameMap.chunks {
		if c == nil {
			continue
		}
		if pX >= (c.XPos-1)*int32(chunkSize) && pX <= (c.XPos+1)*int32(chunkSize) && pZ >= (c.ZPos-1)*int32(chunkSize) && pZ <= (c.ZPos+1)*int32(chunkSize) {
			x := c.XPos * int32(chunkSize)
			z := c.ZPos * int32(chunkSize)
			for _, cube := range c.cubes {
				cX := int32(cube.x) + x
				cY := int32(cube.y)
				cZ := int32(cube.z) + z
				if cY == pY &&
					(cX == pX || cX == pXpC || cX == pXmC) &&
					(cZ == pZ || cZ == pZpC || cZ == pZmC) {
					return true
				}
			}
		}
	}
	return false
}

func (p *player) loadGameMap(mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	x := int32(m.Floor(p.xPos / float64(chunkSize)))
	z := int32(m.Floor(p.zPos / float64(chunkSize)))
	p.loadNewChunk(x, z, collection, mongoSession)
	go p.loopChunkLoader(mongoSession)
}

func (p *player) loopChunkLoader(mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	for {
		x := int32(m.Floor(p.xPos / float64(chunkSize)))
		z := int32(m.Floor(p.zPos / float64(chunkSize)))
		p.removeOldChunks(x, z)
		if !checkForChunk(x, z, p.gameMap.chunks) {
			p.loadNewChunk(x, z, collection, mongoSession)
		}
		if !checkForChunk(x+1, z, p.gameMap.chunks) {
			p.loadNewChunk(x+1, z, collection, mongoSession)
		}
		if !checkForChunk(x-1, z, p.gameMap.chunks) {
			p.loadNewChunk(x-1, z, collection, mongoSession)
		}
		if !checkForChunk(x, z+1, p.gameMap.chunks) {
			p.loadNewChunk(x, z+1, collection, mongoSession)
		}
		if !checkForChunk(x, z-1, p.gameMap.chunks) {
			p.loadNewChunk(x, z-1, collection, mongoSession)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (p *player) removeOldChunks(x, z int32) {
	for i, ch := range p.gameMap.chunks {
		if ch == nil {
			continue
		}
		if ch.XPos == (x-2) || ch.ZPos == (z-2) || ch.XPos == (x+2) || ch.ZPos == (z+2) {
			copy(p.gameMap.chunks[i:], p.gameMap.chunks[i+1:])
			p.gameMap.chunks[len(p.gameMap.chunks)-1] = nil
			p.gameMap.chunks = p.gameMap.chunks[:len(p.gameMap.chunks)-1]
			p.removeOldChunks(x, z)
			break
		}
	}
}

func (p *player) loadNewChunk(x, z int32, collection *mgo.Collection, mongoSession *mgo.Session) {
	ch := Chunk{}
	err := collection.Find(bson.M{"xpos": x, "zpos": z}).One(&ch)
	if err != nil {
		fmt.Printf("Creating Chunk at X %v Z %v\n", x, z)
		genChunk(x, z, mongoSession)
		mongoSession.Fsync(false)
		err = collection.Find(bson.M{"xpos": x, "zpos": z}).One(&ch)
		if err != nil {
			log.Fatalf("Could not create chunk: %s\n", err)
		}
	}
	ch.loadChunk(mongoSession)
	p.gameMap.chunks = append(p.gameMap.chunks, &ch)

}

func checkForChunk(x, z int32, chunks []*Chunk) bool {
	for _, ch := range chunks {
		if ch == nil {
			continue
		}
		if ch.XPos == x && ch.ZPos == z && ch.fullyLoaded {
			return true
		}
	}
	return false
}

func createDatabaseLink() {
	var err error
	user.logFile, err = os.OpenFile("ErrorLog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(user.logFile)
	mongoSession, err = mgo.Dial(mongodb)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
}
func closeMongoSession() {
	user.logFile.Close()
	mongoSession.Close()
}

type Chunk struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	XPos, ZPos  int32
	fullyLoaded bool
	cubes       []chunkCube
}

type chunkCube struct {
	x, y, z, cubeType uint8
}

func (c *Chunk) Update(cubes []Terrain.Cube) {
	for _, cube := range cubes {
		c.cubes = append(c.cubes, chunkCube{cube.XPos, cube.YPos, cube.ZPos, cube.GetCubeType()})
	}
}

func (c *Chunk) loadChunk(mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	cubes := []Terrain.Cube{}
	collection.Find(bson.M{"chunkid": c.ID}).All(&cubes)
	c.Update(cubes)
	c.fullyLoaded = true
}

func (c *Chunk) getCubeArray(cubeType uint8) []float32 {
	positions := []float32{}
	for _, cube := range c.cubes {
		if cube.cubeType == cubeType {
			positions = append(positions, float32(int32(cube.x)+(c.XPos*int32(chunkSize))), float32(cube.y), float32(int32(cube.z)+(c.ZPos*int32(chunkSize))))
		}
	}
	return positions
}

const (
	chunkSize uint8 = 64
	seaLevel  uint8 = 64
)

func genChunk(x, z int32, mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	c := &Chunk{XPos: x, ZPos: z}
	err := collection.Insert(c)
	mongoSession.Fsync(false)
	collection.Find(c).One(c)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
	collection = session.DB("GameDatabase").C("Cubes")

	var cubeType uint8
	var waitGroup sync.WaitGroup
	waitGroup.Add(int(seaLevel))
	for yPos := uint8(0); yPos < seaLevel; yPos++ {
		if (yPos + 1) >= seaLevel {
			cubeType = Terrain.Grass
		} else if (yPos + 5) >= seaLevel {
			cubeType = Terrain.Dirt
		} else if (yPos + 10) >= seaLevel {
			cubeType = Terrain.Gravel
		} else {
			cubeType = Terrain.Stone
		}
		go genLayer(c, yPos, cubeType, collection, &waitGroup)
	}
	waitGroup.Wait()
}

func genLayer(c *Chunk, y, cubeType uint8, collection *mgo.Collection, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	bulk := collection.Bulk()
	for x := uint8(0); x < chunkSize; x++ {
		for z := uint8(0); z < chunkSize; z++ {
			bulk.Insert(Terrain.Cube{ChunkID: c.ID, XPos: x, YPos: y, ZPos: z, CubeType: cubeType})
		}
	}
	_, err := bulk.Run()
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
}

func Render(vao, positionBuffer, textureBuffer uint32) {
	user.gameMap.RenderLevel(vao, positionBuffer, textureBuffer)
}

func (gameMap *level) RenderLevel(vao, positionBuffer, textureBuffer uint32) {

	gl.BindVertexArray(vao)

	for _, c := range gameMap.chunks {
		if c == nil || len(c.cubes) == 0 {
			continue
		}
		for _, gCube := range Terrain.GCubes {
			positions := c.getCubeArray(gCube.Gtype)

			if len(positions) == 0 {
				continue
			}

			gl.BindBuffer(gl.ARRAY_BUFFER, textureBuffer)
			gl.BufferData(gl.ARRAY_BUFFER, len(gCube.Texture)*4, gl.Ptr(gCube.Texture), gl.STATIC_DRAW)

			gl.BindBuffer(gl.ARRAY_BUFFER, positionBuffer)
			gl.BufferData(gl.ARRAY_BUFFER, len(positions)*4, gl.Ptr(positions), gl.STATIC_DRAW)

			instances := int32(len(positions) / 3)
			gl.DrawElementsInstanced(gl.TRIANGLES, 36, gl.UNSIGNED_INT, gl.Ptr(nil), int32(instances))
		}
	}

}
