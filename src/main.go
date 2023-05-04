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

	p := NewPlanet(1.0, 150, 0)
	p.addMoon(.2, 128, 30, 5, mgl32.Vec3{1, 0, 0}, 2)
	p.addMoon(.4, 128, 100, 10, mgl32.Vec3{1, 1, 0}, 0.5)
	p.moons[1].addMoon(0.1, 128, 10, 1.5, mgl32.Vec3{1, 1, 0}, 3)

	skybox := NewSkybox("skybox1", "skybox.shader")

	for !window.ShouldClose() {
		// Update:
		cam.Inputs(window)

		// Draw:
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		p.draw()

		// Draw the skybox LAST
		skybox.draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
