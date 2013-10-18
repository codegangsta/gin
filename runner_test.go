package gin_test

import (
	"github.com/codegangsta/gin"
	"testing"
)

func Test_NewRunner(t *testing.T) {
	bin := "test_fixtures/build_success/build_success"
	runner := gin.NewRunner(bin)

	expect(t, runner.BinInfo().Name(), "build_success")
}

func Test_Runner_Run(t *testing.T) {
	bin := "test_fixtures/build_success/build_success"
	runner := gin.NewRunner(bin)

	command, err := runner.Run()
}

// func Test_Runner_SettingEnvironment
// func Test_Runner_WritingOutput
// func Test_Runner_ThrowingError
// func Test_Runner_RestartingUpdatedBinary
