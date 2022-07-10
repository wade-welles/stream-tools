package terminal

import (
	"golang.org/x/term"
)

// TODO: Model the terminal so we can interact with it while saving a cached
// state.
//type Terminal struct {
//	Height int
//	Width int
//}

func Size() (int, int, err) {
	if term.IsTerminal(0) {
		println("in a term")
	} else {
		println("not in a term")
	}
	return term.GetSize(0)
}

func Height() int {
	_, height, err := Size()
	if err != nil {
		return height
	} else {
		return 0
	}
}

func Width() int {
	width, _, err := Size()
	if err != nil {
		return width
	} else {
		return 0
	}
}
