package gin_test

import (
	"github.com/codegangsta/gin"
	"os"
	"strings"
	"testing"
)

func Test_Builder_Build_Success(t *testing.T) {
	wd := "test_fixtures/build_success/"

	builder := gin.NewBuilder(wd)
	err := builder.Build()
	expect(t, err, nil)

	file, err := os.Open(wd + "build_success")

	if err != nil {
		t.Fatal("File has not been written")
	}

	refute(t, file, nil)
}

func Test_Builder_Build_Failure(t *testing.T) {
	wd := "test_fixtures/build_failure/"

	builder := gin.NewBuilder(wd)
	err := builder.Build()
	refute(t, err, nil)

	expect(t, strings.Contains(builder.Errors(), "./main.go:4: undefined: this"), true)
	expect(t, strings.Contains(builder.Errors(), "./main.go:4: undefined: compile"), true)
}
