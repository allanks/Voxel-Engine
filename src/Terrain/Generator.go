package Terrain

import (
	"log"
	m "math"
	"os"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongodb   string = "localhost:27017"
	chunkSize int    = 32
	seaLevel  int    = 30
)

var (
	gameMap level
)

type level struct {
	chunks []*chunk
}

type chunk struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	XPos, ZPos int
	instances  []int
	colors     []float32
}

func (c *chunk) getPositions() []float32 {
	val := []float32{}
	for i := 0; i < (len(c.instances) / 3); i++ {
		val = append(val, float32((c.XPos*chunkSize)+c.instances[i*3]), float32(c.instances[(i*3)+1]), float32((c.ZPos*chunkSize)+c.instances[(i*3)+2]))
	}
	return val
}

func (c *chunk) getColors() []float32 {
	return c.colors
}

func (c *chunk) Update(cubes []Cube) {
	for _, cube := range cubes {
		c.instances = append(c.instances, cube.XPos, cube.YPos, cube.ZPos)
		gCube := gCubes[cube.CubeType]
		c.colors = append(c.colors, gCube.getColors()...)
	}
}

func GenLevel() {

	f, err := os.OpenFile("LevelGen.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	var mongoSession *mgo.Session
	mongoSession, err = mgo.Dial(mongodb)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	defer mongoSession.Close()

	var count int
	collection := mongoSession.DB("GameDatabase").C("Cubes")
	count, err = collection.Count()
	if err != nil {
		log.Fatalf("Count Failed: %s\n", err)
	}
	if count == 0 {
		var waitGroup sync.WaitGroup
		waitGroup.Add(13)
		go genChunk(0, 0, &waitGroup, mongoSession)
		go genChunk(1, 0, &waitGroup, mongoSession)
		go genChunk(0, 1, &waitGroup, mongoSession)
		go genChunk(1, 1, &waitGroup, mongoSession)
		go genChunk(-1, 0, &waitGroup, mongoSession)
		go genChunk(0, -1, &waitGroup, mongoSession)
		go genChunk(-1, -1, &waitGroup, mongoSession)
		go genChunk(1, -1, &waitGroup, mongoSession)
		go genChunk(-1, 1, &waitGroup, mongoSession)
		go genChunk(2, 0, &waitGroup, mongoSession)
		go genChunk(0, 2, &waitGroup, mongoSession)
		go genChunk(-2, 0, &waitGroup, mongoSession)
		go genChunk(0, -2, &waitGroup, mongoSession)
		waitGroup.Wait()
		mongoSession.Fsync(false)
	}

	loadGameMap(mongoSession)
}

func loadGameMap(mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	err := collection.Find(nil).All(&gameMap.chunks)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(gameMap.chunks))
	for _, c := range gameMap.chunks {
		go loadChunk(&waitGroup, c, mongoSession)
	}
	waitGroup.Wait()

}

func loadChunk(waitGroup *sync.WaitGroup, c *chunk, mongoSession *mgo.Session) {
	defer waitGroup.Done()
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	cubes := []Cube{}
	collection.Find(bson.M{"chunkid": c.ID}).All(&cubes)
	c.Update(cubes)
}

func genChunk(x, z int, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) {
	defer waitGroup.Done()
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	c := &chunk{XPos: x, ZPos: z}
	err := collection.Insert(c)
	mongoSession.Fsync(false)
	collection.Find(c).One(c)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
	collection = session.DB("GameDatabase").C("Cubes")
	var cubeType int
	for yPos := 0; yPos < seaLevel; yPos++ {
		if (yPos + 1) >= seaLevel {
			cubeType = Grass
		} else if (yPos + 5) >= seaLevel {
			cubeType = Dirt
		} else if (yPos + 10) >= seaLevel {
			cubeType = Gravel
		} else {
			cubeType = Stone
		}
		for x := 0; x < chunkSize; x++ {
			for z := 0; z < chunkSize; z++ {
				err := collection.Insert(&Cube{ChunkID: c.ID, XPos: x, YPos: yPos, ZPos: z, CubeType: cubeType})
				if err != nil {
					log.Printf("RunQuery : ERROR : %s\n", err)
				}
			}
		}
	}
}

func IsInCube(xPos, yPos, zPos, collisionDistance float64) bool {
	pX := int(m.Floor(xPos))
	pY := int(m.Floor(yPos))
	pZ := int(m.Floor(zPos))
	pXpC := int(m.Floor(xPos + collisionDistance))
	pXmC := int(m.Floor(xPos - collisionDistance))
	pZpC := int(m.Floor(zPos + collisionDistance))
	pZmC := int(m.Floor(zPos - collisionDistance))
	for _, c := range gameMap.chunks {
		x := c.XPos * chunkSize
		z := c.ZPos * chunkSize
		for i := 0; i < (len(c.instances) / 3); i++ {
			cX := c.instances[i*3] + x
			cY := c.instances[(i*3)+1]
			cZ := c.instances[(i*3)+2] + z
			if (cX == pX || cX == pXpC || cX == pXmC) &&
				(cZ == pZ || cZ == pZpC || cZ == pZmC) &&
				cY == pY {
				return true
			}
		}
	}
	return false
}
