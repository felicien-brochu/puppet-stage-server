package main

import (
	"os/exec"
	"strings"
)

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
