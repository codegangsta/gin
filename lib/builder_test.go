package gin_test

import (
	"fmt"
	"github.com/codegangsta/gin/lib"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_Builder_Build_Success(t *testing.T) {
	wd := filepath.Join("test_fixtures", "build_success")
	bin := "build_success"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	builder := gin.NewBuilder(wd, bin)
	err := builder.Build()
	expect(t, err, nil)

	file, err := os.Open(filepath.Join(wd, bin))
	if err != nil {
		t.Fatalf("File has not been written: %v", err)
	}

	refute(t, file, nil)
}

func Test_Builder_Build_Failure(t *testing.T) {
	wd := filepath.Join("test_fixtures", "build_failure")

	builder := gin.NewBuilder(wd, "bin")
	err := builder.Build()
	refute(t, err, nil)

	expect(t, strings.Contains(builder.Errors(), fmt.Sprintf(".%smain.go:4: undefined: this", string(os.PathSeparator))), true)
	expect(t, strings.Contains(builder.Errors(), fmt.Sprintf(".%smain.go:4: undefined: compile", string(os.PathSeparator))), true)
}
