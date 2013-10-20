package gin_test

import (
	"io"
	"os"
	"os/exec"
)

type MockRunner struct {
}

func (m *MockRunner) Run() (*exec.Cmd, error) {
	return nil, nil
}

func (m *MockRunner) Info() (os.FileInfo, error) {
	return nil, nil
}

func (m *MockRunner) SetWriter(io.Writer) {
}

func (m *MockRunner) Kill() error {
	return nil
}
