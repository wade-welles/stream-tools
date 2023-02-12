package streamkit

import (
	"fmt"
	"strings"
	"time"

	obs "github.com/wade-welles/streamkit/obs"
	x11 "github.com/wade-welles/streamkit/x11"
)

type Toolkit struct {
	// NOTE: Short-poll rate [we will rewrite without short polling after]
	Delay time.Duration
	OBS   *obs.OBS
	X11   *x11.X11
}

// TODO: Could pass the host for OBS and the host for X11 as variadic strings so
// it can be empty, or provide position 1 for obs position 2 for x11 (though x11
// should assumingly always be 127.0.0.1 whereas obs reasonably could be
// different
func New() Toolkit {
	toolkit := Toolkit{
		OBS: &obs.OBS{
			Client: obs.ConnectToOBS("192.168.1.1:4444"),
		},
		X11: &x11.X11{
			Client: x11.ConnectToX11(),
		},
		Delay: 1500 * time.Millisecond,
	}
	// TODO: Pull showname from SceneCollection.ScName

	// TODO: The scenes and show object should be populated based on whatever
	//       the scene collection that is currently active is. but keep in
	//       mind the goal is to abstract awwy some of the less good design
	//       bits into a better logical construct

	toolkit.OBS.Show = &obs.Show{
		OBS:    toolkit.OBS,
		Name:   "she hacked you",
		Scenes: obs.Scenes{}, // TODO: Maybe populate on initialization
	}
	toolkit.OBS.Show.Cache()
	toolkit.X11.InitActiveWindow()
	return toolkit
}

func (t Toolkit) HandleWindowEvents() {

	fmt.Printf("Number of scenes parsed: %v\n", len(t.OBS.Show.Scenes))

	for _, scene := range t.OBS.Show.Scenes {
		fmt.Printf("scene:\n")
		fmt.Printf("  name: %v\n", scene.Name)
		fmt.Printf("  item_count: %v\n", len(scene.Items))
		fmt.Printf("  items:\n")
		for _, item := range scene.Items {
			fmt.Printf("    item:\n")
			fmt.Printf("      name: %v\n", item.Name)
		}
	}

	fmt.Printf("Names of scenes: %v\n", strings.Join(t.OBS.Show.SceneNames(), ", "))

	primaryScene, ok := t.OBS.Show.Scene("Primary")
	if ok {
		// TODO: We need to cache or initialize the items in a given scene
	} else {
		panic(fmt.Errorf("failed to locate primary scene"))
	}
	bumperScene, ok := t.OBS.Show.Scene("Bumper")
	if !ok {
		panic(fmt.Errorf("failed to locate bumper scene"))
	}

	// TODO: This lookup is not connecting with the parsed items when we are
	// running .Cache() on each scene as its parsed in Show.Cache()
	fmt.Printf("# of primaryScene items: %v\n", len(primaryScene.Items))

	vimWindowName := "Primary Terminal (VIM Window)"
	vimWindow, ok := primaryScene.Item(vimWindowName)
	if ok {
		fmt.Printf("found vimWindow item: %v\n", vimWindow)
	} else {
		panic(
			fmt.Errorf(
				"failed to find item '" + vimWindowName + "' within the primary scene",
			),
		)
	}

	consoleWindow, _ := primaryScene.Item("CONSOLE")
	chromiumWindow, _ := primaryScene.Item("CHROMIUM")

	tick := time.Tick(t.Delay)
	for {
		select {
		case <-tick:
			if t.X11.HasActiveWindowChanged() {
				time.Sleep(4 * time.Second)

				currentScene := t.OBS.Show.Current
				activeWindow := t.X11.ActiveWindow()

				if currentScene.HasName("content:primary") {
					switch activeWindow {
					case x11.Primary, x11.Secondary:
						t.X11.CacheActiveWindow()
						if !vimWindow.Visible {
							bumperScene.Transition()
							primaryScene.Transition(4 * time.Second)

							chromiumWindow.Hide().Lock().Update()
							vimWindow.Unhide().Lock().Update()
							consoleWindow.Unhide().Lock().Update()
						}
					case x11.Chromium:
						t.X11.CacheActiveWindow()
						if !chromiumWindow.Visible {
							bumperScene.Transition()
							primaryScene.Transition(4 * time.Second)

							vimWindow.Hide().Lock().Update()
							consoleWindow.Unhide().Lock().Update()
							chromiumWindow.Unhide().Lock().Update()
						}
					default: // UndefinedName
						// TODO: We should never allow this condition to ever occur, and by
						// doing that we optimize it further bc we are not checking conditions
						// that we dont want to begin with
						fmt.Println("[undefined] active window?(%v)", t.X11.ActiveWindow())
					}
				}
			}
			// TODO: Check what the active widnow currently is; then use obs-ws to
			// change the scenes with the bumper in between
		}
	}

	//time.Sleep(2 * time.Second)
	//t.AvatarToggle()
}

// TODO: Put this back together once we ahve scenes and items parsed
//func (t Toolkit) AvatarToggle() {
//	if primaryScene, ok := t.OBS.Show.Scene("content:primary"); ok {
//
//		dynamicAvatar, _ := primaryScene.Item("dynamic avatar")
//		staticAvatar, _ := primaryScene.Item("static avatar")
//
//		//primaryScene.Update()
//		//primaryScene.Cache()
//
//		//dynamicAvatar.Update()
//		//dynamicAvatar.Cache()
//
//		if dynamicAvatar.Visible {
//			staticAvatar.Print()
//			staticAvatar.Unhide().Update()
//			staticAvatar.Print()
//			fmt.Printf("---\n")
//
//			dynamicAvatar.Print()
//			dynamicAvatar.Hide().Update()
//			dynamicAvatar.Print()
//		} else {
//			dynamicAvatar.Print()
//			dynamicAvatar.Unhide().Update()
//			dynamicAvatar.Print()
//
//			fmt.Printf("---\n")
//
//			staticAvatar.Print()
//			staticAvatar.Hide().Update()
//			staticAvatar.Print()
//		}
//	}
//}
