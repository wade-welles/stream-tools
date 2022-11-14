package obs

import (
	"fmt"
	"time"

	// OBS
	goobs "github.com/andreykaipov/goobs"

	events "github.com/andreykaipov/goobs/api/events"
	sceneitems "github.com/andreykaipov/goobs/api/requests/sceneitems"
	scenes "github.com/andreykaipov/goobs/api/requests/scenes"
	//studiomode "github.com/andreykaipov/goobs/api/requests/ui"
	//typedefs "github.com/andreykaipov/goobs/api/typedefs"
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
	Position

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
type BlendMode uint8

const (
	Normal BlendMode = iota
	Additive
	Subtract
	Screen
	Multiply
	Lighten
	Darken
)

func (bm BlendMode) String() string {
	switch bm {
	case Additive:
		return "additive"
	case Subtract:
		return "subtract"
	case Multiply:
		return "multiply"
	case Lighten:
		return "lighten"
	case Darken:
		return "darken"
	default: // Normal
		return "normal"
	}
}

func MarshalBlendMode(mode string) BlendMode {
	switch mode {
	case Additive.String():
		return Additive
	case Subtract.String():
		return Subtract
	case Multiply.String():
		return Multiply
	case Lighten.String():
		return Lighten
	case Darken.String():
		return Darken
	default: // Normal
		return Normal
	}
}

type Item struct {
	Id   int
	Name string
	Type ItemType
	Blend BlendMode

	// TODO: Add
	//        type
	//        folder/relation
	//
	Layer

	// NOTE: All that is needed for a tree structure
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

// TODO: It looks like most of these values are no longer accessible, and the
// individual attributes that are available have to be interacted with
// individually ;_;
func (it Item) Update() (Item, bool) {
	// TODO: Not even exactly sure what "enabled" means but for now will assume
	// visible boolean is the "enabled" value. This went from not great to worse. 
	itemEnabledParams := sceneitems.SetSceneItemEnabledParams{
		SceneName:        it.Scene.Name,
		SceneItemId:      float64(it.Id),
		SceneItemEnabled: &it.Visible,
	}
	// TODO: Technically now we should be checking the returned values for changes
	// to confirm the update was actually successful
	_, err := it.Scene.Show.OBS.SceneItems.SetSceneItemEnabled(&itemEnabledParams)

	itemLockedParams := sceneitems.SetSceneItemLockedParams{
		SceneName:       it.Scene.Name,
		SceneItemId:     float64(it.Id),
		SceneItemLocked: &it.Locked, 
	}
	_, err = it.Scene.Show.OBS.SceneItems.SetSceneItemLocked(&itemLockedParams)

	itemIndexParams := sceneitems.SetSceneItemIndexParams{
		SceneName:       it.Scene.Name,
		SceneItemId:     float64(it.Id),
		SceneItemIndex:  float64(it.Index),
	}
	_, err = it.Scene.Show.OBS.SceneItems.SetSceneItemIndex(&itemIndexParams)

	// TODO: Should eventually create a enumerator to work with the string using
	// ints but it accepts a string so blank until more work on this will do. we
	// just wont submit it. though may be pointless since we almost never use this
	// but then again maybe when we have programmatic control over it, then new
	// uses will become apparent
	itemBlendModeParams := sceneitems.SetSceneItemBlendModeParams{
		SceneName:          it.Scene.Name,
		SceneItemId:        float64(it.Id),
		SceneItemBlendMode: Normal.String(),
	}
	_, err = it.Scene.Show.OBS.SceneItems.SetSceneItemBlendMode(&itemBlendModeParams)

	// TODO: Either
	//   1) will need to create our own transform object
	//   2) possibly better for simplicity, we could store some of these values, 
	//      then generate the transform params necessary to move the item from 
	//      previous position to the new position based on 
	//itemTransformParams := sceneitems.SetSceneItemTransformParams{
	//	SceneName:          it.Scene.Name,
	//	SceneItemId:        it.Id,
	//	SceneItemTransform: &typedefs.SceneItemTransform{
	//		Alignment:       0, // float64
	//		SourceHeight:    0, // float64
	//		SourceWidth:     0, // float64
	//		Width:           0, // float64
	//		Height:          0, // float64
	//		Rotation:        0, // float64
	//		PositionX:       0, // float64
	//		PositionY:       0, // float64
	//		BoundsType:     "", // string
	//		BoundsAlignment: 0, // float64
  //    BoundsHeight:    0, // float64
	//		BoundsWidth:     0, // float64
	//		CropBottom:      0, // float64
	//		CropLeft:        0, // float64
	//		CropRight:       0, // float64
	//		CropTop:         0, // float64
	//		ScaleX:          0, // float64
	//		ScaleY:          0, // float64
	//	},
	//}
	//_, err := it.Scene.Show.OBS.SceneItems.SetSceneItemTransformProperties(
	//	&itemTransformParams,
	//)

	return it, err != nil
}

//func (sh Show) OBSScenes() (obsScenes Scenes) {
//	if apiResponse, err := sh.OBS.Client.Scenes.GetSceneList(); err == nil {
//		sh.Current, _ = sh.Scene(apiResponse.CurrentScene)
//
//		for _, scene := range apiResponse.Scenes {
//			if newScene, ok := sh.NewScene(scene.Name); ok {
//				obsScenes = append(obsScenes, newScene)
//
//			}
//		}
//	}
//
//	return obsScenes
//}

// TODO: Consider passing up the error instead of Item but this
//       currently matches
//func (it *Item) Cache() (*Item, bool) {
//	// TODO: Maybe separate function to get Scene List; and make that
//	//       separate and then call it here to simplify this
//	if apiResponse, err := it.Scene.Show.OBS.Client.Scenes.GetSceneList(); err == nil {
//
//		for _, scene := range apiResponse.Scenes {
//			if it.Scene.HasName(scene.Name) {
//				for _, sceneItem := range scene.Sources {
//					if it.HasName(sceneItem.Name) {
//						if parsedItem, ok := ParseItem(it.Scene, sceneItem); ok {
//							it.Id = parsedItem.Id // May not be necessary if static
//							it.Parent = parsedItem.Parent
//							it.Items = parsedItem.Items
//							it.Scene = parsedItem.Scene
//							return it, true
//						}
//					}
//				}
//			}
//		}
//	}
//	return it, false
//	//return errors.New("no scene found")
//}

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
	return it
}

func (it *Item) Lock() *Item {
	it.Layer.Locked = true
	return it
}

func (it *Item) Unhide() *Item {
	it.Layer.Visible = true
	return it
}

func (it *Item) Hide() *Item {
	it.Layer.Visible = false
	return it
}

func (it Item) Print() Item {
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
	return it
}

//func (it *Item) ParseItem(item typedefs.SceneItem) (*Item, bool) {
//	if parsedItem, ok := ParseItem(it.Scene, item); ok {
//		it.Items = append(it.Items, parsedItem)
//		return parsedItem, true
//	}
//	return nil, false
//}

// NOTE: Alias
//func (sc *Scene) ParseItem(item typedefs.SceneItem) (*Item, bool) {
//	if parsedItem, ok := ParseItem(sc, item); ok {
//		sc.Items = append(sc.Items, parsedItem)
//		return parsedItem, true
//	}
//	return nil, false
//}

//func (it *Item) ToggleVisibility() bool {
//	var ok bool
//	if it.Visible {
//		_, ok = it.Hide().Update()
//	} else {
//		_, ok = it.Unhide().Update()
//	}
//	return ok
//}

//func (it *Item) ToggleLock() bool {
//	var ok bool
//	if it.Locked {
//		_, ok = it.Unlock().Update()
//	} else {
//		_, ok = it.Lock().Update()
//	}
//	return ok
//}

//func ParseItem(scene *Scene, item typedefs.SceneItem) (*Item, bool) {
//	// TODO: Should be validation on if scene/item already exists
//
//	//	Alignment float64 `json:"alignment,omitempty"`
//	//  [x]	Cx float64 `json:"cx,omitempty"`
//	//	[x] Cy float64 `json:"cy,omitempty"`
//	//	    // List of children (if this item is a group)
//	//	[X] GroupChildren []SceneItem `json:"groupChildren,omitempty"`
//	//	    // Scene item ID
//	//	[X] Id int `json:"id,omitempty"`
//	//	    // Whether or not this Scene Item is locked and can't be moved around
//	//	[X] Locked bool `json:"locked,omitempty"`
//	//	    // Whether or not this Scene Item is muted.
//	//	[ ] Muted bool `json:"muted,omitempty"`
//	//	    // The name of this Scene Item.
//	//	[X] Name string `json:"name,omitempty"`
//	//	    // Name of the item's parent (if this item belongs to a group)
//	//	[X] ParentGroupName string `json:"parentGroupName,omitempty"`
//	//	    // Whether or not this Scene Item is set to "visible".
//	//	[X] Render bool `json:"render,omitempty"`
//	//	[X] SourceCx float64 `json:"source_cx,omitempty"`
//	//	[X] SourceCy float64 `json:"source_cy,omitempty"`
//	//	    // Source type. Value is one of the following: "input", "filter", "transition", "scene" or "unknown"
//	//	[X] Type string `json:"type,omitempty"`
//	//	[ ] Volume float64 `json:"volume,omitempty"`
//	//	[X] X float64 `json:"x,omitempty"`
//	//	[X] Y float64 `json:"y,omitempty"`
//
//	cachedItem := &Item{
//		// NOTE: Intentionally left out muted and volume to only keep that logic
//		//       in the audiomixer and its audio sources
//		Scene: scene,
//		Id:    item.Id,
//		Name:  item.Name,
//		Type:  MarshalItemType(item.Type),
//		Layer: Layer{
//			Visible:   item.Render,
//			Locked:    item.Locked,
//			Alignment: int(item.Alignment),
//			Position: Position{
//				X: item.X,
//				Y: item.Y,
//			},
//			Dimensions: Dimensions{
//				Width:  item.Cx,
//				Height: item.Cy,
//			},
//			Source: Dimensions{
//				Width:  item.SourceCx,
//				Height: item.SourceCy,
//			},
//		},
//	}
//
//	if 0 < len(item.ParentGroupName) {
//		cachedItem.Parent, _ = scene.Item(item.ParentGroupName)
//	} else if len(item.ParentGroupName) == 0 {
//		scene.Items = append(scene.Items, cachedItem)
//	} else if 0 < len(item.GroupChildren) {
//		for _, childItem := range item.GroupChildren {
//			parsedChildItem, _ := ParseItem(scene, childItem)
//			cachedItem.Items = append(cachedItem.Items, parsedChildItem)
//		}
//	}
//
//	return cachedItem, false
//}

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
	Index     int
	Visible   bool
	Locked    bool
	Alignment int
	Rotation  float64
	Position
	Dimensions
	Source Dimensions
}

