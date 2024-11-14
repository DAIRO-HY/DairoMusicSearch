package controller

import (
	"DairoMusicSearch/controller/download"
	"DairoMusicSearch/controller/search"
)

func RegistryController() {
	search.Init()
	download.Init()
}
