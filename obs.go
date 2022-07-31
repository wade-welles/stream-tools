package obstools

import (
	"fmt"
	"log"
	"math/rand"
	"time"

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

// NOTE: The point on the source that the item is manipulated from.
//       The sum of 1=Left or 2=Right, and 4=Top or 8=Bottom,
//        or omit to center on that axis.
//  TODO: Probably just not just ignore but totally subtract
type Alignment uint8

const (
	CenterAlign      Alignment = 0
	LeftAlign        Alignment = 1
	RightAlign       Alignment = 2
	TopAlign         Alignment = 4
	TopLeftAlign     Alignment = 5
	TopRightAlign    Alignment = 6
	BottomAlign      Alignment = 8
	BottomLeftAlign  Alignment = 9
	BottomRightAlign Alignment = 10
)

func (a Alignment) String() string {
	switch a {
	case CenterAlign:
		return "center"
	case LeftAlign:
		return "left"
	case RightAlign:
		return "right"
	case TopAlign:
		return "top"
	case TopRightAlign:
		return "top-right"
	case BottomAlign:
		return "bottom"
	case BottomLeftAlign:
		return "bottom-left"
	case BottomRightAlign:
		return "bottom-right"
	default:
		return "undefined"
	}
}

func MarshalAlignment(alignment int) Alignment { return Alignment(alignment) }

type ItemType uint8

// NOTE: From typedefs/scene_item.go from goobs
// Source type. Value is one of the following: "input", "filter", "transition", "scene" or "unknown"
const (
	UnknownType ItemType = iota
	InputType
	FilterType
	TransitionType
	SceneType
)

func (itt ItemType) String() string {
	switch itt {
	case InputType:
		return "input"
	case FilterType:
		return "filter"
	case TransitionType:
		return "transition"
	case SceneType:
		return "scene"
	default: // UnknownType
		return "unknown"
	}
}

func MarshalItemType(itemType string) ItemType {
	switch itemType {
	case InputType.String():
		return InputType
	case FilterType.String():
		return FilterType
	case TransitionType.String():
		return TransitionType
	case SceneType.String():
		return SceneType
	default:
		return UnknownType
	}
}

type Dimensions struct {
	Height float64
	Width  float64
}

type Item struct {
	Id   int
	Name string
	Type ItemType

	// TODO: Add
	//        type
	//        folder/relation
	//

	Layer

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

func PrintItem(item typedefs.SceneItem) {
	fmt.Printf("__item__\n")
	fmt.Printf("  id: %v \n", item.Id)
	fmt.Printf("  type: %v \n", item.Type)
	fmt.Printf("  name: %v \n", item.Name)
	fmt.Printf("  hidden: %v \n", !item.Render)
	fmt.Printf("  locked: %v \n", item.Locked)
	// Group
	fmt.Printf("  _group_\n")
	fmt.Printf("    is_group: %v \n", len(item.GroupChildren) != 0)
	fmt.Printf("    parent_group_name: %v \n", item.ParentGroupName)
	fmt.Printf("    group_children: \n")
	fmt.Printf("      count: %v \n", len(item.GroupChildren))
	fmt.Printf("      children: %v \n", item.GroupChildren)
	// Audio
	fmt.Printf("  _audio_\n")
	fmt.Printf("    muted: %v \n", item.Muted)
	fmt.Printf("    volume: %v \n", item.Volume)
	// TODO: 3 difference coordinates? seems silly
	fmt.Printf("  _position_\n")
	fmt.Printf("    alignment: %v \n", item.Alignment)
	fmt.Printf("    c_x: %v, c_y: %v \n", item.Cx, item.Cy)
	fmt.Printf("    source_c_x: %v, source_c_y: %v \n", item.SourceCx, item.SourceCy)
	fmt.Printf("    X: %v, Y: %v \n", item.X, item.Y)
}

// TODO: Look at OUR Item type and redo the print function using our
//       simplified type
//func (it Item) Print() {
//	fmt.Printf("__item__\n")
//	fmt.Printf("  id: %v \n", it.Id)
//	fmt.Printf("  type: %v \n", it.Type)
//	fmt.Printf("  name: %v \n", it.Name)
//	fmt.Printf("  hidden: %v \n", !it.Render)
//	fmt.Printf("  locked: %v \n", it.Locked)
//	// Group
//	fmt.Printf("  _group_\n")
//	fmt.Printf("    is_group: %v \n", len(it.GroupChildren) != 0)
//	fmt.Printf("    parent_group_name: %v \n", it.ParentGroupName)
//	fmt.Printf("    group_children: \n")
//	fmt.Printf("      count: %v \n", len(it.GroupChildren))
//	fmt.Printf("      children: %v \n", it.GroupChildren)
//	// Audio
//	fmt.Printf("  _audio_\n")
//	fmt.Printf("    muted: %v \n", it.Muted)
//	fmt.Printf("    volume: %v \n", it.Volume)
//	// TODO: 3 difference coordinates? seems silly
//	fmt.Printf("  _position_\n")
//	fmt.Printf("    alignment: %v \n", it.Alignment)
//	fmt.Printf("    c_x: %v, c_y: %v \n", it.Cx, it.Cy)
//	fmt.Printf("    source_c_x: %v, source_c_y: %v \n", it.SourceCx, it.SourceCy)
//	fmt.Printf("    X: %v, Y: %v \n", it.X, it.Y)
//}

func ParseItem(item typedefs.SceneItem) *Item {
	// TODO: Not yet cahcing the scene pointers and show pointer (essentially no
	// relationships at all atm)

	// cX is width of sprite, cY is height

	// x, y is position of the sprite

	return &Item{
		Id:   item.Id,
		Name: item.Name,
		Type: MarshalItemType(item.Type),
		Layer: Layer{
			Visible:   item.Render,
			Locked:    item.Locked,
			Alignment: item.Alignment,
			Position: Position{
				X: item.X,
				Y: item.Y,
			},
			Dimensions: Dimensions{
				Width:  item.Cx,
				Height: item.Cy,
			},
			Source: Dimensions{
				Width:  item.SourceCx,
				Height: item.SourceCy,
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
	Order     int
	Visible   bool
	Locked    bool
	Alignment int
	Position
	Dimensions
	Source Dimensions
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

// TODO: We need a function for ADDING a scene, and this can then also have
//       the validation necessary to prevent duplicates
func (sh *Show) CacheScenes() (bool, error) {
	apiResponse, err := sh.OBS.Client.Scenes.GetSceneList()
	if err != nil {
		return false, err
	}

	scenes := apiResponse.Scenes

	items := Items{}
	for _, scene := range scenes {
		fmt.Printf("___________________________________________\n")
		fmt.Printf("__scene__\n")
		fmt.Printf("  name: %v \n", scene.Name)
		for _, item := range scene.Sources {
			// NOTE: Cover the case were we do need the child items as their own
			//       item(s) with a zed.
			//       We want the ability to store just the child items into a
			//       a group item's Items (sub-items essentially) field.
			//       We will eventually support nested recurisve single item
			//       type representing all items and all folder types and then
			//       we use that abstraction to interact with OBS so we are not
			//       confined to OBS not so good design patterns and data modeling
			//       architecture? .. :(
			//childItems := Items{}
			if 0 < len(item.GroupChildren) {
				for _, childItem := range item.GroupChildren {
					PrintItem(childItem)
					items = append(items, ParseItem(childItem))
				}
			}

			// TODO: Use ParentGroupName to ensure we carry over any existing
			// relationships

			//       Has sub-items
			//       Has scene info
			PrintItem(item)
			items = append(items, ParseItem(item))
			// TODO: Not yet caching the Show object or scene objects
			// TODO: Assign X/Y position value on item
		}

		sh.AddScene(
			scene.Name,
			items,
			len(scene.Name) == len(apiResponse.CurrentScene) &&
				scene.Name == apiResponse.CurrentScene,
		)
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

func (scs Scenes) Index(sceneIndex int) *Scene {
	for index, scene := range scs {
		if sceneIndex == index {
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

// TODO: THis is nice bc we get chaining additions scs.Add(sc1).Add(sc2)
//   okay
//     so
///     maybe like name and items?

func (scs Scenes) Exists(sceneName string) bool {
	return scs.Name(sceneName) != nil
}

//  Add can not be done like this bc of the data type pulled fro0m the API
func (sh *Show) AddScene(sceneName string, sceneItems Items, sceneIsCurrent bool) *Show {
	if sh.Scenes.Exists(sceneName) {
		fmt.Printf("Scene already exists, skipping (should update rly)\n")
		return sh
	}

	fmt.Printf("scene: %v \n", sceneName)
	sh.Scenes = append(sh.Scenes, &Scene{
		Show:      sh,
		Name:      sceneName,
		Items:     sceneItems,
		IsCurrent: sceneIsCurrent,
	})
	return sh
}

// show.Scene("bumper").Transition()
// TODO: Pass in variadic time.Duration, and this can be the time to sleep
//       then preform the transition.
//
//      scene.Transition()
//
func (sc Scene) Transition(sleepDuration ...time.Duration) error {
	if 0 < len(sleepDuration) {
		fmt.Printf("sleeping \n")
		time.Sleep(sleepDuration[0])
	}
	// TODO: What if sc is NIL??????? well it cant be that but what about empty?
	//       lets validate!
	//
	//         (and eventually put validation in standardized methods)
	//
	fmt.Printf("scene\n")
	fmt.Printf("  name: %v\n", sc.Name)

	//if len(sc.Name) == 0 {
	//	return errors.New("scene is undefined")
	//}

	return sc.Show.SetCurrentScene(sc.Name)
}

func (sh Show) SetCurrentScene(sceneName string) error {
	sceneRequest := scenes.SetCurrentSceneParams{
		SceneName: sceneName,
	}
	_, err := sh.OBS.Scenes.SetCurrentScene(&sceneRequest)
	return err
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
		fmt.Printf("No sources!\n")
		// TODO: Why would we exit? And exit 0??? Thats no error
		//os.Exit(0)
	}
}
