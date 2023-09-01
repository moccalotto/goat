package motor

import (
	"goat/util"
)

type BasicLine struct {
	renderer  *BasicRectRenderer
	camera    *Camera
	thickness float32
	color     util.V4
	pos       Position
}

func CreateBasicLine(x1, y1, x2, y2, thickness float32, camera *Camera, renderer *BasicRectRenderer) *BasicLine {
	bl := BasicLine{
		color:    util.OPAQ_WHITE(),
		camera:   camera,
		renderer: renderer,
	}
	bl.Set(
		thickness,
		util.V2{X: x1, Y: y1},
		util.V2{X: x2, Y: y2},
	)

	return &bl
}

func (L *BasicLine) Set(thickness float32, p1, p2 util.V2) {

	/*
	   => set scaleX = distance(p1, p2) [MAYBE ADD THICKNESS]
	   => set sacleY = thickness
	   => set posXY = midpoint(p1, p2)
	   => set Angle = angle(distance(p1, p2))
	*/

	dist := p2.Sub(p1)
	midpoint := p1.Add(dist.Scale(.5))

	L.pos.SetScale(dist.Len(), thickness)
	L.pos.SetXY(midpoint.X, midpoint.Y)

	L.pos.SetAngle(dist.Angle())
}

func (L *BasicLine) Draw() {
	L.renderer.Draw(L.camera.GetMatrix(), L.pos.GetMatrix(), L.color)
}
