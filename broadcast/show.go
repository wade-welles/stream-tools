package broadcast

import (
	show "github.com/wade-welles/streamkit/broadcast/show"
)

// TODO: We will need to restructure the goobs from our show object since it
// could just as easily be ffmpeg of output of window that is assembled from
// v42l and whatever else

// Doing this disentanglement has a lot of cascading benefits that may almost
// certainly missed if not thought long abuot it

type Show struct {
	Id   int
	Name string

	// Season []*Season
	// Episodes []*Episode

	ActiveScene *show.Scene
	StudioScene *show.Scene

	Scenes []*show.Scene
}

func (sh *Show) EmptyScene() *show.Scene {
	return &show.Scene{Name: ""}
}

func (sh *Show) Scene(name string) *show.Scene {
	for _, scene := range sh.Scenes {
		if scene.Name == name {
			return scene
		}
	}
	return nil
}

func (sh *Show) ParseScene(name string, index int) *show.Scene {
	// Validate name & index
	var err error
	if !(0 < len(name) && len(name) < 255) &&
		!(0 <= index && index < 999) {
		panic(err)
	}

	validScene := &show.Scene{
		Index: index,
		Name:  name,
	}

	sh.Scenes = append(sh.Scenes, validScene)

	return validScene
}

// GOOBS TYPEDEF
//
//	type Scene struct {
//		SceneIndex int    `json:"sceneIndex"`
//		SceneName  string `json:"sceneName"`
//	}
