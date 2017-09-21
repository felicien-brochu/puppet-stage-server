package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

type PositionCommand struct {
	servoID int
	P       int
}

func (cmd PositionCommand) String() string {
	return fmt.Sprintf("P%d;%d\r", cmd.servoID, cmd.P)
}

var serialPort *serial.Port

func readSerialLoop() {
	buffer := make([]byte, 128)
	for true {
		n, err := serialPort.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(buffer[:n]))
	}
}

func main() {
	ports, err := ListSerialPorts()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(ports))
	if len(ports) == 0 {
		log.Fatal("No active serial port")
	}
	conf := &serial.Config{Name: ports[0], Baud: 115200}
	serialPort, err = serial.OpenPort(conf)
	if err != nil {
		log.Fatal(err)
	}

	go readSerialLoop()

	time.Sleep(2000 * time.Millisecond)
	servoID := 0
	speed := 700.
	delay := time.Duration(20)
	step := speed / (1. / 0.02)
	max := 470.
	min := 130.

	for i := 0; i < 6; i++ {
		for p := min; p < max; p += step {
			cmd1 := PositionCommand{
				servoID,
				int(p),
			}
			_, err = serialPort.Write([]byte(cmd1.String()))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(delay * time.Millisecond)
		}

		for p := max; p > min; p -= step {
			cmd1 := PositionCommand{
				servoID,
				int(p),
			}
			_, err = serialPort.Write([]byte(cmd1.String()))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond * delay)
		}
	}

	fmt.Print("END##########")

}
