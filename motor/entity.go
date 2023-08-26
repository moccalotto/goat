package motor

type Entity interface {
	Draw()
	Update()
	Clone() Entity
}
