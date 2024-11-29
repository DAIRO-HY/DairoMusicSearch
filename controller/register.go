package controller

import (
	"DairoMusicSearch/controller/download"
	"DairoMusicSearch/controller/search"
	"DairoMusicSearch/controller/set"
)

func RegistryController() {
	search.Init()
	download.Init()
	set.Init()
}
