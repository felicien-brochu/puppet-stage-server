package drivers

import (
	"felicien/puppet-server/model"
	"fmt"
	"time"
)

const (
	commandInterval time.Duration = 20 * time.Millisecond
)

var puppetDrivers = make(map[string]*PuppetDriver)

// PuppetDriver driver for a puppet (multiple boards)
type PuppetDriver struct {
	puppet       model.Puppet
	boardDrivers map[string]*BoardDriver
	started      bool
	ticker       *time.Ticker
	senderTicker chan time.Time
	boardTickers []chan time.Time
	done         chan struct{}
}

// GetPuppetDriver returns the PuppetDriver corresponding to the Puppet
func GetPuppetDriver(puppet model.Puppet) *PuppetDriver {
	if driver, ok := puppetDrivers[puppet.Name]; ok {
		return driver
	}
	return nil
}

// AddPuppetDriver creates and stores a new PuppetDriver
func AddPuppetDriver(puppet model.Puppet) (*PuppetDriver, error) {
	var driver *PuppetDriver
	if driver = GetPuppetDriver(puppet); driver == nil {
		var err error
		driver, err = NewPuppetDriver(puppet)
		if err != nil {
			return nil, err
		}
		puppetDrivers[puppet.Name] = driver
	}

	return driver, nil
}

// NewPuppetDriver creates a new PuppetDriver and returns it
func NewPuppetDriver(puppet model.Puppet) (*PuppetDriver, error) {
	driver := new(PuppetDriver)
	driver.puppet = puppet
	driver.boardDrivers = make(map[string]*BoardDriver)
	driver.boardTickers = make([]chan time.Time, 0)
	driver.started = false

	for id, board := range puppet.Boards {
		boardDriver, err := NewBoardDriver(board)
		if err != nil {
			return nil, err
		}
		driver.boardDrivers[id] = boardDriver
	}

	return driver, nil
}

// Start starts all boards
func (driver *PuppetDriver) Start() error {
	if !driver.started {
		driver.ticker = time.NewTicker(commandInterval)
		driver.boardTickers = make([]chan time.Time, 0)
		driver.senderTicker = make(chan time.Time)
		driver.done = make(chan struct{})

		for _, board := range driver.puppet.Boards {
			err := driver.StartBoard(board.ID)
			if err != nil {
				return err
			}
		}

		go driver.tick()
		driver.started = true
	}
	return nil
}

func (driver *PuppetDriver) tick() {
MainLoop:
	for {
		select {
		case t := <-driver.ticker.C:
			for _, boardTicker := range driver.boardTickers {
				boardTicker <- t
			}
			driver.senderTicker <- t
		case <-driver.done:
			break MainLoop
		}
	}

	close(driver.senderTicker)
	for _, boardTicker := range driver.boardTickers {
		close(boardTicker)
	}
}

// Stop stops all boards
func (driver *PuppetDriver) Stop() {
	driver.ticker.Stop()
	close(driver.done)
	for _, boardDriver := range driver.boardDrivers {
		boardDriver.Stop()
	}
	driver.started = false
}

// StartBoard starts a specific board
func (driver *PuppetDriver) StartBoard(boardID string) error {
	boardDriver, ok := driver.boardDrivers[boardID]
	if !ok {
		return fmt.Errorf("No board for id '%s'", boardID)
	}

	if boardDriver.started {
		return nil
	}

	var boardTicker = make(chan time.Time)
	driver.boardTickers = append(driver.boardTickers, boardTicker)
	return boardDriver.Start(boardTicker)
}

// GetBoardDriver returns a board driver corresponding to board ID.
func (driver *PuppetDriver) GetBoardDriver(boardID string) *BoardDriver {
	boardDriver, ok := driver.boardDrivers[boardID]
	if !ok {
		return nil
	}

	return boardDriver
}

// SetServoPosition set a servo position
func (driver *PuppetDriver) SetServoPosition(servoID string, position int) error {
	var boardID string
	var servo model.Servo

BoardLoop:
	for boardIDKey, board := range driver.puppet.Boards {
		for servoIDKey, servoItem := range board.Servos {
			if servoIDKey == servoID {
				boardID = boardIDKey
				servo = servoItem
				break BoardLoop
			}
		}
	}

	if boardID == "" {
		return fmt.Errorf("No servo for id '%s' in puppet '%s'", servoID, driver.puppet.ID)
	}

	boardDriver := driver.GetBoardDriver(boardID)
	var positionCommand = PositionCommand{
		Addr:     servo.Addr,
		Position: position,
	}
	boardDriver.AddCommand(positionCommand)

	return nil
}

// GetSenderTicker returns a ticker that can be used to be in sync with the boards' buses rythm
func (driver *PuppetDriver) GetSenderTicker() chan time.Time {
	return driver.senderTicker
}
