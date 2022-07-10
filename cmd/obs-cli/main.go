package main

import (
	"fmt"
	"os"

	obs "github.com/wade-welles/obs-tools"

	cli "github.com/multiverse-os/cli"
)

func main() {
	fmt.Println("obs-cli")
	fmt.Println("===========")
	// TODO: Don't pass this value in, just read the name of the scene collection
	// and assign it as the show name
	toolkit := obs.NewToolkit()

	fmt.Println("Toolkit for purpose of building CLI interface: %v", toolkit)

	cmd, initErrors := cli.New(cli.App{
		Name:        "obs-service",
		Description: "A long running obs service toolkit",
		Version:     cli.Version{Major: 0, Minor: 1, Patch: 0},
		Actions: cli.Actions{
			OnStart: func(c *cli.Context) error {
				c.CLI.Log("OnStart action")
				return nil
			},
			Fallback: func(c *cli.Context) error {
				c.CLI.Log("Fallback action")
				return nil
			},
			OnExit: func(c *cli.Context) error {
				c.CLI.Log("on exit action")
				return nil
			},
		},
	})

	if len(initErrors) == 0 {
		cmd.Parse(os.Args).Execute()
	}
}
