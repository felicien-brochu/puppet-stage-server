package players

import (
	"felicien/puppet-server/model"
	"fmt"
)

var puppetPlayers = make(map[string]*PuppetPlayer)

// PuppetPlayer player for a puppet (multiple boards)
type PuppetPlayer struct {
	puppet       model.Puppet
	boardPlayers map[string]*BoardPlayer
}

// GetPuppetPlayer returns the PuppetPlayer corresponding to the Puppet
func GetPuppetPlayer(puppet model.Puppet) *PuppetPlayer {
	if player, ok := puppetPlayers[puppet.Name]; ok {
		return player
	}
	return nil
}

// AddPuppetPlayer creates and stores a new PuppetPlayer
func AddPuppetPlayer(puppet model.Puppet) (*PuppetPlayer, error) {
	var player *PuppetPlayer
	if player = GetPuppetPlayer(puppet); player == nil {
		var err error
		player, err = NewPuppetPlayer(puppet)
		if err != nil {
			return nil, err
		}
		puppetPlayers[puppet.Name] = player
	}

	return player, nil
}

// NewPuppetPlayer creates a new PuppetPlayer and returns it
func NewPuppetPlayer(puppet model.Puppet) (*PuppetPlayer, error) {
	player := new(PuppetPlayer)
	player.puppet = puppet
	player.boardPlayers = make(map[string]*BoardPlayer)

	for id, board := range puppet.Boards {
		boardPlayer, err := NewBoardPlayer(board)
		if err != nil {
			return nil, err
		}
		player.boardPlayers[id] = boardPlayer
	}

	return player, nil
}

// Start starts all boards
func (player *PuppetPlayer) Start() error {
	for _, board := range player.puppet.Boards {
		err := player.StartBoard(board.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// StartBoard starts a specific board
func (player *PuppetPlayer) StartBoard(boardID string) error {
	boardPlayer, ok := player.boardPlayers[boardID]
	if !ok {
		return fmt.Errorf("No board for id '%s'", boardID)
	}

	if boardPlayer.started {
		return nil
	}

	err := boardPlayer.Start()
	return err
}

// GetBoardPlayer returns a board player corresponding to board ID.
func (player *PuppetPlayer) GetBoardPlayer(boardID string) *BoardPlayer {
	boardPlayer, ok := player.boardPlayers[boardID]
	if !ok {
		return nil
	}

	return boardPlayer
}
