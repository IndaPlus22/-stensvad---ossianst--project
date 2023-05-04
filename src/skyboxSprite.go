package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Skybox struct {
	texture CubemapTexture
	shader  Shader

	vb VertexBuffer
	ib IndexBuffer
	va VertexArray
}

func NewSkybox(texturePath string, shaderPath string) Skybox {
	var vertices = []float32{
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

	var indices = []uint32{
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

	s := Skybox{
		NewCubemapTexture(texturePath),
		NewShader(shaderPath),
		VertexBuffer{0},
		IndexBuffer{0, 0},
		VertexArray{0},
	}
	s.shader.bind()
	projection := cam.ProjMatrix()
	view := cam.ViewMatrix().Mat3().Mat4()

	s.shader.setUniformMat4fv("view", view)
	s.shader.setUniformMat4fv("projection", projection)

	s.shader.unbind()

	s.vb = NewVertexBuffer(vertices)
	s.ib = NewIndexBuffer(indices)

	s.vb.bind()
	s.va = NewVertexArray([]int{3})

	return s
}

func (s *SkyboxSprite) draw() {

	s.texture.bind(0)

	view := cam.ViewMatrix().Mat3().Mat4()

	projection := cam.ProjMatrix()

	s.shader.bind()

	s.shader.setUniformMat4fv("view", view)
	s.shader.setUniformMat4fv("projection", projection)

	s.va.bind()
	s.ib.bind()

	gl.DrawElements(gl.TRIANGLES, s.ib.count, gl.UNSIGNED_INT, gl.PtrOffset(0))

	s.va.unbind()
	s.ib.unbind()
	s.shader.unbind()

	gl.DepthFunc(gl.LESS)
}
