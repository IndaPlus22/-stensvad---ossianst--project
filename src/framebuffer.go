package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type FrameBuffer struct {
	id     uint32
	width  uint32
	height uint32
}

func NewFrameBuffer(w uint32, h uint32) FrameBuffer {
	var id uint32
	gl.GenFramebuffers(1, &id)
	fb := FrameBuffer{id, w, h}

	return fb
}

func (fb *FrameBuffer) addColorTexture(slot uint32, texWidth uint32, texHeight uint32, colorAttachment uint32) {
	// TODO, will be implemented if needed
}

func (fb *FrameBuffer) addDepthTexture(slot uint32, texWidth uint32, texHeight uint32) {
	// Create new texture
	var tex uint32
	gl.GenTextures(1, &tex)

	// Bind texture to "slot"
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	// Texture paramaters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32, int32(texWidth), int32(texHeight), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)

	fb.bind()
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, tex, 0)
	fb.unbind()
}

func (fb *FrameBuffer) delete() {
	gl.DeleteFramebuffers(1, &fb.id)
}

func (fb *FrameBuffer) bind() {
	gl.Viewport(0, 0, int32(fb.width), int32(fb.height))

	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.id)
}

func (fb *FrameBuffer) unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
