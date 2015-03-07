package Terrain

import (
	"log"
	"os"
	"sync"

	"gopkg.in/mgo.v2"
)

const (
	mongodb string = "localhost:27017"
)

var (
	cubes []*Cube
)

type ()

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
		go genLayer(0, -1, 0, 5, &waitGroup, mongoSession)
		go genLayer(0, -2, 0, 6, &waitGroup, mongoSession)
		go genLayer(0, -3, 0, 7, &waitGroup, mongoSession)
		go genLayer(0, -4, 0, 8, &waitGroup, mongoSession)
		go genLayer(0, -5, 0, 7, &waitGroup, mongoSession)
		go genLayer(0, -6, 0, 6, &waitGroup, mongoSession)
		go genLayer(0, -7, 0, 5, &waitGroup, mongoSession)
		waitGroup.Wait()
	}

	err = collection.Find(nil).All(&cubes)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

}

func genLayer(xPos, yPos, zPos, size int32, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) {
	defer waitGroup.Done()
	for i := -size; i < size; i++ {
		go genRow(yPos, i, size, mongoSession)
	}
}

func genRow(yPos, zPos, size int32, mongoSession *mgo.Session) {
	session := mongoSession.Copy()
	defer session.Close()
	collection := session.DB("GameDatabase").C("Cubes")

	for i := -size; i < size; i++ {
		log.Printf("Creating Cube %v\n", &Cube{XPos: i, YPos: yPos, ZPos: zPos})
		err := collection.Insert(&Cube{XPos: i, YPos: yPos, ZPos: zPos})
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
