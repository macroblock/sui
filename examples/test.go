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

	ftpConn    *ftp.ServerConn
	numThreads = 1

	root    *sui.RootWindow
	lbFiles *ListBox
)

type ftpItem struct {
	filename string
	active   bool
	working  bool
	byteSent int64
	file     *os.File
	err      error
	done     chan interface{}
}

func isClosed(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func ftpInit() (*ftp.ServerConn, error) {
	c, err := ftp.DialTimeout(ftpHost+":"+strconv.Itoa(ftpPort), 5*time.Second)
	if err != nil {
		return nil, err
	}
	err = c.Login(ftpUser, ftpPassword)
	if err != nil {
		return nil, err
	}
	err = c.ChangeDir("/master/temp")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ftpClose(c *ftp.ServerConn) {
	c.Quit()
}

func ftpJob(item *ftpItem) {
	fmt.Println("sending:", item.filename)
	stat, err := item.file.Stat()
	if err != nil {
		fmt.Println("stat err: ", err)
	} else {
		fmt.Printf("size %d\n", stat.Size)
	}
	item.err = ftpConn.StorFrom("/master/temp/xxx", item.file, 0)
	if item.err != nil {
		fmt.Println("error!!!!:", item.err)
	}
	//close(item.done)
}

func ftpStartJob(item *ftpItem) {
	item.working = true
	//item.done = make(chan interface{})
	file, err := os.Open(item.filename)
	if err != nil {
		panic("preJob: " + fmt.Sprint(err))
	}
	item.file = file
	go ftpJob(item)
	//item.active = false
	//item.working = false
	//item.file = nil
	//file.Close()
	//if item.err != nil {
	//	panic("postJob: " + fmt.Sprint(err))
	//}
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
	filename := "c:\\tools\\src.mpg"
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

func onDraw() {
	o := sui.Sender()
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
}

/*func onEnter(o sui.Widget) bool {
	o.SetClearColor(sui.Palette.BackgroundHi)
	return true
}

func onLeave(o sui.Widget) bool {
	o.SetClearColor(sui.Palette.Background)
	return true
}*/

func onMouseClick() {
	o := sui.Sender()
	fmt.Println("!!!!!!!! MouseClick: ", o)
	o.SetClearColor(sui.Palette.BackgroundLo)
}

func onMouseOver() {
	x := sui.MouseOver()
	if x != nil && x != root {
		x.SetClearColor(sui.Palette.BackgroundHi)
	}
	if sui.PrevMouseOver() != nil && sui.PrevMouseOver() != root {
		sui.PrevMouseOver().SetClearColor(sui.Palette.Background)
	}
}

func onPressMouseDown() {
	o := sui.Sender()
	fmt.Println("MousePressDown: ", o)
}

func onPressMouseUp() {
	o := sui.Sender()
	fmt.Println("MousePressUp: ", o)
}

func onDropFile() {
	item := &ftpItem{}
	item.filename = sui.DropFile()
	item.active = true

	lbFiles.AddItem(fmt.Sprint(item.active, " ", item.filename), item)
	files = append(files, sui.DropFile())
}

func loop() {
	sui.PostUpdate()
	items := lbFiles.Items()
	nJobs := 0
	for i := range items {
		item := items[i].Data.(*ftpItem)
		if item.working {
			nJobs++
		}
		percents := 0
		if item.file != nil && item.working {
			stat, err := item.file.Stat()
			if err != nil {
				fmt.Println("draw stat err: ", err)
			} else {
				fmt.Printf("draw size %d\n", stat.Size())
				item.byteSent, _ = item.file.Seek(0, 1)
				percents = int(item.byteSent * 100 / stat.Size())
			}
		}
		items[i].Name = fmt.Sprint(percents, item.active, item.working, item.byteSent, " ", item.filename)
	}
	if nJobs >= numThreads {
		return
	}
	for i := range items {
		item := items[i].Data.(*ftpItem)
		if !item.working && item.active {
			item.working = true
			ftpStartJob(item)
		}
		items[i].Name = fmt.Sprint(item.active, item.working, item.byteSent, " ", item.filename)
	}

	//fmt.Print("x")
	sui.PostUpdate()
}

func main() {
	fmt.Println(ftpUser + ":" + ftpPassword + "@ftp://" + ftpHost + ":" + strconv.Itoa(ftpPort))
	//ftpTest()
	err := sui.Init()
	defer sui.Close()
	if err != nil {
		panic(err)
	}
	root = sui.NewRootWindow("test", 800, 600)
	root.SetClearColor(sui.Palette.BackgroundLo)
	//root.SetClearColor(sui.Color32(0x00000000))
	root.OnDropFile = onDropFile
	//root.OnDraw = onDraw
	//root.OnEnter = onEnter
	//root.OnLeave = onLeave
	root.OnMouseOver = onMouseOver
	root.OnMouseButtonDown = onPressMouseDown
	root.OnMouseButtonUp = onPressMouseUp
	//root.OnMouseClick = onMouseClick

	btnInc := sui.NewBox(40, 35)
	btnInc.Move(5, 5)
	btnInc.OnMouseOver = onMouseOver
	btnInc.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Inc")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnInc.OnMouseClick = func() {
		numThreads = sui.MinInt(numThreads+1, maxThreads)
		sui.PostUpdate()
	}

	btnDec := sui.NewBox(40, 35)
	btnDec.Move(50, 5)
	btnDec.OnMouseOver = onMouseOver
	btnDec.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Dec")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnDec.OnMouseClick = func() {
		numThreads = sui.MaxInt(numThreads-1, minThreads)
		sui.PostUpdate()
	}

	lblNumThreads := sui.NewBox(40, 35)
	lblNumThreads.Move(95, 5)
	lblNumThreads.OnMouseOver = onMouseOver
	lblNumThreads.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(10, 5), strconv.Itoa(numThreads))
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}

	btnStop := sui.NewBox(50, 35)
	btnStop.Move(140, 5)
	btnStop.OnMouseOver = onMouseOver
	btnStop.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Stop")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnStop.OnMouseClick = func() {
		sui.PostUpdate()
	}

	btnMoveToTop := sui.NewBox(95, 35)
	btnMoveToTop.Move(300, 5)
	btnMoveToTop.OnMouseOver = onMouseOver
	btnMoveToTop.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "to Top")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnMoveToTop.OnMouseClick = func() {
		sui.PostUpdate()
	}

	btnMoveToBottom := sui.NewBox(95, 35)
	btnMoveToBottom.Move(400, 5)
	btnMoveToBottom.OnMouseOver = onMouseOver
	btnMoveToBottom.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "to Bottom")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnMoveToBottom.OnMouseClick = func() {
		sui.PostUpdate()
	}

	btnPlay := sui.NewBox(50, 35)
	btnPlay.Move(195, 5)
	btnPlay.OnMouseOver = onMouseOver
	btnPlay.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Play")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnPlay.OnMouseClick = func() {
		sui.PostUpdate()
	}

	lbFiles = NewListBox(790, 350)
	lbFiles.Move(5, 45)
	lbFiles.OnMouseOver = onMouseOver

	root.AddChild(btnInc)
	root.AddChild(btnDec)
	root.AddChild(lblNumThreads)
	root.AddChild(btnStop)
	root.AddChild(btnPlay)
	root.AddChild(btnMoveToTop)
	root.AddChild(btnMoveToBottom)
	root.AddChild(lbFiles)

	/*panel := sui.NewBox(500, 180)
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
	*/

	ftpConn, err = ftpInit()
	defer ftpClose(ftpConn)
	if err != nil {
		panic(err)
	}

	sui.OnLoop = loop

	sui.Run()

	root.Close()
}
