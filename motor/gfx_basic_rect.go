package motor

import (
	"goat/util"
)

type BasicRect struct {
	Renderer *BasicRectRenderer
	Camera   *Camera
	Deleted  bool
	Position

	Color util.V4
}

func CreateBasicRect(x, y, w, h, a float32, camera *Camera, renderer *BasicRectRenderer) *BasicRect {

	R := BasicRect{
		Camera:   camera,
		Renderer: renderer,
		Color:    util.OPAQ_WHITE(),
	}

	R.SetXY(x, y)
	R.SetScale(w, h)
	R.SetAngle(a)

	return &R
}

func (R *BasicRect) Draw() {
	if R.Deleted {
		return
	}
	camMatrix := R.Camera.GetMatrix()
	thingMatrix := R.GetMatrix()

	R.Renderer.UniColor = R.Color
	R.Renderer.Draw(camMatrix, thingMatrix, R.Color)
}
