package obstools

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	// OBS

	goobs "github.com/andreykaipov/goobs"
	events "github.com/andreykaipov/goobs/api/events"
	requests "github.com/andreykaipov/goobs/api/requests"
	scenes "github.com/andreykaipov/goobs/api/requests/scenes"
	sources "github.com/andreykaipov/goobs/api/requests/sources"
	studiomode "github.com/andreykaipov/goobs/api/requests/studio_mode"
	typedefs "github.com/andreykaipov/goobs/api/typedefs"
)

// TODO: STRONGLY consider the concept of just recursively nesting a generic
// layer-able type that we use to store both items, and scenes, and merge those
// two concepts basically via a common interrface

// TODO: BUT IF NOT; then we should RUN the other way, and create specialized
// audio level tracks stored in their own compartmentalized object to control
// mute/volume **and volume is very important; since we unmute bg music and mute
// our primary microphone that is a single channel spread across two.
type Item struct {
	Id   int
	Name string

	Layer

	Position *Position

	Scene *Scene
	Items *Items
	Show  *Show
}

type Position struct {
	X float64
	Y float64
}

func (it Item) HasName(name string) bool {
	return len(it.Name) == len(name) && it.Name == name
}

func ParseItem(item typedefs.SceneItem) *Item {
	// TODO: Not yet cahcing the scene pointers and show pointer (essentially no
	// relationships at all atm)
	return &Item{
		Id:   item.Id,
		Name: item.Name,
		Layer: Layer{
			Visible: item.Render,
			Locked:  item.Locked,
			Position: Position{
				X: item.X,
				Y: item.Y,
			},
		},
	}
}

type Items []*Item

func (its Items) Name(name string) *Item {
	for _, item := range its {
		if item.HasName(name) {
			return item
		}
	}
	return nil
}

// TODO: SHOULD this be linked list? I feel like it wouild be much nicer if it
// was! I dont want to reimplement a ton of linked list methods

// obs.Client.Scenes.First().ItemCollections.First().Items.Hidden() =>
// obs.Client.Scenes.First().Items.Hidden() => all hidden items in scene

type Layer struct {
	Order   int
	Visible bool
	Locked  bool
	Position
}

type Scene struct {
	Layer
	Show        *Show
	Name        string
	IsCurrent   bool
	IsPreviewed bool
	Items       Items
	OBSObject   typedefs.Scene
}

func (sc *Scene) Unlock() { sc.Locked = false }
func (sc *Scene) Lock()   { sc.Locked = true }

func (it *Item) Unlock() { it.Locked = false }
func (it *Item) Lock()   { it.Locked = true }
func (it *Item) Unhide() { it.Visible = false }
func (it *Item) Hide()   { it.Visible = true }

func (its Items) Unlocked() (unlockedItems Items) {
	for _, item := range its {
		if !item.Locked {
			unlockedItems = append(unlockedItems, item)
		}
	}
	return unlockedItems
}

// items.LockedItems()
// scenes.LockedScenes()

// scenes.Locked()
// items.Locked()
func (its Items) Locked() (lockedItems Items) {
	for _, item := range its {
		if item.Locked {
			lockedItems = append(lockedItems, item)
		}
	}
	return lockedItems
}

func (its Items) Visible() (visible Items) {
	for _, item := range its {
		if item.Visible {
			visible = append(visible, item)
		}
	}
	return visible
}

func (its Items) Hidden() (hidden Items) {
	for _, item := range its {
		if !item.Visible {
			hidden = append(hidden, item)
		}
	}
	return hidden
}

func (its Items) Size() int    { return len(its) }
func (its Items) First() *Item { return its[0] }
func (its Items) Last() *Item  { return its[its.Size()-1] }

// scene.Item("itemName").Hide()
// scene.Items.Name("itemName").Hide()
func (sc Scene) Item(name string) *Item {
	return sc.Items.Name(name)
}

func (sc Scene) HasName(name string) bool {
	return len(sc.Name) == len(name) && sc.Name == name
}

func (scs Scenes) Locked() (locked Scenes) {
	for _, scene := range scs {
		if scene.Locked {
			locked = append(locked, scene)
		}
	}
	return locked
}

// scenes.Unlocked().First()
func (scs Scenes) Unlocked() (unlocked Scenes) {
	for _, scene := range scs {
		if !scene.Locked {
			unlocked = append(unlocked, scene)
		}
	}
	return unlocked
}

// scenes.Visible().Size()
func (scs Scenes) Visible() (visible Scenes) {
	for _, scene := range scs {
		if scene.Visible {
			visible = append(visible, scene)
		}
	}
	return visible
}

func (scs Scenes) Hidden() (hidden Scenes) {
	for _, scene := range scs {
		if !scene.Visible {
			hidden = append(hidden, scene)
		}
	}
	return hidden
}

type Scenes []*Scene

