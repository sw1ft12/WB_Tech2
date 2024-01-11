package patterns

type State interface {
	pressButton()
}

type Computer struct {
	currState State
}

// pressButton - do action depending on curr state
func (c *Computer) pressButton() {
	c.currState.pressButton()
}

// changeState - changes curren state
func (c *Computer) changeState(state State) {
	c.currState = state
}

func newComputer() *Computer {
	computer := Computer{}
	computer.changeState(&TurnedOff{computer: &computer})
	return &computer
}

// TurnedOff State
type TurnedOff struct {
	computer *Computer
}

// pressbutton in turned off state
func (t *TurnedOff) pressButton() {
	t.computer.changeState(&TurnedOn{computer: t.computer})
}

// Asleep state
type Asleep struct {
	computer *Computer
}

// pressbutton in asleep state
func (a *Asleep) pressButton() {
	a.computer.changeState(&TurnedOn{computer: a.computer})
}

type TurnedOn struct {
	computer *Computer
}

// pressbutton in turned on state
func (t *TurnedOn) pressButton() {
	t.computer.changeState(&Asleep{computer: t.computer})
}
