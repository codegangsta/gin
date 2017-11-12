package gin_test

import (
	"path/filepath"
	"testing"

	"github.com/codegangsta/gin/lib"
)

func Test_NewTester(t *testing.T) {
	dir := filepath.Join("test_fixtures", "build_success")

	tester := gin.NewTester(dir, []string{})
	err := tester.Run()
	expect(t, err, nil)
}
