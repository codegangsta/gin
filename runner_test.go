package gin_test

import (
	"bytes"
	"github.com/codegangsta/gin"
	"os"
	"testing"
	"time"
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

func Test_Runner_Kill(t *testing.T) {
	bin := "test_fixtures/writing_output"
	runner := gin.NewRunner(bin)

	cmd1, err := runner.Run()
	expect(t, err, nil)

	cmd2, err := runner.Run()
	expect(t, err, nil)
	expect(t, cmd1, cmd2)

	time.Sleep(time.Second * 1)
	os.Chtimes(bin, time.Now(), time.Now())
	if err != nil {
		t.Fatal("Error with Chtimes")
	}

	cmd3, err := runner.Run()
	expect(t, err, nil)
	refute(t, cmd1, cmd3)
}

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
