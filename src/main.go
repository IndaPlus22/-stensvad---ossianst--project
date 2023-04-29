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

	sphereVertices, sphereIndices := planet.GenPlanet(1.0, 150, 30)

	sphere := NewSprite(sphereVertices, sphereIndices, "square.png", "lighting.shader")

	skybox := NewSkybox("skybox1", "skybox.shader")

	for !window.ShouldClose() {
		// Update:
		cam.Inputs(window)

		// Sends position of the light source and camera to the shader:
		sphere.shader.bind()
		sphere.shader.setUniform3f("lightPos", float32(math.Cos(cam.Time)*5.0), 0.0, float32(math.Sin(cam.Time)*5.0))
		sphere.shader.setUniform3f("camPos", cam.GetPosition().X(), cam.GetPosition().Y(), cam.GetPosition().Z())

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
