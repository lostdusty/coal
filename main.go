package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	appGuibalt    = app.NewWithID("br.dev.lostdusty.guibalt")
	windowGuibalt = appGuibalt.NewWindow("Guibalt")
)

func main() {
	windowGuibalt.CenterOnScreen()
	windowGuibalt.SetMaster()
	windowGuibalt.Resize(fyne.Size{Width: 600, Height: 400})

	startScreen := showMainScreen()
	if !appGuibalt.Preferences().BoolWithFallback("first-run", false) {
		log.Println("showing first run screen...")
		startScreen = showFirstRunScreen()
	}

	windowGuibalt.SetContent(startScreen)
	windowGuibalt.ShowAndRun()
}
