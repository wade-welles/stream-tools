package obs

import (
	"fmt"

	goobs "github.com/andreykaipov/goobs"

	events "github.com/andreykaipov/goobs/api/events"
	ui "github.com/andreykaipov/goobs/api/requests/ui"
	typedefs "github.com/andreykaipov/goobs/api/typedefs"
)

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
//	AverageRenderTime uint
//	FramesLost         uint
//	FramesSkipped      uint
//	FramesDropped      uint
//	DataOutput         uint
//	Bitrate            uint
//}

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

// TODO: I don't see this in the code anymore-- its because this is
// extrapolating a lot so we have a much easier to use API
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

// TODO: So its now clear that ParseShow runs but ParseScene is never ran
func MarshalAlignment(alignment int) Alignment { return Alignment(alignment) }

// TODO: This should either give the full OBS object returned (like an init
// function) or this should be a method on OBS after it is created you connect.
func ConnectToOBS(host string) *goobs.Client {
	client, err := goobs.New(host)
	// TODO: If we fail to connect, we could have streamkit launch obs
	if err != nil {
		panic(err)
	}
	return client
}

// TODO: Item type ideally should have group or folder in here but these are
//
//	the ones from OBS; so itd probably need to be separate from our
//	custom ones and then added together with a further type
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

// NOTE: The API we are creating looks something like:
//          obs.Show.Scenes.First()
//          obs.Show.Scenes.First().Items.First().(Show()|.Hide())

type Client struct {
	WS *goobs.Client

	// TODO: We want OBS to reflect whats running in the application and Show is
	// our local cache of it
	Mode Mode

	//Stats      *Stats
	//AudioMixer *AudioMixer

	//Sources []*goobs.Source
}

func (obs Client) IsMode(mode Mode) bool {
	switch mode {
	case StudioMode, StreamingMode, RecordingMode:
		return true
	default:
		return false
	}
}

func (obs Client) ToggleStudioMode() bool {
	toggledValue := !obs.IsMode(StudioMode)
	studioModeEnabledParams := ui.SetStudioModeEnabledParams{
		StudioModeEnabled: &toggledValue,
	}

	// NOTE: I hate using Ui,.. its an acronym :(
	_, err := obs.WS.Ui.SetStudioModeEnabled(&studioModeEnabledParams)
	return err == nil
}

func (obs Client) StudioMode() bool {
	obs.Mode = StudioMode
	// NOTE: This using pointer to a boolean is incredibly tedious
	studioMode := true
	studioModeEnabledParams := ui.SetStudioModeEnabledParams{
		StudioModeEnabled: &studioMode,
	}

	_, err := obs.WS.Ui.SetStudioModeEnabled(&studioModeEnabledParams)
	return err == nil
}

// go Events() to call this because its meant to be event driven-- but you know
func (obs Client) Events() {
	for event := range obs.WS.IncomingEvents {
		switch e := event.(type) {
		case *events.SceneItemEnableStateChanged:
			fmt.Printf("Scene Item Enabled %-25q (%v): %v\n", e.SceneName, e.SceneItemId, e.SceneItemEnabled)
		default:
			fmt.Printf("Unhandled event: %#v\n", e)
		}
	}
}

// TODO: This returns typedefs, but we intend to abstract all those away so we
// never should be returning them-- at least not as a public func
func (obs Client) Scenes() ([]*typedefs.Scene, error) {
	apiResponse, err := obs.WS.Scenes.GetSceneList()
	return apiResponse.Scenes, err
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
