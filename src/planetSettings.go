package main

import (
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

func EarthSettings() PlanetSettings {
	return PlanetSettings{
		PlanetShape{
			// General:
			1.0, // radius
			100, // resolution
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

		true,
		true,

		"spots.png",
		"normalmap_rocky.png",
		"planet.shader",

		1.0,
		3.0,
	}
}

func MoonSettings() PlanetSettings {
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

		true,
		true,

		"spots.png",
		"normalmap_craters.png",
		"planet.shader",

		3.0,
		2.0,
	}
}

func RandomColors() PlanetColors {
	return PlanetColors{
		mgl32.Vec3{0.98, 0.90, 0.62},
		mgl32.Vec3{0.95, 0.79, 0.41},
		mgl32.Vec3{0.49, 0.62, 0.00},
		mgl32.Vec3{0.26, 0.41, 0.00},
		mgl32.Vec3{0.42, 0.32, 0.24},
		mgl32.Vec3{0.90, 0.90, 0.95},
		mgl32.Vec3{0.50, 0.50, 0.90},
	}
}
