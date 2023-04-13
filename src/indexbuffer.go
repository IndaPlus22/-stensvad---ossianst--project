package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type IndexBuffer struct {
	id    uint32
	count int32
}

func NewIndexBuffer(indices []uint32) IndexBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	ib := IndexBuffer{id, int32(len(indices))}
	ib.bind()
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	ib.unbind()

	return ib
}

func (ib *IndexBuffer) delete() {
	gl.DeleteBuffers(1, &ib.id)
}

func (ib *IndexBuffer) bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.id)
}

func (ib *IndexBuffer) unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}
