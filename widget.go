package sui

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//type func() func(o Widget) bool

type Box struct {
	Widget
	surface    *sdl.Surface
	renderer   *sdl.Renderer
	children   []Widget
	parent     Widget
	clearColor Color
	drawColor  Color
	textColor  Color
	pos        Point
	size       Point
	font       *ttf.Font

	OnDraw            func()
	OnEnter           func()
	OnLeave           func()
	OnMouseButtonDown func()
	OnMouseButtonUp   func()
	OnMouseClick      func()
	OnMouseOver       func()
	OnMouseScroll     func()
	OnKeyPress        func()
	OnResize          func()
}

func NewBox(w, h int) *Box {
	widget := &Box{}
	//widget.init(w, h)
	widget.Resize(w, h)
	widget.SetClearColor(Palette.Background)
	widget.SetColor(Palette.Foreground)
	widget.SetTextColor(Palette.Text)
	widget.SetFont(defaultFont)
	return widget
}

func (o *Box) Close() {
	if o.Parent() != nil {
		o.Parent().RemoveChild(o)
	}
	o.surface.Free()
}

func (o *Box) Pos() Point {
	return o.pos
}

func (o *Box) Size() Point {
	return o.size
}

func (o *Box) Parent() Widget {
	return o.parent
}

func (o *Box) Children() []Widget {
	return o.children
}

func (o *Box) SetParent(parent Widget) {
	o.parent = parent
	PostUpdate()
}

func (o *Box) Surface() *sdl.Surface {
	return o.surface
}

func (o *Box) Renderer() *sdl.Renderer {
	return o.renderer
}

func (o *Box) SetClearColor(color Color) {
	o.clearColor = color
	PostUpdate()
}
func (o *Box) ClearColor() Color {
	return o.clearColor
}

func (o *Box) SetColor(color Color) {
	o.drawColor = color
	PostUpdate()
}

func (o *Box) Color() Color {
	return o.drawColor
}

func (o *Box) SetTextColor(color Color) {
	o.textColor = color
	PostUpdate()
}

func (o *Box) TextColor() Color {
	return o.textColor
}

func (o *Box) Font() *ttf.Font {
	return o.font
}

func (o *Box) SetFont(font *ttf.Font) {
	o.font = font
}

func (o *Box) AddChild(child Widget) {
	o.RemoveChild(child)
	o.children = append(o.children, child)
	child.SetParent(o)
	PostUpdate()
}

func (o *Box) RemoveChild(child Widget) {
	for i, c := range o.children {
		if c == child {
			if i == 0 {
				o.children = o.children[1:]
			} else {
				o.children = append(o.children[:i], o.children[i+1:]...)
			}
			return
		}
	}
	PostUpdate()
}

func (o *Box) TranslateAbsToRel(x, y int) (int, int) {
	x = x - o.Pos().X
	y = y - o.Pos().Y
	if o.Parent() == nil {
		return x, y
	}
	return o.Parent().TranslateAbsToRel(x, y)
}

func (o *Box) enter() {
	callback(o.OnEnter)
}

func (o *Box) leave() {
	callback(o.OnLeave)
}

func (o *Box) mouseButtonDown() {
	callback(o.OnMouseButtonDown)
}

func (o *Box) mouseButtonUp() {
	callback(o.OnMouseButtonUp)
}

func (o *Box) mouseClick() {
	callback(o.OnMouseClick)
}

func (o *Box) mouseOver() {
	callback(o.OnMouseOver)
}
func (o *Box) mouseScroll() {
	callback(o.OnMouseScroll)
}

func (o *Box) keyPress() {
	callback(o.OnKeyPress)
}

func (o *Box) Move(x, y int) {
	o.SetPos(Point{x, y})
}

func (o *Box) SetPos(pos Point) {
	o.pos = pos
	PostUpdate()
}

func (o *Box) Resize(w, h int) {
	o.SetSize(Point{w, h})
}

func (o *Box) SetSize(size Point) {
	size.Max(Point{})
	surface, err := sdl.CreateRGBSurface(0, int32(size.X), int32(size.Y), 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
	__(err)
	renderer, err := sdl.CreateSoftwareRenderer(surface)
	__(err)
	//w, h, err = o.renderer.GetRendererOutputSize()
	//____panic(err)

	o.surface.Free()
	o.surface = surface
	o.renderer.Destroy()
	o.renderer = renderer
	o.size = size

	callback(o.OnResize)

	PostUpdate()
}

func (o *Box) Clear() {
	r := o.Renderer()
	r.SetDrawColor(o.clearColor.RGBA())
	r.Clear()
}

func (o *Box) Fill(rect Rect) {
	r := o.Renderer()
	r.SetDrawColor(o.drawColor.RGBA())
	r.FillRect(rect.Rect())
}

func (o *Box) Rect(rect Rect) {
	r := o.Renderer()
	r.SetDrawColor(o.drawColor.RGBA())
	r.DrawRect(rect.Rect())
}

func (o *Box) Line(a, b Point) {
	r := o.Renderer()
	r.SetDrawColor(o.drawColor.RGBA())
	r.DrawLine(a.X, a.Y, b.X, b.Y)
}

func (o *Box) WriteText(pos Point, str string) Point {
	var solid *sdl.Surface
	var err error

	if str == "" {
		return Point{0, 0}
	}

	if solid, err = o.Font().RenderUTF8_Blended(str, o.textColor.Color); err != nil {
		fmt.Fprint(os.Stderr, "Failed to render text: %s\n", err)
		return Point{0, 0}
	}
	defer solid.Free()
	src := sdl.Rect{0, 0, solid.W, solid.H}
	dst := NewRect(pos, o.Size())
	if err = solid.Blit(&src, o.Surface(), dst.Rect()); err != nil {
		fmt.Fprint(os.Stderr, "Failed to put text on window surface: %s\n", err)
		return Point{}
	}

	return Point{int(solid.W), int(solid.H)}
}

func (o *Box) Draw() {
	if !callback(o.OnDraw) {
		o.Clear()
	}
}

func (o *Box) Repaint() {
	glob.sender = o
	o.Draw()
	glob.sender = nil
	for _, child := range o.children {
		child.Repaint()
		src := NewRect(Point{}, child.Size())
		dst := NewRect(child.Pos(), child.Size())
		child.Surface().Blit(src.Rect(), o.Surface(), dst.Rect())
	}
}
