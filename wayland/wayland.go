package wayland

import (
	"fmt"

	wl "github.com/rajveermalviya/go-wayland"
)

func WaylandTest() {
	fmt.Println("wayland")
}

func Connect(addr string) (*wl.Display, error) {
	return wl.Connect(addr)
}
