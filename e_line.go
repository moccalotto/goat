package main

type ELine struct {
	a *Vector
	b *Vector
}

func (e *ELine) Draw(dm *Drawing) {
	dm.Line(
		e.a.X(),
		e.a.Y(),
		e.b.X(),
		e.b.Y())
}
