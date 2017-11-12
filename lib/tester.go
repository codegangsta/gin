package gin

import (
	"fmt"
	"os/exec"
	"strings"
)

type tester struct {
	dir      string
	testArgs []string
}

func NewTester(dir string, testArgs []string) *tester {
	return &tester{
		dir:      dir,
		testArgs: testArgs,
	}
}

func (t *tester) Run() error {
	testArgs := strings.Join(t.testArgs, " ")

	args := append([]string{"sh", "-c", "go test $(go list ./... | grep -v /vendor/) " + testArgs})

	var command *exec.Cmd

	command = exec.Command(args[0], args[1:]...)
	command.Dir = t.dir
	output, err := command.CombinedOutput()

	fmt.Print(string(output))

	return err
}
