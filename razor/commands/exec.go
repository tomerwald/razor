package commands

import (
	"encoding/json"
	"log"
	"os/exec"
)

type ExecCommand struct {
	ExecutablePath string
	Params         string
}

func (ec *ExecCommand) run() ([]byte, error) {
	log.Printf("running: %s %s\r\n", ec.ExecutablePath, ec.Params)
	return exec.Command(ec.ExecutablePath, ec.Params).Output()
}

func RunExec(payload []byte) ([]byte, error) {
	var m ExecCommand
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		return m.run()
	}
}
