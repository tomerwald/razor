package commands

import (
	"encoding/json"
	"log"
	"os/exec"
	"syscall"
)

type ExecCommand struct {
	ExecutablePath string
	Params         string
}

func (ec *ExecCommand) run() ([]byte, error) {
	log.Printf("running: %s %s\r\n", ec.ExecutablePath, ec.Params)
	cmdInstance := exec.Command(ec.ExecutablePath, ec.Params)
	cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmdInstance.Output()
}

func RunExec(payload []byte) ([]byte, error) {
	var m ExecCommand
	err := json.Unmarshal(payload, &m)
	if err != nil {
		return nil, err
	} else {
		out, err := m.run()
		if err != nil {
			return []byte(string(out) + " " + err.Error()), nil
		} else {
			return out, nil
		}
	}
}
