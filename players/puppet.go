package players

import (
	"felicien/puppet-server/drivers"
	"felicien/puppet-server/model"
	"fmt"
)

type PuppetPlayer struct {
	puppet       model.Puppet
	puppetDriver *drivers.PuppetDriver
	playing      bool
	stagePlayer  stagePlayer
}

func NewPuppetPlayer(puppet model.Puppet, puppetDriver *drivers.PuppetDriver) PuppetPlayer {
	return PuppetPlayer{
		puppet:       puppet,
		puppetDriver: puppetDriver,
		playing:      false,
	}
}

func (player *PuppetPlayer) PlayStage(stage model.Stage, playStart model.Time, stateChan chan string) {
	player.stagePlayer = newStagePlayer(stage, player.puppet, playStart)
	player.playing = true
	player.stagePlayer.play(player, stateChan, player.puppetDriver.GetSenderTicker())
}

func (player *PuppetPlayer) StopStage() {
	player.stagePlayer.stop()
	player.playing = false
}

func (player *PuppetPlayer) playServoPosition(servoID string, position int) {
	err := player.puppetDriver.SetServoPosition(servoID, position)
	if err != nil {
		fmt.Println(err)
	}
}

func (player *PuppetPlayer) PreviewServo(servoID string, value float64) {
	if !player.playing {
		for _, board := range player.puppet.Boards {
			for _, servo := range board.Servos {
				if servo.ID == servoID {
					var position = int((value/100)*float64(servo.Max-servo.Min)) + servo.Min
					player.puppetDriver.SetServoPosition(servoID, position)
					return
				}
			}
		}
	}
}
