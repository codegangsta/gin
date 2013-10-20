package gin

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

type Builder interface {
	Build() error
	Errors() string
}

type builder struct {
	dir    string
	errors string
}

func NewBuilder(dir string) Builder {
	return &builder{dir: dir}
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	command := exec.Command("go", "build")
	command.Dir = b.dir

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	err = command.Start()
	if err != nil {
		return err
	}

	errors, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err
	}

	b.errors = string(errors)
	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}

	return err
}
