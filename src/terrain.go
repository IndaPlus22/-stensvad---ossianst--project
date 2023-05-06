package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Crater struct {
	position mgl32.Vec3
	radius   float32
}

var seed = float32(0.0)

func GenTerrain(points []mgl32.Vec3, shape *PlanetShape) {
	craters := genCraters(shape.numCraters)

	rand.NewSource(time.Now().UnixNano())
	seed = rand.Float32() * 1.0e5

	var wg sync.WaitGroup

	numGoroutines := 20 // Set the desired number of goroutines
	numPoints := len(points)

	concurrency := (numPoints + numGoroutines - 1) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startIndex int) {
			defer wg.Done()
			endIndex := startIndex + concurrency
			if endIndex > numPoints {
				endIndex = numPoints
			}

			for j := startIndex; j < endIndex; j++ {
				points[j] = points[j].Mul(getHeightAtPoint(points[j], shape, craters))
			}
		}(i * concurrency)
	}

	wg.Wait()
}

func getHeightAtPoint(point mgl32.Vec3, shape *PlanetShape, craters []Crater) float32 {
	continentHeight := detailedNoise(point, shape.continentAmplitude, shape.continentFrequency*shape.frequency)

	if continentHeight < 0.0 {
		continentHeight *= shape.oceandepth
	}
	continentHeight = smoothMax(continentHeight, -shape.oceanFloorDepth, shape.oceanSmoothness)

	mountainMask := smoothMax(1e-6, detailedNoise(point, shape.mountainMaskAmplitude, shape.mountainFrequency*shape.frequency*1.1)+shape.mountainMaskOffset, shape.mountainMaskSmoothness)
	mountainHeight := smoothMax(0, ridgidNoise(point, shape.mountainAmplitude, shape.mountainFrequency*shape.frequency), shape.mountainSmoothness)
	mountainHeight = smoothMin(mountainMask, mountainHeight, 0)

	craterHeight := getCraterHeight(point, shape, craters)

	return 1.0 + ((continentHeight+mountainHeight)/3.0+craterHeight)*shape.amplitude
}

func simpleNoise(point mgl32.Vec3, amplitude, frequency float32) float32 {
	x, y, z := point.X()*frequency, point.Y()*frequency, point.Z()*frequency
	return Snoise(x, y, z) * amplitude
}

func detailedNoise(point mgl32.Vec3, amplitude, frequency float32) float32 {
	noiseHeight := float32(0.0)

	for i := 0; i < 5; i++ {
		x, y, z := point.X()*frequency, point.Y()*frequency, point.Z()*frequency
		noiseHeight += Snoise(x, y, z) * float32(amplitude)
		frequency *= 2.0
		amplitude *= 0.5
	}

	return noiseHeight
}

func ridgidNoise(point mgl32.Vec3, amplitude, frequency float32) float32 {
	return amplitude*0.5 - float32(math.Abs(float64(detailedNoise(point, amplitude, frequency))))
}

func getCraterHeight(point mgl32.Vec3, shape *PlanetShape, craters []Crater) float32 {
	craterHeight := float32(0.0)

	for i := 0; i < len(craters); i++ {
		x := mgl32.Vec3.Len(point.Sub(craters[i].position)) / craters[i].radius

		cavity := x*x - 1.0
		rimX := float32(math.Min(float64(x-1.0-shape.craterRimWidth), 0))
		rim := shape.craterRimSteepness * rimX * rimX

		craterShape := smoothMax(cavity, shape.craterFloorHeight, shape.craterSmoothness)
		craterShape = smoothMin(craterShape, float32(rim), shape.craterSmoothness)
		craterHeight += craterShape * craters[i].radius
	}

	return craterHeight
}

func genCraters(numCraters uint32) []Crater {
	craters := make([]Crater, numCraters)

	for i := 0; i < len(craters); i++ {
		position := randomPointOnSphere()
		radius := float32(math.Pow(rand.Float64(), 2) * 0.25)
		craters[i] = Crater{position, radius}
	}
	return craters
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
