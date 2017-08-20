package main

import (
	"fmt"

	"github.com/macroblock/sui"
)

func onDraw(o sui.Widget) bool {
	rect := sui.NewRect(sui.Point{}, o.Size())
	//fmt.Println("rect: ", rect)
	//rect.Extend(-1)
	//fmt.Println("rect2: ", rect)
	//o.SetClearColor(sui.Color32(0xffff0000))
	o.Clear()
	o.SetColor(sui.Color32(0xffffffff))
	o.Rect(rect)
	o.WriteText(sui.Point{10, 10}, "~!@#$%^&*()_+|[]{};:'<>? TTF Test string 0123456789!")
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

func main() {
	err := sui.Init()
	defer sui.Close()
	if err != nil {
		panic(err)
	}
	root := sui.NewRootWindow("test", 800, 600)
	//root.SetClearColor(sui.Color32(0x00000000))
	root.OnDraw = onDraw
	//root.OnEnter = onEnter
	//root.OnLeave = onLeave
	root.OnMouseOver = onMouseOver
	root.OnMouseButtonDown = onPressMouseDown
	root.OnMouseButtonUp = onPressMouseUp

	panel := sui.NewBox(500, 500)
	panel.Move(20, 20)
	//panel.SetClearColor(sui.Color32(0xffff000))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	root.AddChild(panel)
	fmt.Println(root)

	panel = sui.NewBox(250, 250)
	panel.Move(40, 40)
	//panel.SetClearColor(sui.Color32(0xff00ff00))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	root.AddChild(panel)
	fmt.Println(root)

	panel = sui.NewBox(200, 200)
	panel.Move(60, 60)
	//panel.SetClearColor(sui.Color32(0xff0000ff))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	root.AddChild(panel)
	fmt.Println(root)
	//_ = sui.NewSystemWindow("test", 800, 600)

	panel = sui.NewBox(200, 500)
	panel.Move(540, 20)
	//panel.SetClearColor(sui.Color32(0xffff000))
	panel.OnDraw = onDraw
	//panel.OnEnter = onEnter
	//panel.OnLeave = onLeave
	panel.OnMouseOver = onMouseOver
	panel.OnMouseButtonDown = onPressMouseDown
	panel.OnMouseButtonUp = onPressMouseUp
	root.AddChild(panel)
	fmt.Println(root)

	sui.Run()

	root.Close()
}
