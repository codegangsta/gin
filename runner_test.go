package gin_test

import (
	"bytes"
	"github.com/codegangsta/gin"
	"testing"
)

func Test_NewRunner(t *testing.T) {
	bin := "test_fixtures/writing_output"
	runner := gin.NewRunner(bin)

	fi, _ := runner.Info()
	expect(t, fi.Name(), "writing_output")
}

func Test_Runner_Run(t *testing.T) {
	bin := "test_fixtures/writing_output"
	runner := gin.NewRunner(bin)

	cmd, err := runner.Run()
	expect(t, err, nil)
	expect(t, cmd.Process == nil, false)
}

// func Test_Runner_SettingEnvironment(t *testing.T) {
// }

func Test_Runner_SetWriter(t *testing.T) {
	buff := bytes.NewBufferString("")
	expect(t, buff.String(), "")

	bin := "test_fixtures/writing_output"
	runner := gin.NewRunner(bin)
	runner.SetWriter(buff)

	cmd, err := runner.Run()
	cmd.Wait()
	expect(t, err, nil)
	expect(t, buff.String(), "Hello world\n")
}

// func Test_Runner_ThrowingError
func Test_Runner_RestartingUpdatedBinary(t *testing.T) {

}

// func Test_Runner_NotRestartingSameBinary
