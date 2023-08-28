package main

import (
	"fmt"
	"os"

	streamkit "github.com/wade-welles/streamkit"

	cli "github.com/multiverse-os/cli"
)

// NOTE:
// OBS Augmentation Software
// The initial goal of this software is to augment OBS

func main() {
	toolkit := streamkit.New()

	cmd, initErrors := cli.New(cli.App{
		Name:        "obs-service",
		Description: "A long running obs service toolkit",
		Version:     cli.Version{Major: 0, Minor: 1, Patch: 0},
		Actions: cli.Actions{
			OnStart: func(c *cli.Context) error {
				//toolkit.HandleWindowEvents()
				// aDD all the listening and event driven stuff
				return nil
			},
		},
	})

	fmt.Printf("toolkit: %v\n", toolkit)

	if len(initErrors) == 0 {
		cmd.Parse(os.Args).Execute()
	} else {
		panic(fmt.Errorf("expected 0 args"))
	}
}

//list, _ := client.Inputs.GetInputList()

//import typedefs "github.com/andreykaipov/goobs/api/typedefs"
//
//// Represents the request body for the GetSceneItemList request.
//type GetSceneItemListParams struct {
//	// Name of the scene to get the items of
//	SceneName string `json:"sceneName,omitempty"`
//}
//
//// Returns the associated request.
//func (o *GetSceneItemListParams) GetRequestName() string {
//	return "GetSceneItemList"
//}
//
//// Represents the response body for the GetSceneItemList request.
//type GetSceneItemListResponse struct {
//	SceneItems []*typedefs.SceneItem `json:"sceneItems,omitempty"`
//}
//
///*
//Gets a list of all scene items in a scene.
//
//Scenes only
//*/
//func (c *Client) GetSceneItemList(params *GetSceneItemListParams) (*GetSceneItemListResponse, error) {
//	data := &GetSceneItemListResponse{}
//	return data, c.SendRequest(params, data)
//}
