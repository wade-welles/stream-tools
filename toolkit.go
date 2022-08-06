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
	toolkit.OBS.Show.CacheScenes()
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

	fmt.Printf("iterating over SCENES to later do things like transition...\n")
	//scene := map[string]*Scene{}
	//for index, cachedScene := range t.OBS.Show.Scenes {
	//	// TODO: We are getting multiple of each of the scenes, we need to do a
	//	// confirmation of existing scene so we dont have duplicates; can even hash
	//	// the name and maybe some other value to avoid string comparisons
	//	fmt.Printf("scene:\n")
	//	fmt.Printf("  len(scene.Items)=(%v)\n", len(cachedScene.Items))
	//	fmt.Printf("  scene_index: %v\n", index)
	//	fmt.Printf("  scene_name: %v\n", cachedScene.Name)
	//	fmt.Printf("  print neach scene item with the visibility and locked status\n")
	//	cachedSceneName := strings.Split(cachedScene.Name, "content:")
	//	fmt.Printf("  cached_scene_name: %v\n", cachedSceneName)
	//	fmt.Printf("  len(cachedSceneName)=%v\n", len(cachedSceneName))
	//	if 0 < len(cachedSceneName) {
	//		fmt.Printf("  cachedSceneName[1]: %v\n", cachedSceneName[1])

	//		//scene[cachedSceneName[1]] = cachedScene
	//	}
	//}

	//fmt.Printf("how many scenes were loaded into the 'scene' map object?(%v)\n", len(scene))

	//for sceneName, scene := range scene {
	//	fmt.Printf("scene:\n")
	//	fmt.Printf("  key:sceneName: %v\n", sceneName)
	//	fmt.Printf("  value:scene: %v\n", scene)
	//	fmt.Printf("  scene.name: %v\n", scene.Name)
	//	fmt.Printf("  len(scene.Items): %v\n", len(scene.Items))
	//	for index, item := range scene.Items {
	//		fmt.Printf("\n\nitem:\n")
	//		fmt.Printf("  index: %v\n", index)
	//		if item != nil {
	//			fmt.Printf("  name: %v\n", item.Name)
	//		} else {
	//			fmt.Printf(" item == nil :(")
	//		}
	//	}
	//}

	// NOTE: Incorrect way of doing this; but using it to create a functional demo
	// quickly. This short polls and we would obviously want to take advantage of
	// the built in ability that exists within the API to subscribe to events like
	// changes in the active window.
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

			fmt.Printf("current cached window: %v \n", t.X11.CurrentWindow)

			if t.X11.HasActiveWindowChanged() {
				fmt.Printf("The active window is: %v \n", t.X11.ActiveWindow())
				fmt.Printf("ACTIVE window CHANGED!\n")

				switch t.X11.ActiveWindow() {
				case Primary, Secondary:
					fmt.Printf("[primary] active window?(%v) \n", t.X11.ActiveWindow())
					t.X11.CacheActiveWindow()

					time.Sleep(4 * time.Second)
					if bumperScene, ok := t.OBS.Show.Scene("bumper"); ok {
						bumperScene.Transition()
					}

					if primaryScene, ok := t.OBS.Show.Scene("primary"); ok {
						primaryScene.Transition(4 * time.Second)
					}

				case Chromium:
					fmt.Println("[chromium] active window?(%v)", t.X11.ActiveWindow())
					t.X11.CacheActiveWindow()

					time.Sleep(4 * time.Second)
					if bumperScene, ok := t.OBS.Show.Scene("bumper"); ok {
						bumperScene.Transition()
					}

					if primaryScene, ok := t.OBS.Show.Scene("primary"); ok {
						primaryScene.Transition(4 * time.Second)
						// TODO: hide the vim terminal and unhide chrome
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

			t.AvatarToggle()
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
	if primaryScene, ok := t.OBS.Show.Scene("primary"); ok {

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
