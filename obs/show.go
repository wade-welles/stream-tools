package obs

import (
	"fmt"

	goobs "github.com/andreykaipov/goobs"
)

//type API struct {
//	WS *goobs.Client
//	*ShowAPI
//}
//
//type ShowAPI struct {
//	Scenes *scenes.Client
//	Items  *sceneitems.Client
//}

// for adding/removing/lsiting/checking-if-present
type Show struct {
	OBS *goobs.Client

	Name string

	Current *Scene
	Preview *Scene

	Mode []Mode

	Scenes Scenes

	// client.Scenes.GetSceneList()

	// Profile
	// All items?
	// History/Log of bumpers played (bgs used)
}

type Shows []*Show

func (sh *Show) ParseScene(name string) (*Scene, error) {
	// Validation
	// TODO: This may be completely unncessary now that we clear the show.Scenes
	// and then iterate over the apiResponse and rebuild each of them
	//if _, ok := sh.Scene(name); ok {
	//	fmt.Printf("Scene already exists, skipping (should update rly)\n")
	//	return sh.Scenes.Name(name)
	//}

	fmt.Printf("ParseScene ran...\n")

	parsedScene := &Scene{
		Show: sh,
		Name: name,
	}

	fmt.Printf("parsedScene.Name: %v\n", parsedScene.Name)

	//parsedScene.Cache()

	// TODO: Somehwere here we should be iterating over the scene items and adding
	// them before we append the scene to the show. Otherwise we have to iterate
	// over each of the show.Scenes and get their items.

	sh.Scenes = append(sh.Scenes, parsedScene)
	return parsedScene, nil
}

func (sh Show) Scene(sceneName string) (*Scene, bool) {
	fmt.Printf("sh.Scenes.Name(sceneName): and sceneName is %v\n", sceneName)
	// TODO: So the problem is when we hit this, Scenes doesn't exist properly
	if sh.Scenes == nil {
		fmt.Printf("scenes was nil\n")
		sh.Scenes = Scenes{}
	} else {
		fmt.Printf("scenes was nil\n")
	}

	return sh.Scenes.Name(sceneName)
}

func (sh Show) SceneNames() (sceneNames []string) {
	for _, scene := range sh.Scenes {
		sceneNames = append(sceneNames, scene.Name)
	}
	return sceneNames
}

// TODO: We should either make this CacheScenes() or make it cache both Scenes,
// and then their items
func (sh *Show) Cache() (*Show, bool) {

	// NOTE: For simplicity, for now we will just set scenes to empty and then
	// populate with API data. So we set show.Scenes to an empty slice of scenes
	sh.Scenes = Scenes{}

	if apiResponse, err := sh.OBS.Scenes.GetSceneList(); err == nil {

		// apiResponse == type GetSceneListResponse struct {
		// 	CurrentPreviewSceneName string `json:"currentPreviewSceneName,omitempty"`
		// 	CurrentProgramSceneName string `json:"currentProgramSceneName,omitempty"`
		// 	Scenes []*typedefs.Scene `json:"scenes,omitempty"`
		// }

		// TODO: Instead of iterating over the obs scene names and then using that
		// to look it up, lets just see if we can iterate over the scenes in the API
		// and then append the parsed version of those to append

		// TODO: Its a way,... right?? ?? hello ?
		//obsSceneNames := []string{}
		//cachedSceneNames := []string{}
		//for _, cachedScene := range sh.Scenes {
		//	cachedSceneNames = append(cachedSceneNames, cachedScene.Name)
		//}

		for _, scene := range apiResponse.Scenes {
			//obsSceneNames = append(obsSceneNames, scene.SceneName)

			//if cachedScene, ok := sh.Scene(scene.SceneName); ok {
			//	cachedScene.Cache()
			//} else {
			newScene, err := sh.ParseScene(scene.SceneName)
			if err != nil {
				panic(err)
			}

			fmt.Printf("new_scene: \n")
			fmt.Printf("  name: %v\n", newScene.Name)
			fmt.Printf("  item_count: %v\n", len(newScene.Items))
			// NOTE: Here we would create a cached scene from the data that
			//       does exist in the OBS state
			//       And left over scenes would need to be purged
			//         Again, just clearing and rebuilding seems easier;
			//         but we would likely lose data since our models are
			//         more complex than the OBS data models
			//}
		}

		return sh, len(sh.Scenes) == len(apiResponse.Scenes)
	} else {
		return sh, false
	}
}
