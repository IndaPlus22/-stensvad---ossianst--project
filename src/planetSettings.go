package main

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

type PlanetSettings struct {
	shape  PlanetShape
	colors PlanetColors

	hasAtmosphere bool
	hasOcean      bool

	texturePath   string
	normalMapPath string
	shaderPath    string

	textureScale   float32
	normalMapScale float32
}

type PlanetShape struct {
	radius    float32
	res       uint32
	amplitude float32
	frequency float32

	oceandepth      float32
	oceanFloorDepth float32
	oceanSmoothness float32

	continentAmplitude float32
	continentFrequency float32

	mountainAmplitude  float32
	mountainFrequency  float32
	mountainSmoothness float32

	mountainMaskAmplitude  float32
	mountainMaskSmoothness float32
	mountainMaskOffset     float32

	numCraters         uint32
	craterRimWidth     float32
	craterRimSteepness float32
	craterSmoothness   float32
	craterFloorHeight  float32
}

type PlanetColors struct {
	shoreColLow  mgl32.Vec3
	shoreColHigh mgl32.Vec3
	flatColLow   mgl32.Vec3
	flatColHigh  mgl32.Vec3
	steepColLow  mgl32.Vec3
	steepColHigh mgl32.Vec3
	waterCol     mgl32.Vec3
}

func DefaultEarth() PlanetSettings {
	return PlanetSettings{
		PlanetShape{
			// General:
			1.0, // radius
			200, // resolution
			1.0, // amplitude
			1.0, // frequency

			// Ocean:
			7.0, // depth
			0.4, // floor depth
			0.8, // smoothness

			// Continent:
			0.15, // amplitude
			1.0,  // frequency

			// Mountain:
			1.5,  // amplitude
			0.75, // frequency
			0.5,  // smoothness

			// Mountain Mask:
			1.1,  // amplitude
			0.4,  // smoothness
			-0.5, // offset

			// Crater:
			0,    // count
			0.7,  // rim width
			0.4,  // rim steepness
			0.3,  // smoothness
			-0.3, // floor height
		},

		PlanetColors{
			mgl32.Vec3{0.98, 0.90, 0.62},
			mgl32.Vec3{0.95, 0.79, 0.41},
			mgl32.Vec3{0.49, 0.62, 0.00},
			mgl32.Vec3{0.26, 0.41, 0.00},
			mgl32.Vec3{0.42, 0.32, 0.24},
			mgl32.Vec3{0.90, 0.90, 0.95},
			mgl32.Vec3{0.50, 0.50, 0.90},
		},

		true, // has atmosphere
		true, // has oceans

		"spots.png",           // texture
		"normalmap_rocky.png", // normal map
		"planet.shader",       // shader

		1.0, // texture scale
		3.0, // normal map scale
	}
}

func DefaultMoon() PlanetSettings {
	return PlanetSettings{
		PlanetShape{
			// General:
			1.0, // radius
			100, // resolution
			1.0, // amplitude
			1.0, // frequency

			// Ocean:
			1.0,  // depth
			10.0, // floor depth
			0.0,  // smoothness

			// Continent:
			0.15, // amplitude
			1.0,  // frequency

			// Mountain:
			1.1, // amplitude
			0.1, // frequency
			0.1, // smoothness

			// Mountain Mask:
			1.1,  // amplitude
			0.1,  // smoothness
			-0.1, // offset

			// Crater:
			40,   // count
			0.7,  // rim width
			0.4,  // rim steepness
			0.3,  // smoothness
			-0.3, // floor height
		},

		PlanetColors{
			mgl32.Vec3{0.60, 0.60, 0.60},
			mgl32.Vec3{0.60, 0.60, 0.60},
			mgl32.Vec3{0.60, 0.60, 0.60},
			mgl32.Vec3{0.60, 0.60, 0.60},
			mgl32.Vec3{0.60, 0.60, 0.60},
			mgl32.Vec3{0.60, 0.60, 0.60},
			mgl32.Vec3{0.00, 0.00, 0.00},
		},

		false, // has atmosphere
		false, // has oceans

		"spots.png",             // texture
		"normalmap_craters.png", // normal map
		"planet.shader",         // shader

		3.0, // texture scale
		0.5, // normal map scale
	}
}

func DefaultSun() PlanetSettings {
	return PlanetSettings{
		PlanetShape{
			// General:
			2.0, // radius
			50,  // resolution
			0.0, // amplitude
			0.0, // frequency

			// Ocean:
			1.0, // depth
			0.0, // floor depth
			0.0, // smoothness

			// Continent:
			0.0, // amplitude
			0.0, // frequency

			// Mountain:
			0.0, // amplitude
			0.0, // frequency
			0.0, // smoothness

			// Mountain Mask:
			0.0, // amplitude
			0.0, // smoothness
			0.0, // offset

			// Crater:
			0,   // count
			0.0, // rim width
			0.0, // rim steepness
			0.0, // smoothness
			0.0, // floor height
		},

		PlanetColors{},

		true,  // has atmosphere
		false, // has oceans

		"sun.png",    // texture
		"sun.png",    // normal map
		"sun.shader", // shader

		1.2, // texture scale
		0.0, // normal map scale
	}
}

func RandomColors() PlanetColors {
	// Generate random base colors
	shoreCol := mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
	flatCol := mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
	steepCol := mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
	waterCol := mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}

	// Set colors with offsets for additional similar colors
	return PlanetColors{
		shoreCol,
		shoreCol.Add(mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}.Mul(0.4)),
		flatCol,
		flatCol.Add(mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}.Mul(0.4)),
		steepCol,
		steepCol.Add(mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}.Mul(0.4)),
		waterCol,
	}
}
