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
		2, 1, 6,
		1, 5, 6,
		7, 4, 0,
		0, 3, 7,
		6, 5, 4,
		4, 7, 6,
		2, 3, 0,
		0, 1, 2,
		5, 1, 0,
		0, 4, 5,
		7, 3, 6,
		2, 6, 3,
	}

	sphereVertices, sphereIndices := planet.GenPlanet(1.0, 16)

	sphere := NewSprite(sphereVertices, sphereIndices, "lighting.shader")

	skybox := NewSkyboxSprite(skyboxVertices, skyboxIndices, "skybox1", "skybox.shader")

	// TEST
	atmosphere := NewPostProcessingFrame(windowWidth, windowHeight, "atmosphere.shader")
	// TEST END

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
		camPos := cam.GetPosition()

		// Sends position of the light source and camera to the shader:
		sphere.shader.bind()
		sphere.shader.setUniform3f("lightPos", float32(math.Cos(t)*5.0), 0.0, float32(math.Sin(t)*5.0))
		sphere.shader.setUniform3f("camPos", camPos.X(), camPos.Y(), camPos.Z())

		// TEST
		camDir := cam.GetOrientation()
		atmosphere.shader.bind()
		atmosphere.shader.setUniform3f("camDir", camDir.X(), camDir.Y(), camDir.Z())
		atmosphere.shader.setUniform3f("camPos", camPos.X(), camPos.Y(), camPos.Z())
		atmosphere.shader.setUniformMat4fv("viewMatrix", cam.ViewMatrix())
		atmosphere.shader.setUniformMat4fv("projMatrix", cam.ProjMatrix())
		fmt.Println(camPos)
		atmosphere.fb.bind()
		// TEST END

		// Draw:
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)
		gl.Enable(gl.CULL_FACE)

		// Draw the skybox LAST
		skybox.draw()
		sphere.draw()

		// TEST
		gl.Disable(gl.DEPTH_TEST)
		atmosphere.fb.unbind()
		atmosphere.draw()
		// TEST END

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
