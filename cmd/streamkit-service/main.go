package main

import (
	"os"

	streamkit "github.com/wade-welles/streamkit"

	cli "github.com/multiverse-os/cli"
)

// NOTE:
// OBS Augmentation Software
// The initial goal of this software is to augment OBS

func main() {
	streamkit := streamkit.New()

	cmd, initErrors := cli.New(cli.App{
		Name:        "obs-service",
		Description: "A long running obs service toolkit",
		Version:     cli.Version{Major: 0, Minor: 1, Patch: 0},
		Actions: cli.Actions{
			OnStart: func(c *cli.Context) error {
				c.CLI.Log("OnStart action: streamkit.HandleWindowEvents()")
				streamkit.HandleWindowEvents()
				return nil
			},
			//OnExit: func(c *cli.Context) error {
			//	c.CLI.Log("on exit action")
			//	return nil
			//},
		},
	})

	if len(initErrors) == 0 {
		cmd.Parse(os.Args).Execute()
	}
}
