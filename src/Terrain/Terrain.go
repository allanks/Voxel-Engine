package Terrain

import (
	"encoding/gob"
	"fmt"
	"log"
	m "math"
	"net"

	"github.com/allanks/Voxel-Engine/src/Model"
	"github.com/allanks/Voxel-Engine/src/Server/DataType"
	"github.com/go-gl/glow/gl-core/4.5/gl"
)

const (
	chunkSize  int = 16
	renderSize int = 8
	viewSize   int = 32
)

var conn net.Conn

type Level struct {
	chunks []*clientChunk
}

type clientChunk struct {
	DataType.Chunk
	loaded    bool
	drawables []float32
}

func (gameMap *Level) IsInCube(xPos, yPos, zPos, collisionDistance float64) bool {
	pX := int(m.Floor(xPos))
	//pY := int(m.Floor(yPos))
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
		if pX == c.XPos*chunkSize && pZ == c.ZPos*chunkSize {
		}
	}
	return false
}

func (gameMap *Level) RenderLevel(vao, typeBuffer uint32, offset int32) {

	for _, c := range gameMap.chunks {
		if c == nil || len(c.drawables) == 0 {
			continue
		}

		gl.Uniform3f(offset, float32(c.XPos*chunkSize), 0.0, float32(c.ZPos*chunkSize))
		Model.Render(typeBuffer, c.drawables, Model.Cube)
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

func (gameMap *Level) InitChunk(x, z int) {
	c := clientChunk{DataType.Chunk{XPos: x / chunkSize, ZPos: z / chunkSize}, false, []float32{}}
	gameMap.chunks = append(gameMap.chunks, &c)
	gameMap.loadChunkFromServer(&c)
}

func (gameMap *Level) updateChunk(cubes *DataType.CubeChunk) {
	for _, chunk := range gameMap.chunks {
		if chunk.XPos == cubes.XPos && chunk.ZPos == cubes.ZPos {
			chunk.drawables = cubes.Cubes
			return
		}
	}
}

func (gameMap *Level) loadChunkFromServer(c *clientChunk) {
	encoder := gob.NewEncoder(conn)
	encoder.Encode(c.Chunk)
	cubes := &DataType.CubeChunk{}
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(cubes)
	if err == nil {
		gameMap.updateChunk(cubes)
	}
}

func (gameMap *Level) LoopChunkLoader(pX, pZ float64) {

	x := int(m.Floor(pX / float64(chunkSize)))
	z := int(m.Floor(pZ / float64(chunkSize)))
	gameMap.removeOldChunks(x, z)
	for c := 1; c < renderSize; c++ {
		for i := x - c; i < x+c; i++ {
			for j := z - c; j < z+c; j++ {
				if !gameMap.checkForChunk(i, j) {
					c := clientChunk{DataType.Chunk{XPos: i, ZPos: j}, false, []float32{}}
					gameMap.chunks = append(gameMap.chunks, &c)
					gameMap.loadChunkFromServer(&c)
				}
			}
		}
	}
}

func StartConnection() {
	var err error
	conn, err = net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Connection error", err)
	}
}

func CloseConnection() {
	conn.Close()
}

func (gameMap *Level) PrintChunks() {
	for _, chunk := range gameMap.chunks {
		fmt.Printf("chunk %v\n", chunk)
	}
}
