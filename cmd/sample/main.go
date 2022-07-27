package main

import (
	_ "embed"

	"github.com/go-baselib/postgresql-cdc/listener"
)

//go:embed sample.yaml
var conf []byte

func main() {
	if err := listener.Start(conf); err != nil {
		panic(err)
	}
}
