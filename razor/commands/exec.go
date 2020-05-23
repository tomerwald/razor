package commands

import (
	"encoding/json"
	"log"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/net/context"
)

type ExecResponse struct {
	Error  string
	Stdout string
}

func (tr *ExecResponse) toBuffer() []byte {
	out, _ := json.Marshal(tr)
	return out
}

// NewExecResponse is a ExecResponse constructor
func NewExecResponse(execError string, content string) *ExecResponse {
	return &ExecResponse{
		Error:  execError,
		Stdout: content,
	}
}

// ExecCommand is a payload type used to run executables on the machine
type ExecCommand struct {
	ExecutablePath string
	Params         string
	Timeout        int
}

func (ec *ExecCommand) run() ([]byte, error) {
	log.Printf("running: %s %s\r\n", ec.ExecutablePath, ec.Params)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ec.Timeout)*time.Millisecond)
	defer cancel()
	cmdInstance := exec.CommandContext(ctx, ec.ExecutablePath, ec.Params)
	cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmdInstance.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, err
	}
	return output, err
}

func RunExec(payload []byte) []byte {
	var m ExecCommand
	err := json.Unmarshal(payload, &m)
	if err != nil {
		return []byte(commandError{"Failed to unmarshal command"}.Error())
	}
	out, err := m.run()
	if err != nil {
		return NewExecResponse(err.Error(), string(out)).toBuffer()
	}
	return NewExecResponse("", string(out)).toBuffer()
}
