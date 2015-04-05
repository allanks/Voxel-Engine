package Terrain

import (
	"fmt"
	"log"
	m "math"
	"os"
	"sync"

	"github.com/go-gl/glow/gl-core/4.5/gl"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongodb    string = "localhost:27017"
	chunkSize  int    = 16
	maxHeight  int    = 128
	seaLevel   uint8  = 64
	stoneLevel uint8  = 60
)

var (
	logFile      *os.File
	mongoSession *mgo.Session
	cubeLoader   []Cube
)

type Level struct {
	chunks       []*Chunk
	seed         int64
	defaultBlock []*Cube
}

type Chunk struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	XPos, ZPos  int
	fullyLoaded bool
	cubes       [chunkSize][maxHeight][chunkSize]uint8
}

func (gameMap *Level) IsInCube(xPos, yPos, zPos, collisionDistance float64) bool {
	pX := int(m.Floor(xPos))
	pY := int(m.Floor(yPos))
	pZ := int(m.Floor(zPos))
	pXpC := ((int(m.Floor(xPos+collisionDistance)) % chunkSize) + chunkSize) % chunkSize
	pXmC := ((int(m.Floor(xPos-collisionDistance)) % chunkSize) + chunkSize) % chunkSize
	pZpC := ((int(m.Floor(zPos+collisionDistance)) % chunkSize) + chunkSize) % chunkSize
	pZmC := ((int(m.Floor(zPos-collisionDistance)) % chunkSize) + chunkSize) % chunkSize
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("X+ %v X- %v Z+ %v Z- %v\n", pXpC, pXmC, pZpC, pZmC)
		}
	}()
	for _, c := range gameMap.chunks {
		if c == nil {
			continue
		}
		if pX <= (c.XPos+1)*chunkSize && pX >= c.XPos*chunkSize && pZ <= (c.ZPos+1)*chunkSize && pZ >= c.ZPos*chunkSize {
			if c.cubes[pXpC][pY][pZpC] != Empty ||
				c.cubes[pXpC][pY][pZmC] != Empty ||
				c.cubes[pXmC][pY][pZpC] != Empty ||
				c.cubes[pXmC][pY][pZmC] != Empty {
				return true
			}
		}
	}
	return false
}

