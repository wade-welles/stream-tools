package obstools

import (
	"fmt"
	"strings"
	"time"

	x11 "github.com/linuxdeepin/go-x11-client"
	ewmh "github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
)

type Window struct {
	Name    WindowName
	Focused bool // aka Active

}

type Windows []*Window

type WindowName uint8 // 0..255

const (
	UndefinedName WindowName = iota
	Primary
	Secondary
	Chromium
)

func (wt WindowName) String() string {
	switch wt {
	case Primary:
		return "primary"
	case Secondary:
		return "secondary"
	case Chromium:
		return "chromium"
	default: // UndefinedName
		return "undefined"
	}
}

func MarshalWindowName(wt string) WindowName {
	switch strings.ToLower(wt) {
	case Primary.String():
		return Primary
	case Secondary.String():
		return Secondary
	case Chromium.String():
		return Chromium
	default: // UndefinedName
		return UndefinedName
	}
}

////////////////////////////////////////////////////////////////////////////////
type X11 struct {
	Client *x11.Conn

	CurrentWindow WindowName

	ActiveWindowChangedAt time.Time
}

func ConnectToX11() *x11.Conn {
	client, err := x11.NewConn()
	if err != nil {
		panic(err)
	}
	return client
}

// TODO: X11 =(should return a Window() function=> WindowName type
func (x *X11) Window() WindowName {
	x.CacheActiveWindow()
	return x.CurrentWindow
}

func (x *X11) HasActiveWindowChanged() bool {
	// NOTE
	// If we record time last active window change happened, then we can limit
	// the number of times it can change within x amount of time for better
	// reliability under pressure.
	return x.IsCurrentWindow(x.ActiveWindow())
}

func (x *X11) ActiveWindow() WindowName {
	windowName, err := ewmh.GetActiveWindow(x.Client).Reply(x.Client)
	if err != nil {
		fmt.Println("error:", err)
	}

	activeWindowName, err := ewmh.GetWMName(x.Client, windowName).Reply(x.Client)
	if err != nil {
		fmt.Println("error:", err)
	}

	if strings.HasSuffix(strings.ToLower(activeWindowName), Chromium.String()) {
		activeWindowName = strings.ToLower(activeWindowName[(len(activeWindowName) - len(Chromium.String())):])
	}

	return MarshalWindowName(activeWindowName)
}

func (x *X11) CacheActiveWindow() WindowName {
	x.CurrentWindow = x.ActiveWindow()
	x.ActiveWindowChangedAt = time.Now()
	return x.CurrentWindow
}

// TODO: Naming issue here bigger than may originally appear
// x.IsActiveWindow(Primary)
func (x *X11) IsActiveWindow(windowName WindowName) bool {
	return x.ActiveWindow() != windowName
}

func (x *X11) IsCurrentWindow(windowName WindowName) bool {
	return x.CurrentWindow != windowName
}
