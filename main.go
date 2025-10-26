package main

import (
	"bytes"
	"image"
	"image/png"
	"moviestream-gui/gui"
	"moviestream-gui/history"
	"moviestream-gui/queue"
	"moviestream-gui/settings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	
	// Set application icon
	iconImg := gui.CreateAppIcon()
	myWindow.SetIcon(fyne.NewStaticResource("icon.png", iconToBytes(iconImg)))
	
	// Create and show the main GUI
	gui.CreateMainUI(myWindow)
	
	myWindow.Resize(fyne.NewSize(900, 650))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

// iconToBytes converts an image to PNG bytes
func iconToBytes(img image.Image) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}
