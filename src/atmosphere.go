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

	shader    Shader
	normalMap Texture
}

/*
Creates a frame covering the entire screen to apply post processing on, using a
frame buffer with resolution wxh and an associated shader program

Parameters:
- w: the width of the frame buffer
- h: the height of the frame buffer
- shaderPath: the file name of the associated shader program

Returns:
- ppf: a PostProcessingFrame object
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

	normalMap := NewTexture("normalmap_ocean.png")
	normalMap.bind(5)

	// Create framebuffer
	fb := NewFrameBuffer(w, h)
	fb.addColorTexture(2, w, h, gl.COLOR_ATTACHMENT0, gl.RGBA)
	fb.addColorTexture(3, w, h, gl.COLOR_ATTACHMENT1, gl.RGBA32F)
	fb.addDepthTexture(10, w, h)

	ppf := PostProcessingFrame{va, fb, ib, []uint32{}, shader, normalMap}
	ppf.shader.bind()
	ppf.shader.setUniform1i("colorTexture", 2)
	ppf.shader.setUniform1i("depthTexture", 3)
	ppf.shader.setUniform1f("camNear", cam.GetNearPlane())
	ppf.shader.setUniform1f("camFar", cam.GetFarPlane())
	ppf.shader.setUniform1i("oceanNormalMap", 5)

	return ppf
}

/*
Create a Uniform Buffer of vec4 values

Parameters:
- block: the name of the Uniform Block in the shader program
- vectors: vec4 values to add as buffer data
- size: the amount of memory to allocate for the buffer data
*/
func (ppf *PostProcessingFrame) addUniformBufferVec4(block string, vectors []mgl32.Vec4, size int) {
	var id uint32
	gl.GenBuffers(1, &id)
	gl.BindBuffer(gl.UNIFORM_BUFFER, id)

	gl.BufferData(gl.UNIFORM_BUFFER, size, nil, gl.DYNAMIC_DRAW)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, int(unsafe.Sizeof(mgl32.Vec4{}))*len(vectors), unsafe.Pointer(&vectors[0]))

	blockId := gl.GetUniformBlockIndex(ppf.shader.id, gl.Str(block+"\x00"))

	gl.UniformBlockBinding(ppf.shader.id, blockId, 0)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 0, id)

	ppf.ub = append(ppf.ub, id)
}

/*
Update vec4 values of Uniform Buffer

Parameters:
- ub: the id of the uniform buffer to update
- vectors: the updated vec4 values
*/
func (ppf *PostProcessingFrame) updateUniformBufferVec4(ub uint32, vectors []mgl32.Vec4) {
	gl.BindBuffer(gl.UNIFORM_BUFFER, ub)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, int(unsafe.Sizeof(mgl32.Vec4{}))*len(vectors), unsafe.Pointer(&vectors[0]))
}

// Bind necessary buffers and shader programs and render the post processing effects
func (ppf *PostProcessingFrame) draw() {
	ppf.shader.bind()
	ppf.va.bind()
	ppf.ib.bind()

	gl.DrawElements(gl.TRIANGLES, ppf.ib.count, gl.UNSIGNED_INT, gl.PtrOffset(0))

	ppf.shader.unbind()
	ppf.va.unbind()
}
