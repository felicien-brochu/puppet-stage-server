package drivers

import "fmt"

// Command is to be sent by a bus to change the state of a board
type Command interface {
	CommandString() string
}

// PositionCommand changes the position of device
type PositionCommand struct {
	Addr     int
	Position int
}

// CommandString returns the string representation to be sent on a Bus
func (cmd PositionCommand) CommandString() string {
	return fmt.Sprintf("P%d;%d\r", cmd.Addr, cmd.Position)
}
