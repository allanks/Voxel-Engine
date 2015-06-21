package Terrain

import (
	"encoding/gob"
	"fmt"
	"log"
	m "math"
	"net"

	"github.com/allanks/Voxel-Engine/src/Model"
	"github.com/allanks/Voxel-Engine/src/Server/DataType"
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

func (gameMap *Level) GetYCubes(xPos, yPos, zPos, height float64) []float32 {
	pX := int(m.Floor(xPos))
	pY := int(m.Floor(yPos))
	nY := int(m.Floor(yPos - height))
	pZ := int(m.Floor(zPos))
	cubes := []float32{}
	for _, c := range gameMap.chunks {
		if c == nil {
			continue
		}
		if isInRange(pX/chunkSize, c.XPos) && isInRange(pZ/chunkSize, c.ZPos) {
			for i := 0; i < len(c.drawables)/4; i++ {
				if int(m.Floor(float64(c.drawables[(i*4)+3]))) != DataType.Empty &&
					int(m.Floor(float64(c.drawables[i*4])))+(c.XPos*chunkSize) == pX &&
					int(m.Floor(float64(c.drawables[(i*4)+2])))+(c.ZPos*chunkSize) == pZ &&
					(isInRange(int(m.Floor(float64(c.drawables[(i*4)+1]))), pY) ||
						isInRange(int(m.Floor(float64(c.drawables[(i*4)+1]))), nY)) {
					cubes = append(cubes, c.drawables[i*4]+float32(c.XPos*chunkSize), c.drawables[(i*4)+1], c.drawables[(i*4)+2]+float32(c.ZPos*chunkSize))
				}
			}
		}
	}
	return cubes
}

func (gameMap *Level) GetCubes(query []DataType.Pos) []DataType.Pos {
	cubes := []DataType.Pos{}
	for _, cubeAt := range query {
		qx, qy, qz := DataType.FloorToInt(cubeAt.XPos), DataType.FloorToInt(cubeAt.YPos), DataType.FloorToInt(cubeAt.ZPos)
		for _, chunk := range gameMap.chunks {
			if qx/chunkSize == chunk.XPos && qz/chunkSize == chunk.ZPos {
				for i := 0; i < len(chunk.drawables)/4; i++ {
					if DataType.FloorToInt(chunk.drawables[(i*4)+3]) != DataType.Empty &&
						DataType.FloorToInt(chunk.drawables[i*4])+(chunk.XPos*chunkSize) == qx &&
						DataType.FloorToInt(chunk.drawables[(i*4)+1]) == qy &&
						DataType.FloorToInt(chunk.drawables[(i*4)+2])+(chunk.ZPos*chunkSize) == qz {
						cubes = append(cubes, DataType.Pos{
							XPos: chunk.drawables[i*4] + float32(chunk.XPos*chunkSize),
							YPos: chunk.drawables[(i*4)+1],
							ZPos: chunk.drawables[(i*4)+2] + float32(chunk.ZPos*chunkSize)})
					}
				}
			}
		}
	}
	return cubes
}

func (gameMap *Level) GetXZCubes(xPos, yPos, zPos float64) []float32 {
	pX := int(m.Floor(xPos))
	pY := int(m.Floor(yPos))
	pZ := int(m.Floor(zPos))
	cubes := []float32{}
	for _, c := range gameMap.chunks {
		if c == nil {
			continue
		}
		if isInRange(pX/chunkSize, c.XPos) && isInRange(pZ/chunkSize, c.ZPos) {
			for i := 0; i < len(c.drawables)/4; i++ {
				if int(m.Floor(float64(c.drawables[(i*4)+3]))) != DataType.Empty &&
					int(m.Floor(float64(c.drawables[(i*4)+1]))) == pY &&
					isInRange(int(m.Floor(float64(c.drawables[i*4])))+(c.XPos*chunkSize), pX) &&
					isInRange(int(m.Floor(float64(c.drawables[(i*4)+2])))+(c.ZPos*chunkSize), pZ) {
					cubes = append(cubes, c.drawables[i*4]+float32(c.XPos*chunkSize), c.drawables[(i*4)+1], c.drawables[(i*4)+2]+float32(c.ZPos*chunkSize))
				}
			}
		}
	}
	return cubes
}

func isInRange(static, dynamic int) bool {
	return static == dynamic || static == dynamic+1 || static == dynamic-1
}

func (gameMap *Level) RenderLevel() {

	for _, c := range gameMap.chunks {
		if c == nil || len(c.drawables) == 0 {
			continue
		}

		Model.BindBuffers([]float32{float32(c.XPos * chunkSize), 0.0, float32(c.ZPos * chunkSize)}, Model.Cube)
		Model.Render(c.drawables, Model.Cube)
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
