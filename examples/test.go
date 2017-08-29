package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/veandco/go-sdl2/sdl"

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

type ftpItem struct {
	filename  string
	stopped   bool
	working   bool
	done      bool
	bytesSent int64
	file      *os.File
	fileSize  int64
	err       error
}

func NewFtpItem(name string) *ftpItem {
	item := &ftpItem{}
	item.filename = name
	return item
}

func (o *ftpItem) Close() {
	if o.file != nil {
		o.file.Close()
	}
}

func (o *ftpItem) InitLocalFile() {
	if o.err != nil {
		return
	}
	o.bytesSent = 0
	file, err := os.Open(o.filename)
	if err != nil {
		o.err = err
		return
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		o.err = err
		return
	}
	o.fileSize = stat.Size()
	o.file = file
}

func (o *ftpItem) InitRemoteFile() {
	if o.err != nil {
		return
	}
}

func (o *ftpItem) FileSize() int64 {
	if o.file == nil {
		return -1
	}
	return o.fileSize
}

func (o *ftpItem) FilePos() int64 {
	if o.file == nil {
		return 0
	}
	pos, err := o.file.Seek(0, 1)
	if err != nil {
		pos = 0
	}
	return pos
}

func (o *ftpItem) WorkingOne() int {
	if o.working && o.file != nil {
		return 1
	}
	return 0
}

func (o *ftpItem) ReadyToWork() bool {
	if !o.working && !o.stopped {
		return true
	}
	return false
}

func (o *ftpItem) Stor() {
	if o.err != nil {
		return
	}
	o.err = ftpConn.StorFrom("/master/temp/xxx", o.file, 0)
}

func (o *ftpItem) StartJob() {
	if o.err != nil {
		return
	}
	o.working = true
	go o.job()
}

func (o *ftpItem) job() {
	o.stopped = true
	o.InitLocalFile()
	o.InitRemoteFile()
	o.Stor()
}

func loop() {
	sui.PostUpdate()
	items := lbFiles.Items()
	numWorkers := 0
	for i := range items {
		item := items[i].Data.(*ftpItem)
		numWorkers += item.WorkingOne()
		percents := int(item.FilePos() * 100 / item.FileSize())
		items[i].Name = fmt.Sprint(percents, item.stopped, item.working, " ", item.filename)
	}
	if numWorkers >= numThreads {
		return
	}
	for i := range items {
		item := items[i].Data.(*ftpItem)
		if item.ReadyToWork() {
			item.working = true
			//item.StartJob()
			return
		}
		//items[i].Name = fmt.Sprint(item.active, item.working, item.bytesSent, " ", item.filename)
	}
}

func deleteFtpItem(items []listBoxItem, i int) []listBoxItem {
	if i < 0 && i >= len(items) {
		panic("delete index wrong")
	}
	fmt.Println("\nin: ", i, "\n", items)
	if items[i].Data != nil {
		item := items[i].Data.(*ftpItem)
		item.Close()
		items[i].Data = nil
	}
	if i == 0 {
		items = items[1:]
	} else if i == len(items)-1 {
		items = items[:i]
	} else {
		items = append(items[:i], items[i+1:]...)
	}
	//fmt.Println("\nout: ", i, "\n", items)
	return items
}

func onKeyPress() {
	switch sui.KeySym() {
	case sdl.K_DELETE:
		items := lbFiles.Items()
		for i := 0; i < len(items); {
			if items[i].Selected || i == lbFiles.itemIndex {
				items = deleteFtpItem(items, i)
				lbFiles.itemIndex = -1
				lbFiles.items = items
				sui.PostUpdate()
			} else {
				i++
			}
		}
	case sdl.K_HOME, sdl.K_LEFT:
		lbFiles.itemIndex = sui.MinInt(0, len(lbFiles.Items())-1)
		lbFiles.offset = 0
		sui.PostUpdate()
	case sdl.K_END, sdl.K_RIGHT:
		lbFiles.itemIndex = len(lbFiles.Items()) - 1
		lbFiles.offset = sui.MaxInt(0, lbFiles.itemIndex-lbFiles.Size().Y/itemHeight+1)
		sui.PostUpdate()
	case sdl.K_UP:
		lbFiles.itemIndex--
		lbFiles.itemIndex = sui.MaxInt(0, lbFiles.itemIndex)
		lbFiles.itemIndex = sui.MinInt(len(lbFiles.Items())-1, lbFiles.itemIndex)
		lbFiles.CalcOffset()
		sui.PostUpdate()
	case sdl.K_DOWN:
		fmt.Println("key DOWN")
		lbFiles.itemIndex++
		lbFiles.itemIndex = sui.MaxInt(0, lbFiles.itemIndex)
		lbFiles.itemIndex = sui.MinInt(len(lbFiles.Items())-1, lbFiles.itemIndex)
		lbFiles.CalcOffset()
		sui.PostUpdate()
	}
}

func moveToTop(toTop bool) {
	items := lbFiles.Items()
	newItems := []listBoxItem{}
	for i := range items {
		if items[i].Selected || i == lbFiles.itemIndex {
			p := items[i]
			items[i].Data = nil
			//p.Selected = false
			newItems = append(newItems, p)
			lbFiles.itemIndex = -1
			lbFiles.items = items
			sui.PostUpdate()
		}
	}
	for i := 0; i < len(items); {
		if items[i].Data == nil {
			items = deleteFtpItem(items, i)
		} else {
			i++
		}
	}
	itemIndex := 0
	if toTop {
		items = append(newItems, items...)
	} else {
		items = append(items, newItems...)
		itemIndex = len(items) - 1
	}
	lbFiles.items = items
	lbFiles.itemIndex = itemIndex
	lbFiles.CalcOffset()
}

func onDropFile() {
	item := NewFtpItem(sui.DropFile())
	lbFiles.AddItem(fmt.Sprint(item.stopped, " ", item.filename), item)
	lbFiles.itemIndex = len(lbFiles.items) - 1
	lbFiles.CalcOffset()
	sui.PostUpdate()
	//files = append(files, sui.DropFile())
}

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
	root.OnKeyPress = onKeyPress
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
	//lblNumThreads.OnMouseOver = onMouseOver
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
		moveToTop(true)
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
		moveToTop(false)
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

	fInfo := sui.NewBox(790, 195)
	fInfo.Move(5, 400)
	fInfo.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Info")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}

	root.OnResize = func() {
		size := root.Size()
		w := size.X - 10
		h := size.Y - 250
		lbFiles.Resize(w, h)
		x := 5
		y := size.Y - 200
		w = size.X - 10
		h = 195
		fInfo.Move(x, y)
		fInfo.Resize(w, h)
	}

	root.AddChild(btnInc)
	root.AddChild(btnDec)
	root.AddChild(lblNumThreads)
	root.AddChild(btnStop)
	root.AddChild(btnPlay)
	root.AddChild(btnMoveToTop)
	root.AddChild(btnMoveToBottom)
	root.AddChild(lbFiles)
	root.AddChild(fInfo)

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