// TODO: Better track Current scene (its what is transitioned to last, we
//       ideally will have a log like system, or at least dates) and get
//       rid of this bool below bc its bad 4 rlly
type Scene struct {
	Layer

	Name string

	Show  *Show
	Items Items

	IsCurrent   bool
	IsPreviewed bool
}

func (its Items) Unlocked() (unlockedItems Items) {
	for _, item := range its {
		if !item.Locked {
			unlockedItems = append(unlockedItems, item)
		}
	}
	return unlockedItems
}

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
	return len(sc.Name) == len(name) && sc.Name == name
}

type Scenes []*Scene

func (scs Scenes) Size() int     { return len(scs) }
func (scs Scenes) First() *Scene { return scs[0] }
func (scs Scenes) Last() *Scene  { return scs[scs.Size()-1] }

// TODO: I thought Source == SceneItem

// NOTE: Scene Collection == All (Show|Channel) Scenes

type Mode uint8 // 0..255

const (
	UndefinedMode Mode = iota
	StudioMode
	StreamingMode
	RecordingMode
)

func (sm Mode) String() string {
	switch sm {
	case StudioMode:
		return "studio"
	case StreamingMode:
		return "streaming"
	case RecordingMode:
		return "recording"
	default: // UndefinedMode
		return "undefined"
	}
}