func (gameMap *Level) RenderLevel(vao, positionBuffer, textureBuffer uint32) {

	gl.BindVertexArray(vao)

	for _, c := range gameMap.chunks {
		if c == nil || len(c.cubes) == 0 {
			continue
		}
		for _, gCube := range GCubes {
			if gCube.Gtype == 0 {
				continue
			}
			positions := []float32{}
			for x := 0; x < chunkSize; x++ {
				for z := 0; z < chunkSize; z++ {
					for y := 0; y < maxHeight; y++ {
						if c.cubes[x][y][z] == gCube.Gtype {
							positions = append(positions, float32(x+(c.XPos*chunkSize)), float32(y), float32(z+(c.ZPos*chunkSize)))
						}
					}
				}
			}

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

func (ch *Chunk) genChunk(mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	err := collection.Insert(ch)
	mongoSession.Fsync(false)
	collection.Find(ch).One(ch)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
	collection = session.DB("GameDatabase").C("Cubes")

	var cubeType uint8
	var waitGroup sync.WaitGroup
	waitGroup.Add(int(seaLevel))
	for yPos := uint8(0); yPos < seaLevel; yPos++ {
		if (yPos + 1) >= seaLevel {
			cubeType = Grass
		} else if (yPos + 5) >= seaLevel {
			cubeType = Dirt
		} else if (yPos + 10) >= seaLevel {
			cubeType = Gravel
		} else {
			cubeType = Stone
		}
		go genLayer(ch, yPos, cubeType, collection, &waitGroup)
	}
	waitGroup.Wait()
}

func genLayer(c *Chunk, y, cubeType uint8, collection *mgo.Collection, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	bulk := collection.Bulk()
	for x := 0; x < chunkSize; x++ {
		for z := 0; z < chunkSize; z++ {
			bulk.Insert(Cube{ChunkID: c.ID, XPos: uint8(x), YPos: y, ZPos: uint8(z), CubeType: cubeType})
		}
	}
	_, err := bulk.Run()
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
}

func (gameMap *Level) removeOldChunks(x, z int) {
	for i, ch := range gameMap.chunks {
		if ch == nil {
			continue
		}
		if ch.XPos == (x-3) || ch.ZPos == (z-3) || ch.XPos == (x+3) || ch.ZPos == (z+3) {
			copy(gameMap.chunks[i:], gameMap.chunks[i+1:])
			gameMap.chunks[len(gameMap.chunks)-1] = nil
			gameMap.chunks = gameMap.chunks[:len(gameMap.chunks)-1]
			gameMap.removeOldChunks(x, z)
			break
		}
	}
}

func (ch *Chunk) loadNewChunk(x, z int) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	err := collection.Find(bson.M{"xpos": x, "zpos": z}).One(ch)
	if err != nil {
		fmt.Printf("Creating Chunk at X %v Z %v\n", x, z)
		ch.genChunk(mongoSession)
		mongoSession.Fsync(false)
		err = collection.Find(bson.M{"xpos": x, "zpos": z}).One(ch)
		if err != nil {
			log.Fatalf("Could not create chunk: %s\n", err)
		}
	}
	ch.loadChunk(mongoSession)
}

func (c *Chunk) loadChunk(mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	collection.Find(bson.M{"chunkid": c.ID}).All(&cubeLoader)
	c.Update(cubeLoader)
	c.fullyLoaded = true
}

func (c *Chunk) Update(cubes []Cube) {
	for _, cube := range cubes {
		c.cubes[int(cube.XPos)][int(cube.YPos)][int(cube.ZPos)] = cube.GetCubeType()
	}
}

func (gameMap *Level) checkForChunk(x, z int) bool {
	for _, ch := range gameMap.chunks {
		if ch == nil {
			continue
		}
		if ch.XPos == x && ch.ZPos == z {
			return true
		}
	}
	return false
}

func (gameMap *Level) LoadGameMap(pX, pZ float64) {
	if mongoSession == nil {
		createDatabaseLink()
	}
	x := int(m.Floor(pX / float64(chunkSize)))
	z := int(m.Floor(pZ / float64(chunkSize)))
	ch := Chunk{}
	ch.XPos = x
	ch.ZPos = z
	gameMap.chunks = append(gameMap.chunks, &ch)
	ch.loadNewChunk(x, z)
}

func (gameMap *Level) initChunk(x, z int) {
	ch := Chunk{}
	ch.XPos = x
	ch.ZPos = z
	gameMap.chunks = append(gameMap.chunks, &ch)

	ch.loadNewChunk(x, z)
}

func (gameMap *Level) LoopChunkLoader(pX, pZ float64) {

	x := int(m.Floor(pX / float64(chunkSize)))
	z := int(m.Floor(pZ / float64(chunkSize)))
	gameMap.removeOldChunks(x, z)
	if !gameMap.checkForChunk(x, z) {
		gameMap.initChunk(x, z)
	}
	if !gameMap.checkForChunk(x+1, z) {
		gameMap.initChunk(x+1, z)
	}
	if !gameMap.checkForChunk(x-1, z) {
		gameMap.initChunk(x-1, z)
	}
	if !gameMap.checkForChunk(x, z+1) {
		gameMap.initChunk(x, z+1)
	}
	if !gameMap.checkForChunk(x, z-1) {
		gameMap.initChunk(x, z-1)
	}
	if !gameMap.checkForChunk(x+2, z) {
		gameMap.initChunk(x+2, z)
	}
	if !gameMap.checkForChunk(x-2, z) {
		gameMap.initChunk(x-2, z)
	}
	if !gameMap.checkForChunk(x, z+2) {
		gameMap.initChunk(x, z+2)
	}
	if !gameMap.checkForChunk(x, z-2) {
		gameMap.initChunk(x, z-2)
	}
	if !gameMap.checkForChunk(x+1, z+1) {
		gameMap.initChunk(x+1, z+1)
	}
	if !gameMap.checkForChunk(x+1, z-1) {
		gameMap.initChunk(x+1, z-1)
	}
	if !gameMap.checkForChunk(x-1, z+1) {
		gameMap.initChunk(x-1, z+1)
	}
	if !gameMap.checkForChunk(x-1, z-1) {
		gameMap.initChunk(x-1, z-1)
	}
}

func createDatabaseLink() {
	var err error
	logFile, err = os.OpenFile("ErrorLog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
	mongoSession, err = mgo.Dial(mongodb)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
}
func closeMongoSession() {
	logFile.Close()
	mongoSession.Close()
}
