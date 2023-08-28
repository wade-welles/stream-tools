package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	x11 "github.com/wade-welles/streamkit/x11"
)

type Application struct {
	Name  string
	X11   *x11.X11
	Delay time.Duration
	Paths map[PathType]Path
}

// TODO: Initialize some Paths for our basic application, which would include,
// config (~/.config/$APP_NAME/*) => so config file could be the existance of a
// .conf file in the folder.
// Local is the local cache stored in (~/.local/share/$APP_NAME/*)
// System is only necessary if the application is capable or even desirable to
// ever be ran as root.
// Log is the path to the logs, with the ability to output to a variety of log
// files (see how CLI framework works)
type PathType uint8

const (
	Config PathType = iota
	Data
)

// NOTE: Advantage here is we can make methods on this path, like write log, or
// load or init config.
type Path string

func main() {
	// TODO: Check if not root, if root just close, don't bother with system use
	// of this tool
	userHome, _ := os.UserHomeDir()

	x11App := Application{
		Name:  "x11-cli",
		X11:   x11.X11Connect("10.101.101.1:0"),
		Delay: 2 * time.Second,
	}

	x11App.Paths = map[PathType]Path{
		Config: Path(fmt.Sprintf("%v/.config/%v", userHome, x11App.Name)),
		Data:   Path(fmt.Sprintf("%v/.local/share/%v", userHome, x11App.Name)),
	}

	fmt.Printf("%v\n", x11App.Name)
	fmt.Printf("%v\n", strings.Repeat("=", len(x11App.Name)))

	fmt.Printf("Looking for configuration: %v\n", x11App.Paths[Config])
	fmt.Printf("Storing local data: %v\n", x11App.Paths[Data])

	//x11App.X11.InitActiveWindow()

	// TODO: Probably want to load some settings from a YAML config to make things
	// easier

	//fmt.Printf("x11App:\n")

	//tick := time.Tick(x11App.Delay)
	//for {
	//	select {
	//	case <-tick:
	//		if x11App.X11.HasActiveWindowChanged() {
	//			fmt.Printf("HasActiveWindowChanged(): true\n")

	//			activeWindow := x11App.X11.ActiveWindow()
	//			fmt.Printf("  active_window_title: %s\n", activeWindow.Title)

	//			fmt.Printf("  x11.ActiveWindowTitle: %v\n", x11App.X11.ActiveWindowTitle)

	//			// NOTE: This worked to prevent it from repeating
	//			// HasActiveWindowChanged() over and over
	//			x11App.X11.CacheActiveWindow()

	//		} else {
	//			fmt.Printf("tick,...\n")
	//			fmt.Printf("  x11.ActiveWindowTitle: %v\n", x11App.X11.ActiveWindowTitle)
	//			fmt.Printf(
	//				"  x11.ActiveWindow().Type.String(): %v\n",
	//				x11App.X11.ActiveWindow().Type.String(),
	//			)
	//		}
	//	}
	//}

}
