package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/macroblock/sui"
)

const (
	minThreads = 1
	maxThreads = 50
)

var (
	ftpHost     = ""
	ftpPort     = -1
	ftpUser     = ""
	ftpPassword = ""

	files []string

	numThreads = 1
)

func isClosed(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func ftpTest() {
	c, err := ftp.DialTimeout(ftpHost+":"+strconv.Itoa(ftpPort), 5*time.Second)
	if err != nil {
		panic(err)
	}

	err = c.Login(ftpUser, ftpPassword)
	if err != nil {
		panic(err)
	}

	err = c.ChangeDir("/master/temp")
	if err != nil {
		panic(err)
	}
	filename := "\\\\Rackstation\\RACKSTATION\\FILMS_RT\\_series_eng\\Komnata_104\\s01\\sd_2017_Komnata_104_s01__q0w0_trailer.mpg"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	ch := make(chan interface{})
	go func() {
		fmt.Println("sending...")
		c.StorFrom("/master/temp/xxx", file, 0)
		close(ch)
	}()

	fmt.Println("wait")
	for !isClosed(ch) {
		pos, _ := file.Seek(0, 1)
		stat, _ := file.Stat()
		len := stat.Size()
		fmt.Printf("uploaded: %d%%\r", int(float64(pos)/float64(len)*100))
	}
	fmt.Println("done.")

	// err = c.Logout()
	// if err != nil {
	// 	panic(err)
	// }

	c.Quit()

}

func onDraw(o sui.Widget) bool {
	rect := sui.NewRect(sui.Point{}, o.Size())
	//fmt.Println("rect: ", rect)
	//rect.Extend(-1)
	//fmt.Println("rect2: ", rect)
	//o.SetClearColor(sui.Color32(0xffff0000))
	o.Clear()
	o.SetColor(sui.Color32(0xffffffff))
	o.Rect(rect)
	pos := sui.NewPoint(10, 10)
	ofs := o.WriteText(pos, "~!@#$%^&*()_+|[]{};:'<>? TTF Test string 0123456789!")
	for _, fileName := range files {
		pos.X = 10
		pos.Y += ofs.Y
		ofs = o.WriteText(pos, fileName)
	}
	return true
}

/*func onEnter(o sui.Widget) bool {
	o.SetClearColor(sui.Palette.BackgroundHi)
	return true
}

func onLeave(o sui.Widget) bool {
	o.SetClearColor(sui.Palette.Background)
	return true
}*/

func onMouseClick(o sui.Widget) bool {
	fmt.Println("!!!!!!!! MouseClick: ", o)
	o.SetClearColor(sui.Palette.BackgroundLo)
	return true
}

func onMouseOver(o sui.Widget) bool {
	if o != nil {
		o.SetClearColor(sui.Palette.BackgroundHi)
	}
	if sui.PrevMouseOver() != nil {
		sui.PrevMouseOver().SetClearColor(sui.Palette.Background)
	}
	return true
}

func onPressMouseDown(o sui.Widget) bool {
	fmt.Println("MousePressDown: ", o)
	return true
}

func onPressMouseUp(o sui.Widget) bool {
	fmt.Println("MousePressUp: ", o)
	return true
}

func onDropFile(o sui.Widget) bool {
	files = append(files, sui.DropFile())
	return true
}

func main() {
	fmt.Println(ftpUser + ":" + ftpPassword + "@ftp://" + ftpHost + ":" + strconv.Itoa(ftpPort))
	ftpTest()
	err := sui.Init()
	defer sui.Close()
	if err != nil {
		panic(err)
	}
	root := sui.NewRootWindow("test", 800, 600)
	root.SetClearColor(sui.Palette.BackgroundLo)
	//root.SetClearColor(sui.Color32(0x00000000))
	root.OnDropFile = onDropFile
	//root.OnDraw = onDraw
	//root.OnEnter = onEnter
	//root.OnLeave = onLeave
	//root.OnMouseOver = onMouseOver
	root.OnMouseButtonDown = onPressMouseDown
	root.OnMouseButtonUp = onPressMouseUp
	//root.OnMouseClick = onMouseClick

	btnInc := sui.NewBox(40, 35)
	btnInc.Move(5, 5)
	btnInc.OnDraw = func(o sui.Widget) bool {
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Inc")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
		return true
	}
	btnInc.OnMouseClick = func(o sui.Widget) bool {
		numThreads = sui.MinInt(numThreads+1, maxThreads)
		sui.PostUpdate()
		return true
	}

	btnDec := sui.NewBox(40, 35)
	btnDec.Move(50, 5)
	btnDec.OnDraw = func(o sui.Widget) bool {
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Dec")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
		return true
	}
	btnDec.OnMouseClick = func(o sui.Widget) bool {
		numThreads = sui.MaxInt(numThreads-1, minThreads)
		sui.PostUpdate()
		return true
	}

	lblNumThreads := sui.NewBox(40, 35)
	lblNumThreads.Move(95, 5)
	lblNumThreads.OnDraw = func(o sui.Widget) bool {
		o.Clear()
		o.WriteText(sui.NewPoint(10, 5), strconv.Itoa(numThreads))
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
		return true
	}

	btnStop := sui.NewBox(50, 35)
	btnStop.Move(140, 5)
	btnStop.OnDraw = func(o sui.Widget) bool {
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Stop")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
		return true
	}
	btnStop.OnMouseClick = func(o sui.Widget) bool {
		//numThreads = sui.MaxInt(numThreads-1, minThreads)
		sui.PostUpdate()
		return true
	}

	btnPlay := sui.NewBox(50, 35)
	btnPlay.Move(195, 5)
	btnPlay.OnDraw = func(o sui.Widget) bool {
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Play")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
		return true
	}
	btnPlay.OnMouseClick = func(o sui.Widget) bool {
		//numThreads = sui.MaxInt(numThreads-1, minThreads)
		sui.PostUpdate()
		return true
	}

	lbFiles := sui.NewBox(790, 350)
	lbFiles.Move(5, 45)
	lbFiles.OnDraw = func(o sui.Widget) bool {
		o.Clear()
		pos := sui.NewPoint(10, 10)
		for _, fileName := range files {
			ofs := o.WriteText(pos, fileName)
			pos.X = 10
			pos.Y += ofs.Y
		}
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
		return true
	}

	root.AddChild(btnInc)
	root.AddChild(btnDec)
	root.AddChild(lblNumThreads)
	root.AddChild(btnStop)
	root.AddChild(btnPlay)
	root.AddChild(lbFiles)

	panel := sui.NewBox(500, 180)
	panel.Move(20, 400)
	//panel.SetClearColor(sui.Color32(0xffff000))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	panel.OnMouseClick = onMouseClick
	root.AddChild(panel)
	fmt.Println(root)

	panel = sui.NewBox(250, 140)
	panel.Move(40, 420)
	//panel.SetClearColor(sui.Color32(0xff00ff00))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	panel.OnMouseClick = onMouseClick
	root.AddChild(panel)
	fmt.Println(root)

	panel = sui.NewBox(200, 100)
	panel.Move(60, 440)
	//panel.SetClearColor(sui.Color32(0xff0000ff))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	panel.OnMouseClick = onMouseClick
	root.AddChild(panel)
	fmt.Println(root)
	//_ = sui.NewSystemWindow("test", 800, 600)

	panel = sui.NewBox(200, 180)
	panel.Move(540, 400)
	//panel.SetClearColor(sui.Color32(0xffff000))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	panel.OnMouseClick = onMouseClick
	root.AddChild(panel)
	fmt.Println(root)

	sui.Run()

	root.Close()
}
