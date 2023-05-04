package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"
	"stensvad-ossianst-melvinbe-project/src/camera"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800 * 3
const windowHeight = 600 * 3

var cam = camera.NewCamera(windowWidth, windowHeight, mgl32.Vec3{0.0, 0.0, 2.0})

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

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.34, 0.32, 0.45, 1.0)

	atmosphere := NewPostProcessingFrame(windowWidth, windowHeight, "atmosphere.shader")

	p := NewPlanet(1.0, 150, 0)
	p.addMoon(.2, 128, 30, 5, mgl32.Vec3{1, 0, 0}, 2)
	p.addMoon(.4, 128, 100, 10, mgl32.Vec3{1, 1, 0}, 0.5)
	p.moons[1].addMoon(0.1, 128, 10, 1.5, mgl32.Vec3{1, 1, 0}, 3)

	skybox := NewSkybox("skybox1", "skybox.shader")

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
		atmosphere.shader.setUniformMat4fv("viewMatrix", cam.ViewMatrix())
		atmosphere.shader.setUniformMat4fv("projMatrix", cam.ProjMatrix())

		// Send planet properties to post processing shader:
		var planetOrigin mgl32.Vec3 = mgl32.Vec3{0.0, 0.0, 0.0}
		var atmosphereScale float32 = 1.3
		var planetRadius float32 = 1.0
		atmosphere.shader.setUniform3f("planetOrigin", planetOrigin.X(), planetOrigin.Y(), planetOrigin.Z())
		atmosphere.shader.setUniform1f("planetRadius", planetRadius)
		atmosphere.shader.setUniform1f("atmosphereScale", atmosphereScale)

		// Bind the framebuffer for postprocessing before drawing:
		atmosphere.fb.bind()

		// Draw:
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)
		//gl.Enable(gl.CULL_FACE)

		p.draw()

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
