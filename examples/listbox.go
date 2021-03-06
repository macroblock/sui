package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/macroblock/sui"
)

const (
	itemHeight = 20
	xOffs      = 5
	yOffs      = 0
)

type (
	listBoxItem struct {
		Selected bool
		Name     string
		Data     interface{}
	}

	ListBox struct {
		sui.Box
		items     []listBoxItem
		itemIndex int
		offset    int
	}
)

func (o *ListBox) Repaint() {
	sui.SetSender(o)
	o.Draw()
	sui.SetSender(nil)
	for _, child := range o.Children() {
		child.Repaint()
		src := sui.NewRect(sui.Point{}, child.Size())
		dst := sui.NewRect(child.Pos(), child.Size())
		child.Surface().Blit(src.Rect(), o.Surface(), dst.Rect())
	}
}

func NewListBox(w, h int) *ListBox {
	lb := ListBox{
		Box:       *sui.NewBox(w, h),
		itemIndex: -1,
	}
	lb.OnDraw = draw
	lb.OnMouseClick = mouseClick
	lb.OnMouseScroll = mouseScroll
	return &lb
}

func (o *ListBox) AddItem(str string, data interface{}) {
	item := listBoxItem{
		Selected: false,
		Name:     str,
		Data:     data,
	}
	o.items = append(o.items, item)
	sui.PostUpdate()
}

func (o *ListBox) Items() []listBoxItem {
	return o.items
}

func (o *ListBox) CalcOffset() {
	if o.offset > o.itemIndex {
		o.offset = sui.MaxInt(0, o.itemIndex)
	}
	if o.offset+o.Size().Y/itemHeight-1 < o.itemIndex {
		o.offset = o.itemIndex - o.Size().Y/itemHeight + 1
		o.offset = sui.MinInt(len(o.items)-1, o.offset)
	}
}

func mouseClick() {
	o := sui.Sender().(*ListBox)
	isShift := sui.ModShift() != 0
	//fmt.Println(isShift)
	index := o.offset + sui.MousePos().Y/itemHeight
	if index >= len(o.items) {
		return
	}
	if isShift && o.itemIndex >= 0 {
		a := sui.MinInt(o.itemIndex, index)
		b := sui.MaxInt(o.itemIndex, index)
		for i := a; i <= b; i++ {
			o.items[i].Selected = true
		}
	} else {
		for i := range o.items {
			o.items[i].Selected = false
		}
	}
	o.itemIndex = index
	sui.PostUpdate()
}

func mouseScroll() {
	o := sui.Sender().(*ListBox)
	scroll := sui.MouseScroll()
	//fmt.Println(scroll)
	k := 2
	o.offset -= scroll.X*k + scroll.Y*k
	o.offset = sui.MaxInt(0, o.offset)
	bound := sui.MaxInt(0, len(o.items)-o.Size().Y/itemHeight)
	o.offset = sui.MinInt(bound, o.offset)
	sui.PostUpdate()
}

func bpsToStr(bps int) string {
	const kilo = 1000
	postfix := "B/s"
	val := float64(bps)
	if val > kilo {
		val /= kilo
		postfix = "KB/s"
		if val > kilo {
			val /= kilo
			postfix = "MB/s"
		}
	}

	format := "%.0f"
	if val < 100 {
		format = "%.1f"
	}
	return fmt.Sprintf(format, val) + " " + postfix
}

func secondsToStr(sec int) string {
	ret := strconv.Itoa(sec)
	ss := sec % 60
	sec /= 60
	mm := sec % 60
	hh := sec / 60
	ret = fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	return ret
}

func drawItem(rect sui.Rect, item *ftpItem) {
	o := sui.Sender().(*ListBox)
	textColor := o.TextColor()
	_ = textColor
	pos := rect.Pos

	if item.err == nil && item.oldErr == nil {
		o.SetTextColor(sui.Palette.Info)
	} else if item.stopped {
		o.SetTextColor(sui.Palette.Error)
	} else {
		o.SetTextColor(sui.Palette.Warning)
	}

	pos.X += 5
	if item.working {
		percent := -1
		if item.fileSize != 0 {
			percent = int(item.bytesSent * 100 / item.fileSize)
		}
		o.WriteText(pos, strconv.Itoa(percent)+"%")
	}

	pos.X += 45
	bps := 0
	if item.working {
		bps = item.Bps()
		o.WriteText(pos, bpsToStr(bps))
	}
	pos.X += 100
	if item.working && bps != 0 {
		o.WriteText(pos, secondsToStr(int(item.fileSize-item.bytesSent)/bps))
	}

	pos.X += 80
	//o.SetTextColor(textColor)
	o.WriteText(pos, filepath.Base(item.filename))
}

func draw() {
	o := sui.Sender().(*ListBox)
	//o.SetClearColor(sui.Color32(0xffffffff))
	o.Clear()
	clearColor := o.ClearColor()
	drawColor := o.Color()
	textColor := o.TextColor()
	pos := sui.NewPoint(xOffs, yOffs)
	for i := o.offset; i < len(o.items); i++ {
		if pos.Y >= o.Size().Y {
			break
		}
		rect := sui.NewRect(sui.NewPoint(0, pos.Y-yOffs), sui.NewPoint(o.Size().X, itemHeight+1))
		if i == o.itemIndex {
			//o.SetColor(drawColor)
			//o.Fill(rect)
			//o.SetColor(clearColor)
			//o.Rect(rect)
			//o.SetTextColor(clearColor)
			//o.WriteText(pos, o.items[i].Name)
			o.SetColor(sui.Palette.Accent)
			o.Fill(rect)
			o.SetColor(sui.Palette.Normal)
			o.Rect(rect)
			drawItem(rect, o.items[i].Data.(*ftpItem))
		} else if o.items[i].Selected {
			// o.SetColor(sui.Palette.Select)
			// o.Fill(rect)
			// o.SetColor(clearColor)
			// //o.Rect(rect)
			// o.SetTextColor(clearColor)
			//o.WriteText(pos, o.items[i].Name)
			o.SetColor(sui.Palette.Select)
			o.Fill(rect)
			o.SetColor(sui.Palette.Normal)
			o.Rect(rect)
			drawItem(rect, o.items[i].Data.(*ftpItem))
		} else {
			// o.SetColor(clearColor)
			// o.Fill(rect)
			// o.SetColor(drawColor)
			// //o.Rect(rect)
			// o.SetTextColor(textColor)
			// //o.WriteText(pos, o.items[i].Name)
			o.SetColor(clearColor)
			o.Fill(rect)
			o.SetColor(sui.Palette.Normal)
			o.Rect(rect)
			drawItem(rect, o.items[i].Data.(*ftpItem))
		}
		pos.Y += itemHeight
	}
	o.SetClearColor(clearColor)
	o.SetColor(drawColor)
	o.SetTextColor(textColor)
	o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
}
