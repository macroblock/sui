package main

import (
	"fmt"

	"github.com/macroblock/sui"
)

func onDraw(o sui.Widgeter) bool {
	rect := sui.NewRect(sui.Point{}, o.Size())
	fmt.Println("rect: ", rect)
	o.Renderer().SetDrawColor(255, 77, 77, 255)
	//o.Surface().FillRect(&rect, 0xffff8888)
	o.Renderer().FillRect(&rect)
	rect.X++
	rect.Y++
	rect.W -= 2
	rect.H -= 2
	fmt.Println("rect2: ", rect)
	o.Renderer().SetDrawColor(255, 77, 255, 77)
	//o.Surface().FillRect(&rect, 0xff88ff88)
	o.Renderer().FillRect(&rect)
	o.WriteText(sui.Point{10, 10}, "~!@#$%^&*()_+|[]{};:'<>? TTF Test string 0123456789!", 0xffffffff)
	return true
}

func main() {
	err := sui.Init()
	defer sui.Close()
	if err != nil {
		panic(err)
	}
	root := sui.NewSystemWindow("test", 800, 600)

	panel := sui.NewWidget(500, 500)
	panel.Move(20, 20)
	panel.SetColor(0xffff0000)
	panel.SetOnDraw(onDraw)
	root.AddChild(panel)
	fmt.Println(root)

	panel = sui.NewWidget(250, 250)
	panel.Move(40, 40)
	panel.SetColor(0xff00ff00)
	panel.SetOnDraw(onDraw)
	root.AddChild(panel)
	fmt.Println(root)

	panel = sui.NewWidget(200, 200)
	panel.Move(60, 60)
	panel.SetColor(0xff0000ff)
	panel.SetOnDraw(onDraw)
	root.AddChild(panel)
	fmt.Println(root)
	//_ = sui.NewSystemWindow("test", 800, 600)

	sui.Run()

	root.Close()
}
