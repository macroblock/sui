package sui

type palette struct {
	Background     Color
	BackgroundHi   Color
	BackgroundLo   Color
	Foreground     Color
	Text           Color
	SelectedItemBg Color
}

var Palette = palette{
	BackgroundHi:   Color32(0xffeedddd),
	Background:     Color32(0xffddcccc),
	BackgroundLo:   Color32(0xffccbbbb),
	Foreground:     Color32(0xff443333),
	Text:           Color32(0xff443333),
	SelectedItemBg: Color32(0xff887777),
}
