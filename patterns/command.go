package patterns

type ToggleOnCommand struct {
	receiver *Receiver
}

// Execute command.
func (c *ToggleOnCommand) Execute() string {
	return c.receiver.ToggleOn()
}

type ToggleOffCommand struct {
	receiver *Receiver
}

// Execute command.
func (c *ToggleOffCommand) Execute() string {
	return c.receiver.ToggleOff()
}

type Receiver struct{}

func (r *Receiver) ToggleOn() string {
	return "Toggle On"
}

func (r *Receiver) ToggleOff() string {
	return "Toggle Off"
}

// Command provides a command interface.
type Command interface {
	Execute() string
}

type Invoker struct {
	commands []Command
}

func (i *Invoker) StoreCommand(command Command) {
	i.commands = append(i.commands, command)
}

// UnStoreCommand removes command.
func (i *Invoker) UnStoreCommand() {
	if len(i.commands) != 0 {
		i.commands = i.commands[:len(i.commands)-1]
	}
}

// Execute all commands.
func (i *Invoker) Execute() string {
	var result string
	for _, command := range i.commands {
		result += command.Execute() + "\n"
	}
	return result
}
