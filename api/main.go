package main

import (
	"net/http"

	"github.com/jessevdk/go-flags"
	"github.com/maxwellhealth/bongo"
	"github.com/zenazn/goji"
)

var cfg struct {
	Name string `short:"n" long:"name" description:"A name" required:"true"`
}

var connection bongo.Connection

func main() {
	initConfig()
	initDb()
}

func initConfig() {
	_, err := flags.Parse(&cfg)
	if err != nil {
		panic(err)
	}
}

func initDb() {
	var err error

	config := &bongo.Config{
		ConnectionString: "localhost",
		Database:         "bongotest",
	}

	connection, err = bongo.Connect(config)
	if err != nil {
		panic(err)
	}
}
