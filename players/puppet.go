package players

import (
	"felicien/puppet-server/drivers"
	"felicien/puppet-server/model"
	"log"
)

// PuppetPlayer is used to play a stage on its puppet or preview a stage on its puppet.
type PuppetPlayer struct {
	puppet       model.Puppet
	puppetDriver *drivers.PuppetDriver
	playing      bool
	stagePlayer  *stagePlayer
}

// NewPuppetPlayer creates a new puppet player
func NewPuppetPlayer(puppet model.Puppet, puppetDriver *drivers.PuppetDriver) *PuppetPlayer {
	return &PuppetPlayer{
		puppet:       puppet,
		puppetDriver: puppetDriver,
		playing:      false,
	}
}

// PlayStage creates a stagePlayer and start playing the stage. Reports the state of playing by stateChan.
func (player *PuppetPlayer) PlayStage(stage model.Stage, playStart model.Time, stateChan chan string) {
	player.stagePlayer = newStagePlayer(stage, playStart)
	player.playing = true
	player.stagePlayer.play(player, stateChan, player.puppetDriver.GetSenderTicker())
}

// StopStage stops playing the stage
func (player *PuppetPlayer) StopStage() {
	player.stagePlayer.stop()
	player.playing = false
}

func (player *PuppetPlayer) playServoPosition(servoID string, position int) {
	err := player.puppetDriver.SetServoPosition(servoID, position)
	if err != nil {
		log.Println(err)
	}
}

// PreviewStage preview the servo positions of the stage at t time.
func (player *PuppetPlayer) PreviewStage(stage model.Stage, t model.Time) {
	playFrame(stage, t, player, true)
	<-player.puppetDriver.GetSenderTicker()
}
