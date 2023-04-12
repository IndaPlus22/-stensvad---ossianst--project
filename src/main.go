package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"
	"stensvad-ossianst-melvinbe-project/src/camera"
	"stensvad-ossianst-melvinbe-project/src/planet"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800 * 2
const windowHeight = 600 * 2

var cam = camera.NewCamera(windowWidth, windowHeight, mgl32.Vec3{0.0, 0.0, 3.0})

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
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Planet Generator", nil, nil)
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
	gl.Enable(gl.CULL_FACE)

	var vertices, indices = planet.Gen(50)

	cube := NewSprite(vertices, indices, "square.png", "simple.shader")

	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		// Calculate deltatime
		time := glfw.GetTime()
		deltatime := time - previousTime
		previousTime = time

		// Update:
		cam.Inputs(window)

		cube.rotation = cube.rotation.Add(mgl32.Vec3{0, float32(deltatime * 0.0), 0})

		// Draw:
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		cube.draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
