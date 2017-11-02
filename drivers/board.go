package drivers

import (
	"errors"
	"felicien/puppet-server/model"
	"time"
)

// BoardDriver represents a driver associated with a Board
type BoardDriver struct {
	servos      map[string]model.Servo
	bus         Bus
	commandSink chan Command
	started     bool
	done        chan struct{}
}

// NewBoardDriver creates a new BoardDriver
func NewBoardDriver(board model.Board) (*BoardDriver, error) {
	bus, err := NewBus(board.BusType)
	if err != nil {
		return nil, err
	}

	var driver = BoardDriver{
		servos:      board.Servos,
		bus:         bus,
		commandSink: make(chan Command),
		started:     false,
		done:        make(chan struct{}),
	}

	return &driver, nil
}

// Start starts the BoardDriver
func (driver *BoardDriver) Start(ticker chan time.Time) error {
	err := driver.bus.Open()
	if err != nil {
		return err
	}

	driver.commandSink = make(chan Command)
	driver.done = make(chan struct{})
	driver.started = true
	go driver.play(ticker)
	return nil
}

func (driver *BoardDriver) play(ticker chan time.Time) {
	defer driver.bus.Close()
	defer close(driver.commandSink)

	var otherCommands []Command
	var positionCommands = make(map[int]PositionCommand)

MainLoop:
	for {
		select {
		case command := <-driver.commandSink:
			if positionCommand, ok := command.(PositionCommand); ok {
				positionCommands[positionCommand.Addr] = positionCommand
			} else {
				otherCommands = append(otherCommands, command)
			}
		case <-ticker:
			if len(otherCommands) > 0 || len(positionCommands) > 0 {
				err := driver.sendCommands(positionCommands, otherCommands)
				if err != nil {
					break MainLoop
				}
				otherCommands = make([]Command, 0)
				positionCommands = make(map[int]PositionCommand)
			}
		case <-driver.done:
			break MainLoop
		}
	}
	driver.started = false
}

// Stop stops the bus and driver
func (driver *BoardDriver) Stop() {
	if driver.started {
		close(driver.done)
	}
}

// AddCommand adds a command for the next push
func (driver *BoardDriver) AddCommand(command Command) error {
	if !driver.started {
		return errors.New("BoardDriver is not started")
	}
	driver.commandSink <- command
	return nil
}

func (driver *BoardDriver) sendCommands(positionCommands map[int]PositionCommand, otherCommands []Command) error {
	for _, command := range otherCommands {
		_, err := driver.bus.Writer().Write([]byte(command.CommandString()))
		if err != nil {
			return err
		}
	}

	for _, command := range positionCommands {
		_, err := driver.bus.Writer().Write([]byte(command.CommandString()))
		if err != nil {
			return err
		}
	}
	return nil
}
