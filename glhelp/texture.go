package glhelp

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Texture struct {
	handle uint32
	target uint32 // same target as gl.BindTexture(<this param>, ...)
	unit   uint32 // Texture unit. You can bundle many textures together in one unit, and allow the shader to work on all of them
	wrapR  int32
	wrapS  int32
	W      int32
	H      int32
	pix    []uint8
}

func CreateTextureFromFile(filePath string, wrapR, wrapS int32) (*Texture, error) {

	img, err := LoadImage(filePath)

	if err != nil {
		return nil, err
	}

	return CreateTexture(img, wrapR, wrapS)
}

func CreateTexture(img image.Image, wrapR, wrapS int32) (*Texture, error) {
	imgRgba := image.NewRGBA(img.Bounds())
	draw.Draw(imgRgba, imgRgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	stride := imgRgba.Rect.Size().X * 4 // bytes per horizontal line
	if imgRgba.Stride != stride {
		// this really shouldnt happen
		return nil, fmt.Errorf("only 32-bit colors supported. Stride is %d, but should be %d", imgRgba.Stride, stride)
	}
	return &Texture{
			handle: 0,
			target: uint32(gl.TEXTURE_2D),
			wrapS:  wrapS,
			wrapR:  wrapR,
			W:      int32(imgRgba.Rect.Size().X),
			H:      int32(imgRgba.Rect.Size().Y),
			pix:    imgRgba.Pix,
		},
		nil
}

func (T *Texture) Initialize() {

	gl.GenTextures(1, &T.handle)

	T.Bind(gl.TEXTURE0)

	// All of these hardcoded params could be set dynamic,
	// for instance as properties on the T struct.
	// Also, we could lazy-exec all of these settings
	gl.TexParameteri(T.target, gl.TEXTURE_WRAP_R, T.wrapR)
	gl.TexParameteri(T.target, gl.TEXTURE_WRAP_S, T.wrapS)
	gl.TexParameteri(T.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(T.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter

	gl.TexImage2D(
		T.target,         // Target. In our case hardcoded to TEXTURE_2D
		0,                // quality. 0 is best (I think)
		gl.SRGB_ALPHA,    // internal format
		T.W,              // width
		T.H,              // height
		0,                // border. Must be zero.
		gl.RGBA,          // pixes stored as R, G, B, A
		gl.UNSIGNED_BYTE, // one byte at a time
		gl.Ptr(T.pix),    // pointer to first pixel
	)

	gl.GenerateMipmap(T.target)
	AssertGLOK()
}

func (T *Texture) Bind(unit uint32) {
	T.unit = unit
	gl.ActiveTexture(unit)
	gl.BindTexture(T.target, T.handle)
}

func (T *Texture) Unbind() {
	if T.unit == 0 {
		return
	}

	T.unit = 0
	gl.BindTexture(T.target, 0)
}

func (T *Texture) Setuniform(uniformLoc int32) error {
	if T.unit == 0 {
		return errors.New("texture not bound")
	}
	gl.Uniform1i(uniformLoc, int32(T.unit-gl.TEXTURE0))
	return nil
}