func MarshalMode(modeName string) Mode {
	switch modeName {
	case StudioMode.String():
		return StudioMode
	case StreamingMode.String():
		return StreamingMode
	case RecordingMode.String():
		return RecordingMode
	default: // Undefined
		return UndefinedMode
	}
}

// TODO: Populate Modes, studio, streaming, recording; and have methods
//       for adding/removing/lsiting/checking-if-present
type Show struct {
	OBS *OBS

	Name string

	Current *Scene
	Preview *Scene

	Mode []Mode

	Scenes Scenes

	// Profile
	// All items?
	// History/Log of bumpers played (bgs used)
}

//type Source struct {
//
//}

//type Usage struct {
//	CPU    int
//	Disk   int
//	Memory int
//}
//
//type Stats struct {
//	Usage     Usage
//	Streaming bool
//	Recording bool
//
//	FramesPerSecond uint8
//
//	Time               uint
//	AvgerageRenderTime uint
//	FramesLost         uint
//	FramesSkipped      uint
//	FramesDropped      uint
//	DataOutput         uint
//	Bitrate            uint
//}

// NOTE: The API we are creating looks something like:
//          obs.Show.Scenes.First()
//          obs.Show.Scenes.First().Items.First().(Show()|.Hide())

type OBS struct {
	*goobs.Client

	Show *Show

	//Stats      *Stats
	//AudioMixer *AudioMixer

	//Sources []*goobs.Source
}

