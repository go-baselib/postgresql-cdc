package config

import (
	"testing"

	_ "embed"
)

//go:embed test.yaml
var test []byte

func TestInitConfig(t *testing.T) {
	InitConfig(test)
}

func TestGetConfig(t *testing.T) {
	t.Log(GetConfig())
}
