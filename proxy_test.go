package gin_test

import (
	"github.com/codegangsta/gin"
	"testing"
)

func Test_NewProxy(t *testing.T) {
	proxy := gin.NewProxy(builder, runner)

	expect(t, proxy != nil, true)
}
