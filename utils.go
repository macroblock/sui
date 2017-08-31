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
	sign int
	sdl.Color
}

func (o Color) Hi() Color {
	o.R += uint8(Palette.incHi.sign * int(Palette.incHi.R))
	o.G += uint8(Palette.incHi.sign * int(Palette.incHi.G))
	o.B += uint8(Palette.incHi.sign * int(Palette.incHi.B))
	o.A += uint8(Palette.incHi.sign * int(Palette.incHi.A))
	return o
}

func (o Color) Lo() Color {
	o.R += uint8(Palette.incLo.sign * int(Palette.incHi.R))
	o.G += uint8(Palette.incLo.sign * int(Palette.incHi.G))
	o.B += uint8(Palette.incLo.sign * int(Palette.incHi.B))
	o.A += uint8(Palette.incLo.sign * int(Palette.incHi.A))
	return o
}

func Color32(c int64) Color {
	sign := 1
	if c < 0 {
		sign = -1
	}
	return Color{sign: sign,
		Color: sdl.Color{
			R: uint8((c >> 0) & 0xff),
			G: uint8((c >> 8) & 0xff),
			B: uint8((c >> 16) & 0xff),
			A: uint8((c >> 24) & 0xff),
		},
	}
}

func Color32b(c int64) Color {
	sign := 1
	if c < 0 {
		sign = -1
	}
	return Color{sign: sign,
		Color: sdl.Color{
			R: uint8((c >> 24) & 0xff),
			G: uint8((c >> 16) & 0xff),
			B: uint8((c >> 8) & 0xff),
			A: uint8((c >> 0) & 0xff),
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

func callback(fn func()) bool {
	if fn != nil {
		fn()
		return true
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
