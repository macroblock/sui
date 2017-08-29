package sui

import (
	"github.com/veandco/go-sdl2/sdl"
)

// RootWindow ...
type RootWindow struct {
	Box
	window     *sdl.Window
	OnDropFile func()
}

// NewRootWindow ...
func NewRootWindow(title string, width, height int) *RootWindow {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	__(err)

	root := &RootWindow{
		Box:    Box{},
		window: window,
	}
	w, h := window.GetSize()
	root.Resize(w, h)
	root.SetClearColor(Palette.Background)
	root.SetColor(Palette.Foreground)
	root.SetTextColor(Palette.Text)
	root.SetFont(defaultFont)

	sdl.EventState(sdl.DROPFILE, sdl.ENABLE)

	glob.rootWindows = append(glob.rootWindows, root)

	return root
}

func (o *RootWindow) UpdateSurface() {
	o.window.UpdateSurface()
}

func (o *RootWindow) Resize(w, h int) {
	o.SetSize(NewPoint(w, h))
}

func (o *RootWindow) SetSize(_ Point) {
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
	//fmt.Printf("getRendererOutputSize: id [%d], size: %dx%d\n", o.window.GetID(), sizew, sizeh)

	o.size = Point{sizew, sizeh}

	callback(o.OnResize)

	PostUpdate()
}

func (o *RootWindow) Close() {
	o.Box.Close()
	//o.renderer.Destroy()
	//o.surface.Free()
	o.window.Destroy()
}

func (o *RootWindow) dropFile() {
	callback(o.OnDropFile)
}
