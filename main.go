package main

import (
	"moviestream-gui/gui"
	"moviestream-gui/history"
	"moviestream-gui/queue"
	"moviestream-gui/settings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	// Use NewWithID to fix preferences API warning
	myApp := app.NewWithID("com.moviestream.app")
	
	// Apply Warp Dark theme
	myApp.Settings().SetTheme(gui.GetCurrentTheme())
	
	// Initialize settings system
	settings.Initialize(myApp)
	
	// Initialize queue and history systems
	queue.Initialize(myApp)
	history.Initialize(myApp)
	
	myWindow := myApp.NewWindow("MovieStream")
	
	// Set application icon using theme icon
	myWindow.SetIcon(theme.MediaVideoIcon())
	
	// Create and show the main GUI
	gui.CreateMainUI(myWindow)
	
	myWindow.Resize(fyne.NewSize(900, 650))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
