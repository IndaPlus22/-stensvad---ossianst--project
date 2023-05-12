package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Global variables
var windowWidth = 800 * 2
var windowHeight = 600 * 2

var cam = NewCamera(windowWidth, windowHeight, mgl32.Vec3{0.0, 0.0, 5.0})
var planets = []*Planet{}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Planets!", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	fbWidth, fbHeight := window.GetFramebufferSize()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.34, 0.32, 0.45, 1.0)

	// Create planets
	earthSettings := DefaultEarth()
	moonSettings := DefaultMoon()

	sun := NewPlanet(DefaultSun())

	earthSettings.shape.radius = 1.5
	p1 := NewPlanet(earthSettings)

	earthSettings.shape.radius = 1.0
	earthSettings.colors = RandomColors()
	p2 := NewPlanet(earthSettings)

	earthSettings.shape.radius = 0.75
	earthSettings.colors = RandomColors()
	p3 := NewPlanet(earthSettings)

	moonSettings.shape.radius = 0.75
	m1 := NewPlanet(moonSettings)

	moonSettings.shape.radius = 0.5
	moonSettings.colors = RandomColors()
	m2 := NewPlanet(moonSettings)

	moonSettings.shape.radius = 0.3
	moonSettings.colors = RandomColors()
	m3 := NewPlanet(moonSettings)

	// Set orbits of planets
	p1.addOrbital(m1, 6.0, mgl32.Vec3{0.0, 1.0, 0.1}, -1.25)
	p1.addOrbital(m2, 5.0, mgl32.Vec3{0.5, 1.0, 0.0}, 1.5)
	p2.addOrbital(m3, 4.0, mgl32.Vec3{0.0, 1.0, 0.2}, -1.75)

	sun.addOrbital(p1, 12.0, mgl32.Vec3{0.1, 1.0, 0.1}, 0.75)
	sun.addOrbital(p2, 18.0, mgl32.Vec3{0.2, 1.0, 0.0}, -1.1)
	sun.addOrbital(p3, 21.0, mgl32.Vec3{0.0, 1.0, 0.3}, 1.25)

	// Create atmospheres
	// Send planet positions to uniform buffer
	planetPositions := []mgl32.Vec4{}
	for _, planet := range planets {
		// First three are planet coordinates, fourth is planet scale
		if planet.hasAtmosphere {
			p := planet.position
			planetPositions = append(planetPositions, mgl32.Vec4{p.X(), p.Y(), p.Z(), planet.scale})
		}
	}
	atmosphere := NewPostProcessingFrame(uint32(fbWidth), uint32(fbHeight), "atmosphere.shader")
	atmosphere.addUniformBufferVec4("PlanetPositions", planetPositions, int(unsafe.Sizeof(mgl32.Vec4{})*10))

	// Create skybox
	skybox := NewSkybox("skybox2", "skybox.shader")

	for !window.ShouldClose() {
		// Update:
		cam.Inputs(window)
		camPos := cam.GetPosition()

		// Send the world position, direction, projection matrix and view matrix of the camera
		// as well as the position of the light to the atmosphere shader:
		camDir := cam.GetOrientation()
		atmosphere.shader.bind()
		atmosphere.shader.setUniform3f("camDir", camDir.X(), camDir.Y(), camDir.Z())
		atmosphere.shader.setUniform3f("camPos", camPos.X(), camPos.Y(), camPos.Z())
		atmosphere.shader.setUniform3f("lightPos", sun.position.X(), sun.position.Y(), sun.position.Z())
		atmosphere.shader.setUniformMat4fv("viewMatrix", cam.ViewMatrix())
		atmosphere.shader.setUniformMat4fv("projMatrix", cam.ProjMatrix())

		// Send planet properties to post processing shader:
		for i := range planetPositions {
			p := planets[i].position
			planetPositions[i] = mgl32.Vec4{p.X(), p.Y(), p.Z(), planets[i].scale}
		}

		atmosphere.updateUniformBufferVec4(atmosphere.ub[0], planetPositions)

		// Bind the framebuffer for postprocessing before drawing:
		atmosphere.fb.bind()

		// Draw:
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)

		sun.Draw()

		// Draw the skybox LAST
		skybox.draw()

		// Disable depth testing and apply post processing:
		gl.Disable(gl.DEPTH_TEST)
		atmosphere.fb.unbind()
		atmosphere.draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
