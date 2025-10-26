package main

import (
	"moviestream-gui/gui"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("MovieStream - Movies & TV Shows")
	
	// Create and show the main GUI
	gui.CreateMainUI(myWindow)
	
	myWindow.Resize(fyne.NewSize(1000, 700))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

