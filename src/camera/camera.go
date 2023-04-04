package camera

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position    mgl32.Vec3
	Orientation mgl32.Vec3
	Up          mgl32.Vec3

	firstClick bool

	width  int
	height int

	speed       float32
	sensitivity float32

	savedMouseX     float64
	savedMouseY     float64
	lastFrameMouseX float64
	lastFrameMouseY float64
}

// Constructor
func NewCamera(width int, height int, position mgl32.Vec3) Camera {
	c := Camera{position, mgl32.Vec3{0.0, 0.0, -1.0}, mgl32.Vec3{0.0, 1.0, 0.0}, true, width, height, 0.1, 30, 0, 0, 0, 0}

	return c
}

// Tell the cameara ho
func (c *Camera) Matrix(FOVdeg float32, nearPlane float32, farPlane float32, shader *uint32, uniform string) {
	view := mgl32.Ident4()
	projection := mgl32.Ident4()
	center := mgl32.Vec3{}
	center = center.Add(c.Position)
	center = center.Add(c.Orientation)

	view = mgl32.LookAtV(c.Position, center, c.Up)
	projection = mgl32.Perspective(mgl32.DegToRad(FOVdeg), float32(c.width/c.height), nearPlane, farPlane)

	projview := projection
	projview = projview.Mul4(view)

	gl.UniformMatrix4fv(gl.GetUniformLocation(*shader, gl.Str(uniform)), 1, false, &projview[0])
}

func (c *Camera) Inputs(window *glfw.Window) {

	if window.GetKey(glfw.KeyW) == glfw.Press {
		temp := c.Orientation
		c.Position = c.Position.Add(temp.Mul(c.speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		temp := c.Orientation
		temp = temp.Cross(c.Up).Normalize()
		c.Position = c.Position.Add(temp.Mul(c.speed * -1))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		temp := c.Orientation
		c.Position = c.Position.Add(temp.Mul(c.speed * -1))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		temp := c.Orientation
		temp = temp.Cross(c.Up).Normalize()
		c.Position = c.Position.Add(temp.Mul(c.speed))
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		temp := c.Up
		c.Position = c.Position.Add(temp.Mul(c.speed))
	}
	if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		temp := c.Up
		c.Position = c.Position.Add(temp.Mul(c.speed * -1))
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		c.speed = 0.4
	} else if window.GetKey(glfw.KeyLeftShift) == glfw.Release {
		c.speed = 0.1
	}

	if window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

		//Make the camera not jump when starting to look around
		if c.firstClick {
			mouseX, mouseY := window.GetCursorPos()
			c.savedMouseX = mouseX
			c.savedMouseY = mouseY
			c.lastFrameMouseX = mouseX
			c.lastFrameMouseY = mouseY
			c.firstClick = false
		}

		mouseX, mouseY := window.GetCursorPos()

		fmt.Println(mouseX, mouseY)

		rotX := c.sensitivity * float32((mouseY-c.lastFrameMouseY)/float64(c.height))
		rotY := c.sensitivity * float32((mouseX-c.lastFrameMouseX)/float64(c.height))

		rotationMatrix := mgl32.HomogRotate3D(mgl32.DegToRad(-rotX), c.Orientation.Cross(c.Up).Normalize())
		newOrientation := rotationMatrix.Mul4x1(mgl32.Vec4{c.Orientation.X(), c.Orientation.Y(), c.Orientation.Z(), 1})

		c.Orientation = mgl32.Vec3{c.Orientation.X(), newOrientation.Y(), -1}

		rotationMatrix = mgl32.HomogRotate3D(mgl32.DegToRad(-rotY), c.Up)
		newOrientation = rotationMatrix.Mul4x1(mgl32.Vec4{c.Orientation.X(), c.Orientation.Y(), c.Orientation.Z(), 1})
		c.Orientation = mgl32.Vec3{newOrientation.X(), c.Orientation.Y(), -1}

		//Update new last frames
		c.lastFrameMouseX = mouseX
		c.lastFrameMouseY = mouseY
	} else if window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Release && c.firstClick == false {
		window.SetCursorPos(c.savedMouseX, c.savedMouseY)
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		c.firstClick = true
	}
}
