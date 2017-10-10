package players

import (
	"felicien/puppet-server/model"
	"time"
)

// BoardPlayer represents a player associated with a Board
type BoardPlayer struct {
	servos      map[string]model.Servo
	bus         Bus
	commandSink chan model.Command
	started     bool
}

const (
	commandInterval time.Duration = 40 * time.Millisecond
)

// NewBoardPlayer creates a new BoardPlayer
func NewBoardPlayer(board model.Board) (*BoardPlayer, error) {
	bus, err := NewBus(board.BusType)
	if err != nil {
		return nil, err
	}

	var player = BoardPlayer{
		servos:      board.Servos,
		bus:         bus,
		commandSink: make(chan model.Command),
		started:     false,
	}

	return &player, nil
}

// Start starts the BoardPlayer
func (player *BoardPlayer) Start() error {
	err := player.bus.Open()
	if err != nil {
		return err
	}

	player.started = true
	go player.play()

	return nil
}

func (player *BoardPlayer) play() {
	defer player.Stop()

	var otherCommands []model.Command
	var positionCommands = make(map[int]model.PositionCommand)
	var ticker = time.NewTicker(commandInterval)

	for {
		select {
		case command := <-player.commandSink:
			if positionCommand, ok := command.(model.PositionCommand); ok {
				positionCommands[positionCommand.Addr] = positionCommand
			} else {
				otherCommands = append(otherCommands, command)
			}
		case <-ticker.C:
			if len(otherCommands) > 0 || len(positionCommands) > 0 {
				player.sendCommands(positionCommands, otherCommands)
				otherCommands = make([]model.Command, 0)
				positionCommands = make(map[int]model.PositionCommand)
			}
		}
	}
}

// Stop stops the bus and player
func (player *BoardPlayer) Stop() {
	player.bus.Close()
	player.started = false
}

// AddCommand adds a command for the next push
func (player *BoardPlayer) AddCommand(command model.Command) {
	player.commandSink <- command
}

func (player *BoardPlayer) sendCommands(positionCommands map[int]model.PositionCommand, otherCommands []model.Command) {
	for _, command := range otherCommands {
		player.bus.Writer().Write([]byte(command.CommandString()))
	}

	for _, command := range positionCommands {
		player.bus.Writer().Write([]byte(command.CommandString()))
	}
}
