package Terrain

import (
	"fmt"
	"log"
	m "math"
	"os"

	"github.com/go-gl/glow/gl-core/4.5/gl"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongodb    string = "localhost:27017"
	chunkSize  int    = 16
	maxHeight  int    = 128
	seaLevel   int    = 64
	renderSize int    = 8
	viewSize   int    = 32
)

var (
	logFile      *os.File
	mongoSession *mgo.Session
	cubeLoader   []Cube
)

type Level struct {
	chunks []*Chunk
	noise  *simplexNoise
}

type Chunk struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	XPos, ZPos  int
	fullyLoaded bool
	cubes       [chunkSize * chunkSize][maxHeight]float32
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
			if c.cubes[pXpC*chunkSize+pZpC][pY] != Empty ||
				c.cubes[pXpC*chunkSize+pZmC][pY] != Empty ||
				c.cubes[pXmC*chunkSize+pZpC][pY] != Empty ||
				c.cubes[pXmC*chunkSize+pZmC][pY] != Empty {
				return true
			}
		}
	}
	return false
}

func (gameMap *Level) RenderLevel(vao, typeBuffer uint32, chunkPosition int32) {

	gl.BindVertexArray(vao)
	for _, c := range gameMap.chunks {
		if c == nil {
			continue
		}
		for x := 0; x < chunkSize; x++ {
			for z := 0; z < chunkSize; z++ {
				gl.Uniform2f(chunkPosition, float32((c.XPos*chunkSize)+x), float32((c.ZPos*chunkSize)+z))

				slice := c.cubes[x*chunkSize+z][:]

				gl.BindBuffer(gl.ARRAY_BUFFER, typeBuffer)
				gl.BufferData(gl.ARRAY_BUFFER, len(slice)*4, gl.Ptr(slice), gl.STATIC_DRAW)

				gl.DrawElementsInstanced(gl.TRIANGLES, 36, gl.UNSIGNED_INT, gl.Ptr(nil), int32(maxHeight))
			}
		}
	}
}

func (gameMap *Level) genChunk(ch *Chunk) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	err := collection.Insert(ch)
	mongoSession.Fsync(false)
	collection.Find(ch).One(ch)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}

	for x := 0; x < chunkSize; x++ {
		for z := 0; z < chunkSize; z++ {
			gameMap.genColumn(ch, x, z)
		}
	}
	go ch.persistChunk()
}

func (gameMap *Level) genColumn(ch *Chunk, x, z int) {
	n := (gameMap.noise.getNoise(float64((ch.XPos*chunkSize)+x), float64((ch.ZPos*chunkSize)+z)) + 1.0) / 2.0
	h := int((n * 4) + 60)
	//fmt.Printf("Got a Noise of %v\n", n)

	var cubeType uint8
	for y := 0; y < h; y++ {
		if (y + 1) >= h {
			cubeType = Grass
		} else if (y + 5) >= h {
			cubeType = Dirt
		} else if (y + 10) >= h {
			cubeType = Gravel
		} else {
			cubeType = Stone
		}
		ch.cubes[x*chunkSize+z][y] = float32(cubeType)
	}

}

func (ch *Chunk) persistChunk() {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	bulk := collection.Bulk()
	for x := 0; x < chunkSize; x++ {
		for z := 0; z < chunkSize; z++ {
			for y := 0; y < maxHeight; y++ {
				bulk.Insert(Cube{ChunkID: ch.ID, XPos: uint8(x), YPos: uint8(y), ZPos: uint8(z), CubeType: uint8(ch.cubes[x*chunkSize+z][y])})
			}
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
		if ch.XPos == (x-(renderSize+1)) || ch.ZPos == (z-(renderSize+1)) || ch.XPos == (x+(renderSize+1)) || ch.ZPos == (z+(renderSize+1)) {
			copy(gameMap.chunks[i:], gameMap.chunks[i+1:])
			gameMap.chunks[len(gameMap.chunks)-1] = nil
			gameMap.chunks = gameMap.chunks[:len(gameMap.chunks)-1]
			gameMap.removeOldChunks(x, z)
			break
		}
	}
}

func (gameMap *Level) loadNewChunk(ch *Chunk, x, z int) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	err := collection.Find(bson.M{"xpos": x, "zpos": z}).One(ch)
	if err != nil {
		fmt.Printf("Creating Chunk at X %v Z %v\n", x, z)
		gameMap.genChunk(ch)
		mongoSession.Fsync(false)
	} else {
		ch.loadChunk(mongoSession)
	}
}

func (c *Chunk) loadChunk(mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	collection.Find(bson.M{"chunkid": c.ID}).All(&cubeLoader)
	c.Update(cubeLoader)
	c.fullyLoaded = true
	//fmt.Printf("Cubes: %v\n", c.cubes[:])
}

func (c *Chunk) Update(cubes []Cube) {
	for _, cube := range cubes {
		c.cubes[int(cube.XPos)*chunkSize+int(cube.ZPos)][int(cube.YPos)] = float32(cube.GetCubeType())
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
	gameMap.noise = createSimplexNoise(200, 255.0, 0.5)
	x := int(m.Floor(pX / float64(chunkSize)))
	z := int(m.Floor(pZ / float64(chunkSize)))
	ch := Chunk{}
	ch.XPos = x
	ch.ZPos = z
	gameMap.chunks = append(gameMap.chunks, &ch)
	gameMap.loadNewChunk(&ch, x, z)
}

func (gameMap *Level) initChunk(x, z int) {
	ch := Chunk{}
	ch.XPos = x
	ch.ZPos = z
	gameMap.chunks = append(gameMap.chunks, &ch)

	gameMap.loadNewChunk(&ch, x, z)
}

func (gameMap *Level) LoopChunkLoader(pX, pZ float64) {

	x := int(m.Floor(pX / float64(chunkSize)))
	z := int(m.Floor(pZ / float64(chunkSize)))
	gameMap.removeOldChunks(x, z)
	for c := 1; c < renderSize; c++ {
		for i := x - c; i < x+c; i++ {
			for j := z - c; j < z+c; j++ {
				if !gameMap.checkForChunk(i, j) {
					gameMap.initChunk(i, j)
				}
			}
		}
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
	index := mgo.Index{
		Key: []string{"chunkid"},
	}
	collection := mongoSession.DB("GameDatabase").C("Cubes")
	collection.EnsureIndex(index)

}
func closeMongoSession() {
	logFile.Close()
	mongoSession.Close()
}