func (scs Scenes) Size() int     { return len(scs) }
func (scs Scenes) First() *Scene { return scs[0] }
func (scs Scenes) Last() *Scene  { return scs[scs.Size()-1] }

// TODO: I thought Source == SceneItem

// NOTE: Scene Collection == All (Show|Channel) Scenes
type Show struct {
	OBS    *OBS
	Name   string
	Scenes Scenes

	// NOTE: Ideally we keep the typedefs.Scene in a linked list for better
	// interactions
	//Scenes       *list.List
}

//type Source struct {
//
//}

type Usage struct {
	CPU    int
	Disk   int
	Memory int
}

type Stats struct {
	Usage     Usage
	Streaming bool
	Recording bool

	FramesPerSecond uint8

	Time               uint
	AvgerageRenderTime uint
	FramesLost         uint
	FramesSkipped      uint
	FramesDropped      uint
	DataOutput         uint
	Bitrate            uint
}

// NOTE: The API we are creating looks something like:
//          obs.Show.Scenes.First()
//          obs.Show.Scenes.First().Items.First().(Show()|.Hide())
type OBS struct {
	*goobs.Client

	Stats      *Stats
	AudioMixer *AudioMixer
	Show       *Show

	//Sources []*goobs.Source
}

func (obs OBS) StudioMode() bool {
	response, err := obs.Client.StudioMode.GetStudioModeStatus()
	if err != nil {
		return false
	} else {
		return response.StudioMode
	}
}

//  TODO: This would be better as obs.Shows.Active(); but that means creating a
//  Shows type (POTENTIALLY, but ideally no)

func (obs OBS) ActiveShow() *Show {
	return nil
}

// TODO: This should either give the full OBS object returned (like an init
// function) or this should be a method on OBS after it is created you connect.
func ConnectToOBS() *goobs.Client {
	client, err := goobs.New("127.0.0.1:4444", goobs.WithDebug(true))
	if err != nil {
		panic(err)
	}
	return client
}

// go Events()
func (obs OBS) Events() {
	for event := range obs.Client.IncomingEvents {
		switch e := event.(type) {
		case *events.SourceVolumeChanged:
			fmt.Printf("Volume changed for %-25q: %f\n", e.SourceName, e.Volume)
		default:
			log.Printf("Unhandled event: %#v", e.GetUpdateType())
		}
	}
}

// TODO: Switch to Scene

// TODO: Toggle Visibility on element(source) in scene

// TODO: List Scenes

/////////////////////////////////////////////////////////////
//func (obs OBS) Scenes() ([]*typedefs.Scene, error) {
//	apiResponse, err := obs.Client.Scenes.GetSceneList()
//	return apiResponse.Scenes, err
//}

// TODO: Not a fan of this name, the general functionality; at least all obs
// caching should be if separated also combined in a single public method
func (sh Show) CacheScenes() (bool, error) {
	apiResponse, err := sh.OBS.Client.Scenes.GetSceneList()
	if err != nil {
		return false, err
	}

	scenes := apiResponse.Scenes

	// TODO: Iterate through the scene.Sources ([]SceneItem) and add them by
	// creating our item objects and adding them to the new scene.
	// TODO: Look at the apiresponse object returned from GetSceneList to
	// determine if only current and a list is given or if preview is also
	// provided or possibly other variables
	items := Items{}
	cachedScenes := Scenes{}
	for _, scene := range scenes {
		// TODO: Initialize the scene object, add it to a slice of scenes and make
		// sure that the scene and items are tied together so one can easily move
		// between them
		// Items: scene.Sources

		for _, item := range scene.Sources {
			fmt.Printf("item: ", item)

			childItems := Items{}
			if 0 < len(item.GroupChildren) {
				// TODO: This item has children and this will need to be checked
				// recursively perhaps to guarantee we dont miss any
				for _, childItem := range item.GroupChildren {
					if 0 < len(childItem.GroupChildren) {
						fmt.Printf("is this even possible? bc if so, ...\n")
					}
					// TODO: This item has children and this will need to be checked
					childItems = append(childItems, ParseItem(childItem))
				}
			}

			// TODO: Use ParentGroupName to ensure we carry over any existing
			// relationships

			//       Has sub-items
			//       Has scene info
			items = append(items, ParseItem(item))
			// TODO: Not yet caching the Show object or scene objects
			// TODO: Assign X/Y position value on item

			cachedScenes = append(cachedScenes, &Scene{
				Name:  scene.Name,
				Items: items,
			})

			fmt.Println("scene:", scene.Name)
			sh.Scenes = append(sh.Scenes, &Scene{
				Name:  scene.Name,
				Items: items,
				IsCurrent: len(scene.Name) == len(apiResponse.CurrentScene) &&
					scene.Name == apiResponse.CurrentScene,
			})

		}
	}

	return 0 < sh.Scenes.Size(), err
}

