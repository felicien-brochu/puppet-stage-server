package players

import (
	"errors"
	"felicien/puppet-server/model"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/tarm/serial"
)

// Bus describes a communication bus (serial COM bus, wifi etc.)
type Bus interface {
	Open() error
	Close() error
	Reader() io.Reader
	Writer() io.Writer
}

// NewBus creates a new bus according to given bus type
func NewBus(busType model.BusType) (Bus, error) {
	if busType == model.BusTypeSerial {
		return DefaultSerialBus()
	}
	return nil, errors.New("BusType not supported")
}

// SerialBus describes a serial bus
type SerialBus struct {
	portConfig *serial.Config
	port       *serial.Port
}

// NewSerialBus constructor
func NewSerialBus(name string, baud int) *SerialBus {
	serialBus := new(SerialBus)
	serialBus.portConfig = &serial.Config{Name: name, Baud: baud}
	return serialBus
}

// DefaultSerialBus returns the first available SerialBus, error if none
func DefaultSerialBus() (*SerialBus, error) {
	ports, err := ListSerialPorts()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("No active serial port")
	}

	serialBus := new(SerialBus)
	serialBus.portConfig = &serial.Config{Name: ports[0], Baud: 115200}

	return serialBus, err
}

// Open opens the bus
func (serialBus *SerialBus) Open() error {
	if serialBus.portConfig == nil {
		return errors.New("SerialBus.Connect() error: no serial portConfig in SerialBus")
	}

	port, err := serial.OpenPort(serialBus.portConfig)
	if err != nil {
		return err
	}

	serialBus.port = port
	return nil
}

// Close closes the SerialBus
func (serialBus *SerialBus) Close() error {
	if serialBus.port == nil {
		return nil
	}
	return serialBus.port.Close()
}

// Reader returns a reader to the SerialBus
func (serialBus *SerialBus) Reader() io.Reader {
	return serialBus.port
}

// Writer returns a writer to the SerialBus
func (serialBus *SerialBus) Writer() io.Writer {
	return serialBus.port
}

// ListSerialPorts returns a list of active serial ports
func ListSerialPorts() ([]string, error) {
	cmd := exec.Command("powershell", "[System.IO.Ports.SerialPort]::GetPortNames()")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	n, err := stdout.Read(buffer)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(buffer[:n]), "\n"), nil
}
