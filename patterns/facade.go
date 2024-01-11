package patterns

type Composition struct{}

type Reproducer struct {
	comp *Composition
}

// Reproduce - Plays sounds
func (*Reproducer) Reproduce() {
	//...
}

type VolumeController struct {
	volume int
}

// Increase - increase volume by value
func (vc *VolumeController) Increase(value int) {
	vc.volume += value
}

// Decrease - decrease volume by value
func (vc *VolumeController) Decrease(value int) {
	vc.volume -= value
}

// Player - facade for playing music
type Player struct {
	*Reproducer
	*VolumeController
}
