package gin

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Runner interface {
	Run() (*exec.Cmd, error)
	Info() (os.FileInfo, error)
	SetWriter(io.Writer)
}

type runner struct {
	bin    string
	writer io.Writer
}

func NewRunner(bin string) Runner {
	return &runner{
		bin:    bin,
		writer: ioutil.Discard,
	}
}

func (r *runner) Run() (*exec.Cmd, error) {
	command := exec.Command(r.bin)
	stdout, err := command.StdoutPipe()
	if err != nil {
		return command, err
	}

	err = command.Start()
	if err != nil {
		return command, err
	}

	go io.Copy(r.writer, stdout)
	return command, err
}

func (r *runner) Info() (os.FileInfo, error) {
	return os.Stat(r.bin)
}

func (r *runner) SetWriter(writer io.Writer) {
	r.writer = writer
}
