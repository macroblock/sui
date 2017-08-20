package sui

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var glob struct {
	sysWindows []*SystemWindow
	needUpdate bool
}

// SystemWindow ...
type SystemWindow struct {
	Widget
	window *sdl.Window
}

// NewSystemWindow ...
func NewSystemWindow(title string, width, height int) *SystemWindow {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	__(err)

	sysWindow := &SystemWindow{
		Widget: Widget{},
		window: window,
	}
	w, h := window.GetSize()
	sysWindow.Resize(w, h)
	sysWindow.SetColor(0xff00dddd)
	sysWindow.SetFont(defaultFont)

	glob.sysWindows = append(glob.sysWindows, sysWindow)

	return sysWindow
}

func (o *SystemWindow) UpdateSurface() {
	o.window.UpdateSurface()
}

func (o *SystemWindow) Resize(w, h int) {
	//fmt.Printf("sys resize: id [%d], size: %dx%d\n", o.window.GetID(), w, h)
	sizew, sizeh := o.window.GetSize()
	//fmt.Printf("getSize: id [%d], size: %dx%d\n", o.window.GetID(), sizew, sizeh)
	err := error(nil)
	o.surface.Free()
	o.surface, err = o.window.GetSurface()
	__(err)
	o.renderer.Destroy()
	o.renderer, err = sdl.CreateSoftwareRenderer(o.surface)
	__(err)
	sizew, sizeh, err = o.renderer.GetRendererOutputSize()
	__(err)
	fmt.Printf("getRendererOutputSize: id [%d], size: %dx%d\n", o.window.GetID(), sizew, sizeh)

	o.size = Point{sizew, sizeh}
	PostUpdate()
}

func (o *SystemWindow) Close() {
	o.Widget.Close()
	o.renderer.Destroy()
	o.surface.Free()
}

// Init ...
func Init() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	__(err)

	err = ttf.Init()
	__(err)

	InitFonts()

	return nil
}

// Close ...
func Close() {
	//for _, root := range glob.sysWindows {
	//	root.Close()
	//}
	glob.sysWindows = nil

	ttf.Quit()
	sdl.Quit()
}

func PostUpdate() {
	glob.needUpdate = true
}

// Run ...
func Run() int {
	quit := false
	PostUpdate()
	for event := sdl.PollEvent(); event != nil || !quit; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			quit = true
		case *sdl.WindowEvent:
			switch t.Event {
			case sdl.WINDOWEVENT_RESIZED:
				fmt.Printf("WINDOWEVENT RESIZED: id [%d], size: %dx%d\n", t.WindowID, t.Data1, t.Data2)
				for _, root := range glob.sysWindows {
					root.Resize(int(t.Data1), int(t.Data2))
				}
				PostUpdate()

				//case sdl.WINDOWEVENT_SIZE_CHANGED:
				//	fmt.Printf("size changed: id [%d], size: %dx%d\n", t.WindowID, t.Data1, t.Data2)
			}
		case *sdl.KeyDownEvent:
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			if t.Keysym.Sym == sdl.K_ESCAPE {
				quit = true
			}
		case *sdl.KeyUpEvent:
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		}

		if quit {
			break
		}
		//fmt.Println("->")
		if glob.needUpdate {
			for _, root := range glob.sysWindows {
				//fmt.Println(root.W(), root.H())
				root.Repaint()
				root.UpdateSurface()
			}
			glob.needUpdate = false
		}
		//fmt.Println("))")

		time.Sleep(1)
	}

	fmt.Printf("done.")
	return 0
}
