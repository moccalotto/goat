package main

type ELine struct {
	id    uint64
	x1    float64
	y1    float64
	x2    float64
	y2    float64
	color Color
}

func CreateELine(x1, y1, x2, y2 float64) *ELine {
	return &ELine{
		id:    0,
		x1:    x1,
		y1:    y1,
		x2:    x2,
		y2:    y2,
		color: Color{R: 0, G: 0, B: 0, A: 255},
	}
}

func CreateELineX(centX, centY, length, radians float64) *ELine {

	half := CreatePolarVector(radians, length/2)

	return CreateELine(
		centX+half.x,
		centY+half.y,
		centX-half.x,
		centX-half.y,
	)
}

func (e *ELine) Draw(dm *Drawing) {

	dm.Push()
	dm.fgColor = e.color
	dm.Line(e.x1, e.y1, e.x2, e.y2)
	dm.Pop()
}

func (e *ELine) SetId(id uint64) {
	// has the id been set previously?
	if e.id != 0 {
		panic("you cannot set the id of an entity more than once!")
	}

	e.id = id
}

func (e *ELine) GetId() uint64 {
	return e.id
}

func (e *ELine) SetCoords(x1, y1, x2, y2 float64) {
	e.x1 = x1
	e.y1 = y1
	e.x2 = x2
	e.y2 = y2
}

func (e *ELine) Color(c ...uint8) (uint8, uint8, uint8, uint8) {
	switch len(c) {
	case 0:
		return e.color.R, e.color.G, e.color.B, e.color.A
	case 1:
		return e.Color(c[0], c[0], c[0], 255)
	case 3:
		return e.Color(c[0], c[1], c[2], 255)
	case 4:
		e.color = Color{R: c[0], G: c[1], B: c[2], A: c[3]}
		return e.Color()
	default:
		panic("Color() takes 1, 3, or 4 arguments.")
	}
}

func (e *ELine) Offset(offset Vector) {

	e.x1 += offset.X()
	e.x2 += offset.X()

	e.y1 += offset.Y()
	e.y2 += offset.Y()
}

func (e *ELine) Rotate(angle float64) {
	e.Pivot(
		e.x1+(e.x2-e.x1)/2,
		e.y1+(e.y2-e.y1)/2,
		angle)
}

func (e *ELine) Pivot(pivotX, pivotY, radians float64) {

	PivotVec := CreateVector(pivotX, pivotY)

	vec1 := CreateVector(e.x1, e.y1).Sub(PivotVec).Rotate(radians).Add(PivotVec)
	vec2 := CreateVector(e.x2, e.y2).Sub(PivotVec).Rotate(radians).Add(PivotVec)

	e.x1 = vec1.x
	e.y1 = vec1.y
	e.x2 = vec2.x
	e.y2 = vec2.y
}
