package gin_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/codegangsta/gin/lib"
)

func Test_NewRunner(t *testing.T) {
	filename := "writing_output"
	if runtime.GOOS == "windows" {
		filename += ".bat"
	}
	bin := filepath.Join("test_fixtures", filename)

	runner := gin.NewRunner(bin)

	fi, _ := runner.Info()
	expect(t, fi.Name(), filename)
}

func Test_Runner_Run(t *testing.T) {
	bin := filepath.Join("test_fixtures", "writing_output")
	if runtime.GOOS == "windows" {
		bin += ".bat"
	}
	runner := gin.NewRunner(bin)

	cmd, err := runner.Run()
	expect(t, err, nil)
	expect(t, cmd.Process == nil, false)
}

// func Test_Runner_SettingEnvironment(t *testing.T) {
// }

func Test_Runner_Kill(t *testing.T) {
	bin := filepath.Join("test_fixtures", "writing_output")
	if runtime.GOOS == "windows" {
		bin += ".bat"
	}

	runner := gin.NewRunner(bin)

	cmd1, err := runner.Run()
	expect(t, err, nil)

	_, err = runner.Run()
	expect(t, err, nil)

	time.Sleep(time.Second * 1)
	os.Chtimes(bin, time.Now(), time.Now())
	if err != nil {
		t.Fatal("Error with Chtimes")
	}

	cmd3, err := runner.Run()
	expect(t, err, nil)

	if runtime.GOOS != "windows" {
		// does not seem to work as expected on windows
		refute(t, cmd1, cmd3)
	}
}

func Test_Runner_SetWriter(t *testing.T) {
	buff := bytes.NewBufferString("")
	expect(t, buff.String(), "")

	bin := filepath.Join("test_fixtures", "writing_output")
	if runtime.GOOS == "windows" {
		bin += ".bat"
	}

	runner := gin.NewRunner(bin)
	runner.SetWriter(buff)

	cmd, err := runner.Run()
	cmd.Wait()
	expect(t, err, nil)

	if runtime.GOOS == "windows" {
		expect(t, buff.String(), "Hello world\r\n")
	} else {
		expect(t, buff.String(), "Hello world\n")
	}
}
