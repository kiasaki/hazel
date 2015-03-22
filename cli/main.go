package main

import (
	"log"
	"os"
	"path"

	"github.com/mitchellh/cli"
)

var (
	appname       string = "hazel-cli"
	version              = "0.1.0"
	tokenFilename        = path.Join(os.Getenv("HOME"), ".hazel_token")
)

func main() {
	ui := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorBlue,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader: os.Stdin,
			Writer: os.Stdout,
		},
	}

	c := cli.NewCLI(appname, version)
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"apps": func() (cli.Command, error) {
			return &ApplicationsCommand{ui}, nil
		},
		"login": func() (cli.Command, error) {
			return &LoginCommand{ui}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
