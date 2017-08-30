package sui

type palette struct {
	Background     Color
	BackgroundHi   Color
	BackgroundLo   Color
	Foreground     Color
	Text           Color
	SelectedItemBg Color

	DangerAlert  Color
	WarningAlert Color
	InfoAlert    Color
}

var Palette = palette{
	BackgroundHi:   Color32(0xffeedddd),
	Background:     Color32(0xffddcccc),
	BackgroundLo:   Color32(0xffccbbbb),
	Foreground:     Color32(0xff443333),
	Text:           Color32(0xff443333),
	SelectedItemBg: Color32(0xff887777),

	InfoAlert:    Color32(0xffff0000),
	WarningAlert: Color32(0xffff0000),
	DangerAlert:  Color32(0xffff0000),

	// Primary: Color32(0xff007bff),
	// PrimaryLo: Color32(0xff0062cc),
	// PrimaryHi: Color32(0xff007bff), // TEMP
	// Secondary: Color32(0xff868e96),
	// SecondaryLo: Color32(0xff6c757d),
	// SecondaryHi: Color32(0xff868e96), // TEMP
	// Success: Color32(0xff28a745),
	// SuccessLo: Color32(0xff1e7e34),
	// SuccessHi: Color32(0xff28a745), // TEMP
	// Info: Color32(0xff17a2b8),
	// InfoLo: Color32(0xff117a8b),
	// InfoHi: Color32(0xff17a2b8), // TEMP
	// Warning: Color32(0xffffc107),
	// WarningLo: Color32(0xffd39e00),
	// WarningHi: Color32(0xffffc107), // TEMP
	// Danger: Color32(0xffdc3545),
	// DangerLo: Color32(0xffbd2130),
	// DangerHi: Color32(0xffdc3545), // TEMP
	// Light: Color32(0xfff8f9fa),
	// LightLo: Color32(0xffdae0e5),
	// LightHi: Color32(0xfff8f9fa), // TEMP
	// Dark: Color32(0xff343a40),
	// DarkLo: Color32(0xff1d2124),
	// DarkHi: Color32(0xff343a40), // TEMP
	// White: Color32(0xfffffff),
	// Black: Color32(0xff000000),
}
