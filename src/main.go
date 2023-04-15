package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"
	"runtime"
	"stensvad-ossianst-melvinbe-project/src/camera"
	"stensvad-ossianst-melvinbe-project/src/planet"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800 * 2
const windowHeight = 600 * 2

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
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
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

	var skyboxVertices = []float32{
		// Positions
		-1, -1, 1,
		1, -1, 1,
		1, -1, -1,
		-1, -1, -1,
		-1, 1, 1,
		1, 1, 1,
		1, 1, -1,
		-1, 1, -1,
	}

	var skyboxIndices = []uint32{
		1, 2, 6,
		6, 5, 1,
		0, 4, 7,
		7, 3, 0,
		4, 5, 6,
		6, 7, 4,
		0, 3, 2,
		2, 1, 0,
		0, 1, 5,
		5, 4, 0,
		3, 7, 6,
		6, 2, 3,
	}

	sphereVertices, sphereIndices := planet.GenPlanet(1.0, 8)

	sphere := NewSprite(sphereVertices, sphereIndices, "lighting.shader")

	skybox := NewSkyboxSprite(skyboxVertices, skyboxIndices, "skybox1", "skybox.shader")

	previousTime := glfw.GetTime()

	t := 0.0

	for !window.ShouldClose() {
		// Calculate deltatime
		time := glfw.GetTime()
		deltatime := time - previousTime
		previousTime = time

		t += deltatime

		// Update:
		cam.Inputs(window)

		sphere.shader.bind()
		sphere.shader.setUniform3f("lightPos", float32(math.Cos(t)*5.0), 0.0, float32(math.Sin(t)*5.0))

		// Draw:
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		sphere.draw()

		// Draw the skybox LAST
		skybox.draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
