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
	}
}

func (r *runner) Run() (*exec.Cmd, error) {
	if r.needsRefresh() {
    println("refreshing")
		// r.Kill()
	}

	if r.command == nil {
		err := r.runBin()
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
		//r.command.Process.Release()
		r.command.Process.Kill()
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

	//r.starttime = time.Now()

	go io.Copy(r.writer, stdout)
	return nil
}

func (r *runner) needsRefresh() bool {
  _, err := r.Info()
  if err != nil {
    return false
  } else {
    return true
  }
}
