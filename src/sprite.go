package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	position mgl32.Vec3
	rotation mgl32.Vec3

	texture Texture
	shader  Shader

	vb VertexBuffer
	ib IndexBuffer
	va VertexArray
}

func NewSprite(vertices []float32, indices []uint32, texturePath string, shaderPath string) Sprite {
	s := Sprite{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 0, 0},
		NewTexture(texturePath),
		NewShader(shaderPath),
		VertexBuffer{0},
		IndexBuffer{0, 0},
		VertexArray{0},
	}

	s.shader.bind()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	s.shader.setUniformMat4fv("projection", projection)
	s.shader.setUniform1i("tex", 0)
	s.shader.unbind()

	s.updateMesh(vertices, indices)

	return s
}

func (s *Sprite) updateMesh(vertices []float32, indices []uint32) {
	s.vb = NewVertexBuffer(vertices)
	s.ib = NewIndexBuffer(indices)

	s.vb.bind()
	s.va = NewVertexArray([]int{3, 2, 3})
}

func (s *Sprite) draw() {
	s.texture.bind(0)

	model := mgl32.Translate3D(s.position.X(), s.position.Y(), s.position.Z())
	model = model.Mul4(mgl32.HomogRotate3D(float32(s.rotation.X()), mgl32.Vec3{1, 0, 0}))
	model = model.Mul4(mgl32.HomogRotate3D(float32(s.rotation.Y()), mgl32.Vec3{0, 1, 0}))
	model = model.Mul4(mgl32.HomogRotate3D(float32(s.rotation.Z()), mgl32.Vec3{0, 0, 1}))

	view := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	s.shader.bind()

	s.shader.setUniformMat4fv("model", model)
	s.shader.setUniformMat4fv("camera", view)

	s.va.bind()
	s.ib.bind()

	gl.DrawElements(gl.TRIANGLES, s.ib.count, gl.UNSIGNED_INT, gl.PtrOffset(0))

	s.va.unbind()
	s.ib.unbind()
	s.shader.unbind()
}
