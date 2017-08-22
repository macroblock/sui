package sui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Widget interface {
	Close()
	Parent() Widget
	SetParent(Widget)
	Children() []Widget
	AddChild(child Widget)
	RemoveChild(child Widget)
	Pos() Point
	SetPos(pos Point)
	Move(x, y int)
	Size() Point
	SetSize(size Point)
	Resize(width, height int)
	Font() *ttf.Font
	Surface() *sdl.Surface
	Renderer() *sdl.Renderer
	Repaint()
	SetClearColor(color Color)
	SetColor(color Color)
	SetTextColor(color Color)
	Clear()
	Fill(rect Rect)
	Rect(rect Rect)
	Line(a, b Point)
	WriteText(pos Point, str string) Point
	//IsInside(x, y int32) bool
	TranslateAbsToRel(x, y int) (int, int)
	Draw()
	enter()
	leave()
	mouseButtonDown()
	mouseButtonUp()
	mouseClick()
	//MouseDoubleClick(x, y int32)
	mouseOver()
	mouseScroll()
	//KeyDown(key sdl.Keycode, mod uint16)
	//KeyUp(key sdl.Keycode, mod uint16)
	//DragMove(x, y int32, payload DragPayload)
	//DragEnter(x, y int32, payload DragPayload)
	//DragLeave(payload DragPayload)
	//DragDrop(x, y int32, payload DragPayload) bool
}
