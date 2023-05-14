package show

import (
	scene "github.com/wade-welles/streamkit/broadcast/show/scene"
)

type Scene struct {
	Id   string
	Name string

	Items []*scene.Item
}

func (sc *Scene) HasName(name string) bool {
	return sc != nil && sc.Name == name
}