//func (obs *OBS) CurrentScene() (string, error) {
//	apiResponse, err := obs.Client.Scenes.GetCurrentScene()
//	return apiResponse.CurrentScene, err
//}

//func (obs *OBS) CacheScene() (string, error) {
//	apiResponse, err := obs.Client.Scenes.GetCurrentScene()
//	if err == nil {
//		obs.Show.CurrentScene = apiResponse.CurrentScene
//	}
//	//return obs.Show.CurrentScene
//	return apiResponse.CurrentScene, err
//
//}

// NOTE API we should aim for should use linked lists and then we can make our
// own scene type that ahs methods like transition.
// obs.Scenes.First().Transition()

// obs.Scenes.Last().Items.First().Hide()

// NOTE: Because our scene object would have a current field, then we can pull
// it out instead of saving the name of it. There would be a scenes object like
//    type Scenes *list.List
// and with that we could add methods like Current() and the base methods that
// would be included would be container/list or, in other words lniked list
//
// obs.Scenes.Current().Name

// obs.Scene("name").Transition()

func (scs Scenes) Current() *Scene {
	for _, scene := range scs {
		if scene.IsCurrent {
			return scene
		}
	}
	return nil
}

//func (scs Scenes) Random() *Scene {
//	return nil
//}

func (scs Scenes) Preview() *Scene {
	for _, scene := range scs {
		if scene.IsPreviewed {
			return scene
		}
	}
	return nil
}

func (scs Scenes) Name(name string) *Scene {
	for _, scene := range scs {
		if scene.HasName(name) {
			return scene
		}
	}
	return nil
}

// show.Scene("bumper").Transition()
func (sc Scene) Transition() error {
	_, err := sc.Show.OBS.Client.Scenes.SetCurrentScene(
		&scenes.SetCurrentSceneParams{
			SceneName: sc.Name,
		},
	)
	return err
}

// show.Transition("bumper")
func (sh Show) Transition(sceneName string) error {
	_, err := sh.OBS.Client.Scenes.SetCurrentScene(
		&scenes.SetCurrentSceneParams{
			SceneName: sceneName,
		},
	)
	return err
}

func (sh Show) Scene(name string) *Scene {
	return sh.Scenes.Name(name)
}

// obs.PreviewScene("name")

// obs.Scenes.First().Preview()
// obs.Scene("name").Transition()

// obs.Scenes.First().Items.First().Hide()

// obs.Scenes.Preview() => Scene that is in the preview of studiomode

//type Params requests.ParamsBasic
//type Response requests.ResponseBasic

type Params struct {
	*requests.ParamsBasic
	Name  string
	Value string
}

// SceneName string `json:"scene-name,omitempty"`

// obs.Scenes.First().Preview() => sets the preview in studiomode
//func (sc Scene) Preview() error {
//	_, err := sc.Show.OBS.Client.StudioMode.SetPreviewScene(
//		studiomode.SetPreviewSceneParams{
//			SceneName: sc.Name,
//		},
//	)
//	return err
//}

func (sc Scene) Preview() error {
	for index, scene := range sc.Show.Scenes {
		sc.Show.Scenes[index].IsPreviewed = false
		if scene.HasName(sc.Name) {
			err := sc.Show.Preview(sc.Name)
			sc.Show.Scenes[index].IsPreviewed = (err == nil)
		}
	}
	return nil
}

func (sh Show) Preview(sceneName string) error {
	apiRequest := studiomode.SetPreviewSceneParams{
		SceneName: sceneName,
	}
	_, err := sh.OBS.Client.StudioMode.SetPreviewScene(&apiRequest)
	return err
}

func (scs Scenes) Previewed() *Scene {
	for _, scene := range scs {
		if scene.IsPreviewed {
			return scene
		}
	}
	return nil
}

type AudioMixer []*Audio

type Audio struct {
	Name   string
	Volume int
	Muted  bool
}

//////////////////////////////////

func (obs OBS) Volume() {
	//// TODO: List must be obtained from OBS object, and preferably it should be
	////       cached in OBS object.
	//name := list.Sources[0].Name

	//apiResponse, _ := obs.Client.Sources.GetVolume(
	//	&sources.GetVolumeParams{
	//		Source: name,
	//	},
	//)
	//fmt.Printf("%s is muted? %t\n", name, apiResponse.Muted)
}

// TODO: Toggle Mute On Source
//        * Mute + reduce volume to 0 [hide completely]
//        * Mute

func (obs OBS) Sources() {
	list, _ := obs.Client.Sources.GetSourcesList()

	for _, v := range list.Sources {
		if _, err := obs.Client.Sources.SetVolume(&sources.SetVolumeParams{
			Source: v.Name,
			Volume: rand.Float64(),
		}); err != nil {
			panic(err)
		}
	}

	if len(list.Sources) == 0 {
		fmt.Println("No sources!")
		os.Exit(0)
	}
}
