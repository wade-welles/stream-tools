package scene

import (
	obs "github.com/wade-welles/streamkit/obs"
)

type Scene struct {
	Id   string
	Name string

	Items *Items

	showId string
}

func (sc Scene) Show() *obs.Show {
	// Iterate over showIds stored in some slice
	// and find the show and output that object after
	return &obs.Show{}
}
func (sc Scene) Item(itemId string) *obs.Show {
	return &obs.Show{}
}
