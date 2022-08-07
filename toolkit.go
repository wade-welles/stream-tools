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

			//     has it changed to a valid value? *not just has the value chagned*

			// TODO:
			// on item static avatar it has show as NIL

			if t.X11.HasActiveWindowChanged() {
				switch t.X11.ActiveWindow() {
				case Primary, Secondary:
					fmt.Printf("[primary+secondary] active window?(%v)\n", t.X11.ActiveWindow())
					t.X11.CacheActiveWindow()
					time.Sleep(4 * time.Second)

					fmt.Printf("t: %v\n", t)
					fmt.Printf("t.OBS: %v\n", t.OBS)
					fmt.Printf("t.OBS.Show: %v\n", t.OBS.Show)
					fmt.Printf("t.OBS.Show.Scenes: %v\n", t.OBS.Show.Scenes)
					fmt.Printf("len(t.OBS.Show.Scenes): %v\n", len(t.OBS.Show.Scenes))

					if bumperScene, ok := t.OBS.Show.Scene("content:bumper"); ok {
						fmt.Printf("bumperScene.Name: %v\n", bumperScene.Name)
						bumperScene.Transition()
					} else {
						fmt.Printf("failed to transition to bumper\n")
					}

					if primaryScene, ok := t.OBS.Show.Scene("content:primary"); ok {
						primaryScene.Transition(4 * time.Second)
						fmt.Printf("attempting to transition to primary\n")
					} else {
						fmt.Printf("failed to transition to primary\n")
					}
				case Chromium:
					fmt.Printf("[chromium] active window?(%v)\n", t.X11.ActiveWindow())
					t.X11.CacheActiveWindow()

					time.Sleep(4 * time.Second)
					if bumperScene, ok := t.OBS.Show.Scene("bumper:content"); ok {
						fmt.Printf("attempting to transition to bumper\n")
						bumperScene.Transition()
					} else {
						fmt.Printf("failed to transition to bumper\n")
					}

					if primaryScene, ok := t.OBS.Show.Scene("content:primary"); ok {
						fmt.Printf("attempting to transition to primary\n")
						primaryScene.Transition(4 * time.Second)
						// TODO: hide the vim terminal and unhide chrome
					} else {
						fmt.Printf("failed to transition to primary\n")
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

//func (sh Show) ToggleItemVisibility(item *Item) (error, bool) {
//	// TODO: Pretty sure we actually dont need to do this since
//	//       we are interacting with
//	cachedItem := sh.Scene(item.Scene.Name).Item(item.Name)
//
//	if cachedItem.Visible {
//		return cachedItem.Hide()
//	} else {
//		return cachedItem.Unhide()
//	}
//}

//func (sh Show) ToggleItemLock(item *Item) (error, bool) {
//	cachedItem := sh.Scene(item.Scene.Name).Item(item.Name)
//
//	if cachedItem.Locked {
//		return cachedItem.Unlock()
//	} else {
//		return cachedItem.Lock()
//	}
//}

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
			staticAvatar.Unhide()
			staticAvatar.Print()
			fmt.Printf("---\n")

			dynamicAvatar.Print()
			dynamicAvatar.Hide()
			dynamicAvatar.Print()
		} else {
			dynamicAvatar.Print()
			dynamicAvatar.Unhide()
			dynamicAvatar.Print()

			fmt.Printf("---\n")

			staticAvatar.Print()
			staticAvatar.Hide()
			staticAvatar.Print()
		}
	}
}
