package main

import (
	"log"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	coalApp    = app.NewWithID("br.dev.lostdusty.coal")
	coalWindow = coalApp.NewWindow("coal")
	regexPaste = regexp.MustCompile(`https?\:\/\/[^\s]+`)
)

func main() {
	coalWindow.CenterOnScreen()
	coalWindow.SetMaster()
	coalWindow.Resize(fyne.Size{Width: 600, Height: 400})

	startScreen := showMainScreen()
	if !coalApp.Preferences().BoolWithFallback("first-run", false) {
		log.Println("showing first run screen...")
		startScreen = showFirstRunScreen()
	}

	coalWindow.SetContent(startScreen)
	coalWindow.ShowAndRun()
}
