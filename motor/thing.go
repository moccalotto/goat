package motor

type Thing interface {
	Draw()
	Update()
	Clone() Thing
}
