package players

import (
	"felicien/puppet-server/model"
	"math"
	"strconv"
	"time"
)

const frameInterval = 10 * time.Millisecond

type stagePlayer struct {
	stage     model.Stage
	puppet    model.Puppet
	playStart model.Time
	started   bool
	startTime time.Time
	done      chan struct{}
	stateChan chan string
	ticker    chan time.Time
}

func newStagePlayer(stage model.Stage, puppet model.Puppet, playStart model.Time) stagePlayer {
	return stagePlayer{
		stage:     stage,
		puppet:    puppet,
		playStart: playStart,
		started:   false,
	}
}

func (player *stagePlayer) play(puppetPlayer *PuppetPlayer, stateChan chan string, ticker chan time.Time) {
	if player.started {
		player.stop()
	}
	player.done = make(chan struct{})
	player.startTime = time.Now()
	player.started = true
	player.stateChan = stateChan
	player.ticker = ticker

	go player.playRoutine(puppetPlayer)
}

func (player *stagePlayer) stop() {
	if player.started {
		close(player.done)
		player.started = false
	}
}

func drainTicker(ticker chan time.Time) {
	for {
		select {
		case <-ticker:
		default:
			return
		}
	}
}

func (player *stagePlayer) playRoutine(puppetPlayer *PuppetPlayer) {
	player.stateChan <- "start"
	drainTicker(player.ticker)

	player.playFrame(puppetPlayer, player.getCurrentTime())
	endTime := model.Time(player.stage.Duration)

MainLoop:
	for {
		select {
		case <-player.ticker:
			t := player.getCurrentTime()

			if t.Before(endTime) {
				player.stateChan <- strconv.Itoa(int(t))
				player.playFrame(puppetPlayer, t)
			} else {
				t = endTime
				player.stateChan <- strconv.Itoa(int(t))
				player.playFrame(puppetPlayer, t)
				break MainLoop
			}
		case <-player.done:
			break MainLoop
		}
	}
	player.stateChan <- "stop"
	close(player.stateChan)
}

func (player *stagePlayer) playFrame(puppetPlayer *PuppetPlayer, t model.Time) {
	var frame = player.stage.GetFrameAt(t)

	for servoID, value := range frame {
		var driverSequence model.DriverSequence
		for _, driverSequenceItem := range player.stage.Sequences {
			if driverSequenceItem.ServoID == servoID {
				driverSequence = driverSequenceItem
				break
			}
		}

		var servo model.Servo
		for _, board := range player.puppet.Boards {
			if servoItem, ok := board.Servos[servoID]; ok {
				servo = servoItem
				break
			}
		}

		var position int
		if math.IsNaN(value) || !driverSequence.PlayEnabled {
			position = servo.DefaultPosition
		} else {
			position = int((value/100)*float64(servo.Max-servo.Min)) + servo.Min
		}

		puppetPlayer.playServoPosition(servoID, position)
	}
}

func (player *stagePlayer) getCurrentTime() model.Time {
	return model.Time(time.Since(player.startTime)) + player.playStart
}
