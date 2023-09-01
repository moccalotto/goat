package motor

type Behavior interface {
	Update(e *Sprite)
}

type SimpleBehavior struct {
	updateFunc func(*Sprite)
}

func (SB *SimpleBehavior) Update(e *Sprite) {
	SB.updateFunc(e)
}

func CreateSimpleBehavior(f func(e *Sprite)) *SimpleBehavior {
	return &SimpleBehavior{
		updateFunc: f,
	}
}
