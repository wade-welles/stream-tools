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
	toolkit.OBS.Show = &Show{
		OBS:    toolkit.OBS,
		Name:   "she hacked you",
		Scenes: Scenes{},
	}
	toolkit.OBS.Show.CacheScenes()
	toolkit.X11.CacheActiveWindow()
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

	// TODO: defer a close of the close
	//defer t.Connection.Close()
	//if show := t.OBS.Show; show != nil {
	//	fmt.Printf("Show found!(%v)\n", show)
	//	t.OBS.Show.CacheScenes()
	//} else {
	//	fmt.Printf("No Show found!\n")
	//}
	//	panic(err)
	//}

	// NOTE: Incorrect way of doing this; but using it to create a functional demo
	// quickly. This short polls and we would obviously want to take advantage of
	// the built in ability that exists within the API to subscribe to events like
	// changes in the active window.
	tick := time.Tick(t.Delay)
	for {
		select {
		case <-tick:
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

			//     has it changed to a valid value? *not just has the value chagned*
			if t.X11.HasActiveWindowChanged() {
				switch t.X11.ActiveWindow() {
				case Primary:
					fmt.Println("[primary] active window?(%v)", t.X11.ActiveWindow())
					// TODO:
					// 	* Show [A RANDOM BUMPER] for X(5?) seconds

					//	bumpers := []*Scene{}
					// TODO: This was never actually checked if it was working or
					//       if a real attempt on this functionality
					//activeScene := t.OBS.Show.Scene("bumper+holdcard")

					//activeScene.Transition()
					//activeScene.Item("avatar").Hide()

					//  * Switch to the primary content
					//  * Mute any background music, and unmuting the mic

					//	_, _ = client.Sources.ToggleMute(&sources.ToggleMuteParams{Source: name})

					t.X11.CacheActiveWindow()
				case Secondary:
					fmt.Println("[secondary] active window?(%v)", t.X11.ActiveWindow())
					// TODO:
					// 	* Show [A RANDOM BUMPER] for X(5?) seconds
					//  * Switch to the primary content
					//  * Mute any background music, and unmuting the mic

					t.X11.CacheActiveWindow()
				case Chromium:
					fmt.Println("[chromium] active window?(%v)", t.X11.ActiveWindow())
					// TODO:
					// 	* Show [A RANDOM BUMPER] for X(5?) seconds
					//  * Switch to the chrome window
					//  * Mute any background music, and unmuting the mic
					t.X11.CacheActiveWindow()
				default: // UndefinedName
					fmt.Println("[undefined] active window?(%v)", t.X11.ActiveWindow())
				}
				// TODO: Check what the active widnow currently is; then use obs-ws to
				// change the scenes with the bumper in between
			}
		}
	}
}
