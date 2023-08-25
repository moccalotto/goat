package glhelp

import (
	"github.com/go-gl/mathgl/mgl32"
)

// A camera can transform world-sizes into opengl sizes [-1, 1]A
// it also has a position so you can chose to only see some of your work area
type Camera struct {
	wPosX        float32 // The X position the camera is looking at - in world coordinates
	wPosY        float32 // The Y position the camera is looking at - in world coordinates
	wFrameWidth  float32 // The width of the camera view - in world coordinates
	wFrameHeight float32 // The height of the camera view - in world coordinates
	wRotation    float32 // the rotation, relative to the world's X-axis, of the camera - in radians

	transMatrixCache mgl32.Mat3
	cacheValid       bool
}

func CreateCamera() *Camera {
	return &Camera{
		wPosX:        0,
		wPosY:        0,
		wFrameWidth:  2,
		wFrameHeight: 2,
		wRotation:    0.0,
	}
}

func (C *Camera) GetMatrix() mgl32.Mat3 {

	if C.cacheValid {
		return C.transMatrixCache
	}

	scale := mgl32.Scale2D(2/C.wFrameWidth, 2/C.wFrameHeight)
	rotate := mgl32.HomogRotate2D(C.wRotation)
	translate := mgl32.Translate2D(C.wPosX, C.wPosY)

	// Scale, Rotate, Translate: reverse order as when transforming models
	// this is because a camera can be considered an "inverse" model.
	C.transMatrixCache = MatMulX3(scale, rotate, translate)
	C.cacheValid = true

	return C.transMatrixCache
}

func (C *Camera) Rotate(radians float32) {
	C.cacheValid = false
	C.wRotation -= radians // cam movement must be inverted to behave as expected
}

func (C *Camera) SetRotation(radians float32) {
	C.cacheValid = false
	C.wRotation = -radians // cam rotation must be inverted to behave as expected
}

func (C *Camera) SetPosition(x, y float32) {
	C.cacheValid = false
	C.wPosX = -x // camera movement must be negative to
	C.wPosY = -y // behave as expected
}

func (C *Camera) Move(x, y float32) {
	C.cacheValid = false
	C.wPosX -= x // camera movement must be negative to
	C.wPosY -= y // behave as expected
}

func (C *Camera) SetFrameSize(w, h float32) {
	C.cacheValid = false
	C.wFrameWidth = w
	C.wFrameHeight = h
}

// Zoom(10) = increase zoom 10x
func (C *Camera) Zoom(amount float32) {
	C.cacheValid = false
	C.wFrameWidth /= amount
	C.wFrameHeight /= amount
}

// the number of length units in each direction the camera can see.
func (C *Camera) GetFrameSize() (float32, float32) {
	return C.wFrameWidth, C.wFrameHeight
}

// Get frame size as a vector
func (C *Camera) GetFrameSizeV() mgl32.Vec2 {
	return mgl32.Vec2{C.wFrameWidth, C.wFrameHeight}
}
