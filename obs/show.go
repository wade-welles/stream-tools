package obs

import (
	"fmt"

	goobs "github.com/andreykaipov/goobs"
	sceneitems "github.com/andreykaipov/goobs/api/requests/sceneitems"
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
func (sh *Show) Cache() bool {

	// NOTE: For simplicity, for now we will just set scenes to empty and then
	// populate with API data. So we set show.Scenes to an empty slice of scenes
	sh.Scenes = Scenes{}

	apiScenesResponse, err := sh.OBS.Scenes.GetSceneList()
	if err != nil {
		fmt.Printf("error(%v)\n", err)
	}
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
	fmt.Printf("show:\n")
	fmt.Printf("  scenes:\n")
	//
	for _, scene := range apiScenesResponse.Scenes {
		apiResponse, err := sh.OBS.SceneItems.GetSceneItemList(
			&sceneitems.GetSceneItemListParams{
				SceneName: scene.SceneName,
			})
		if err != nil {
			fmt.Printf("error(%v)\n", err)
		}

		fmt.Printf("      scene:\n")
		fmt.Printf("        object: %v\n", scene)
		//fmt.Printf("        item_count: %v\n", len(apiResponse.SceneItems))
		fmt.Printf("        items:\n")
		fmt.Printf("          object: %v\n", apiResponse)
		for _, sceneItem := range apiResponse.SceneItems {
			fmt.Printf("          item: \n")
			fmt.Printf("            object: %v\n", sceneItem)
		}
	}
	//return sh, len(sh.Scenes) == len(apiResponse.Scenes)

	return false
}

// what goes in/ what goes out?
// in show we keep our data object without any logic related to interacting with
// goobs separate, then its easier to swap it out

func (sh *Show) PrintDebug() {
	fmt.Printf("show: \n")
	fmt.Printf("  object: %v\n", sh)
	// Breaks here because no scenes, its not even set to empty
	fmt.Printf("  scenes: %v\n", sh.Scenes)
	fmt.Printf("  scene_count: %v\n", len(sh.Scenes))
	for _, scene := range sh.Scenes {
		fmt.Printf("    scene:\n")
		fmt.Printf("      name: %v\n", scene.Name)
		fmt.Printf("      item_count: %v\n", len(scene.Items))
		fmt.Printf("      items:\n")
		for _, item := range scene.Items {
			fmt.Printf("      item:\n")
			fmt.Printf("        name: %v\n", item.Name)
		}
	}

}
