package Terrain

import (
	m "math"
	"math/rand"
	"time"
)

type simplexNoise struct {
	octaves                 []noise
	frequencies, amplitudes []float64
}

func createSimplexNoise(seed int64, height, persistence float64) *simplexNoise {
	numOfOctaves := int(m.Ceil(m.Log2(height)))
	simplex := simplexNoise{}

	simplex.octaves = make([]noise, numOfOctaves)
	simplex.amplitudes = make([]float64, numOfOctaves)
	simplex.frequencies = make([]float64, numOfOctaves)

	if seed == 0 {
		seed = rand.New(rand.NewSource(time.Now().Unix())).Int63()
	}

	random := rand.New(rand.NewSource(seed))

	for i := 0; i < numOfOctaves; i++ {
		simplex.octaves[i] = initNoise(random.Int63())
		simplex.frequencies[i] = m.Pow(float64(2), float64(i))
		simplex.amplitudes[i] = m.Pow(persistence, float64(numOfOctaves-i))
	}

	return &simplex
}

func (simplex *simplexNoise) getNoise(x, z float64) float64 {
	result := float64(0)

	for i := 0; i < len(simplex.octaves); i++ {
		freq := simplex.frequencies[i]
		result = result + (simplex.octaves[i].generateNoise(x/freq, z/freq) * simplex.amplitudes[i])
	}
	return result
}

// Noise Function

const (
	NUM_OF_SWAP int = 400
)

var grad3 [][]int = [][]int{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
}
var p_supply []int = []int{151, 160, 137, 91, 90, 15, //this contains all the numbers between 0 and 255, these are put in a random order depending upon the seed
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
}

type noise struct {
	perm, permMod12 [512]int
}

func initNoise(seed int64) noise {
	n := noise{}
	for i := 0; i < len(n.perm); i++ {
		n.perm[i] = i & 255
	}

	if seed == 0 {
		seed = rand.New(rand.NewSource(time.Now().Unix())).Int63()
	}

	p := make([]int, 256)
	copy(p, p_supply)

	random := rand.New(rand.NewSource(seed))

	for i := 0; i < NUM_OF_SWAP; i++ {
		from := random.Intn(len(p_supply))
		to := random.Intn(len(p_supply))

		tmp := p[from]
		p[from] = p[to]
		p[to] = tmp
	}

	for i := 0; i < 512; i++ {
		n.perm[i] = p[i&255]
		n.permMod12[i] = n.perm[i] % 12
	}
	return n
}

func (n *noise) generateNoise(x, z float64) float64 {
	skew := 0.5 * (m.Sqrt(3) - 1)
	fac := (x + z) * skew
	i := m.Floor(x + fac)
	j := m.Floor(z + fac)

	unSkew := (3.0 - m.Sqrt(3.0)) / 6.0
	t := (i + j) * unSkew
	x0 := x - i + t
	z0 := z - i + t
	var tri int
	if x0 > z0 {
		tri = 1
	} else {
		tri = 0
	}
	trj := 1 - tri

	x1 := x0 - float64(tri) + unSkew
	z1 := z0 - float64(trj) + unSkew
	x2 := x0 - 1 + (2 * unSkew)
	z2 := z0 - 1 + (2 * unSkew)

	ii := int(i) & 255
	jj := int(j) & 255
	gi0 := n.permMod12[ii+n.perm[jj]] % 12
	gi1 := n.permMod12[ii+tri+n.perm[jj+trj]] % 12
	gi2 := n.permMod12[ii+1+n.perm[jj+1]] % 12

	t0 := 0.5 - (x0 * x0) - (z0 * z0)
	var n0 float64
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 = t0 * t0
		n0 = t0 * t0 * dot(grad3[gi0], x0, z0)
	}

	t1 := 0.5 - (x1 * x1) - (z1 * z1)
	var n1 float64
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 = t1 * t1
		n1 = t1 * t1 * dot(grad3[gi1], x1, z1)
	}

	t2 := 0.5 - x2*x2 - z2*z2
	var n2 float64
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 = t2 * t2
		n2 = t2 * t2 * dot(grad3[gi2], x2, z2)
	}

	return 70.0 * (n0 + n1 + n2)
}

func dot(g []int, x, y float64) float64 {
	return float64(g[0])*x + float64(g[1])*y
}
