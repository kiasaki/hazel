package main

import (
	"flag"
	"io/ioutil"
	"strings"

	"github.com/kiasaki/batbelt/rest"
	"github.com/mitchellh/cli"
)

type LoginCommand struct {
	Ui cli.Ui
}

func (c LoginCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("login", flag.ContinueOnError)
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
	email, err := c.Ui.Ask("Enter email:")
	if err != nil {
		return 1
	}
	password, err := c.Ui.Ask("Enter password (not hidden):")
	if err != nil {
		return 1
	}

	token, err := client.Auth().Login(email, password)
	if err != nil {
		if err == rest.ErrStatusNotFound {
			c.Ui.Error("Error: no user associated to this email address")
		} else if err == rest.ErrStatusBadRequest {
			c.Ui.Error("Error: email/password combination is invalid")
		} else {
			c.Ui.Error("Error: " + err.Error())
		}
		return 1
	}

	err = ioutil.WriteFile(tokenFilename, []byte(token), 0600)
	if err != nil {
		c.Ui.Error("Error saving token to disk: " + err.Error())
		return 1
	}

	c.Ui.Info("Login successful!")
	return 0
}

func (c LoginCommand) Help() string {
	helpText := `
Usage: hazel login

  Log in with your Hazel credentials. Input is accepted by typing
  on the terminal.

Options:

  -api-url=http://localhost:6201 Hazel API URL.
`
	return strings.TrimSpace(helpText)
}

func (c LoginCommand) Synopsis() string {
	return "Login with your Hazel credentials"
}
