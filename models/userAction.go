package models

type UserActions struct {
	SkipChan chan bool
	StopChan chan bool

	Stopped bool
}

func NewActions() *UserActions {
	return &UserActions{
		SkipChan: make(chan bool, 1),
		StopChan: make(chan bool, 1),
	}
}

func (a *UserActions) Stop() {
	a.Stopped = true
	a.StopChan <- true
}

func (a *UserActions) Skip() {
	a.SkipChan <- true
}