package sui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Widgeter interface {
	AddChild(child Widgeter)
	RemoveChild(child Widgeter)
	Move(x, y int)
	Resize(width, height int)
	SetPos(pos Point)
	SetSize(size Point)
	Surface() *sdl.Surface
	Renderer() *sdl.Renderer
	Children() []Widgeter
	Parent() Widgeter
	SetParent(Widgeter)
	Repaint()
	Pos() Point
	Size() Point
	//MouseDoubleClick(x, y int32)
	//MousePressDown(x, y int32, button uint8)
	//MousePressUp(x, y int32, button uint8)
	//MouseMove(x, y, xrel, yrel int32)
	//KeyDown(key sdl.Keycode, mod uint16)
	//KeyUp(key sdl.Keycode, mod uint16)
	//TranslateXYToWidget(globalX, globalY int32) (x, y int32)
	//IsInside(x, y int32) bool
	//HasFocus(focus bool)
	Font() *ttf.Font
	WriteText(pos Point, str string, color uint32) Point
	Close()
	//DragMove(x, y int32, payload DragPayload)
	//DragEnter(x, y int32, payload DragPayload)
	//DragLeave(payload DragPayload)
	//DragDrop(x, y int32, payload DragPayload) bool
}
