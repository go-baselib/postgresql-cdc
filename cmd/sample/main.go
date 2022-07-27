package main

import (
	_ "embed"

	"github.com/go-baselib/postgresql-cdc/internal/config"
	"github.com/go-baselib/postgresql-cdc/internal/listener"
)

//go:embed sample.yaml
var conf []byte

func main() {
	config.InitConfig(conf)

	if err := listener.Start(); err != nil {
		panic(err)
	}
}
