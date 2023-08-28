package x11

import (
	"time"

	x11 "github.com/linuxdeepin/go-x11-client"
)

type X11 struct {
	Client *x11.Conn // 	xdisplay       *x.Conn

	// TODO: Populated by previous active windows
	//          (THIS REQUIRES UNIQUE-NESS CHECK)
	//Desktops
	//Windows Windows

	// TODO
	// When needed bother to store the history of active windows but that
	// isn't needed quite yet, so there is about ZERO point in implementing
	// it.

	// TODO: Maybe just cache the active window name so we do simple name
	// comparison, but this leads to a bug where two windows with the same name
	// are considered the name window
	ActiveWindowTitle     string
	ActiveWindowChangedAt time.Time
}

// TODO: If we can move some of these to be methods of Window struct, it would
// be better organized but there will be obvious limitations we have to work
// through
//func (x *x11) CurrentWindowTitle() string {
//	return x.ActiveWindow().Title
//}
//
//func (x *x11) HasActiveWindowChanged() bool {
//	return x.ActiveWindowTitle != x.ActiveWindow().Title
//}
//
//func (x *x11) ActiveWindow() *Window {
//	//	var err error
//	//	// TODO: This returns an x.Window object which can get all sorts of
//	//	// information beyond just the name, like the PID. We shouldn't need a second
//	//	// call at all to get the title of the window, thats absurdist.
//	activeWindow, err := ewmh.GetActiveWindow(x.Client).Reply(x.Client)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("active_window: %v\n", activeWindow)
//
//	// TODO: Do we actually need to do GetWMName? Shouldn't we actually do the
//	// GetWindowInfo thing so we get it and much more information we could cache
//	activeWindowTitle, err := ewmh.GetWMName(
//		x.Client,
//		activeWindow,
//	).Reply(x.Client)
//	if err != nil {
//		fmt.Printf("error(%v)\n", err)
//	}
//
//	pid, err := ewmh.GetWMPid(x.Client, activeWindow).Reply(x.Client)
//	if err != nil {
//		fmt.Printf("error(%v)\n", err)
//	} else {
//		fmt.Printf("\tPid:%v\n", pid)
//		data, _ := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
//		fmt.Printf("\t\tCmdline: %s\n", data)
//	}
//
//	// TODO: Maybe have a cache window data or some such func
//	cachedWindow := &Window{
//		Title: activeWindowTitle,
//		PID:   pid,
//	}
//
//	return cachedWindow
//}
//
//func (x *x11) InitActiveWindow() *Window {
//	activeWindow := x.ActiveWindow()
//	x.ActiveWindowTitle = activeWindow.Title
//	x.ActiveWindowChangedAt = time.Now()
//	return activeWindow
//}
//
//func (x *x11) CacheActiveWindow() *Window {
//	activeWindow := x.ActiveWindow()
//	x.ActiveWindowTitle = x.ActiveWindow().Title
//	x.ActiveWindowChangedAt = time.Now()
//	return activeWindow
//}
