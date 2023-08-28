package show

import (
	scene "github.com/wade-welles/streamkit/broadcast/show/scene"
)

type Scene struct {
	Index int
	Name  string

	Items []*scene.Item
}

func (sc *Scene) HasName(name string) bool {
	return (sc != nil || len(sc.Name) != len(name) || len(name) == 0)
}

func (sc *Scene) Item(name string) *scene.Item {
	return &scene.Item{}
}

func (sc *Scene) ParseItem(name string, index int) *scene.Item {

	return &scene.Item{}
}
