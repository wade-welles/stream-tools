package obs

import (
	"fmt"
	"time"

	sceneitems "github.com/andreykaipov/goobs/api/requests/sceneitems"
	scenes "github.com/andreykaipov/goobs/api/requests/scenes"
	typedefs "github.com/andreykaipov/goobs/api/typedefs"
)

// TODO: Move scene content here
type Scenes []*Scene

func (scs Scenes) Size() int     { return len(scs) }
func (scs Scenes) First() *Scene { return scs[0] }
func (scs Scenes) Last() *Scene  { return scs[scs.Size()-1] }

func (scs Scenes) IsEmpty() bool { return scs.Size() == 0 }

// TODO: Add reverse to get order in the OBS GUI

func (scs Scenes) Name(name string) (*Scene, bool) {
	for _, scene := range scs {
		if scene.HasName(name) {
			return scene, true
		}
	}
	return nil, false
}

// TODO: Better track Current scene (its what is transitioned to last, we
//
//	ideally will have a log like system, or at least dates) and get
//	rid of this bool below bc its bad 4 rlly
//

// NOTE: We will store order of scenes in position in the Scenes slice
type Scene struct {
	Name string

	Show *Show
	// NOTE/TODO:
	// vs scene.show.obs.client (3 pointer lookup of objects) taking the bit of
	// extra memory to directly save OBSClient in scene, so it can access it
	// without the jumps. Then we will benchmark these two options to demonstrate
	// something about Go that is important
	OBS   *ShowAPI
	Items Items

	IsCurrent   bool
	IsPreviewed bool
}

// pass config?
func NewEmptyScene(show *Show) (*Scene, error) {
	// Needs to be minimal but still functional
	return &Scene{
		Name: "",
		Show: show,
		OBS:  new(ShowAPI),
	}, nil
}

func (sc Scene) Item(name string) (*Item, bool) {
	return sc.Items.Name(name)
}

func (sc Scene) ItemNameContains(searchText string) (*Item, bool) {
	return sc.Items.NameContains(searchText)
}

func (sc *Scene) HasName(name string) bool {
	return len(sc.Name) == len(name) && sc.Name == name
}

// TODO: The OBS ws api should be interacted with through Update() alone
//
//	and not scattered through the code so its really hard to maintain
func (sc Scene) Update() bool {
	// TODO: No real easy way to do this unless perhaps updating scene
	//       list at once or deleting and re-creating?
	//         Almost certainly want to use delete and recreate since
	//         no clear edit of Scene; just item
	return false
}

//func (sc *Scene) Preview() (*Scene, bool) {
//	apiRequest := studiomode.SetPreviewSceneParams{
//		SceneName: sc.Name,
//	}
//
//	_, err := sc.OBS.StudioMode.SetPreviewScene(&apiRequest)
//	if err == nil {
//		sc.IsPreviewed = true
//		sc.Show.Preview = sc
//	}
//
//	return sc, err == nil
//
//}

// obs.Scenes.First().Preview() => sets the preview in studiomode
//func (sc Scene) Preview() error {
//	_, err := sc.OBS.StudioMode.SetPreviewScene(
//		studiomode.SetPreviewSceneParams{
//			SceneName: sc.Name,
//		},
//	)
//	return err
//}

