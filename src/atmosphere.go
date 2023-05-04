package main

import "github.com/go-gl/gl/v4.1-core/gl"

type PostProcessingFrame struct {
	va VertexArray
	fb FrameBuffer
	rb RenderBuffer
	ib IndexBuffer

	shader Shader
}

/*
Initialize a framebuffer with an attached renderbuffer with the dimensions wxh, and an associated
shader program with the filename shaderPath.
*/
func NewPostProcessingFrame(w uint32, h uint32, shaderPath string) PostProcessingFrame {
	// Vertices and indices for the postprocessing rectangle
	var vertices = []float32{
		1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, 0.0, 1.0,
	}

	var indices = []uint32{
		3, 2, 1,
		3, 1, 0,
	}

	// Create VAO and VBO for rectangle covering the screen
	vb := NewVertexBuffer(vertices)
	ib := NewIndexBuffer(indices)
	vb.bind()

	va := NewVertexArray([]int{2, 2})
	va.bind()
	shader := NewShader(shaderPath)

	// Create framebuffer
	fb := NewFrameBuffer(w, h)
	fb.addColorTexture(1, w, h, gl.COLOR_ATTACHMENT0)

	// Create renderbuffer and attach to the framebuffer
	rb := NewRenderBuffer(w, h)
	fb.addRenderBuffer(rb.id)

	ppf := PostProcessingFrame{va, fb, rb, ib, shader}
	ppf.shader.bind()
	ppf.shader.setUniform1i("colorTexture", 1)

	return ppf
}

// Bind necessary buffers and shader programs and render the post processing effects.
func (ppf *PostProcessingFrame) draw() {
	ppf.shader.bind()
	ppf.va.bind()
	ppf.ib.bind()

	gl.DrawElements(gl.TRIANGLES, ppf.ib.count, gl.UNSIGNED_INT, gl.PtrOffset(0))

	ppf.shader.unbind()
	ppf.va.unbind()
}
