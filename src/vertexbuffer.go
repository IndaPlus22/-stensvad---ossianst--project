package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexBuffer struct {
	id uint32
}

func NewVertexBuffer(vertices []float32) VertexBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	vb := VertexBuffer{id}
	vb.bind()
	// Calculate and set size of buffer
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	vb.unbind()

	return vb
}

func (vb *VertexBuffer) delete() {
	gl.DeleteBuffers(1, &vb.id)
}

func (vb *VertexBuffer) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vb.id)
}

func (vb *VertexBuffer) unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
