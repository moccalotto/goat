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
	handle uint32
	// unit   uint32 // for now we always work on unit 0
	wrapR int32
	wrapS int32
	W     int32
	H     int32
	pix   []uint8 // We should be able to reuse pix across many texture objects. OR be able to reuse textures again and again
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
			handle: 0,
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
	gl.BindTexture(gl.TEXTURE_2D, T.handle)
	defer gl.BindTexture(gl.TEXTURE_2D, 0)

	// All of these hardcoded params could be set dynamic,
	// for instance as properties on the T struct.
	// Also, we could lazy-exec all of these settings
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, T.wrapR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, T.wrapS)

	//
	// https://gregs-blog.com/2008/01/17/opengl-texture-filter-parameters-explained/
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_LINEAR) // minification filter
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST) // minification filter
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)               // magnification filter

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.SRGB_ALPHA,    // internal format
		T.W,              // width
		T.H,              // height
		0,                // border. Must be zero.
		gl.RGBA,          // pixes stored as R, G, B, A
		gl.UNSIGNED_BYTE, // one byte at a time
		gl.Ptr(T.pix),    // pointer to first pixel
	)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	AssertGLOK()
}

func (T *Texture) GetTextureUnit() int32 {
	return 0
}

func (T *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, T.handle)
	gl.ActiveTexture(uint32(T.GetTextureUnit()))
}

func (T *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (T *Texture) Destroy() {
	T.Unbind()
	gl.DeleteTextures(1, &T.handle)
}
