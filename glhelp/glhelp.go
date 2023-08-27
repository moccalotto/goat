package glhelp

/**
Wrapper for all calls to opengl.

- make importing opengl less stupid
- streamline calls so they return better errors
- wrap stuff in sane data structures
*/

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	Tau           = FullCircle
	FullCircle    = math.Pi * 2
	HalfCircle    = math.Pi
	QuarterCircle = math.Pi / 2

	Degrees = FullCircle / 360

	Float32Size = 4
	Uint32Size  = 4
	Int32Size   = 4
)

func LoadImage(filePath string) (*image.RGBA, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, GlProbablePanic(fmt.Errorf("could not open file '%s' - %v", filePath, err))
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, GlProbablePanic(err)
	}

	imgRgba := image.NewRGBA(img.Bounds())
	draw.Draw(imgRgba, imgRgba.Bounds(), img, image.Pt(0, 0), draw.Src)

	// we must have 4 8-bit colors per pixel
	if imgRgba.Stride != imgRgba.Rect.Size().X*4 {
		return nil, GlProbablePanic(
			fmt.Errorf("only 32-bit colors supported. Stride is %d, but should be %d", imgRgba.Stride, imgRgba.Rect.Size().X*4),
		)
	}

	return imgRgba, nil
}

/**
 * Return a pointer to the given string
 */
func GlStr(str string) *uint8 {
	/**
	 * TODO: are we doing something dangerous here?
	 * What if this temporary string gets garbage collected before OpenGL gets to use it ?
	 */
	return unsafe.StringData(str + "\x00")
}

func ClearScreenF(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func ClearScreenI(r, g, b, a uint8) {
	ClearScreenF(
		float32(r)/255.0,
		float32(g)/255.0,
		float32(b)/255.0,
		float32(a)/255.0,
	)
}

func ReadFile(filename string) (string, error) {

	bytes, err := os.ReadFile(filename)

	if err != nil {
		return "", GlProbablePanic(err)
	}

	return string(bytes), nil
}

func MatMulMany(mats ...mgl32.Mat3) mgl32.Mat3 {
	result := mats[0]

	for i := 1; i < len(mats); i++ {
		result = result.Mul3(mats[i])
	}

	return result
}

func MatMulX3(m1, m2, m3 mgl32.Mat3) mgl32.Mat3 {
	return m1.Mul3(m2).Mul3(m3)
}

func EnableBlending() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func Lerp(start, end, amt float32) float32 {
	diff := end - start

	return start + diff*amt
}

func Sincos(radians float32) (float32, float32) {
	sin, cos := math.Sincos(float64(radians))

	return float32(sin), float32(cos)
}

func Enable2D() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)

}

func Wireframe(on bool) {
	if on {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		return
	}

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}

	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}

	return b
}

func Sign(f float32) float32 {
	if f >= 0 {
		return 1.0
	}

	return -1.0
}

func Hypot(a, b float32) float32 {
	return float32(math.Hypot(float64(a), float64(b)))
}

func NormalizeAngle(r float32) float32 {
	for r > Tau {
		r -= Tau
	}

	for r < 0 {
		r += Tau
	}

	return r
}
