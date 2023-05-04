package main

import (
	"stensvad-ossianst-melvinbe-project/src/planet"

	"github.com/go-gl/mathgl/mgl32"
)

type Planet struct {
	PlanetSprite     Sprite
	axisAroundParent mgl32.Vec3
	moons            []Planet
	orbitTime        float64
}

/*
NewPlanet generates a new planet and retruns it.

Parameters:
- radius: the radius of the planet
- res: the resolution of the planet (number of vertices in circumference / 4)
- numCraters: the amount of craters on the planet

Returns:
- planet: a Planet Struct

Example usage:

	p := GenPLanet(20, 128, 30)
*/
func NewPlanet(radius float32, res uint32, numCraters uint32) Planet {
	planetVertices, planetIndices := planet.GenPlanet(radius, res, numCraters)

	p := Planet{NewSprite(planetVertices, planetIndices, "square.png", "lighting.shader"), mgl32.Vec3{}, nil, 0}

	p.PlanetSprite.shader.bind()
	p.PlanetSprite.shader.setUniform3f("lightPos", float32(5.0), 0.0, float32(5.0))
	return p
}

// Draws the planet that calls this function and all its moon with it
func (p *Planet) draw() {
	p.PlanetSprite.rotation = mgl32.Vec3{0, float32(cam.TimeTot), 0}
	p.PlanetSprite.shader.bind()
	p.PlanetSprite.shader.setUniform3f("camPos", cam.GetPosition().X(), cam.GetPosition().Y(), cam.GetPosition().Z())
	p.PlanetSprite.draw()

	//Draw and rotate all the moons around this planet
	for i := range p.moons {
		rotM := mgl32.HomogRotate3D(float32(cam.TimeDiff*p.moons[i].orbitTime), p.moons[i].axisAroundParent)
		p.moons[i].PlanetSprite.position = rotM.Mul4x1(p.moons[i].PlanetSprite.position.Vec4(0)).Vec3().Add(p.PlanetSprite.position)

		p.moons[i].draw()
		p.moons[i].PlanetSprite.position = p.moons[i].PlanetSprite.position.Sub(p.PlanetSprite.position)
	}
}

// Adds a moon around the planet that calls this function
func (p *Planet) addMoon(radius float32, res uint32, numCraters uint32, distance float32, axis mgl32.Vec3, timeToOrbit float64) {
	m := NewPlanet(radius, res, numCraters)
	m.axisAroundParent = axis.Normalize()
	m.orbitTime = timeToOrbit
	m.PlanetSprite.position = m.axisAroundParent.Cross(mgl32.Vec3{1, 1, 1}).Normalize().Mul(distance)
	p.moons = append(p.moons, m)
}