//func (obs OBS) StudioMode() bool {
//	response, err := obs.Client.StudioMode.GetStudioModeStatus()
//	if err != nil {
//		return false
//	} else {
//		return response.StudioMode
//	}
//}

//  TODO: This would be better as obs.Shows.Active(); but that means creating a
//  Shows type (POTENTIALLY, but ideally no)

//func (obs OBS) ActiveShow() *Show {
//	return nil
//}

// TODO: This should either give the full OBS object returned (like an init
// function) or this should be a method on OBS after it is created you connect.
func ConnectToOBS(host string) *goobs.Client {
	client, err := goobs.New(host)
	if err != nil {
		panic(err)
	}
	return client
}

// go Events()
func (obs OBS) Events() {
	for event := range obs.Client.IncomingEvents {
		switch e := event.(type) {
		case *events.SceneItemEnableStateChanged:
			fmt.Printf("Scene Item Enabled %-25q (%v): %v\n", e.SceneName, e.SceneItemId, e.SceneItemEnabled)
		default:
			fmt.Printf("Unhandled event: %#v\n", e)
		}
	}
}

// TODO: Switch to Scene

// TODO: List Scenes

/////////////////////////////////////////////////////////////
//func (obs OBS) Scenes() ([]*typedefs.Scene, error) {
//	apiResponse, err := obs.Client.Scenes.GetSceneList()
//	return apiResponse.Scenes, err
//}

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
func (sc Scene) Update() bool {
	// TODO: No real easy way to do this unless perhaps updating scene
	//       list at once or deleting and re-creating?
	//         Almost certainly want to use delete and recreate since
	//         no clear edit of Scene; just item
	return false
}

