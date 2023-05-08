package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray struct {
	id uint32
}

/*
NewVertexArray generates a new vertex array

Parameters:
- elementCounts: how many floats every element in a vertex contains

Returns:
- va: the new vertex array object

Example usage:

	// Generate a vertex with room for position(3 floats), uv coordinates(2 floats) and normal(3 floats), for example
	va := NewVertexArray([]int{3, 2, 3})
*/
func NewVertexArray(elementCounts []int) VertexArray {
	var id uint32
	gl.GenVertexArrays(1, &id)
	va := VertexArray{id}
	va.bind()

	// Calculate size of one vertex
	stride := 0
	for i := 0; i < len(elementCounts); i++ {
		stride += elementCounts[i] * 4
	}

	// Set vertex array layout
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
