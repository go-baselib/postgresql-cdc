package main

import (
	_ "embed"

	"github.com/1005281342/postgresql-cdc/internal/config"
	"github.com/1005281342/postgresql-cdc/internal/listener"
)

//go:embed sample.yaml
var conf []byte

func main() {
	config.InitConfig(conf)

	if err := listener.Start(); err != nil {
		panic(err)
	}
}
