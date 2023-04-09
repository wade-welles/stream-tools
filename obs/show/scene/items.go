package scene

import (
	"fmt"

	obs "github.com/wade-welles/streamkit/obs"
)

type Item struct {
	Id   string
	Name string

	sceneId string
	showId  string
}

type Items []*Item

func (i Item) TestFunction() *Item {
	fmt.Printf("streamkit/show/scenes\n")
	return i
}

func (it Item) Show() *obs.Show {
	return &obs.Show{}
}
