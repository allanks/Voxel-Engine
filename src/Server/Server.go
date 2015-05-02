package Server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/allanks/Voxel-Engine/src/Server/DataType"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongodb   string = "localhost:27017"
	chunkSize int    = 16
	maxHeight int    = 128
	seaLevel  int    = 64
)

var (
	logFile      *os.File
	mongoSession *mgo.Session
	cubeLoader   []cube
)

var serverNoise = []*SimplexNoise{
	&SimplexNoise{},
}

type cube struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	ChunkID          bson.ObjectId
	XPos, YPos, ZPos int8
	CubeType         uint8
	Visible          bool
}

func InitServer() {
	createDatabaseLink()
	fmt.Println("Listening")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			fmt.Printf("Error accepting connection %v\n", err)
			continue
		}
		go getChunk(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}

func getChunk(conn net.Conn) {
	fmt.Println("Serving Connection")

	for {
		dec := gob.NewDecoder(conn)
		c := &DataType.Chunk{}
		err := dec.Decode(c)
		if err == nil {
			loadChunk(c, conn)
		} else if neterr, ok := err.(net.Error); ok {
			fmt.Printf("Recieved error %v\n", neterr)
			break
		}
	}
}

func loadChunk(c *DataType.Chunk, conn net.Conn) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	encoder := gob.NewEncoder(conn)

	err := collection.Find(bson.M{"xpos": c.XPos, "zpos": c.ZPos}).One(c)
	if err != nil {
		encoder.Encode(genChunk(c))
	} else {
		encoder.Encode(fetchChunk(c))
	}
}

func genChunk(c *DataType.Chunk) *DataType.CubeChunk {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Chunks")
	err := collection.Insert(c)
	mongoSession.Fsync(false)
	collection.Find(c).One(c)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}

	var cubes []*cube
	for x := 0; x < chunkSize; x++ {
		for z := 0; z < chunkSize; z++ {
			n := (serverNoise[0].GetNoise(float64((c.XPos*chunkSize)+x), float64((c.ZPos*chunkSize)+z)) + 1.0) / 2.0
			h := int((n * 4) + 60)

			var cubeType uint8
			for y := -1; y < h+1; y++ {
				if y >= h || y < 0 {
					cubeType = DataType.Empty
				} else if (y + 1) >= h {
					cubeType = DataType.Grass
				} else if (y + 5) >= h {
					cubeType = DataType.Dirt
				} else if (y + 10) >= h {
					cubeType = DataType.Gravel
				} else {
					cubeType = DataType.Stone
				}
				cubes = append(cubes, &cube{ChunkID: c.ID, XPos: int8(x), YPos: int8(y), ZPos: int8(z), CubeType: uint8(cubeType)})
			}
		}
	}
	filteredCubes := filter(cubes)

	fmt.Printf("Created Chunk at X %v Z %v\n", c.XPos, c.ZPos)
	return &DataType.CubeChunk{XPos: c.XPos, ZPos: c.ZPos, Cubes: filteredCubes}
}

func filter(cubes []*cube) []float32 {
	drawables := []float32{}
	for _, current := range cubes {
		current.Visible = false
		if current.CubeType == DataType.Empty {
			continue
		}
		for _, other := range cubes {
			xDiff := (other.XPos-current.XPos) == 1 || (other.XPos-current.XPos) == -1
			yDiff := (other.YPos-current.YPos) == 1 || (other.YPos-current.YPos) == -1
			zDiff := (other.ZPos-current.ZPos) == 1 || (other.ZPos-current.ZPos) == -1
			xSame := (other.XPos - current.XPos) == 0
			ySame := (other.YPos - current.YPos) == 0
			zSame := (other.ZPos - current.ZPos) == 0
			if other.CubeType == DataType.Empty && ((xDiff && ySame && zSame) || (xSame && yDiff && zSame) || (xSame && ySame && zDiff)) {
				current.Visible = true
				drawables = append(drawables, float32(current.XPos), float32(current.YPos), float32(current.ZPos), float32(current.CubeType))
				break
			}
		}
	}
	go persistChunk(cubes)
	return drawables
}

func persistChunk(cubes []*cube) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	bulk := collection.Bulk()
	for _, current := range cubes {
		bulk.Insert(current)
	}
	_, err := bulk.Run()
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
	}
}

func fetchChunk(c *DataType.Chunk) *DataType.CubeChunk {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")
	collection.Find(bson.M{"chunkid": c.ID}).All(&cubeLoader)
	drawables := []float32{}
	for _, current := range cubeLoader {
		if current.Visible {
			drawables = append(drawables, float32(current.XPos), float32(current.YPos), float32(current.ZPos), float32(current.CubeType))
		}
	}
	fmt.Printf("Loaded Chunk at X %v Z %v\n", c.XPos, c.ZPos)
	return &DataType.CubeChunk{XPos: c.XPos, ZPos: c.ZPos, Cubes: drawables}
}

func LoadGameMap() {
	if mongoSession == nil {
		createDatabaseLink()
	}
	serverNoise[0] = CreateSimplexNoise(200, 255.0, 0.5)
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
