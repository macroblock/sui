package sui

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var glob struct {
	sysWindows []SystemWindow
}

// SystemWindow ...
type SystemWindow struct {
	title         string
	width, height int32
	window        *sdl.Window
	surface       *sdl.Surface
}

// NewSystemWindow ...
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
	InitFonts()

	return nil
}

// Close ...
func Close() {
	ttf.Quit()
	sdl.Quit()
}

// Run ...
func Run() int {

	quit := false
	for event := sdl.PollEvent(); event != nil || !quit; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			quit = true
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

		time.Sleep(1)
	}

	fmt.Printf("done.")
	return 0
}
