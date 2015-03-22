package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var appname = "hazel-cli"
var version = "0.1.0"

func main() {
	c := cli.NewCLI(appname, version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"foo": fooCommandFactory,
		"bar": barCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
