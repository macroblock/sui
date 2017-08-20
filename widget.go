package sui

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type CallbackFn func(o Widgeter) bool

type Widget struct {
	Widgeter
	surface  *sdl.Surface
	renderer *sdl.Renderer
	children []Widgeter
	parent   Widgeter
	bgColor  uint32
	pos      Point
	size     Point
	font     *ttf.Font
	onDraw   CallbackFn
}

func NewWidget(w, h int) *Widget {
	widget := &Widget{}
	//widget.init(w, h)
	widget.Resize(w, h)
	widget.SetColor(0xff00dddd)
	widget.SetFont(defaultFont)
	return widget
}

func (o *Widget) Close() {
	if o.Parent() != nil {
		o.Parent().RemoveChild(o)
	}
	o.surface.Free()
}

func (o *Widget) Pos() Point {
	return o.pos
}

func (o *Widget) Size() Point {
	return o.size
}

func (o *Widget) SetColor(color uint32) {
	o.bgColor = color
	PostUpdate()
}

func (o *Widget) Parent() Widgeter {
	return o.parent
}

func (o *Widget) Children() []Widgeter {
	return o.children
}

func (o *Widget) SetParent(parent Widgeter) {
	o.parent = parent
	PostUpdate()
}

func (o *Widget) Surface() *sdl.Surface {
	return o.surface
}

func (o *Widget) Renderer() *sdl.Renderer {
	return o.renderer
}

func (o *Widget) Font() *ttf.Font {
	return o.font
}

func (o *Widget) SetFont(font *ttf.Font) {
	o.font = font
}

func (o *Widget) SetOnDraw(fn CallbackFn) {
	o.onDraw = fn
}

func (o *Widget) AddChild(child Widgeter) {
	o.RemoveChild(child)
	o.children = append(o.children, child)
	child.SetParent(o)
	PostUpdate()
}

func (o *Widget) RemoveChild(child Widgeter) {
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

func (o *Widget) Move(x, y int) {
	o.SetPos(Point{x, y})
}

func (o *Widget) SetPos(pos Point) {
	o.pos = pos
	PostUpdate()
}

func (o *Widget) Resize(w, h int) {
	o.SetSize(Point{w, h})
}

func (o *Widget) SetSize(size Point) {
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
	PostUpdate()
}

func (o *Widget) FillRect(x, y, w, h int, c uint32) {
	surface := o.Surface()
	rect := NewRect(Point{x, y}, Point{w, h})
	surface.FillRect(&rect, c)
}

func (o *Widget) WriteText(pos Point, str string, color uint32) Point {
	var solid *sdl.Surface
	var err error

	if str == "" {
		return Point{0, 0}
	}

	if solid, err = o.Font().RenderUTF8_Blended(str, uint32toColor(color)); err != nil {
		fmt.Fprint(os.Stderr, "Failed to render text: %s\n", err)
		return Point{0, 0}
	}
	defer solid.Free()
	rectSrc := sdl.Rect{0, 0, solid.W, solid.H}
	rectDst := NewRect(pos, o.Size())
	if err = solid.Blit(&rectSrc, o.Surface(), &rectDst); err != nil {
		//fmt.Fprint(os.Stderr, "Failed to put text on window surface: %s\n", err)
		return Point{}
	}

	return Point{int(solid.W), int(solid.H)}
}

func (o *Widget) Draw() bool {
	return callback(o.onDraw, o)
}

func (o *Widget) Repaint() {
	if !o.Draw() {
		o.FillRect(0, 0, o.size.X, o.size.Y, o.bgColor)
	}
	for _, child := range o.children {
		child.Repaint()
		src := NewRect(Point{}, child.Size())
		dst := NewRect(child.Pos(), child.Size())
		child.Surface().Blit(&src, o.Surface(), &dst)
	}
}
