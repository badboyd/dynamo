package main

import (
	"flag"

	"github.com/badboyd/dynamo/config"
	"github.com/badboyd/dynamo/internal/server"
)

var (
	configFile = flag.String("config", "config/default.yml", "Configuration file path")
)

func init() {
	flag.Parse()
}

func main() {
	conf, err := config.Load(*configFile)
	if err != nil {
		panic(err)
	}

	server := server.New(conf)
	server.Start()
}
