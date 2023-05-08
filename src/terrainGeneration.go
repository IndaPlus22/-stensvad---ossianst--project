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

/*
GenTerrain generates the points of a planet as described in a given planet shape struct

Parameters:
- points: the planet points of sphere before fancy terrain generation
- shape: the planet shape struct containing a recipe for the planets shape

Example usage:

	// Generate points of sphere first
	points, indices := genOctahedron(100)
	normalizePointDistances(points)

	earthSettings := DefaultEarth()

	GenTerrain(points, earthSettings.shape)
	// The sphere is now a planet
*/
func GenTerrain(points []mgl32.Vec3, shape PlanetShape) {
	craters := genCraters(shape.numCraters)

	// Generate a random seed for every planet
	rand.NewSource(time.Now().UnixNano())
	seed = rand.Float32() * 1.0e5

	var wg sync.WaitGroup

	numGoroutines := 20 // How to parallelize
	numPoints := len(points)

	// Amount of points to calculate per goroutine
	concurrency := (numPoints + numGoroutines - 1) / numGoroutines

	// Calculate points concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startIndex int) {
			defer wg.Done()
			endIndex := startIndex + concurrency
			if endIndex > numPoints {
				endIndex = numPoints
			}

			for j := startIndex; j < endIndex; j++ {
				points[j] = points[j].Mul(getHeightAtPoint(points[j], &shape, craters))
			}
		}(i * concurrency)
	}

	wg.Wait()
}

// Calculate the height of a single point
func getHeightAtPoint(point mgl32.Vec3, shape *PlanetShape, craters []Crater) float32 {
	// Generate the general bumpyness of the planet surface and locations of the oceans
	continentHeight := detailedNoise(point, shape.continentAmplitude, shape.continentFrequency*shape.frequency)

	// Deepen the deep areas of the surface to form oceans
	if continentHeight < 0.0 {
		continentHeight *= shape.oceandepth
	}
	// Raise the deepest areas to to ocean floor
	continentHeight = smoothMax(continentHeight, -shape.oceanFloorDepth, shape.oceanSmoothness)

	// Generate a mask for the mountains to keep some areas free from mountains
	mountainMask := smoothMax(1e-6, detailedNoise(point, shape.mountainMaskAmplitude, shape.mountainFrequency*shape.frequency*1.1)+shape.mountainMaskOffset, shape.mountainMaskSmoothness)
	// Generate the actual mountains
	mountainHeight := smoothMax(0, ridgeNoise(point, shape.mountainAmplitude, shape.mountainFrequency*shape.frequency), shape.mountainSmoothness)
	// Limit the mountains to stay within the mask
	mountainHeight = smoothMin(mountainMask, mountainHeight, 0)

	// Add craters
	craterHeight := getCraterHeight(point, shape, craters)

	// Add everything together and scale by general amplitude
	return 1.0 + ((continentHeight+mountainHeight)/3.0+craterHeight)*shape.amplitude
}

// Calls the Snoise function with a specified amplitude and freqency
func simpleNoise(point mgl32.Vec3, amplitude, frequency float32) float32 {
	x, y, z := point.X()*frequency, point.Y()*frequency, point.Z()*frequency
	return Snoise(x, y, z) * amplitude
}

// Repeatadly calls the Snoise function with decreasing amplitude amplitude and increasing freqency
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

// Same as detailed noise but uses negative absolute values to form sharp edges
func ridgeNoise(point mgl32.Vec3, amplitude, frequency float32) float32 {
	return amplitude*0.5 - float32(math.Abs(float64(detailedNoise(point, amplitude, frequency))))
}

// Add together the effect of every crater on a specified point
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

	return craterHeight * shape.amplitude
}

// Randomly generates the positions and radi of every crater
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

// Like the min function, but smooth
func smoothMin(a, b, k float32) float32 {
	h := (b - a + k) / (2 * k)
	h = float32(math.Min(math.Max(float64(h), 0.0), 1.0))
	return a*h + b*(1-h) - k*h*(1-h)
}

// Like the max function, but smooth
func smoothMax(a, b, k float32) float32 {
	return smoothMin(a, b, -k)
}
