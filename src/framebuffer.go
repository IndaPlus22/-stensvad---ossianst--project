package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type FrameBuffer struct {
	id     uint32
	width  uint32
	height uint32

	m_DrawBuffers []uint32
}

/*
Creates a new FrameBuffer with resolution wxh and returns it

Parameters:
- w: the width of the frame buffer
- h: the height of the frame buffer

Returns:
- p: the new FrameBuffer object

Example usage:

	earthSettings := DefaultEarth()
	planet := NewPlanet(earthSettings)
*/
func NewFrameBuffer(w uint32, h uint32) FrameBuffer {
	var id uint32
	gl.GenFramebuffers(1, &id)
	fb := FrameBuffer{id, w, h, []uint32{}}

	return fb
}

/*
Creates a color texture and attaches it to the FrameBuffer

Parameters:
- slot: the texture slot to store the texture in
- texWidth: the width of the texture
- texHeight: the height of the texture
- colorAttachment: what color attachment the texture will use
- pixelSize: the size of the internal format
*/
func (fb *FrameBuffer) addColorTexture(slot uint32, texWidth uint32, texHeight uint32, colorAttachment uint32, pixelSize int32) {
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

	gl.TexImage2D(gl.TEXTURE_2D, 0, pixelSize, int32(texWidth), int32(texHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	fb.bind()
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, colorAttachment, gl.TEXTURE_2D, tex, 0)

	fb.m_DrawBuffers = append(fb.m_DrawBuffers, colorAttachment)
	gl.DrawBuffers(int32(len(fb.m_DrawBuffers)), &fb.m_DrawBuffers[0])
	fb.unbind()
}

/*
Creates a depth texture and attaches it to the FrameBuffer

Parameters:
- slot: the texture slot to store the texture in
- texWidth: the width of the texture
- texHeight: the height of the texture
*/
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

func (fb *FrameBuffer) addRenderBuffer(rb uint32) {
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rb)
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
