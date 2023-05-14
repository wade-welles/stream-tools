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
	ActiveScene *show.Scene
	StudioScene *show.Scene

	Scenes []*show.Scene
}

func (sh *Show) EmptyScene() *show.Scene {
	return &show.Scene{
		Name: "",
	}
}
