package model

import "fmt"

// Command is to be sent by a bus to change the state of a board
type Command interface {
	commandString() string
}

// PositionCommand changes the position of device
type PositionCommand struct {
	addr     int
	position int
}

func (cmd PositionCommand) commandString() string {
	return fmt.Sprintf("P%d;%d\r", cmd.addr, cmd.position)
}