// TODO:
//   type GetSceneItemListParams (goobs) to request the []*SceneItem
func (sc *Scene) Cache() (*Scene, bool) {
	// GetSceneItemPropertiesParams represents the params body for the "GetSceneItemProperties" request.
	// Gets the scene specific properties of the specified source item.

	if apiResponse, err := sc.Show.OBS.Client.Scenes.GetSceneList(); err == nil {

		if currentScene, ok := sc.Show.Scene(apiResponse.CurrentProgramSceneName); ok {
			sc.Show.Current = currentScene
		}

		for _, scene := range apiResponse.Scenes {
			// TODO: Scene no longer comes with its sources
			if sc.HasName(scene.SceneName) {
				fmt.Printf("local scene cache, still exists in obs...\n")
				fmt.Printf("but doing nothing because scenes no longer contain their sources")
				//sc.Items = Items{}
				//for _, item := range scene.Sources {
				//	sc.ParseItem(item)
				//}
				//return sc, true
			}
		}
	}
	return sc, false
}

func (sc *Scene) Transition(sleepDuration ...time.Duration) (*Scene, bool) {
	if 0 < len(sleepDuration) {
		fmt.Printf("sleeping \n")
		time.Sleep(sleepDuration[0])
	}

	_, err := sc.Show.OBS.Scenes.SetCurrentProgramScene(
		&scenes.SetCurrentProgramSceneParams{
			SceneName: sc.Name,
		},
	)

	if err == nil {
		sc.IsCurrent = true
		sc.Show.Current = sc
	}

	return sc, err == nil
}

func (sh Show) SceneNames() (sceneNames []string) {
	for _, scene := range sh.Scenes {
		sceneNames = append(sceneNames, scene.Name)
	}
	return sceneNames
}

func (sh *Show) Cache() (*Show, bool) {
	//sh.Scenes = Scenes{}
	if apiResponse, err := sh.OBS.Client.Scenes.GetSceneList(); err == nil {

		// TODO: Its a way,... right?? ?? hello ?
		obsSceneNames := []string{}
		cachedSceneNames := []string{}
		for _, cachedScene := range sh.Scenes {
			cachedSceneNames = append(cachedSceneNames, cachedScene.Name)
		}

		for _, scene := range apiResponse.Scenes {
			obsSceneNames = append(obsSceneNames, scene.SceneName)

			if cachedScene, ok := sh.Scene(scene.SceneName); ok {
				cachedScene.Cache()
			} else {
				if newScene, ok := sh.NewScene(scene.SceneName); ok {
					newScene.Cache()
				}
				// NOTE: Here we would create a cached scene from the data that
				//       does exist in the OBS state
				//       And left over scenes would need to be purged
				//         Again, just clearing and rebuilding seems easier;
				//         but we would likely lose data since our models are
				//         more complex than the OBS data models
			}
		}

		return sh, len(sh.Scenes) == len(apiResponse.Scenes)
	} else {
		return sh, false
	}
}

func (sh Show) Scene(sceneName string) (*Scene, bool) {
	return sh.Scenes.Name(sceneName)
}

//type Params requests.ParamsBasic
//type Response requests.ResponseBasic

//type Params struct {
//	*requests.ParamsBasic
//
//	Name  string
//	Value string
//}

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

//func (sc *Scene) Preview() (*Scene, bool) {
//	apiRequest := studiomode.SetPreviewSceneParams{
//		SceneName: sc.Name,
//	}
//
//	_, err := sc.Show.OBS.Client.StudioMode.SetPreviewScene(&apiRequest)
//	if err == nil {
//		sc.IsPreviewed = true
//		sc.Show.Preview = sc
//	}
//
//	return sc, err == nil
//
//}

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

// ?? Uses own SceneItem (WITH unique info? different?)
//	sceneItemsParams := sceneitems.GetSceneItemListParams{
//		SceneName: sc.Name,
//	}
//	apiResponse, err := sc.Show.OBS.SceneItems.GetSceneItemProperties(&sceneItemsParams)
//	if err != nil {
//		return err
//	}
