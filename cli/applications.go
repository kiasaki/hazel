package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kiasaki/batbelt/rest"
	"github.com/mitchellh/cli"
)

type ApplicationsCommand struct {
	Ui cli.Ui
}

func (c ApplicationsCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("apps", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	apiUrl := ApiUrlFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	client, err := NewApiClient(*apiUrl)
	if err != nil {
		c.Ui.Error("Error creating Hazel client")
		return 1
	}
	response, err := client.Applications().All()
	if err == rest.ErrStatusUnauthorized {
		c.Ui.Error("Error logging in, run 'hazel login' to authenticate")
		return 1
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error contacting Hazel api: %s", err))
		return 1
	}
	for _, app := range response.Applications {
		name := app.Name
		for i := len(name); i <= 20; i++ {
			name = name + " "
		}
		c.Ui.Output(name + app.GitURL)
	}
	if len(response.Applications) == 0 {
		c.Ui.Output("You have no applications!")
	}

	return 0
}

func (c ApplicationsCommand) Help() string {
	helpText := `
Usage: hazel apps

  Lists applications you own

Options:

  -api-url=http://localhost:6201 Hazel API URL.
`
	return strings.TrimSpace(helpText)
}

func (c ApplicationsCommand) Synopsis() string {
	return "lists your applications"
}
