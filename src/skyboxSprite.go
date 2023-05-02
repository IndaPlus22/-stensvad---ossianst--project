package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type SkyboxSprite struct {
	position mgl32.Vec3
	rotation mgl32.Vec3

	texture CubemapTexture
	shader  Shader

	vb VertexBuffer
	ib IndexBuffer
	va VertexArray
}

func NewSkyboxSprite(vertices []float32, indices []uint32, texturePath string, shaderPath string) SkyboxSprite {
	s := SkyboxSprite{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 0, 0},
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

	s.updateMesh(vertices, indices)

	return s
}

func (s *SkyboxSprite) updateMesh(vertices []float32, indices []uint32) {
	s.vb = NewVertexBuffer(vertices)
	s.ib = NewIndexBuffer(indices)

	s.vb.bind()
	s.va = NewVertexArray([]int{3})
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
