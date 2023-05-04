package main

import "github.com/go-gl/gl/v4.1-core/gl"

type RenderBuffer struct {
	id     uint32
	width  uint32
	height uint32
}

// Initialize a render buffer of dimensions wxh and return it as a RenderBuffer object
func NewRenderBuffer(w uint32, h uint32) RenderBuffer {
	var id uint32

	gl.GenRenderbuffers(1, &id)
	gl.BindRenderbuffer(gl.RENDERBUFFER, id)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(w), int32(h))

	rb := RenderBuffer{id, w, h}
	return rb
}

func (rb *RenderBuffer) bind() {
	gl.BindRenderbuffer(gl.RENDERBUFFER, rb.id)
}

func (rb *RenderBuffer) unbind() {
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
}
