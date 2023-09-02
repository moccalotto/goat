package tractor

import (
	"goat/shed"
)

type BasicLine struct {
	renderer  *BasicRectRenderer
	camera    *Camera
	thickness float32
	color     shed.V4
	pos       Position
}

func CreateBasicLine(x1, y1, x2, y2, thickness float32, camera *Camera, renderer *BasicRectRenderer) *BasicLine {
	bl := BasicLine{
		color:     shed.OPAQ_WHITE(),
		camera:    camera,
		renderer:  renderer,
		thickness: thickness,
	}
	bl.SetCoords(x1, y1, x2, y2)

	return &bl
}

func (L *BasicLine) SetColor(rgba shed.V4) {
	L.color = rgba
}

func (L *BasicLine) SetCoords(x1, y1, x2, y2 float32) {
	L.SetPoints(
		shed.V2{X: x1, Y: y1},
		shed.V2{X: x2, Y: y2},
	)
}

func (L *BasicLine) SetPoints(p1, p2 shed.V2) {

	/*
	   => set scaleX = distance(p1, p2) + thickness
	   => set sacleY = thickness
	   => set posXY = midpoint(p1, p2)
	   => set Angle = angle(distance(p1, p2))
	*/

	dist := p2.Minus(p1)
	midpoint := p1.Plus(dist.Scaled(.5))

	L.pos.SetScale(dist.Len()+L.thickness, L.thickness)
	L.pos.SetXY(midpoint.X, midpoint.Y)

	L.pos.SetAngle(dist.Angle())
}

func (L *BasicLine) Draw() {
	L.renderer.Draw(L.camera.GetMatrix(), L.pos.GetMatrix(), L.color)
}
