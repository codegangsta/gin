package gin

import (
	"os"
	"os/exec"
)

type Runner interface {
	Run() (*exec.Cmd, error)
	Info() (os.FileInfo, error)
}

type runner struct {
	bin string
}

func NewRunner(bin string) Runner {
	return &runner{bin: bin}
}

func (r *runner) Run() (*exec.Cmd, error) {
	command := exec.Command(r.bin)
	//stdout, err := command.StdoutPipe()
	err := command.Start()

	return command, err
}

func (r *runner) Info() (os.FileInfo, error) {
	return os.Stat(r.bin)
}
