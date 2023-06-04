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
	fmt.Sprintf("name we are looking for: %v\n", name)
	for _, scene := range scs {
		fmt.Sprintf("  scene: %v \n", scene)
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

	Broadcast *Broadcast
	// NOTE/TODO:
	// vs scene.show.obs.client (3 pointer lookup of objects) taking the bit of
	// extra memory to directly save OBSClient in scene, so it can access it
	// without the jumps. Then we will benchmark these two options to demonstrate
	// something about Go that is important
	Items Items

	IsCurrent   bool
	IsPreviewed bool
}

func NewEmptyScene(bc *Broadcast) *Scene {
	return &Scene{
		Name:      "",
		Broadcast: bc,
	}
}

func (sc Scene) Item(name string) (*Item, bool) {
	return sc.Items.Name(name)
}

func (sc Scene) ItemNameContains(searchText string) (*Item, bool) {
	return sc.Items.NameContains(searchText)
}

func (sc *Scene) HasName(name string) bool {
	return sc != nil && sc.Name == name
}

func (sc *Scene) ParseItem(item *typedefs.SceneItem) (*Item, error) {
	//func ParseItem(scene *Scene, item *typedefs.SceneItem) (*Item, bool) {
	// TODO: Should be validation on if scene/item already exists
	parsedItem := &Item{
		// NOTE: Intentionally left out muted and volume to only keep that logic
		//       in the audiomixer and its audio sources
		// TODO: Store index because its the layer position and will be important
		Broadcast: sc.Broadcast,
		//OBS: &ShowAPI{
		//	Scenes: &scenes.Client{Client: sc.Show.OBS.WS.Client},
		//	Items:  &sceneitems.Client{Client: sc.Show.OBS.WS.Client},
		//},
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

	return parsedItem, nil
}

type Items []*Item

func (sc *Scene) Cache() (*Scene, bool) {
	fmt.Printf("caching scene, and its associated items...")
	apiResponse, err := sc.Broadcast.Client.SceneItems.GetSceneItemList(
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

	return sc, false
}

func (sc *Scene) Transition(sleepDuration ...time.Duration) (*Scene, bool) {
	if 0 < len(sleepDuration) {
		fmt.Printf("sleeping \n")
		time.Sleep(sleepDuration[0])
	}

	_, err := sc.Broadcast.Client.Scenes.SetCurrentProgramScene(
		&scenes.SetCurrentProgramSceneParams{
			SceneName: sc.Name,
		},
	)

	if err == nil {
		sc.IsCurrent = true
		sc.Broadcast.Program = sc
	}

	return sc, err == nil
}
