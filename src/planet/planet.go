package planet

import (
	"github.com/go-gl/mathgl/mgl32"
)

func GeneratePlanetVertices(resolution int, center mgl32.Vec3) []float32 {
	vertices := []float32{}
	for i := 0; i < resolution; i++ {
		for j := 0; j < resolution; j++ {
			vertices = append(vertices, float32(i), float32(j), 0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
		}
	}
	for i := 0; i < resolution; i++ {
		for j := 0; j < resolution; j++ {
			vertices = append(vertices, float32(i), 0.0, float32(j), 0.0, 0.0, 0.0, 0.0, 0.0)
		}
	}
	for i := 0; i < resolution; i++ {
		for j := 0; j < resolution; j++ {
			vertices = append(vertices, 0.0, float32(i), float32(j), 0.0, 0.0, 0.0, 0.0, 0.0)
		}
	}
	for i := 0; i < resolution; i++ {
		for j := 0; j < resolution; j++ {
			vertices = append(vertices, float32(i), float32(resolution-1), float32(j), 0.0, 0.0, 0.0, 0.0, 0.0)
		}
	}
	for i := 0; i < resolution; i++ {
		for j := 0; j < resolution; j++ {
			vertices = append(vertices, float32(resolution-1), float32(i), float32(j), 0.0, 0.0, 0.0, 0.0, 0.0)
		}
	}
	for i := 0; i < resolution; i++ {
		for j := 0; j < resolution; j++ {
			vertices = append(vertices, float32(i), float32(j), float32(resolution-1), 0.0, 0.0, 0.0, 0.0, 0.0)
		}
	}

	return vertices
}

func GenerateIndices(resolution uint32) []uint32 {
	indices := []uint32{}
	var mult uint32 = 0
	for k := 0; k < 6; k++ {
		for i := uint32(0); i < resolution-1; i++ {
			for j := uint32(0); j < resolution-1; j++ {
				// never trust prioriteringsregler
				index := uint32(j + (i * resolution))
				indices = append(indices, mult+index, mult+(index+resolution+1), mult+(index+resolution))
				indices = append(indices, mult+index, mult+(index+1), mult+(index+resolution+1))
			}
		}
		mult += resolution * resolution
	}

	return indices
}

/*func NormalizeFromCenter(x float32, y float32, z float32, center mgl32.Vec3) mgl32.Vec3 {
	// 1. Normera vektorn: norm = center - coord
	// 2. Placera koordinaten på center + norm
}*/
