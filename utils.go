package sui

import "github.com/veandco/go-sdl2/sdl"

type Point struct {
	X, Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

func (o *Point) Min(point Point) {
	if o.X > point.X {
		o.X = point.X
	}
	if o.Y > point.Y {
		o.Y = point.Y
	}
}

func (o *Point) Max(point Point) {
	if o.X < point.X {
		o.X = point.X
	}
	if o.Y < point.Y {
		o.Y = point.Y
	}
}

type Color struct {
	sdl.Color
}

func Color32(c uint32) Color {
	return Color{
		sdl.Color{
			uint8((c >> 0) & 0xff),
			uint8((c >> 8) & 0xff),
			uint8((c >> 16) & 0xff),
			uint8((c >> 24) & 0xff),
		},
	}
}

func (o Color) RGBA() (byte, byte, byte, byte) {
	return o.R, o.G, o.B, o.A
}

type Rect struct {
	Pos  Point
	Size Point
}

func NewRect(pos, size Point) Rect {
	return Rect{pos, size}
}

func (o Rect) XYWH() (int32, int32, int32, int32) {
	return int32(o.Pos.X), int32(o.Pos.Y), int32(o.Size.X), int32(o.Size.Y)
}

func (o Rect) Rect() *sdl.Rect {
	return &sdl.Rect{int32(o.Pos.X), int32(o.Pos.Y), int32(o.Size.X), int32(o.Size.Y)}
}

func (o *Rect) Extend(d int) {
	o.Pos.X -= d
	o.Pos.Y -= d
	o.Size.X += d << 1
	o.Size.Y += d << 1
}

func callback(fn CallbackFn, o Widget) bool {
	if fn != nil {
		return fn(o)
	}
	return false
}

func __(err error) {
	if err != nil {
		panic(err)
	}
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
