package glhelp

import (
	// "github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	wPosX        float32 // The X position the camera is looking at - in world coordinates
	wPosY        float32 // The Y position the camera is looking at - in world coordinates
	wFrameWidth  float32 // The width of the camera view - in world coordinates
	wFrameHeight float32 // The height of the camera view - in world coordinates
	wRotation    float32 // the rotation, relative to the world's X-axis, of the camera - in radians
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

func (C *Camera) GetTransformationMatrix() mgl32.Mat3 {
	scale := mgl32.Scale2D(2/C.wFrameWidth, 2/C.wFrameHeight)
	rotate := mgl32.HomogRotate2D(C.wRotation)
	translate := mgl32.Translate2D(C.wPosX, C.wPosY)

	// return MatMulMany(scale, rotate, translate)
	return MatMulX3(scale, rotate, translate)
}

func (C *Camera) Rotate(radians float32) {
	C.wRotation -= radians // cam movement must be inverted to behave as expected
}

func (C *Camera) SetRotation(radians float32) {
	C.wRotation = -radians // cam rotation must be inverted to behave as expected
}

func (C *Camera) SetPosition(x, y float32) {
	C.wPosX = -x // camera movement must be negative to
	C.wPosY = -y // behave as expected
}

func (C *Camera) Move(x, y float32) {
	C.wPosX -= x // camera movement must be negative to
	C.wPosY -= y // behave as expected
}

func (C *Camera) SetFrameSize(w, h float32) {
	C.wFrameWidth = w
	C.wFrameHeight = h
}

// Zoom(10) = increase zoom 10x
func (C *Camera) Zoom(amount float32) {
	C.wFrameWidth /= amount
	C.wFrameHeight /= amount
}