// NOTE: Alias
// TODO: Maybe pass ID through and use that to lookup via *api.Client access to
// Items of the scene
func (sc *Scene) ParseItem(item *typedefs.SceneItem) (*Item, error) {

	//func ParseItem(scene *Scene, item *typedefs.SceneItem) (*Item, bool) {
	// TODO: Should be validation on if scene/item already exists
	parsedItem := &Item{
		// NOTE: Intentionally left out muted and volume to only keep that logic
		//       in the audiomixer and its audio sources
		// TODO: Store index because its the layer position and will be important
		OBS: &ShowAPI{
			Scenes: &scenes.Client{Client: sc.Show.OBS.WS.Client},
			Items:  &sceneitems.Client{Client: sc.Show.OBS.WS.Client},
		},
		Scene: sc,
		Id:    item.SceneItemID,
		Name:  item.SourceName,
		Type:  MarshalItemType(item.SourceType),
		Layer: Layer{
			// TODO: Not sure if this is what was render
			Visible:   item.SceneItemEnabled,
			Locked:    item.SceneItemLocked,
			Alignment: int(item.SceneItemTransform.Alignment),
			Rectangle: Rectangle{
				Position: Position{
					X: item.SceneItemTransform.PositionX,
					Y: item.SceneItemTransform.PositionY,
				},
				Width:  item.SceneItemTransform.Width,
				Height: item.SceneItemTransform.Height,
			},
			Source: Rectangle{
				Width:  item.SceneItemTransform.SourceWidth,
				Height: item.SceneItemTransform.SourceHeight,
			},
		},
	}

	parsedItem.Cache()

	// TODO We only get if its a group, then presumably we do a GetSources or
	// GetSceneItems on this specific scene item because its a group, meaning it
	// has items that need to be parsed

	// TODO: There is a GetGroupList() now too; look into that

	//if item.IsGroup {
	//	parsedItem.Parent, _ = scene.Item(item.ParentGroupName)
	//} else if len(item.ParentGroupName) == 0 {
	//	scene.Items = append(scene.Items, parsedItem)
	//} else if 0 < len(item.GroupChildren) {
	//	for _, childItem := range item.GroupChildren {
	//		parsedChildItem, _ := ParseItem(scene, childItem)
	//		parsedItem.Items = append(parsedItem.Items, parsedChildItem)
	//	}
	//}

	return parsedItem, nil
}

type Items []*Item

// NOTE: A recursive calling of Name to check child items is preferred
//       but OBS folders/grouping only supports 1 level so this is
//       adequate
//       And OBS does not support duplicate item naming (or scene)

// TODO: Add ability to pull item based on a search term so we can pull out
// something with an overly complex name like "Primary Terminal (VIM Window)"
// but we want to just be able to check if it has for example "VIM" in the name

// TODO:
//
//	type GetSceneItemListParams (goobs) to request the []*SceneItem
func (sc *Scene) Cache() (*Scene, bool) {
	fmt.Printf("caching scene, and its associated items...")
	// GetSceneItemPropertiesParams represents the params body for the "GetSceneItemProperties" request.
	// Gets the scene specific properties of the specified source item.

	//if apiResponse, err := sc.OBS.Scenes.GetSceneList(); err == nil {

	//if currentScene, ok := sc.Show.Scene(apiResponse.CurrentProgramSceneName); ok {
	//	sc.Show.Current = currentScene
	//}

	//for _, scene := range apiResponse.Scenes {
	// TODO: Scene no longer comes with its sources
	//if sc.HasName(scene.SceneName) {
	//	fmt.Printf("local scene cache, still exists in obs...\n")
	//	fmt.Printf("but doing nothing because scenes no longer contain their sources")
	//}

	// TODO This is not *api.Client but *goobs.Client which makes the typecast not
	// work
	// NOTE: Lets benchmark this against a direct OBSClient object
	apiResponse, err := sc.OBS.Items.GetSceneItemList(
		&sceneitems.GetSceneItemListParams{
			SceneName: sc.Name,
		},
	)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("apiResponse: %v\n", apiResponse)
		fmt.Printf("apiResponse.SceneItems: %v\n", len(apiResponse.SceneItems))
	}

	// NOTE: Clear existing scene items
	sc.Items = Items{}

	// NOTE: Repopulate scene items

	for _, item := range apiResponse.SceneItems {
		sc.ParseItem(item)
	}

	//return sc, true

	//}
	//}
	return sc, false
}

func (sc *Scene) Transition(sleepDuration ...time.Duration) (*Scene, bool) {
	if 0 < len(sleepDuration) {
		fmt.Printf("sleeping \n")
		time.Sleep(sleepDuration[0])
	}

	_, err := sc.OBS.Scenes.SetCurrentProgramScene(
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
