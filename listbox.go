package sui

const (
	itemHeight = 20
	xOffs      = 5
	yOffs      = 5
)

type (
	listBoxItem struct {
		Selected bool
		Name     string
		Data     interface{}
	}

	ListBox struct {
		Box
		items     []listBoxItem
		itemIndex int
		offset    int
	}
)

func NewListBox(w, h int) *ListBox {
	lb := ListBox{
		Box:       *NewBox(w, h),
		itemIndex: -1,
	}
	return &lb
}

func (o *ListBox) AddItem(str string, data interface{}) {
	item := listBoxItem{
		Selected: false,
		Name:     str,
		Data:     data,
	}
	o.items = append(o.items, item)
	PostUpdate()
}

func (o *ListBox) mouseClick() {
	if callback(o.OnMouseClick) {
		return
	}
	index := MousePos().Y / itemHeight
	if index >= len(o.items) {
		return
	}
	o.itemIndex = index
	o.items[index].Selected = !o.items[index].Selected
	PostUpdate()
}

func (o *ListBox) draw() {
	if callback(o.OnDraw) {
		return
	}
	o.SetClearColor(Color32(0xffffffff))
	o.Clear()
	clearColor := o.clearColor
	drawColor := o.drawColor
	textColor := o.textColor
	pos := NewPoint(xOffs, yOffs)
	for i := o.offset; i < len(o.items); i++ {
		if pos.Y >= o.Size().Y {
			break
		}
		rect := NewRect(NewPoint(0, pos.Y-yOffs), NewPoint(o.Size().X, itemHeight))
		if i != o.itemIndex {
			o.SetColor(clearColor)
			//o.Fill(rect)
			o.SetColor(drawColor)
			o.Rect(rect)
			o.SetColor(textColor)
			o.WriteText(pos, o.items[i].Name)
		} else {
			o.SetColor(drawColor)
			//o.Fill(rect)
			o.SetColor(clearColor)
			o.Rect(rect)
			o.SetColor(clearColor)
			o.WriteText(pos, o.items[i].Name)
		}
		pos.Y += itemHeight
	}
	o.clearColor = clearColor
	o.drawColor = drawColor
	o.textColor = textColor
	o.Rect(NewRect(NewPoint(0, 0), o.Size()))
}
