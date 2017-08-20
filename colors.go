package sui

type palette struct {
	Background   Color
	BackgroundHi Color
	BackgroundLo Color
	Foreground   Color
	Text         Color
}

var Palette = palette{
	Background:   Color32(0xffccbbbb),
	BackgroundHi: Color32(0xffddcccc),
	BackgroundLo: Color32(0xffbbaaaa),
	Foreground:   Color32(0xff443333),
	Text:         Color32(0xff443333),
}
