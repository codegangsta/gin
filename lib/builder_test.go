package gin_test

import (
	"github.com/codegangsta/gin/lib"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_Builder_Build_Success(t *testing.T) {
	wd := filepath.Join("test_fixtures", "build_success")
	bin := "build_success"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	builder := gin.NewBuilder(wd, bin, false)
	err := builder.Build()
	expect(t, err, nil)

	file, err := os.Open(filepath.Join(wd, bin))
	if err != nil {
		t.Fatalf("File has not been written: %v", err)
	}

	refute(t, file, nil)
}
