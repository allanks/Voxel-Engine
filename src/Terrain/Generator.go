package Terrain

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/allanks/third-game/src/Player"

	"gopkg.in/mgo.v2"
)

const (
	mongodb string = "localhost:27017"
)

var (
	cubes []*Cube
)

func CheckPlayerCollisions() {
	for _, cube := range cubes {
		x, y, z := Player.GetPosition()
		Player.SetPosistion(cube.CheckCollision(x, y, z, Player.GetPlayerSpeed()))
		xPos, yPos, zPos := cube.CheckCollision(x, y-Player.Height, z, Player.GetPlayerSpeed())
		Player.SetPosistion(xPos, yPos+Player.Height, zPos)
	}
}

func GenLevel(xPos, yPos, zPos int32) {

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
		waitGroup.Add(7)
		go genLayer(0, -1, 0, 5, Grass, &waitGroup, mongoSession)
		go genLayer(0, -2, 0, 6, Grass, &waitGroup, mongoSession)
		go genLayer(0, -3, 0, 7, Grass, &waitGroup, mongoSession)
		go genLayer(0, -4, 0, 8, Grass, &waitGroup, mongoSession)
		go genLayer(0, -5, 0, 7, Grass, &waitGroup, mongoSession)
		go genLayer(0, -6, 0, 6, Grass, &waitGroup, mongoSession)
		go genLayer(0, -7, 0, 5, Grass, &waitGroup, mongoSession)
		waitGroup.Wait()
	}

	err = collection.Find(nil).All(&cubes)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

}

func genLayer(xPos, yPos, zPos, size int32, cubeType float64, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) {
	defer waitGroup.Done()
	for i := -size; i < size; i++ {
		go genRow(yPos, i, size, cubeType, mongoSession)
	}
}

func genRow(yPos, zPos, size int32, cubeType float64, mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")

	for i := -size; i < size; i++ {
		log.Printf("Creating Cube %v\n", &Cube{XPos: float64(i), YPos: float64(yPos), ZPos: float64(zPos), CubeType: cubeType})
		err := collection.Insert(&Cube{XPos: float64(i), YPos: float64(yPos), ZPos: float64(zPos), CubeType: cubeType})
		if err != nil {
			log.Printf("RunQuery : ERROR : %s\n", err)
		}
	}
}

func RenderLevel(vertAttrib, texCoordAttrib uint32, translateUniform int32) {

	for _, cube := range cubes {
		Render(cube, vertAttrib, texCoordAttrib, translateUniform)
	}
}

func PrintCubePos() {
	for _, cube := range cubes {
		fmt.Printf("%v\n", cube)
	}
}
