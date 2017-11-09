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
	playStart model.Time
	started   bool
	startTime time.Time
	done      chan struct{}
	stateChan chan string
	ticker    chan time.Time
}

func newStagePlayer(stage model.Stage, playStart model.Time) *stagePlayer {
	return &stagePlayer{
		stage:     stage,
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

	endTime := model.Time(player.stage.Duration)

MainLoop:
	for {
		select {
		case <-player.ticker:
			t := player.getCurrentTime()

			if t.Before(endTime) {
				player.stateChan <- strconv.Itoa(int(t))
				err := playFrame(player.stage, t, puppetPlayer, false)
				if err != nil {
					break MainLoop
				}
			} else {
				t = endTime
				player.stateChan <- strconv.Itoa(int(t))
				playFrame(player.stage, t, puppetPlayer, false)
				break MainLoop
			}
		case <-player.done:
			break MainLoop
		}
	}
	player.stateChan <- "stop"
	close(player.stateChan)
}

func playFrame(stage model.Stage, t model.Time, puppetPlayer *PuppetPlayer, preview bool) error {
	var frame = stage.GetFrameAt(t, preview)

	for servoID, value := range frame {
		var servo model.Servo
		for _, board := range puppetPlayer.puppet.Boards {
			if servoItem, ok := board.Servos[servoID]; ok {
				servo = servoItem
				break
			}
		}

		var position int
		if math.IsNaN(value) {
			position = servo.DefaultPosition
		} else {
			min := servo.Min
			max := servo.Max
			if servo.Inverted {
				min = servo.Max
				max = servo.Min
			}
			position = int((value/100)*float64(max-min)) + min
		}
		err := puppetPlayer.playServoPosition(servoID, position)
		if err != nil {
			return err
		}
	}
	return nil
}

func (player *stagePlayer) getCurrentTime() model.Time {
	return model.Time(time.Since(player.startTime)) + player.playStart
}
