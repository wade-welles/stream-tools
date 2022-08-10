package obstools

import (
	"fmt"
	"time"

	obs "github.com/wade-welles/obs-tools/obs"
	x11 "github.com/wade-welles/obs-tools/x11"
)

type Toolkit struct {
	// NOTE: Short-poll rate [we will rewrite without short polling after]
	Delay time.Duration
	OBS   *obs.OBS
	X11   *x11.X11
}

func NewToolkit() Toolkit {
	toolkit := Toolkit{
		OBS: &obs.OBS{
			Client: obs.ConnectToOBS(),
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
		Scenes: obs.Scenes{},
	}
	toolkit.OBS.Show.Cache()
	toolkit.X11.InitActiveWindow()
	return toolkit
}

func (t Toolkit) HandleWindowEvents() {
	primaryScene, _ := t.OBS.Show.Scene("content:primary")
	bumperScene, _ := t.OBS.Show.Scene("content:bumper")

	vimWindow, _ := primaryScene.Item("VIM")
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
				switch activeWindow {
				case x11.Primary, x11.Secondary:
					t.X11.CacheActiveWindow()
					if currentScene.HasName("content:primary") {
						if !vimWindow.Visible {
							bumperScene.Transition()
							primaryScene.Transition(4 * time.Second)

							chromiumWindow.Hide().Lock().Update()
							vimWindow.Unhide().Lock().Update()
							consoleWindow.Unhide().Lock().Update()
						}
					}
				case x11.Chromium:
					t.X11.CacheActiveWindow()
					if currentScene.HasName("content:primary") {
						if !chromiumWindow.Visible {
							bumperScene.Transition()
							primaryScene.Transition(4 * time.Second)

							vimWindow.Hide().Lock().Update()
							consoleWindow.Unhide().Lock().Update()
							chromiumWindow.Unhide().Lock().Update()
						}
					}

				default: // UndefinedName
					// TODO: We should never allow this condition to ever occur, and by
					// doing that we optimize it further bc we are not checking conditions
					// that we dont want to begin with
					fmt.Println("[undefined] active window?(%v)", t.X11.ActiveWindow())
				}
				// TODO: Check what the active widnow currently is; then use obs-ws to
				// change the scenes with the bumper in between
			} else {
				fmt.Printf("no ACTIVE window change\n")
			}

			//time.Sleep(2 * time.Second)
			//t.AvatarToggle()
		}
	}
}

func (t Toolkit) AvatarToggle() {
	if primaryScene, ok := t.OBS.Show.Scene("content:primary"); ok {

		dynamicAvatar, _ := primaryScene.Item("dynamic avatar")
		staticAvatar, _ := primaryScene.Item("static avatar")

		//primaryScene.Update()
		//primaryScene.Cache()

		//dynamicAvatar.Update()
		//dynamicAvatar.Cache()

		if dynamicAvatar.Visible {
			staticAvatar.Print()
			staticAvatar.Unhide().Update()
			staticAvatar.Print()
			fmt.Printf("---\n")

			dynamicAvatar.Print()
			dynamicAvatar.Hide().Update()
			dynamicAvatar.Print()
		} else {
			dynamicAvatar.Print()
			dynamicAvatar.Unhide().Update()
			dynamicAvatar.Print()

			fmt.Printf("---\n")

			staticAvatar.Print()
			staticAvatar.Hide().Update()
			staticAvatar.Print()
		}
	}
}
