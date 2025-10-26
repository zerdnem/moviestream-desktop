package main

import (
	"moviestream-gui/gui"
	"moviestream-gui/settings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Use NewWithID to fix preferences API warning
	myApp := app.NewWithID("com.moviestream.app")
	
	// Initialize settings system
	settings.Initialize(myApp)
	
	myWindow := myApp.NewWindow("MovieStream - Movies & TV Shows")
	
	// Create and show the main GUI
	gui.CreateMainUI(myWindow)
	
	myWindow.Resize(fyne.NewSize(1000, 700))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

