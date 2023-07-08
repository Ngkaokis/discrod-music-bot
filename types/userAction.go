package types

func (a *UserActions) Stop() {
	a.Stopped = true
	a.StopChan <- true
}

func (a *UserActions) Skip() {
	a.SkipChan <- true
}