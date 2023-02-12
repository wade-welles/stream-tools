package obs

import (
	"fmt"
	"time"

	goobs "github.com/andreykaipov/goobs"
	sceneitems "github.com/andreykaipov/goobs/api/requests/sceneitems"
	scenes "github.com/andreykaipov/goobs/api/requests/scenes"
	typedefs "github.com/andreykaipov/goobs/api/typedefs"
)

// TODO: Move scene content here
type Scenes []*Scene

func (scs Scenes) Size() int     { return len(scs) }
func (scs Scenes) First() *Scene { return scs[0] }
func (scs Scenes) Last() *Scene  { return scs[scs.Size()-1] }

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
	OBSClient *goobs.Client
	Items     Items

	IsCurrent   bool
	IsPreviewed bool
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
//	_, err := sc.Show.OBS.Client.StudioMode.SetPreviewScene(&apiRequest)
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
//	_, err := sc.Show.OBS.Client.StudioMode.SetPreviewScene(
//		studiomode.SetPreviewSceneParams{
//			SceneName: sc.Name,
//		},
//	)
//	return err
//}

// NOTE: Alias
func (sc *Scene) ParseItem(item *typedefs.SceneItem) (*Item, bool) {
	if parsedItem, ok := ParseItem(sc, item); ok {
		sc.Items = append(sc.Items, parsedItem)
		return parsedItem, true
	}
	return nil, false
}

// TODO:
//
//	type GetSceneItemListParams (goobs) to request the []*SceneItem
func (sc *Scene) Cache() (*Scene, bool) {
	fmt.Printf("caching scene, and its associated items...")
	// GetSceneItemPropertiesParams represents the params body for the "GetSceneItemProperties" request.
	// Gets the scene specific properties of the specified source item.

	//if apiResponse, err := sc.Show.OBS.Client.Scenes.GetSceneList(); err == nil {

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
	apiResponse, err := sc.OBSClient.SceneItems.GetSceneItemList(
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

	_, err := sc.OBSClient.Scenes.SetCurrentProgramScene(
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
