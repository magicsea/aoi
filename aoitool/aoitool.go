package main

import (
	"fmt"
	g "github.com/magicsea/gosprite"
)


const (
	screenW = 1024
	screenH = 768
)

func main() {
	fmt.Println("start")
	err := g.Start(new(ToolScane),screenW, screenH, "AOITool")
	if err != nil {
		fmt.Println("run error:", err)
	}
}
