package obstools

import (
	"fmt"
	"time"
)

type Toolkit struct {
	// NOTE: Short-poll rate [we will rewrite without short polling after]
	Delay time.Duration
	OBS   *OBS
	X11   *X11
}

func NewToolkit() Toolkit {
	toolkit := Toolkit{
		OBS: &OBS{
			Client: ConnectToOBS(),
		},
		X11: &X11{
			Client: ConnectToX11(),
		},
		Delay: 1500 * time.Millisecond,
	}
	// TODO: Pull showname from SceneCollection.ScName

	// TODO: The scenes and show object should be populated based on whatever
	//       the scene collection that is currently active is. but keep in
	//       mind the goal is to abstract awwy some of the less good design
	//       bits into a better logical construct

	toolkit.OBS.Show = &Show{
		OBS:    toolkit.OBS,
		Name:   "she hacked you",
		Scenes: Scenes{},
	}
	toolkit.OBS.Show.Cache()
	toolkit.X11.InitActiveWindow()
	return toolkit
}

func (t Toolkit) HandleWindowEvents() {
	// TODO:
	//   BUG:
	//     our short polling loop is repeating actions that should not be repeated
	//     based on checking if previously active window is the same as the current
	//     active window (NO CHANGE)
	//   BUG:
	//     secondary is not being detected, but primary and chromium is

	tick := time.Tick(t.Delay)
	for {
		select {
		case <-tick:
			//fmt.Printf("ticky tocky...\n")
			//activeWindowName := MarshalWindowName(name)

			//fmt.Printf("marshalled window name (%v)\n", activeWindowName)
			//fmt.Printf("marshalled window name as string (%v)\n", activeWindowName.String())

			// TODO: Completely disable undefined window recognization and so we
			// should never even cache an undefnied window type
			//   * It shouldn't ONLY check if the current window is different than
			//     than the active window-- but it should bypass undefined entirely
			//     AND it should be like if the active windiow is chromium-- it
			//     should only be checking against the other options, maybe even
			//     have like state-machine style pre-defined transitions

			//  * has it changed to a valid value? *not just has the value chagned*

			// TODO:
			// on item static avatar it has show as NIL

			if t.X11.HasActiveWindowChanged() {
				time.Sleep(4 * time.Second)

				currentScene := t.OBS.Show.Current

				primaryScene, _ := t.OBS.Show.Scene("content:primary")

				vimWindow, _ := primaryScene.Item("VIM")
				consoleWindow, _ := primaryScene.Item("CONSOLE")
				chromiumWindow, _ := primaryScene.Item("CHROMIUM")

				bumperScene, _ := t.OBS.Show.Scene("content:bumper")
				// TODO: We would put any significant items cached here
				//       then we can rebuild this switch and checks below
				//       to get shortest and fastest

				// Make naming of windows and focus is "activewindiow"

				// Filter:
				// 		* Only follow window IF its not currently in special
				//      HOLD CARD style bumper OR
				//    * Follow window:
				//					IF not HOLD card
				//					IF not END card
				//          IF not BUMPER card (?)
				//            OR
				//       Alternative strategy; Follow window:
				//          IF window CHROMIUM, PRIMARY (vim), and SECONDARY (Console)

				activeWindow := t.X11.ActiveWindow()
				switch activeWindow {
				case Primary, Secondary:
					t.X11.CacheActiveWindow()
					if currentScene.HasName("content:primary") {
						if !vimWindow.Visible {
							bumperScene.Transition()
							primaryScene.Transition(4 * time.Second)

							vimWindow.Unhide().Lock().Update()
							consoleWindow.Unhide().Lock().Update()
							chromiumWindow.Hide().Update()
						}
					}
				case Chromium:
					t.X11.CacheActiveWindow()
					if currentScene.HasName("content:primary") {
						if !chromiumWindow.Visible {
							bumperScene.Transition()
							primaryScene.Transition(4 * time.Second)

							vimWindow.Hide().Lock().Update()

							consoleWindow.Unhide().Lock().Update()
							chromiumWindow.Unhide().Update()
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
