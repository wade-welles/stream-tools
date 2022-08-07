package obstools

import (
	"fmt"
	"log"
	"time"

	// OBS
	goobs "github.com/andreykaipov/goobs"
	events "github.com/andreykaipov/goobs/api/events"
	requests "github.com/andreykaipov/goobs/api/requests"
	sceneitems "github.com/andreykaipov/goobs/api/requests/scene_items"
	scenes "github.com/andreykaipov/goobs/api/requests/scenes"
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
	FolderType
	TransitionType
	SceneType
)

func (itt ItemType) String() string {
	switch itt {
	case InputType:
		return "input"
	case FolderType:
		return "folder"
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

// TODO: Item type ideally should have group or folder in here but these are
//       the ones from OBS; so itd probably need to be separate from our
//       custom ones and then added together with a further type
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

// TODO: We could use len of Items being greater than one to determine
//       if a given item is a folder (or group using OBS naming) but
//       an item can be a folder (or group) without any items and is
//       still a groupppp
//
//       but keep in mind we didnt want to have to implement a special
//       enumerator or worse have a bool thats like IsGroup or IsFolder
//
//       ideally we keep all that logix with just using our nesting of
//       items recursively
//

type Item struct {
	Id   int
	Name string
	Type ItemType

	// TODO: Add
	//        type
	//        folder/relation
	//
	Layer

	Parent *Item
	Items  Items

	Scene *Scene
}

type Position struct {
	X float64
	Y float64
}

func (it Item) HasName(name string) bool {
	return len(it.Name) == len(name) && it.Name == name
}

// TODO: One issue with this is that something can be a folder without
//       having any items within it; and while that distinction may be
//       useless to us, it may cause issues in the future since there
//       is no simple way to switch between types beyond deleting and
//       recreating

// NOTE: Reconcile the fact that this exists and the type exists
//       and they wont be the same

func (it Item) IsFolder() bool { return (0 < len(it.Items)) }

// Bounds *typedefs.Bounds `json:"bounds,omitempty"`
/// The crop specification for the object (source, scene item, etc).
// Crop *typedefs.Crop `json:"crop,omitempty"`
/// The item specification for this object.
// Item *typedefs.Item `json:"item,omitempty"`
/// The new locked status of the source. 'true' keeps it in its current position, 'false' allows movement.
// Locked *bool `json:"locked,omitempty"`
/// The position of the object (source, scene item, etc).
// Position *typedefs.Position `json:"position,omitempty"`
/// The new clockwise rotation of the item in degrees.
// Rotation float64 `json:"rotation,omitempty"`
/// The scaling specification for the object (source, scene item, etc).
// Scale *typedefs.Scale `json:"scale,omitempty"`
/// Name of the scene the source item belongs to. Defaults to the current scene.
// SceneName string `json:"scene-name,omitempty"`
// // The new visibility of the source. 'true' shows source, 'false' hides source.
// Visible *bool `json:"visible,omitempty"`

func (it Item) Update() (Item, bool) {
	itemParams := sceneitems.SetSceneItemPropertiesParams{
		SceneName: it.Scene.Name,
		Item:      &typedefs.Item{Name: it.Name},
		Locked:    &it.Locked,
		Visible:   &it.Visible,
		// TODO: Eventually Update these values too
		//Bounds:    resp.Bounds,
		//Crop:      resp.Crop,
		//Position:  resp.Position,
		//Rotation:  resp.Rotation,
		//Scale:     resp.Scale,
	}

	_, err := it.Scene.Show.OBS.SceneItems.SetSceneItemProperties(&itemParams)
	return it, err != nil
}

// TODO: Consider passing up the error instead of Item but this
//       currently matches
func (it *Item) Cache() (*Item, bool) {
	if apiResponse, err := it.Scene.Show.OBS.Client.Scenes.GetSceneList(); err == nil {
		for _, scene := range apiResponse.Scenes {
			if it.Scene.HasName(scene.Name) {
				for _, sceneItem := range scene.Sources {
					if it.HasName(sceneItem.Name) {
						if parsedItem, ok := ParseItem(it.Scene, sceneItem); ok {
							it.Id = parsedItem.Id // May not be necessary if static
							it.Parent = parsedItem.Parent
							it.Items = parsedItem.Items
							it.Scene = parsedItem.Scene
							return it, true
						}
					}
				}
			}
		}
	}
	return it, false
	//return errors.New("no scene found")
}

/////////////////////////////////////////////////////
// TODO: Rewrite Item.Cache() with this below??? Bc our going through
//       scenes is kinda silly
// GET / READ
//itemParams := sceneitems.GetSceneItemPropertiesParams{
//	Item:      &typedefs.Item{Name: itemName},
//	SceneName: sceneName,
//}

//apiResponse, err := sh.OBS.SceneItems.GetSceneItemProperties(&itemParams)

//cachedItem := sh.Scene(sceneName).Item(itemName)

//Bounds:    resp.Bounds,
//Crop:      resp.Crop,
//Position:  resp.Position,
//Rotation:  resp.Rotation,

// TODO: Take the API response and use it to update the local cache of
//       the item using our more complex and useful abstraction of
//       source or scene item

// TODO: This update requires us to do a write against the OBS WS API
//       so the change would be reflected within OBS
func (it *Item) Unlock() *Item {
	it.Layer.Locked = false
	it.Update()
	return it
}

func (it *Item) Lock() *Item {
	it.Layer.Locked = true
	it.Update()
	return it
}

func (it *Item) Unhide() *Item {
	it.Layer.Visible = true
	it.Update()
	return it
}

func (it *Item) Hide() *Item {
	it.Layer.Visible = false
	it.Update()
	return it
}

func (it Item) Print() {
	fmt.Printf("item: \n")
	fmt.Printf("  id: %v \n", it.Id)
	fmt.Printf("  name: %v \n", it.Name)
	//fmt.Printf("  type: %v \n", item.Type.String())
	// Its the same without bc by default it tries String()
	fmt.Printf("  type: %v \n", it.Type)
	fmt.Printf("  locked: %v \n", it.Locked)
	fmt.Printf("  visible: %v \n", it.Visible)
	fmt.Printf("  len(items): %v \n", len(it.Items))
	fmt.Printf("  scene(*): %v \n", it.Scene)
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

func (it *Item) ParseItem(item typedefs.SceneItem) (*Item, bool) {
	if parsedItem, ok := ParseItem(it.Scene, item); ok {
		it.Items = append(it.Items, parsedItem)
		return parsedItem, true
	}
	return nil, false
}

// NOTE: Alias
func (sc *Scene) ParseItem(item typedefs.SceneItem) (*Item, bool) {
	if parsedItem, ok := ParseItem(sc, item); ok {
		sc.Items = append(sc.Items, parsedItem)
		return parsedItem, true
	}
	return nil, false
}

func ParseItem(scene *Scene, item typedefs.SceneItem) (*Item, bool) {
	// TODO: Should be validation on if scene/item already exists

	//	Alignment float64 `json:"alignment,omitempty"`
	//  [x]	Cx float64 `json:"cx,omitempty"`
	//	[x] Cy float64 `json:"cy,omitempty"`
	//	    // List of children (if this item is a group)
	//	[X] GroupChildren []SceneItem `json:"groupChildren,omitempty"`
	//	    // Scene item ID
	//	[X] Id int `json:"id,omitempty"`
	//	    // Whether or not this Scene Item is locked and can't be moved around
	//	[X] Locked bool `json:"locked,omitempty"`
	//	    // Whether or not this Scene Item is muted.
	//	[ ] Muted bool `json:"muted,omitempty"`
	//	    // The name of this Scene Item.
	//	[X] Name string `json:"name,omitempty"`
	//	    // Name of the item's parent (if this item belongs to a group)
	//	[X] ParentGroupName string `json:"parentGroupName,omitempty"`
	//	    // Whether or not this Scene Item is set to "visible".
	//	[X] Render bool `json:"render,omitempty"`
	//	[X] SourceCx float64 `json:"source_cx,omitempty"`
	//	[X] SourceCy float64 `json:"source_cy,omitempty"`
	//	    // Source type. Value is one of the following: "input", "filter", "transition", "scene" or "unknown"
	//	[X] Type string `json:"type,omitempty"`
	//	[ ] Volume float64 `json:"volume,omitempty"`
	//	[X] X float64 `json:"x,omitempty"`
	//	[X] Y float64 `json:"y,omitempty"`

	cachedItem := &Item{
		// NOTE: Intentionally left out muted and volume to only keep that logic
		//       in the audiomixer and its audio sources
		Scene: scene,
		Id:    item.Id,
		Name:  item.Name,
		Type:  MarshalItemType(item.Type),
		Layer: Layer{
			Visible:   item.Render,
			Locked:    item.Locked,
			Alignment: int(item.Alignment),
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

	if 0 < len(item.ParentGroupName) {
		cachedItem.Parent, _ = scene.Item(item.ParentGroupName)
	} else if len(item.ParentGroupName) == 0 {
		scene.Items = append(scene.Items, cachedItem)
	} else if 0 < len(item.GroupChildren) {
		for _, childItem := range item.GroupChildren {
			parsedChildItem, _ := ParseItem(scene, childItem)
			cachedItem.Items = append(cachedItem.Items, parsedChildItem)
		}
	}

	return cachedItem, false
}

type Items []*Item

// NOTE: A recursive calling of Name to check child items is preferred
//       but OBS folders/grouping only supports 1 level so this is
//       adequate
//       And OBS does not support duplicate item naming (or scene)
func (its Items) Name(name string) (*Item, bool) {
	fmt.Printf("len(its)=(%v)\n", len(its))

	for index, item := range its {
		if item == nil {
			fmt.Printf("%v) item is somehow nil in loop (%v)\n", index, item)
		} else {
			if item.HasName(name) {
				return item, true
			} else {
				if item.IsFolder() {
					for _, childItem := range item.Items {
						if childItem.HasName(name) {
							return childItem, true
						}
					}
				}
			}
		}
	}
	return nil, false
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

func (sc *Scene) Unlock() {
	sc.Locked = false
	sc.Update()
}

func (sc *Scene) Lock() {
	sc.Locked = true
	sc.Update()
}

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
func (its Items) Folders() (folderItems Items) {
	for _, item := range its {
		if item.IsFolder() {
			folderItems = append(folderItems, item)
		}
	}
	return folderItems
}

func (sc Scene) Item(name string) (*Item, bool) {
	return sc.Items.Name(name)
}

func (sc *Scene) HasName(name string) bool {
	return sc != nil && len(sc.Name) == len(name) && sc.Name == name
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

//func (obs OBS) ActiveShow() *Show {
//	return nil
//}

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

//func (obs *OBS) CurrentScene() (string, error) {
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

// TODO: This can be relied upon unless CurrentScene() above functions to
//       cache it from remote OBS instance over OBS ws API
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

func (scs Scenes) Name(name string) (*Scene, bool) {
	for _, scene := range scs {
		if scene.HasName(name) {
			return scene, true
		}
	}
	return nil, false
}

// TODO: THis is nice bc we get chaining additions scs.Add(sc1).Add(sc2)
//   okay
//     so
///     maybe like name and items?

func (sh *Show) NewScene(name string) (*Scene, bool) {
	if _, ok := sh.Scene(name); ok {
		fmt.Printf("Scene already exists, skipping (should update rly)\n")
		return sh.Scenes.Name(name)
	}
	newScene := &Scene{
		Show:      sh,
		Name:      name,
		IsCurrent: false,
	}

	sh.Scenes = append(sh.Scenes, newScene)
	return newScene, true
}

// TODO: The OBS ws api should be interacted with through Update() alone
//       and not scattered through the code so its really hard to maintain
func (sc *Scene) Current() error {
	sc.IsCurrent = true
	r := scenes.SetCurrentSceneParams{
		SceneName: sc.Name,
	}
	_, err := sc.Show.OBS.Scenes.SetCurrentScene(&r)
	return err
}

func (sc Scene) Update() (*Scene, bool) {
	// TODO: No real easy way to do this unless perhaps updating scene
	//       list at once or deleting and re-creating?
	return sc, false
}

// TODO:
//   type GetSceneItemListParams (goobs) to request the []*SceneItem
func (sc *Scene) Cache() (*Scene, bool) {
	if apiResponse, err := sc.Show.OBS.Client.Scenes.GetSceneList(); err == nil {
		for _, scene := range apiResponse.Scenes {
			if sc.HasName(scene.Name) {
				sc.Items = Items{}
				for _, item := range scene.Sources {
					sc.ParseItem(item)
				}
				return sc, true
			}
		}
	}
	return sc, false
}

// TODO: Need to do an UPDATE_SCENE function obvio
//func (sh *Show) UpdateScene(sceneName string) (*Scene, bool) {
//	// TODO: Update the OBS scene via the ws API using the current
//	//       state of the cached scene (and whatever changes it has)
//
//	var cachedScene *Scene
//	//var cachedItem *Item
//	return cachedScene, false
//}

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

// TODO: Build update from OBS (since OBS has a UI for the time being)
//func (sh Show) ItemAttributes() string {
//

func (sh Show) Scene(sceneName string) (*Scene, bool) {
	return sh.Scenes.Name(sceneName)
}

// TODO: If we don't pass up error then we will fail to provide
//       developers using our library useful errors

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

//func (obs OBS) Volume() {
//	//// TODO: List must be obtained from OBS object, and preferably it should be
//	////       cached in OBS object.
//	//name := list.Sources[0].Name
//
//	//apiResponse, _ := obs.Client.Sources.GetVolume(
//	//	&sources.GetVolumeParams{
//	//		Source: name,
//	//	},
//	//)
//	//fmt.Printf("%s is muted? %t\n", name, apiResponse.Muted)
//}

// TODO: Toggle Mute On Source
//        * Mute + reduce volume to 0 [hide completely]
//        * Mute

//func (obs OBS) Sources() {
//	list, _ := obs.Client.Sources.GetSourcesList()
//
//	for _, v := range list.Sources {
//		if _, err := obs.Client.Sources.SetVolume(&sources.SetVolumeParams{
//			Source: v.Name,
//			Volume: rand.Float64(),
//		}); err != nil {
//			panic(err)
//		}
//	}
//
//	if len(list.Sources) == 0 {
//		fmt.Printf("No sources!\n")
//		// TODO: Why would we exit? And exit 0??? Thats no error
//		//os.Exit(0)
//	}
//}
