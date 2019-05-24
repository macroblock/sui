package sui

import (
	"golang.org/x/image/font/gofont/goregular"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	defaultFont *ttf.Font
	// LatoRegular20 ...
	//LatoRegular20 *ttf.Font
	// LatoRegular24 ...
	//LatoRegular24 *ttf.Font
	// LatoRegular14 ...
	//LatoRegular14 *ttf.Font
	// LatoRegular12 ...
	//LatoRegular12 *ttf.Font
)

// InitFonts ...
func InitFonts() {
	// rwops := sdl.RWFromMem(unsafe.Pointer(&goregular.TTF[0]), len(goregular.TTF))
	rwops, err := sdl.RWFromMem(goregular.TTF)
	if err != nil {
		panic(err)
	}
	defaultFont, _ = ttf.OpenFontRW(rwops, 1, 16)
	// LatoRegular20, _ = ttf.OpenFontRW(rwops, 1, 20)
	// LatoRegular24, _ = ttf.OpenFontRW(rwops, 1, 24)
	// LatoRegular14, _ = ttf.OpenFontRW(rwops, 1, 14)
	// LatoRegular12, _ = ttf.OpenFontRW(rwops, 1, 12)
}
