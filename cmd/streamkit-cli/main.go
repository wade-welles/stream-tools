package main

import (
	"fmt"
	"os"

	streamkit "github.com/wade-welles/streamkit"

	cli "github.com/multiverse-os/cli"
)

func main() {
	fmt.Println("streamkit-cli")
	fmt.Println("=============")
	// TODO: Don't pass this value in, just read the name of the scene collection
	// and assign it as the show name
	toolkit := streamkit.New()

	cmd, initErrors := cli.New(cli.App{
		Name:        "streamkit-service",
		Description: "A long running streaming service toolkit",
		Version:     cli.Version{Major: 0, Minor: 1, Patch: 0},
		Commands: cli.Commands(
			cli.Command{
				Name:        "obs",
				Alias:       "o",
				Description: "show and item in the list",
				Subcommands: cli.Commands(
					cli.Command{
						Name:        "scene",
						Alias:       "s",
						Description: "interaction with new scene object",
						Subcommands: cli.Commands(
							cli.Command{
								Name:  "list",
								Alias: "l",
								Action: func(c *cli.Context) error {

									// TODO NEED TO FIX THIS
									//toolkit.Show.PrintDebug()
									//toolkit.Show.Cache()
									// TODO: So initialize Show and obs.Client in the NewToolkit
									// then cache it there too

									//for _, scene := range toolkit.OBS.Show.SceneNames() {

									//	fmt.Printf("scene: name(%v)\n", scene)
									//	// TODO: Now we need the scene to have items
									//}

									// TODO should do Actions:
									// now we have a simpler tool to test our stupid abstraction
									// god i <3 myself obvio 1
									return nil
								},
							},
						),
					},
				),
			},
		),
		//Actions: cli.Actions{
		//	OnStart: func(c *cli.Context) error {
		//		c.CLI.Log("OnStart action")
		//		//toolkit.
		//		return nil
		//	},
		//	Fallback: func(c *cli.Context) error {
		//		c.CLI.Log("Fallback action")
		//		return nil
		//	},
		//	OnExit: func(c *cli.Context) error {
		//		c.CLI.Log("on exit action")
		//		return nil
		//	},
		//},
	})

	if len(initErrors) == 0 {
		cmd.Parse(os.Args).Execute()
	}
}
