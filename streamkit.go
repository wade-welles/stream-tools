package streamkit

import (
	"time"

	broadcast "github.com/wade-welles/streamkit/broadcast"
	obs "github.com/wade-welles/streamkit/broadcast/obs"
	show "github.com/wade-welles/streamkit/broadcast/show"

	x11 "github.com/wade-welles/streamkit/x11"
)

type Toolkit struct {
	// NOTE: Short-poll rate [we will rewrite without short polling after]
	Delay time.Duration
	OBS   *obs.Client
	X11   *x11.X11
	// TODO: Our local copy of the show is entirely separate from obs.Client so we
	// can change that out while maintaining logic and a data object
	Show   *broadcast.Show
	Config map[string]string
}

// TODO: Could pass the host for OBS and the host for X11 as variadic strings so
// it can be empty, or provide position 1 for obs position 2 for x11 (though x11
// should assumingly always be 127.0.0.1 whereas obs reasonably could be
// different
func New() (toolkit *Toolkit) {
	// TODO: Show should be from config, and obs and x11 information. Logically
	// stored in ~/.config/$APP_NAME and the local data should be
	// ~/.local/share/$APP_NAME

	showConfig := map[string]string{
		"name":     "she hacked you",
		"obs_host": "192.168.1.1:4444",
	}

	wsAPI := obs.ConnectToOBS(showConfig["obs_host"])

	toolkit = &Toolkit{
		Config: showConfig,
		Show: &broadcast.Show{
			Scenes: make([]*show.Scene, 0),
		},
		OBS: &obs.Client{
			WS: wsAPI,
			//Mode: this is studio vs direct stream which is USELESS
			// ui concept only really
		},
		X11: &x11.X11{
			Client: x11.ConnectToX11(),
		},
		Delay: 1500 * time.Millisecond,
	}

	// TODO: We have to pull the scenes and cache them now

	// TODO: Pull showname from SceneCollection.ScName

	// TODO: The wscenes and show object should be populated based on whatever
	//       the scene collection that is currently active is. but keep in
	//       mind the goal is to abstract awwy some of the less good design
	//       bits into a better logical construct

	//toolkit.OBS.Show.OBS.Scenes = &scenes.Client{Client: t.OBS.WS}
	//toolkit.OBS.Show.OBS.Items =
	//toolkit.Show.Client = toolkit.OBS.WS
	//toolkit.Show.Cache()

	toolkit.X11.InitActiveWindow()
	return toolkit
}

//func (t Toolkit) HandleWindowEvents() {
//	parsedShow := t.OBS.Broadcast
//
//	fmt.Printf(
//		"number of scenes parsed from parsedShow object: %v\n",
//		parsedShow.Scenes,
//	)
//
//	scene, _ := parsedShow.ParseScene("Primary")
//	fmt.Printf("scene(%v)\n", scene)
//
//	parsedShow.Cache()
//
//	parsedShow.PrintDebug()
//
//	// this is returning
//	// sh.Scenes.Name(sceneName): and sceneName is Primary, so we get 1 out of 4
//	// of the scene names loaded and no scenes?
//	fmt.Printf(".SceneNames(): %s\n", strings.Join(t.OBS.Show.SceneNames(), ", "))
//	fmt.Printf("len(.Scenenames()) %v\n", len(t.OBS.Show.SceneNames()))
//
//	// TODO Primary was based off a hardcoded window type
//
//	// TODO: THIS =========== LINE is where we are failing, it doesnt look up
//	// Primary so it panics
//
//	primaryScene, ok := parsedShow.Scene("Primary")
//	if ok {
//		// TODO: We need to cache or initialize the items in a given scene
//	} else {
//		panic(fmt.Errorf("failed to locate primary scene"))
//	}
//	bumperScene, ok := parsedShow.Scene("Bumper")
//	if !ok {
//		panic(fmt.Errorf("failed to locate bumper scene"))
//	}
//
//	// TODO: This lookup is not connecting with the parsed items when we are
//	// running .Cache() on each scene as its parsed in Show.Cache()
//	fmt.Printf("# of primaryScene items: %v\n", len(primaryScene.Items))
//	fmt.Printf("# of bumperScene items: %v\n", len(bumperScene.Items))
//
//	//vimWindowName := "Primary Terminal (VIM Window)"
//	//vimWindow, ok := primaryScene.Item(vimWindowName)
//	//if ok {
//	//	fmt.Printf("found vimWindow item: %v\n", vimWindow)
//	//} else {
//	//	panic(
//	//		fmt.Errorf(
//	//			"failed to find item '" + vimWindowName + "' within the primary scene",
//	//		),
//	//	)
//	//}
//
//	//consoleWindow, _ := primaryScene.Item("CONSOLE")
//	//chromiumWindow, _ := primaryScene.Item("CHROMIUM")
//
//	tick := time.Tick(t.Delay)
//	for {
//		select {
//		case <-tick:
//			if t.X11.HasActiveWindowChanged() {
//				// TODO:
//				time.Sleep(4 * time.Second)
//
//				currentScene := t.OBS.Show.Current
//				activeWindow := t.X11.ActiveWindow()
//
//				if currentScene.HasName("Primary") {
//					switch activeWindow.Type {
//					case x11.Terminal:
//						t.X11.CacheActiveWindow()
//						//if !vimWindow.Visible { // is it a terminal, and termainl ID or hash of some combo of things
//						//	bumperScene.Transition()
//						//	primaryScene.Transition(4 * time.Second)
//
//						//	chromiumWindow.Hide().Lock().Update()
//						//	vimWindow.Unhide().Lock().Update()
//						//	consoleWindow.Unhide().Lock().Update()
//						//}
//					case x11.Browser: // TODO: Change to is it a browser
//						t.X11.CacheActiveWindow()
//						//if !chromiumWindow.Visible {
//						//	bumperScene.Transition()
//						//	primaryScene.Transition(4 * time.Second)
//
//						//	vimWindow.Hide().Lock().Update()
//						//	consoleWindow.Unhide().Lock().Update()
//						//	chromiumWindow.Unhide().Lock().Update()
//						//}
//					default: // UndefinedName
//						// TODO: This will error out if for example you select a window from
//						// a parent VM in Multiverse so kinda annoying
//						// TODO: We should never allow this condition to ever occur, and by
//						// doing that we optimize it further bc we are not checking conditions
//						// that we dont want to begin with
//						fmt.Println("[undefined] active window?(%v)", t.X11.ActiveWindow())
//					}
//				} else {
//					fmt.Printf("current scene is not set to primary")
//
//				}
//			}
//			// TODO: Check what the active window currently is; then use obs-ws to
//			// change the scenes with the bumper in between
//		}
//	}
//
//	//time.Sleep(2 * time.Second)
//	//t.AvatarToggle()
//}

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
