package Terrain

import (
	"fmt"
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
	cubes      []Cube
	instances  []float32
	colors     []float32
}

func (c *chunk) getPositions() []float32 {
	return c.instances
}

func (c *chunk) getColors() []float32 {
	return c.colors
}

func (c *chunk) Update() {
	for _, cube := range c.cubes {
		c.instances = append(c.instances, float32(cube.XPos), float32(cube.YPos), float32(cube.ZPos))
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
		waitGroup.Add(1)
		go genChunk(0, 0, &waitGroup, mongoSession)
		/*go genChunk(1, 0, &waitGroup, mongoSession)
		go genChunk(0, 1, &waitGroup, mongoSession)
		go genChunk(-1, 0, &waitGroup, mongoSession)
		go genChunk(0, -1, &waitGroup, mongoSession)*/
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
	fmt.Printf("%v%v\n", "GameMap: ", gameMap)
	for _, c := range gameMap.chunks {
		fmt.Printf("%v%v\n", "Chunk: ", len(c.instances))
	}

}

func loadChunk(waitGroup *sync.WaitGroup, c *chunk, mongoSession *mgo.Session) {
	defer waitGroup.Done()
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	collection.Find(bson.M{}).All(&c.cubes)
	c.Update()
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
	go genLayer(c.ID, c.XPos, c.ZPos, mongoSession)
}

func genLayer(chunkID bson.ObjectId, xPos, zPos int, mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
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
				err := collection.Insert(&Cube{ChunkID: chunkID, XPos: float64(x), YPos: float64(yPos), ZPos: float64(z), CubeType: cubeType})
				if err != nil {
					log.Printf("RunQuery : ERROR : %s\n", err)
				}
			}
		}
	}
}

func FindNearestCubes(xPos, yPos, zPos float64) []Cube {
	var nearCubes = []Cube{}
	for _, c := range gameMap.chunks {
		x := c.XPos * chunkSize
		z := c.ZPos * chunkSize
		for _, cube := range c.cubes {
			cX := cube.XPos + float64(x)
			cZ := cube.ZPos + float64(z)
			if (cX == xPos || cX == (xPos+1) || cX == (xPos-1)) &&
				(cube.YPos == yPos || cube.YPos == (yPos+1) || cube.YPos == (yPos-1)) &&
				(cZ == zPos || cZ == (zPos+1) || cZ == (zPos-1)) {
				nearCubes = append(nearCubes, cube)
			}
		}
	}
	return nearCubes
}

func IsInCube(xPos, yPos, zPos, collisionDistance float64) bool {
	for _, c := range gameMap.chunks {
		x := c.XPos * chunkSize
		z := c.ZPos * chunkSize
		for _, cube := range c.cubes {
			cX := cube.XPos + float64(x)
			cZ := cube.ZPos + float64(z)
			if (cX == m.Floor(xPos) || cX == m.Floor(xPos+collisionDistance) || cX == m.Floor(xPos-collisionDistance)) &&
				(cZ == m.Floor(zPos) || cZ == m.Floor(zPos+collisionDistance) || cZ == m.Floor(zPos-collisionDistance)) &&
				cube.YPos == m.Floor(yPos) {
				return true
			}
		}
	}
	return false
}
