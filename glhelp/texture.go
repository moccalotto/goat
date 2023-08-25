package glhelp

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Texture struct {
	handle      uint32
	unit        uint32 // number between 0 and GL_MAX_COMBINED_TEXTURE_IMAGE_UNITS - GL_TEXTURE0
	typ         uint32
	wrapR       int32
	wrapS       int32
	magFilter   int32
	minFilter   int32
	w           int32
	h           int32
	pix         []uint8 // We should be able to reuse pix across many texture objects. OR be able to reuse textures again and again
	initialized bool
}

func CreateTextureFromFile(filePath string, wrapR, wrapS int32) (*Texture, error) {

	img, err := LoadImage(filePath)

	if err != nil {
		return nil, GlProbablePanic(err)
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
			handle:    0,
			unit:      0,
			wrapS:     wrapS,
			wrapR:     wrapR,
			typ:       gl.TEXTURE_2D,
			minFilter: gl.NEAREST,
			magFilter: gl.NEAREST,

			w:   int32(imgRgba.Rect.Size().X),
			h:   int32(imgRgba.Rect.Size().Y),
			pix: imgRgba.Pix,
		},
		nil
}

func (T *Texture) Finalize() {

	if T.initialized {
		GlLog("Texture already initialzied")
		return
	}

	gl.GenTextures(1, &T.handle)
	gl.BindTexture(T.typ, T.handle)
	defer gl.BindTexture(T.typ, 0)

	gl.TexParameteri(T.typ, gl.TEXTURE_WRAP_R, T.wrapR)
	gl.TexParameteri(T.typ, gl.TEXTURE_WRAP_S, T.wrapS)
	gl.TexParameteri(T.typ, gl.TEXTURE_MIN_FILTER, T.minFilter) // minification filter
	gl.TexParameteri(T.typ, gl.TEXTURE_MAG_FILTER, T.magFilter) // magnification filter
	// https://gregs-blog.com/2008/01/17/opengl-texture-filter-parameters-explained/

	gl.TexImage2D(
		T.typ,            // Most likely T.typ
		0,                // quality level (0 is best)
		gl.SRGB_ALPHA,    // internal format
		T.w,              // width
		T.h,              // height
		0,                // border. Must be zero.
		gl.RGBA,          // pixes stored as R, G, B, A
		gl.UNSIGNED_BYTE, // one byte at a time
		gl.Ptr(T.pix),    // pointer to first pixel
	)

	gl.GenerateMipmap(T.typ)

	T.pix = nil
	T.initialized = true
	AssertGLOK()
}

func (T *Texture) GetTextureUnit() uint32 {
	return T.unit
}

func (T *Texture) Bind() {
	gl.BindTexture(T.typ, T.handle)
	AssertGLOK("BindTexture")
	gl.ActiveTexture(gl.TEXTURE0 + T.unit)
	AssertGLOK("BindTexture")
}

func (T *Texture) Unbind() {
	gl.BindTexture(T.typ, 0)
}

func (T *Texture) Destroy() {
	if !T.initialized {
		return
	}
	T.Unbind()
	gl.DeleteTextures(1, &T.handle)
	AssertGLOK()
}

func (T *Texture) SetMagFilter(magFilter int32) {
	if T.initialized {
		GlPanic(fmt.Errorf("cannot change attributes of a texture that has been Initialized()"))
	}
	T.magFilter = magFilter
}

func (T *Texture) GetMagFilter() int32 {
	return T.magFilter
}

func (T *Texture) SetMinFilter(minFilter int32) {
	if T.initialized {
		GlPanic(fmt.Errorf("cannot change attributes of a texture that has been Initialized()"))
	}
	T.minFilter = minFilter
}

func (T *Texture) GetMinFilter() int32 {
	return T.minFilter
}

func (T *Texture) GetSize() (int32, int32) {
	return T.w, T.h
}

func (T *Texture) GetWrapR() int32 {
	return T.wrapR
}
func (T *Texture) SetWrapR(wrapR int32) {
	if T.initialized {
		GlPanic(fmt.Errorf("cannot change attributes of a texture that has been Initialized()"))
	}
	T.wrapR = wrapR
}
func (T *Texture) SetRepeatR() {
	T.SetWrapR(gl.REPEAT)
}

func (T *Texture) GetWrapS() int32 {
	return T.wrapS
}
func (T *Texture) SetWrapS(wrapS int32) {
	if T.initialized {
		GlPanic(fmt.Errorf("cannot change attributes of a texture that has been Initialized()"))
	}
	T.wrapS = wrapS
}
func (T *Texture) SetRepeatS() {
	T.SetWrapS(gl.REPEAT)
}
