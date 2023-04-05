package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray struct {
	id uint32
}

func NewVertexArray(elementCounts []int) VertexArray {
	var id uint32
	gl.GenVertexArrays(1, &id)
	va := VertexArray{id}
	va.bind()

	stride := 0
	for i := 0; i < len(elementCounts); i++ {
		stride += elementCounts[i] * 4
	}

	offset := 0
	for i := 0; i < len(elementCounts); i++ {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointer(uint32(i), int32(elementCounts[i]), gl.FLOAT, false, int32(stride), gl.PtrOffset(offset))

		offset += elementCounts[i] * 4
	}
	va.unbind()

	return va
}

func (va *VertexArray) delete() {
	gl.DeleteVertexArrays(1, &va.id)
}

func (va *VertexArray) bind() {
	gl.BindVertexArray(va.id)
}

func (va *VertexArray) unbind() {
	gl.BindVertexArray(0)
}
