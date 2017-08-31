package sui

type palette struct {
	// Primary     Color
	// PrimaryLo   Color
	// PrimaryHi   Color // TEMP
	// Secondary   Color
	// SecondaryLo Color
	// SecondaryHi Color // TEMP
	// Success     Color
	// SuccessLo   Color
	// SuccessHi   Color // TEMP
	// Info        Color
	// InfoLo      Color
	// InfoHi      Color // TEMP
	// Warning     Color
	// WarningLo   Color
	// WarningHi   Color // TEMP
	// Danger      Color
	// DangerLo    Color
	// DangerHi    Color // TEMP
	// Light       Color
	// LightLo     Color
	// LightHi     Color // TEMP
	// Dark        Color
	// DarkLo      Color
	// DarkHi      Color // TEMP
	// White       Color
	// Black       Color

	// Base03  Color
	// Base02  Color
	// Base01  Color
	// Base00  Color
	// Base0   Color
	// Base1   Color
	// Base2   Color
	// Base3   Color
	// Yellow  Color
	// Orange  Color
	// Red     Color
	// Magenta Color
	// Violet  Color
	// Blue    Color
	// Cyan    Color
	// Green   Color

	incHi Color
	incLo Color

	Passive Color
	Active  Color
	Accent  Color
	Select  Color

	Normal  Color
	Info    Color
	Warning Color
	Error   Color
}

var Palette = palette{
	// Primary:     Color32b(0x007bffff),
	// PrimaryLo:   Color32b(0x0062ccff),
	// PrimaryHi:   Color32b(0x007bffff), // TEMP
	// Secondary:   Color32b(0x868e96ff),
	// SecondaryLo: Color32b(0x6c757dff),
	// SecondaryHi: Color32b(0x868e96ff), // TEMP
	// Success:     Color32b(0x28a745ff),
	// SuccessLo:   Color32b(0x1e7e34ff),
	// SuccessHi:   Color32b(0x28a745ff), // TEMP
	// Info:        Color32b(0x17a2b8ff),
	// InfoLo:      Color32b(0x117a8bff),
	// InfoHi:      Color32b(0x17a2b8ff), // TEMP
	// Warning:     Color32b(0xffc107ff),
	// WarningLo:   Color32b(0xd39e00ff),
	// WarningHi:   Color32b(0xffc107ff), // TEMP
	// Danger:      Color32b(0xdc3545ff),
	// DangerLo:    Color32b(0xbd2130ff),
	// DangerHi:    Color32b(0xdc3545ff), // TEMP
	// Light:       Color32b(0xf8f9faff),
	// LightLo:     Color32b(0xdae0e5ff),
	// LightHi:     Color32b(0xf8f9faff), // TEMP
	// Dark:        Color32b(0x343a40ff),
	// DarkLo:      Color32b(0x1d2124ff),
	// DarkHi:      Color32b(0x343a40ff), // TEMP
	// White:       Color32b(0xffffffff),
	// Black:       Color32b(0x000000ff),
	// Base03:  Color32b(0x002b36ff),
	// Base02:  Color32b(0x073642ff),
	// Base01:  Color32b(0x0d4e55ff),
	// Base00:  Color32b(0x657b83ff),
	// Base0:   Color32b(0x839496ff),
	// Base1:   Color32b(0x93a1a1ff),
	// Base2:   Color32b(0xeee8d5ff),
	// Base3:   Color32b(0xfdf6e3ff),
	// Yellow:  Color32b(0xb58900ff),
	// Orange:  Color32b(0xcb4b16ff),
	// Red:     Color32b(0xdc322fff),
	// Magenta: Color32b(0xd33682ff),
	// Violet:  Color32b(0x6c71c4ff),
	// Blue:    Color32b(0x268bd2ff),
	// Cyan:    Color32b(0x2aa198ff),
	// Green:   Color32b(0x859900ff),
	incHi: Color32b(0x04040400),
	incLo: Color32b(-0x04040400),

	Passive: Color32b(0x073642ff),
	Active:  Color32b(0x002b36ff),
	Accent:  Color32b(0xfdf6e3ff),
	Select:  Color32b(0x000000ff),

	Normal:  Color32b(0x657b83ff),
	Info:    Color32b(0x268bd2ff),
	Warning: Color32b(0xcb4b16ff),
	Error:   Color32b(0xdc322fff),
}
