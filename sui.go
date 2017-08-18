package sui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var glob struct {
	sysWindows []SystemWindow
}

type SystemWindow struct {
	title         string
	width, height int32
	window        *sdl.Window
	surface       *sdl.Surface
}

func NewSystemWindow(title string, width, height int) (*SystemWindow, error) {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return nil, err
	}

	surface, err := window.GetSurface()

	sysWindow := SystemWindow{
		title:   title,
		width:   int32(width),
		height:  int32(height),
		window:  window,
		surface: surface,
	}
	glob.sysWindows = append(glob.sysWindows, sysWindow)

	return &sysWindow, err
}

// Init ...
func Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	if err := ttf.Init(); err != nil {
		return err
	}
	return nil
}

// Close ...
func Close() {
	ttf.Quit()
	sdl.Quit()
}
