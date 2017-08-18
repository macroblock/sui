package main

import "github.com/macroblock/sui"

func main() {
	err := sui.Init()
	defer sui.Close()
	if err != nil {
		panic(err)
	}
	_, err = sui.NewSystemWindow("test", 800, 600)
	if err != nil {
		panic(err)
	}

	_, err = sui.NewSystemWindow("test", 800, 600)
	if err != nil {
		panic(err)
	}
}
