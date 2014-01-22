package gin

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

type Builder interface {
	Build() error
	Binary() string
	Errors() string
}

type builder struct {
	dir    string
	binary string
	errors string
}

func NewBuilder(dir string, bin string) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") { // check if it already has the .exe extension
			bin += ".exe"
		}
	}

	return &builder{dir: dir, binary: bin}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	command := exec.Command("go", "build", "-o", b.binary)
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
