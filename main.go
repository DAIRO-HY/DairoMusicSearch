package main

import (
	"DairoMusicSearch/Application"
	"DairoMusicSearch/WebHandle"
)

func main() {

	//程序初始化操作
	Application.Init()
	WebHandle.Start()
}
