package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	texture   Texture
	normalMap Texture

	shader Shader

	vb VertexBuffer
	ib IndexBuffer
	va VertexArray
}

/*
NewSprite generates a new sprite model object

Parameters:
- vertices: sprite vertices
- indices: sprite indices
- texturePath: path to texture file from "textures" folder
- normalMapPath: path to normalmap file from "textures" folder
- shaderPath: path to shader file map from "shaders" folder
- textureScale: How often the texture wraps the model surface
- normalMapScale: How often the normal map wraps the model surface

Returns:
- s: the new sprite object

Example usage:

	vertices, indices := GenPlanet(DefaultEarth)
	s := NewSprite(vertices, indices, "spots.png", "normalmap_rocky.png", "planet.shader", 1.0, 2.0)
*/
func NewSprite(vertices []float32, indices []uint32, texturePath, normalMapPath, shaderPath string, textureScale, normalMapScale float32) Sprite {
	s := Sprite{
		NewTexture(texturePath),
		NewTexture(normalMapPath),
		NewShader(shaderPath),
		VertexBuffer{0},
		IndexBuffer{0, 0},
		VertexArray{0},
	}

	// Generate index buffer and vertex array
	s.vb = NewVertexBuffer(vertices)
	s.ib = NewIndexBuffer(indices)

	s.vb.bind()
	s.va = NewVertexArray([]int{3, 3}) // A vertex contains a position(3 floats) and a normal(3 floats)

	// Set constant uniforms once
	s.shader.bind()

	s.shader.setUniform1i("mainTexture", 0)
	s.shader.setUniform1i("normalMap", 1)
	s.shader.setUniform1f("texScale", textureScale)
	s.shader.setUniform1f("nMapScale", normalMapScale)

	s.shader.setUniform3f("lightPos", 0.0, 0.0, 0.0)
	s.shader.setUniform3f("lightColor", 1.0, 1.0, 1.0)

	s.shader.setUniformMat4fv("projection", cam.ProjMatrix())

	s.shader.unbind()

	return s
}

func (s *Sprite) draw(position, rotation mgl32.Vec3, scale float32) {
	// Calculate model matrix as sprites transformation
	model := mgl32.Translate3D(position.X(), position.Y(), position.Z())
	model = model.Mul4(mgl32.HomogRotate3D(float32(rotation.X()), mgl32.Vec3{1, 0, 0}))
	model = model.Mul4(mgl32.HomogRotate3D(float32(rotation.Y()), mgl32.Vec3{0, 1, 0}))
	model = model.Mul4(mgl32.HomogRotate3D(float32(rotation.Z()), mgl32.Vec3{0, 0, 1}))
	model = model.Mul4(mgl32.Scale3D(scale, scale, scale))

	// Get view matrix from camera
	view := cam.ViewMatrix()

	// Projection matrix is already set as it does not change

	s.shader.bind()
	s.texture.bind(0)
	s.normalMap.bind(1)

	s.shader.setUniformMat4fv("model", model)
	s.shader.setUniformMat4fv("view", view)

	s.shader.setUniform3f("camPos", cam.GetPosition().X(), cam.GetPosition().Y(), cam.GetPosition().Z())

	s.va.bind()
	s.ib.bind()

	gl.DrawElements(gl.TRIANGLES, s.ib.count, gl.UNSIGNED_INT, gl.PtrOffset(0))
	//gl.DrawElements(gl.LINES, s.ib.count, gl.UNSIGNED_INT, gl.PtrOffset(0))
	//gl.PointSize(14)
	//gl.DrawElements(gl.POINTS, s.ib.count, gl.UNSIGNED_INT, gl.PtrOffset(0))

	s.va.unbind()
	s.ib.unbind()
	s.shader.unbind()
}
