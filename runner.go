package gin

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type Runner interface {
	Run() (*exec.Cmd, error)
	Info() (os.FileInfo, error)
	SetWriter(io.Writer)
	Kill() error
}

type runner struct {
	bin       string
	writer    io.Writer
	command   *exec.Cmd
	starttime time.Time
}

func NewRunner(bin string) Runner {
	return &runner{
		bin:       bin,
		writer:    ioutil.Discard,
		starttime: time.Now(),
	}
}

func (r *runner) Run() (*exec.Cmd, error) {
	if r.needsRefresh() {
		r.Kill()
	}

	if r.command == nil {
		err := r.runBin()
		time.Sleep(250 * time.Millisecond)
		return r.command, err
	} else {
		return r.command, nil
	}
}

func (r *runner) Info() (os.FileInfo, error) {
	return os.Stat(r.bin)
}

func (r *runner) SetWriter(writer io.Writer) {
	r.writer = writer
}

func (r *runner) Kill() error {
	if r.command != nil && r.command.Process != nil {
		if err := r.command.Process.Kill(); err != nil {
			return err
		}
		r.command = nil
	}

	return nil
}

func (r *runner) runBin() error {
	r.command = exec.Command(r.bin)
	stdout, err := r.command.StdoutPipe()
	if err != nil {
		return err
	}

	err = r.command.Start()
	if err != nil {
		return err
	}

	r.starttime = time.Now()

	go io.Copy(r.writer, stdout)
	return nil
}

func (r *runner) needsRefresh() bool {
	info, err := r.Info()
	if err != nil {
		return false
	} else {
		return info.ModTime().After(r.starttime)
	}
}
