package gin_test

import (
	"io"
	"os"
	"os/exec"
)

type MockRunner struct {
	DidRun bool
}

func NewMockRunner() *MockRunner {
	return &MockRunner{
		DidRun: false,
	}
}

func (m *MockRunner) Run() (*exec.Cmd, error) {
	m.DidRun = true
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

type MockBuilder struct {
	MockErrors string
}

func NewMockBuilder() *MockBuilder {
	return &MockBuilder{}
}

func (m *MockBuilder) Binary() string {
	return "bin"
}

func (m *MockBuilder) Build() error {
	return nil
}

func (m *MockBuilder) Errors() string {
	return m.MockErrors
}
