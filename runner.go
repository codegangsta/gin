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
	return nil, nil
}

func (r *runner) Info() (os.FileInfo, error) {
	return os.Stat(r.bin)
}
