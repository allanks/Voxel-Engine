package ObjectLoader

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func LoadObjFile(fileName string) ([]float32, []float32, []float32) {

	objFile, err := os.Open("resource/models/" + fileName)
	if err != nil {
		panic(err)
	}
	defer objFile.Close()

	fileReader := bufio.NewScanner(objFile)

	vertices := []float32{}
	textures := []float32{}
	normals := []float32{}

	fileVertices := []float32{}
	fileTextures := []float32{}
	fileNormals := []float32{}

	for fileReader.Scan() {
		line := fileReader.Text()
		split := strings.Split(line, " ")
		switch split[0] {
		case "v":
			x, _ := strconv.ParseFloat(split[1], 64)
			y, _ := strconv.ParseFloat(split[2], 64)
			z, _ := strconv.ParseFloat(split[3], 64)
			fileVertices = append(fileVertices, float32(x), float32(y), float32(z))
		case "vt":
			x, _ := strconv.ParseFloat(split[1], 64)
			y, _ := strconv.ParseFloat(split[2], 64)
			fileTextures = append(fileTextures, float32(x), float32(y))
		case "vn":
			x, _ := strconv.ParseFloat(split[1], 64)
			y, _ := strconv.ParseFloat(split[2], 64)
			z, _ := strconv.ParseFloat(split[3], 64)
			fileNormals = append(fileNormals, float32(x), float32(y), float32(z))
		case "f":
			face := strings.Split(split[1], "/")
			vertices = append(vertices, parseFaceTriple(face[0], fileVertices)...)
			textures = append(textures, parseFaceDouble(face[1], fileTextures)...)
			normals = append(normals, parseFaceTriple(face[2], fileNormals)...)
			face = strings.Split(split[2], "/")
			vertices = append(vertices, parseFaceTriple(face[0], fileVertices)...)
			textures = append(textures, parseFaceDouble(face[1], fileTextures)...)
			normals = append(normals, parseFaceTriple(face[2], fileNormals)...)
			face = strings.Split(split[3], "/")
			vertices = append(vertices, parseFaceTriple(face[0], fileVertices)...)
			textures = append(textures, parseFaceDouble(face[1], fileTextures)...)
			normals = append(normals, parseFaceTriple(face[2], fileNormals)...)
		}
	}

	return vertices, normals, textures
}

func parseFaceDouble(index string, slice []float32) []float32 {
	if index != "" {
		x, _ := strconv.ParseInt(index, 10, 64)
		return []float32{slice[((x - 1) * 2)], slice[((x-1)*2)+1]}
	} else {
		return []float32{}
	}
}

func parseFaceTriple(index string, slice []float32) []float32 {
	if index != "" {
		x, _ := strconv.ParseInt(index, 10, 64)
		return []float32{slice[((x - 1) * 3)], slice[((x-1)*3)+1], slice[((x-1)*3)+2]}
	} else {
		return []float32{}
	}
}
