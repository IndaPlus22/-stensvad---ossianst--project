package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type PostProcessingFrame struct {
	va VertexArray
	fb FrameBuffer

	ib IndexBuffer
	ub []uint32

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
	fb.addColorTexture(2, w, h, gl.COLOR_ATTACHMENT0)
	fb.addDepthTexture(3, w, h)

	// Create renderbuffer and attach to the framebuffer
	//rb := NewRenderBuffer(w, h)
	//fb.addRenderBuffer(rb.id)

	ppf := PostProcessingFrame{va, fb, ib, []uint32{}, shader}
	ppf.shader.bind()
	ppf.shader.setUniform1i("colorTexture", 2)
	ppf.shader.setUniform1i("depthTexture", 3)
	ppf.shader.setUniform1f("camNear", cam.GetNearPlane())
	ppf.shader.setUniform1f("camFar", cam.GetFarPlane())

	return ppf
}

func (ppf *PostProcessingFrame) addUniformBufferVec3(block string, vectors []mgl32.Vec4) {
	var id uint32
	gl.GenBuffers(1, &id)
	gl.BindBuffer(gl.UNIFORM_BUFFER, id)
	// potentiellt fel
	gl.BufferData(gl.UNIFORM_BUFFER, int(unsafe.Sizeof(mgl32.Vec4{})*3), nil, gl.DYNAMIC_DRAW)
	// potentiellt fel
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, int(unsafe.Sizeof(mgl32.Vec4{}))*len(vectors), unsafe.Pointer(&vectors[0]))
	// potentiellt fel
	blockId := gl.GetUniformBlockIndex(ppf.shader.id, gl.Str(block+"\x00"))
	gl.UniformBlockBinding(ppf.shader.id, blockId, 0)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 0, id)
	ppf.ub = append(ppf.ub, id)
}

func (ppf *PostProcessingFrame) updateUB(vectors []mgl32.Vec4) {
	for _, ubo := range ppf.ub {
		gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
		gl.BufferSubData(gl.UNIFORM_BUFFER, 0, int(unsafe.Sizeof(mgl32.Vec4{}))*len(vectors), unsafe.Pointer(&vectors[0]))
	}
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
