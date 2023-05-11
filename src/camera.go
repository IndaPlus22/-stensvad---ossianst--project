package main

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position    mgl32.Vec3
	orientation mgl32.Vec3
	up          mgl32.Vec3

	firstClick bool

	width  int
	height int

	speed     float32
	fovDEG    float32
	nearPlane float32
	farPlane  float32

	sensitivity     float64
	yaw             float64
	pitch           float64
	lastFrameMouseX float64
	lastFrameMouseY float64
	timeFactor      float64
	TimeDiff        float64
	TimeTot         float64
}

var previousTime = 0.0

/*
NewCamera generates a new camera and retruns it.

Parameters:
- width: the width of the window
- height: the width of the window
- posistion: the position where the camera should be placed in world space

Returns:
- camera: a Camera Struct

Example usage:

	cam := GenCamera(1080, 720, mgl32.Vec3{0, 0, 0})
*/
func NewCamera(width int, height int, position mgl32.Vec3) Camera {
	c := Camera{
		position,
		mgl32.Vec3{0.0, 0.0, -1.0},
		mgl32.Vec3{0.0, 1.0, 0.0},

		true,

		width, height,

		0.05, 45.0, 0.001, 100,

		0.1, -90, 0, 0, 0, 0.1, 0.0, 0.0}

	return c
}

// Generate a new view matrix based on position and rotation
func (c *Camera) ViewMatrix() mgl32.Mat4 {
	center := c.position
	center = center.Add(c.orientation)

	view := mgl32.LookAtV(c.position, center, c.up)

	return view
}

// Generate a new projection matrix based on camera settings
func (c *Camera) ProjMatrix() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(c.fovDEG), float32(c.width)/float32(c.height), c.nearPlane, c.farPlane)
}

/*
Returns the position of the camera in world space
*/
func (c *Camera) GetPosition() mgl32.Vec3 {
	return c.position
}

func (c *Camera) GetOrientation() mgl32.Vec3 {
	return c.orientation
}

func (c *Camera) GetNearPlane() float32 {
	return c.nearPlane
}

func (c *Camera) GetFarPlane() float32 {
	return c.farPlane
}

// Takes inputs from the user allowing them to controll the camera
func (c *Camera) Inputs(window *glfw.Window) {

	//Update time
	time := glfw.GetTime()
	deltatime := time - previousTime
	previousTime = time

	c.TimeDiff = deltatime * c.timeFactor
	c.TimeTot += c.TimeDiff

	//Positioning of the camera
	if window.GetKey(glfw.KeyW) == glfw.Press {
		temp := c.orientation
		c.position = c.position.Add(temp.Mul(c.speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		temp := c.orientation
		temp = temp.Cross(c.up).Normalize()
		c.position = c.position.Add(temp.Mul(c.speed * -1))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		temp := c.orientation
		c.position = c.position.Add(temp.Mul(c.speed * -1))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		temp := c.orientation
		temp = temp.Cross(c.up).Normalize()
		c.position = c.position.Add(temp.Mul(c.speed))
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		temp := c.up
		c.position = c.position.Add(temp.Mul(c.speed))
	}
	if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		temp := c.up
		c.position = c.position.Add(temp.Mul(c.speed * -1))
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		c.speed = 0.4
	} else if window.GetKey(glfw.KeyLeftShift) == glfw.Release {
		c.speed = 0.03
	}

	//Makes it possible to control what direction the camera is looking
	if window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

		mouseX, mouseY := window.GetCursorPos()

		//Make the camera not jump when starting to look around
		if c.firstClick {
			c.lastFrameMouseX = mouseX
			c.lastFrameMouseY = mouseY
			c.firstClick = false
		}

		xOffset := mouseX - c.lastFrameMouseX
		yOffset := mouseY - c.lastFrameMouseY
		c.lastFrameMouseX = mouseX
		c.lastFrameMouseY = mouseY

		xOffset *= c.sensitivity
		yOffset *= c.sensitivity

		c.yaw += xOffset
		c.pitch += yOffset

		//Stops the user from being able to rotate fullt upwards and downwards
		if c.pitch > 89 {
			c.pitch = 89
		} else if c.pitch < -89 {
			c.pitch = -89
		}

		//Calculate the new orientation of the camera
		frontX := float32(math.Cos(float64(mgl32.DegToRad(float32(c.yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(c.pitch)))))
		frontY := float32(math.Sin(float64(mgl32.DegToRad(float32(-c.pitch)))))
		frontZ := float32(math.Sin(float64(mgl32.DegToRad(float32(c.yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(c.pitch)))))
		front := mgl32.Vec3{frontX, frontY, frontZ}

		c.orientation = front.Normalize()

	} else if window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Release && c.firstClick == false {
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		c.firstClick = true
	}

	// Time-control, keypresses controls the speed of time
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		if c.timeFactor < 0.005 && c.timeFactor > 0 {
			c.timeFactor = 0
		}
		if c.timeFactor > 0 {
			c.timeFactor *= 0.99
		} else if c.timeFactor < 0 {
			c.timeFactor *= 1.01
		} else {
			c.timeFactor = -0.01
		}
	}

	if window.GetKey(glfw.KeyRight) == glfw.Press {
		if c.timeFactor > -0.005 && c.timeFactor < 0 {
			c.timeFactor = 0
		}
		if c.timeFactor < 0 {
			c.timeFactor *= 0.99
		} else if c.timeFactor > 0 {
			c.timeFactor *= 1.01
		} else {
			c.timeFactor = 0.01
		}
	}

	if window.GetKey(glfw.KeyM) == glfw.Press {
		//Update the timefactor if a number key is pressed to the negative of the pressed number
		for i := 0; i < 10; i++ {
			if window.GetKey(glfw.Key0+glfw.Key(i)) == glfw.Press {
				c.timeFactor = -float64(i)
			}
		}
	} else {
		//Update the timefactor if a number key is pressed to the pressed number
		for i := 0; i < 10; i++ {
			if window.GetKey(glfw.Key0+glfw.Key(i)) == glfw.Press {
				c.timeFactor = float64(i)
			}
		}
	}
}
