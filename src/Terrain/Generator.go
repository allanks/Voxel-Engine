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
	layers     [128]*layer
}

type layer struct {
	cubes []Cube
	yPos  int
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
	for j, c := range gameMap.chunks {
		var waitGroup sync.WaitGroup
		waitGroup.Add(len(c.layers))
		for i := range c.layers {
			c.layers[i] = &layer{yPos: i + 1}

			go loadLayer(&waitGroup, c.layers[i], mongoSession)
		}
		waitGroup.Wait()
		gameMap.chunks[j] = c
	}

}

func loadLayer(waitGroup *sync.WaitGroup, l *layer, mongoSession *mgo.Session) {
	defer waitGroup.Done()
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	collection.Find(bson.M{"ypos": l.yPos}).All(&l.cubes)
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
	for i := range c.layers {
		if i < seaLevel {
			go genLayer(c.XPos, i+1, c.ZPos, mongoSession)
		}
	}
}

func genLayer(xPos, yPos, zPos int, mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	var cubeType int
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
			err := collection.Insert(&Cube{XPos: float64((xPos * chunkSize) + x), YPos: float64(yPos), ZPos: float64((zPos * chunkSize) + z), CubeType: cubeType})
			if err != nil {
				log.Printf("RunQuery : ERROR : %s\n", err)
			}
		}
	}
}

func FindNearestCubes(xPos, yPos, zPos float64) []Cube {
	var nearCubes = []Cube{}
	for _, c := range gameMap.chunks {
		for _, l := range c.layers {
			if float64(l.yPos) == yPos || float64(l.yPos) == (yPos+1) || float64(l.yPos) == (yPos-1) {
				for _, cube := range l.cubes {
					if (cube.XPos == xPos || cube.XPos == (xPos+1) || cube.XPos == (xPos-1)) &&
						(cube.ZPos == zPos || cube.ZPos == (zPos+1) || cube.ZPos == (zPos-1)) {
						nearCubes = append(nearCubes, cube)
					}
				}
			}
		}
	}
	return nearCubes
}

func IsInCube(xPos, yPos, zPos, collisionDistance float64) bool {
	for _, c := range gameMap.chunks {
		for _, l := range c.layers {
			if float64(l.yPos) == m.Floor(yPos) {
				for _, cube := range l.cubes {
					if (cube.XPos == m.Floor(xPos) || cube.XPos == m.Floor(xPos+collisionDistance) || cube.XPos == m.Floor(xPos-collisionDistance)) &&
						(cube.ZPos == m.Floor(zPos) || cube.ZPos == m.Floor(zPos+collisionDistance) || cube.ZPos == m.Floor(zPos-collisionDistance)) {
						return true
					}
				}
			}
		}
	}
	return false
}
