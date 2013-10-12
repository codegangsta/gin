package gin_test

import (
	"github.com/codegangsta/gin"
	"testing"
)

func Test_LoadConfig(t *testing.T) {
	config, err := gin.LoadConfig("test_fixtures/config.json")

	expect(t, err, nil)
	expect(t, config.Port, 5678)
	expect(t, config.Server.Port, 3000)
}

func Test_LoadConfig_WithNonExistantFile(t *testing.T) {
	_, err := gin.LoadConfig("im/not/here.json")

	refute(t, err, nil)
}

func Test_LoadConfig_WithMalformedFile(t *testing.T) {
	_, err := gin.LoadConfig("test_fixtures/bad_config.json")

	refute(t, err, nil)
}
