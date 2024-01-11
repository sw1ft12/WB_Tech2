package patterns

import "fmt"

// Class - an abstract interface for mmo ingame character class
type Class interface {
	getClassName() string
}

// Assasin - mmoClass
type Assassin struct {
	className string
	charName  string
}

// newAssasin - constructor
func newAssasin(name string) *Assassin {
	return &Assassin{
		className: "Assassin",
		charName:  name,
	}
}

// getClassName - getter for struct name
func (a *Assassin) getClassName() string {
	return a.className
}

// Warrior - mmoClass
type Warrior struct {
	className string
	charName  string
}

// newWarrior - constructor
func newWarrior(name string) *Warrior {
	return &Warrior{
		className: "Warrior",
		charName:  name,
	}
}

// getClassName - getter for struct name
func (w *Warrior) getClassName() string {
	return w.className
}

// getClass - Factory method for creating new struct
func getClass(className, charName string) (Class, error) {
	switch className {
	case "Warrior":
		return newWarrior(charName), nil
	case "Assassin":
		return newAssasin(charName), nil
	}
	return nil, fmt.Errorf("Wrong class name passed")
}
