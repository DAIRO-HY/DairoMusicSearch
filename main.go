package main

import (
	"DairoMusicSearch/Application"
	"DairoMusicSearch/WebHandle"
	"DairoMusicSearch/controller"
)

func main() {

	//程序初始化操作
	Application.Init()
	controller.RegistryController()
	WebHandle.Start()
}
