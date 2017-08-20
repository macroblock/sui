package sui

import "github.com/veandco/go-sdl2/sdl"

type Point struct {
	X, Y int
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

func NewRect(pos, size Point) sdl.Rect {
	return sdl.Rect{int32(pos.X), int32(pos.Y), int32(size.X), int32(size.Y)}

}

func callback(fn CallbackFn, o Widgeter) bool {
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

func uint32toColor(c uint32) sdl.Color {
	return sdl.Color{
		uint8((c >> 0) & 0xff),
		uint8((c >> 8) & 0xff),
		uint8((c >> 16) & 0xff),
		uint8((c >> 24) & 0xff),
	}
}
