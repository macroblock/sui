package main

import (
	"fmt"

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
	fmt.Println(isShift)
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
			o.SetColor(drawColor)
			o.Fill(rect)
			o.SetColor(clearColor)
			o.Rect(rect)
			o.SetTextColor(clearColor)
			o.WriteText(pos, o.items[i].Name)
		} else if o.items[i].Selected {
			o.SetColor(sui.Palette.SelectedItemBg)
			o.Fill(rect)
			o.SetColor(clearColor)
			o.Rect(rect)
			o.SetTextColor(clearColor)
			o.WriteText(pos, o.items[i].Name)
		} else {
			o.SetColor(clearColor)
			o.Fill(rect)
			o.SetColor(drawColor)
			o.Rect(rect)
			o.SetTextColor(textColor)
			o.WriteText(pos, o.items[i].Name)
		}
		pos.Y += itemHeight
	}
	o.SetClearColor(clearColor)
	o.SetColor(drawColor)
	o.SetTextColor(textColor)
	o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
}
