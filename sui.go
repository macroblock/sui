package sui

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

import "C"

var OnLoop func()

var glob struct {
	rootWindows        []*RootWindow
	needUpdate         bool
	sender             Widget
	focus              Widget
	prevFocus          Widget
	mouseOver          Widget
	prevMouseOver      Widget
	x, y               int
	mouseScroll        Point
	mouseButtonPressed bool
	dropFile           string
	modShift           int
	mouseButton        int
	keysym             int
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
	glob.rootWindows = nil

	ttf.Quit()
	sdl.Quit()
}

func findWidget(x, y int, root Widget) (target Widget, xTarget, yTarget int) {
	target = nil
	x -= root.Pos().X
	y -= root.Pos().Y

	for _, child := range root.Children() {
		xMin := child.Pos().X
		yMin := child.Pos().Y
		xMax := child.Pos().X + child.Size().X
		yMax := child.Pos().Y + child.Size().Y
		xMin = MaxInt(xMin, 0)
		yMin = MaxInt(yMin, 0)
		xMax = MinInt(xMax, root.Size().X)
		yMax = MinInt(yMax, root.Size().Y)
		if x >= xMin && x < xMax && y >= yMin && y < yMax {
			target = child
		}
	}

	if target != nil {
		return findWidget(x, y, target)
	}
	if x >= 0 && x < root.Size().X && y >= 0 && y < root.Size().Y {
		return root, x, y
	}
	return nil, -1, -1
}

func PostUpdate() {
	glob.needUpdate = true
}

func MousePos() Point {
	return NewPoint(glob.x, glob.y)
}

func SetSender(o Widget) {
	glob.sender = o
}

func Sender() Widget {
	return glob.sender
}

func MouseOver() Widget {
	return glob.mouseOver
}

func PrevMouseOver() Widget {
	return glob.prevMouseOver
}

func MouseScroll() Point {
	return glob.mouseScroll
}

func ModShift() int {
	return glob.modShift
}

func MouseButton() int {
	return glob.mouseButton
}

func KeySym() int {
	return glob.keysym
}

func DropFile() string {
	return glob.dropFile
}

// Run ...
func Run() int {
	quit := false
	PostUpdate()
	//for event := sdl.PollEvent(); event != nil || !quit; event = sdl.PollEvent() {
	for {
		event := sdl.PollEvent()
		/*if event == nil {
			time.Sleep(1)
			continue
		}*/

		switch t := event.(type) {
		case *sdl.QuitEvent:
			quit = true

		case *sdl.MouseButtonEvent:
			if t.Type == sdl.MOUSEBUTTONDOWN {
				glob.mouseButtonPressed = true
				glob.focus, glob.x, glob.y = findWidget(int(t.X), int(t.Y), glob.rootWindows[0])
				if glob.prevFocus != glob.focus {
					if glob.prevFocus != nil {
						glob.sender = glob.prevFocus
						glob.prevFocus.leave()
						fmt.Println("leave")
					}
					if glob.focus != nil {
						glob.sender = glob.focus
						glob.focus.enter()
						fmt.Println("enter: ", glob.focus)
					}
					glob.prevFocus = glob.focus
				}
				if glob.focus != nil {
					glob.sender = glob.focus
					glob.focus.mouseButtonDown()
				}
			}
			if t.Type == sdl.MOUSEBUTTONUP {
				glob.mouseButtonPressed = false
				glob.sender, glob.x, glob.y = findWidget(int(t.X), int(t.Y), glob.rootWindows[0])
				if glob.sender != nil {
					glob.sender.mouseButtonUp()
					if glob.focus == glob.sender {
						glob.sender = glob.focus
						glob.mouseButton = int(t.Button)
						glob.focus.mouseClick()
					}
				}
			}
			glob.sender = nil

		case *sdl.MouseMotionEvent:
			if t.Type == sdl.MOUSEMOTION {
				if glob.mouseButtonPressed == false {
					glob.mouseOver, glob.x, glob.y = findWidget(int(t.X), int(t.Y), glob.rootWindows[0])

					if glob.mouseOver != glob.prevMouseOver {
						if glob.mouseOver != nil {
							glob.sender = glob.mouseOver
							glob.mouseOver.mouseOver()
						}
						glob.prevMouseOver = glob.mouseOver
					}
				}
			}
			glob.sender = nil

		case *sdl.MouseWheelEvent:
			if t.Type == sdl.MOUSEWHEEL {
				if glob.mouseOver != nil {
					glob.sender = glob.mouseOver
					glob.mouseScroll.X = int(t.X)
					glob.mouseScroll.Y = int(t.Y)
					glob.mouseOver.mouseScroll()
				}
			}
			glob.sender = nil

		case *sdl.WindowEvent:
			switch t.Event {
			case sdl.WINDOWEVENT_RESIZED:
				fmt.Printf("WINDOWEVENT RESIZED: id [%d], size: %dx%d\n", t.WindowID, t.Data1, t.Data2)
				for _, root := range glob.rootWindows {
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
			if t.Keysym.Sym == sdl.K_LSHIFT || t.Keysym.Sym == sdl.K_RSHIFT {
				glob.modShift = 1
			}
			//fmt.Println("Down")
			if t.Repeat == 0 {
				glob.sender = glob.rootWindows[0]
				glob.keysym = int(t.Keysym.Sym)
				glob.sender.keyPress()
				glob.sender = nil
			}

		case *sdl.KeyUpEvent:
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%x\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			if t.Keysym.Sym == sdl.K_LSHIFT || t.Keysym.Sym == sdl.K_RSHIFT {
				glob.modShift = 0
			}

		case *sdl.DropEvent:
			glob.dropFile = C.GoString((*C.char)(t.File))
			fmt.Println(glob.dropFile)
			glob.rootWindows[0].dropFile()
			glob.dropFile = ""
			PostUpdate()
		}

		if quit {
			break
		}

		callback(OnLoop)

		if glob.needUpdate {
			for _, root := range glob.rootWindows {
				root.Repaint()
				root.UpdateSurface()
			}
			glob.needUpdate = false
		}

		time.Sleep(1)
	}

	fmt.Printf("done.")
	return 0
}
