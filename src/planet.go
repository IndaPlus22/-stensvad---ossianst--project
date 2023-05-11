package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Planet struct {
	sprite Sprite

	position mgl32.Vec3
	rotation mgl32.Vec3
	scale    float32

	axisAroundParent mgl32.Vec3
	orbital          []*Planet
	orbitTime        float64

	hasAtmosphere bool
}

/*
NewPlanet generates a new planet struct and returns it

Parameters:
- settings: the planet settings struct containing a recipe for what the planet will be

Returns:
- p: the new planet object

Example usage:

	earthSettings := DefaultEarth()
	planet := NewPlanet(earthSettings)
*/
func NewPlanet(settings PlanetSettings) *Planet {
	// Generate the planet sprite model
	planetVertices, planetIndices := GenPlanet(settings.shape)

	sprite := NewSprite(
		planetVertices,
		planetIndices,
		settings.texturePath,
		settings.normalMapPath,
		settings.shaderPath,
		settings.textureScale,
		settings.normalMapScale,
	)

	p := &Planet{
		sprite,

		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 0.0},
		settings.shape.radius,

		mgl32.Vec3{},
		nil,
		0,
		settings.hasAtmosphere,
	}

	p.setColors(settings.colors)
	planets = append(planets, p)

	return p
}

// Binds the color uniforms of the "planet.shader" shader to the planets colors
func (p *Planet) setColors(c PlanetColors) {
	p.sprite.shader.bind()

	p.sprite.shader.setUniform3f("shoreColLow", c.shoreColLow.X(), c.shoreColLow.Y(), c.shoreColLow.Z())
	p.sprite.shader.setUniform3f("shoreColHigh", c.shoreColHigh.X(), c.shoreColHigh.Y(), c.shoreColHigh.Z())
	p.sprite.shader.setUniform3f("flatColLow", c.flatColLow.X(), c.flatColLow.Y(), c.flatColLow.Z())
	p.sprite.shader.setUniform3f("flatColHigh", c.flatColHigh.X(), c.flatColHigh.Y(), c.flatColHigh.Z())
	p.sprite.shader.setUniform3f("steepColLow", c.steepColLow.X(), c.steepColLow.Y(), c.steepColLow.Z())
	p.sprite.shader.setUniform3f("steepColHigh", c.steepColHigh.X(), c.steepColHigh.Y(), c.steepColHigh.Z())
	p.sprite.shader.setUniform3f("waterCol", c.waterCol.X(), c.waterCol.Y(), c.waterCol.Z())
}

// Add an orbital to this planet
func (p *Planet) addOrbital(planet *Planet, distance float32, axis mgl32.Vec3, timeToOrbit float64) {
	planet.axisAroundParent = axis.Normalize()
	planet.orbitTime = timeToOrbit
	planet.position = planet.axisAroundParent.Cross(mgl32.Vec3{1, 1, 1}).Normalize().Mul(distance)
	p.orbital = append(p.orbital, planet)
}

// Draws planet and its orbitals
func (p *Planet) Draw() {
	p.sprite.draw(p.position, p.rotation, p.scale)

	p.rotation = mgl32.Vec3{0, float32(cam.TimeTot), 0}

	// Draw and rotate all orbitals around this planet
	for i := range p.orbital {
		rotM := mgl32.HomogRotate3D(float32(cam.TimeDiff*p.orbital[i].orbitTime), p.orbital[i].axisAroundParent)
		p.orbital[i].position = rotM.Mul4x1(p.orbital[i].position.Vec4(0)).Vec3().Add(p.position)

		p.orbital[i].Draw()

		p.orbital[i].position = p.orbital[i].position.Sub(p.position)
	}
}
