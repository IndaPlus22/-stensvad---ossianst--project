package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type CubemapTexture struct {
	id uint32
}

// Creates a new cubemap texture and returns it. Panics if something goes wrong while loading the texture.
func NewCubemapTexture(filePath string) CubemapTexture {
	id, err := genCubemapTexture(filePath)
	if err != nil {
		panic(err)
	}

	return CubemapTexture{id}
}

func genCubemapTexture(filePath string) (uint32, error) {
	filePath = "../res/textures/" + filePath

	var cubemapTexture uint32
	gl.GenTextures(1, &cubemapTexture)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubemapTexture)

	imageNames := []string{"/right.png", "/left.png", "/top.png", "/bottom.png", "/front.png", "/back.png"}

	for i, imagePath := range imageNames {
		imgFile, err := os.Open(filePath + imagePath)
		if err != nil {
			return 0, fmt.Errorf("texture %q not found on disk: %v", filePath, err)
		}
		img, _, err := image.Decode(imgFile)
		if err != nil {
			return 0, err
		}

		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return 0, fmt.Errorf("unsupported stride")
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		gl.TexImage2D(
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)
	}

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return cubemapTexture, nil
}

func (t *CubemapTexture) bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
}

func (t *CubemapTexture) unbind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
}
