package planet

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

var seed = float32(0.0)

// Crater settings
const rimWidth = 0.7
const rimSteepness = 0.4
const smoothness = 0.3
const floorHeight = -0.3

type Crater struct {
	position mgl32.Vec3
	radius   float32
}

var craters = []Crater{}

func GenTerrain(points []mgl32.Vec3, radius float32, numCraters uint32) {
	genCraters(numCraters)

	seed = rand.Float32() * 1.0e5

	// TODO: calculate point for slices seperately using goroutines
	for i := 0; i < len(points); i++ {
		points[i] = points[i].Mul(getHeightAtPoint(points[i]))
		points[i] = points[i].Mul(radius)
	}
}

func getHeightAtPoint(point mgl32.Vec3) float32 {
	height := float32(1.0)

	height += getDetailNoiseHeight(point, 0.075, 0.75)

	height += getRidgeNoiseHeight(point, 0.075, 1.75)

	height += getCraterHeight(point)

	return height
}

func getDetailNoiseHeight(point mgl32.Vec3, amplitude, frequency float32) float32 {
	noiseHeight := float32(0.0)

	for i := 0; i < 5; i++ {
		x, y, z := point.X()*frequency, point.Y()*frequency, point.Z()*frequency
		noiseHeight += Snoise(x, y, z, seed) * float32(amplitude)
		frequency *= 2.0
		amplitude *= 0.5
	}

	return noiseHeight
}

func getRidgeNoiseHeight(point mgl32.Vec3, amplitude, frequency float32) float32 {
	x, y, z := point.X()*frequency, point.Y()*frequency, point.Z()*frequency
	return 0.5 - float32(math.Abs(float64(Snoise(x, y, z, seed))))*amplitude
}

func getCraterHeight(point mgl32.Vec3) float32 {
	craterHeight := float32(0.0)

	for i := 0; i < len(craters); i++ {
		x := mgl32.Vec3.Len(point.Sub(craters[i].position)) / craters[i].radius

		cavity := x*x - 1.0
		rimX := math.Min(float64(x-1.0-rimWidth), 0)
		rim := rimSteepness * rimX * rimX

		craterShape := smoothMax(cavity, floorHeight, smoothness)
		craterShape = smoothMin(craterShape, float32(rim), smoothness)
		craterHeight += craterShape * craters[i].radius
	}

	return craterHeight
}

func genCraters(numCraters uint32) {
	craters = make([]Crater, numCraters)

	for i := 0; i < len(craters); i++ {
		position := randomPointOnSphere()
		radius := float32(math.Pow(rand.Float64(), 2) * 0.25)
		craters[i] = Crater{position, radius}
	}
}

func randomPointOnSphere() mgl32.Vec3 {
	theta := rand.Float64() * 2.0 * math.Pi
	phi := rand.Float64() * math.Pi
	x := math.Cos(theta) * math.Sin(phi)
	y := math.Sin(theta) * math.Sin(phi)
	z := math.Cos(phi)

	return mgl32.Vec3{float32(x), float32(y), float32(z)}
}

func smoothMin(a, b, k float32) float32 {
	h := (b - a + k) / (2 * k)
	h = float32(math.Min(math.Max(float64(h), 0.0), 1.0))
	return a*h + b*(1-h) - k*h*(1-h)
}

func smoothMax(a, b, k float32) float32 {
	return smoothMin(a, b, -k)
}
