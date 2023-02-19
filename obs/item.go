package obs

import (
	"fmt"
	"strings"
	"time"

	sceneitems "github.com/andreykaipov/goobs/api/requests/sceneitems"
)

// TODO: Move the item content here
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

// TODO: This might be better if we dont have unknown and default to Scene type;
// but I'm not sure at the point of writing this if unknown is a possible option
// because sometimes in the API it is
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

type Item struct {
	Id   int
	Name string
	Type ItemType

	// TODO: These are the fields on the API SceneItem object, but remember we are
	// not doing a direct translation but an abstraction to simplify these sort of
	// things; but reconcile needing to fully support the API and wanting to
	// simplfiy both storage and interaction with it.
	// InputKind  string
	// IsGroup    bool

	Layer
	// TODO: Almost certainly blendmode should be inside of Layer but it also only
	// applies to items and not scenes, the idea behind it is that layer holds
	// information pertaining to both
	Blend BlendMode
	// SceneItemEnabled bool (visiblity?)
	// SceneItemIndex   int
	// SceneItemLocked  bool
	// SceneItemTransform scenes.SceneItemTransform (oh, current transform)

	// NOTE: All that is needed for a tree structure
	Parent *Item
	Items  Items

	Scene *Scene
	OBS   *ShowAPI

	LastCachedAt time.Time
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
//
//	so the change would be reflected within OBS
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

// TODO: This should be possible to pull multiple items, it should return
// []*Item, and we should be adding all of the matches to the slice and
// returning the slice
func (its Items) NameContains(searchText string) (*Item, bool) {
	fmt.Printf("len(its)=(%v)\n", len(its))

	for _, item := range its {
		if item != nil {
			if item.NameContains(searchText) {
				return item, true
			} else {
				if item.IsFolder() {
					for _, childItem := range item.Items {
						if childItem.NameContains(searchText) {
							return childItem, true
						}
					}
				}
			}
		}
	}
	return nil, false
}

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
	Rectangle
	Source Rectangle
}

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

// TODO: One issue with this is that something can be a folder without
//       having any items within it; and while that distinction may be
//       useless to us, it may cause issues in the future since there
//       is no simple way to switch between types beyond deleting and
//       recreating

// NOTE: Reconcile the fact that this exists and the type exists
//       and they wont be the same

func (it Item) IsFolder() bool { return (0 < len(it.Items)) }

func (it Item) NameContains(searchText string) bool {
	return strings.Contains(it.Name, searchText)
}

func (it Item) HasName(name string) bool {
	return len(it.Name) == len(name) && it.Name == name
}

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
	_, err := it.OBS.Items.SetSceneItemEnabled(&itemEnabledParams)

	itemLockedParams := sceneitems.SetSceneItemLockedParams{
		SceneName:       it.Scene.Name,
		SceneItemId:     float64(it.Id),
		SceneItemLocked: &it.Locked,
	}
	_, err = it.OBS.Items.SetSceneItemLocked(&itemLockedParams)

	// NOTE: This is really awkward, the scenes.SceneItem object stores both Id
	// and Index as int, but expects to pass it as float64.
	itemIndexParams := sceneitems.SetSceneItemIndexParams{
		SceneName:      it.Scene.Name,
		SceneItemId:    float64(it.Id),
		SceneItemIndex: float64(it.Index),
	}
	_, err = it.OBS.Items.SetSceneItemIndex(&itemIndexParams)

	itemBlendModeParams := sceneitems.SetSceneItemBlendModeParams{
		SceneName:          it.Scene.Name,
		SceneItemId:        float64(it.Id),
		SceneItemBlendMode: Normal.String(),
	}
	_, err = it.OBS.Items.SetSceneItemBlendMode(&itemBlendModeParams)

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
	//_, err := it.OBS.SceneItems.SetSceneItemTransformProperties(
	//	&itemTransformParams,
	//)

	return it, err == nil
}

// TODO: Consider passing up the error instead of Item but this
//
//	currently matches
func (it *Item) Cache() (*Item, bool) {
	// TODO: Maybe separate function to get Scene List; and make that
	//       separate and then call it here to simplify this
	if apiResponse, err := it.OBS.Scenes.GetSceneList(); err == nil {

		for _, scene := range apiResponse.Scenes {
			if it.Scene.HasName(scene.SceneName) {
				apiResponse, err := it.OBS.Items.GetSceneItemList(&sceneitems.GetSceneItemListParams{SceneName: scene.SceneName})
				if err == nil {
					for _, item := range apiResponse.SceneItems {
						if it.HasName(item.SourceName) {
							fmt.Printf("  item:\n")
							fmt.Printf("    name: %v\n", item.SourceName)
						}
					}
					// TODO: Decide to parse if does not exist or does this just cache
					// things we already have a copy of? Parse is for initialization
				}

				fmt.Printf("  item_count: %v\n", len(apiResponse.SceneItems))

				// TODO:[!] This wont work anymore because scene no longer contains scene
				// items
				//for _, sceneItem := range scene.SceneItems {
				//	if it.HasName(sceneItem.Name) {
				//		fmt.Printf("went through scene list and then iterated through that scenes items to find the correct one\n")
				//		//if parsedItem, ok := ParseItem(it.Scene, sceneItem); ok {
				//		//	it.Id = parsedItem.Id // May not be necessary if static
				//		//	it.Parent = parsedItem.Parent
				//		//	it.Items = parsedItem.Items
				//		//	it.Scene = parsedItem.Scene
				//		//	return it, true
				//		//}
				//	}
				//}
			}
		}
	}
	return it, false
	//return errors.New("no scene found")
}

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

func (its Items) Size() int    { return len(its) }
func (its Items) First() *Item { return its[0] }
func (its Items) Last() *Item  { return its[its.Size()-1] }

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
