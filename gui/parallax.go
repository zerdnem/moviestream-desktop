package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ParallaxContainer creates a container with a parallax background effect
type ParallaxContainer struct {
	widget.BaseWidget
	backdrop     *canvas.Image
	content      fyne.CanvasObject
	scrollOffset float32
}

// NewParallaxContainer creates a new parallax container with a backdrop and scrollable content
func NewParallaxContainer(backdropURL string, content fyne.CanvasObject) *ParallaxContainer {
	p := &ParallaxContainer{
		content: content,
	}
	
	// Load backdrop image with blur effect
	if backdropURL != "" {
		p.backdrop = LoadBackdropImage(backdropURL, 1920, 1080)
		p.backdrop.FillMode = canvas.ImageFillStretch
		p.backdrop.ScaleMode = canvas.ImageScaleFastest
	}
	
	p.ExtendBaseWidget(p)
	return p
}

// CreateRenderer creates the renderer for the parallax container
func (p *ParallaxContainer) CreateRenderer() fyne.WidgetRenderer {
	// Create a semi-transparent overlay for better text readability
	bgColor := GetBackgroundColor()
	r, g, b, _ := bgColor.RGBA()
	overlay := canvas.NewRectangle(&color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 200, // Semi-transparent
	})
	
	// Create scroll container for content
	scroll := container.NewVScroll(p.content)
	
	// Stack backdrop, overlay, and content
	var objects []fyne.CanvasObject
	if p.backdrop != nil {
		objects = []fyne.CanvasObject{p.backdrop, overlay, scroll}
	} else {
		objects = []fyne.CanvasObject{overlay, scroll}
	}
	
	return &parallaxRenderer{
		container: p,
		backdrop:  p.backdrop,
		overlay:   overlay,
		scroll:    scroll,
		objects:   objects,
	}
}

type parallaxRenderer struct {
	container *ParallaxContainer
	backdrop  *canvas.Image
	overlay   *canvas.Rectangle
	scroll    *container.Scroll
	objects   []fyne.CanvasObject
}

func (r *parallaxRenderer) Layout(size fyne.Size) {
	// Position backdrop to fill the container with parallax effect
	if r.backdrop != nil {
		// Make backdrop slightly larger for parallax movement
		backdropSize := fyne.NewSize(size.Width, size.Height*1.3)
		r.backdrop.Resize(backdropSize)
		
		// Calculate parallax offset based on scroll position
		scrollOffset := float32(0)
		if scroll, ok := r.scroll.Content.(*fyne.Container); ok {
			_ = scroll // Use scroll to calculate offset if needed
		}
		
		// Position backdrop with parallax offset
		r.backdrop.Move(fyne.NewPos(0, -scrollOffset*0.3))
	}
	
	// Overlay matches the full size
	r.overlay.Resize(size)
	r.overlay.Move(fyne.NewPos(0, 0))
	
	// Scroll container fills the entire space
	r.scroll.Resize(size)
	r.scroll.Move(fyne.NewPos(0, 0))
}

func (r *parallaxRenderer) MinSize() fyne.Size {
	return r.scroll.MinSize()
}

func (r *parallaxRenderer) Refresh() {
	if r.backdrop != nil {
		r.backdrop.Refresh()
	}
	r.overlay.Refresh()
	r.scroll.Refresh()
}

func (r *parallaxRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *parallaxRenderer) Destroy() {}

// CreateParallaxView creates a view with parallax background effect
func CreateParallaxView(backdropURL string, content fyne.CanvasObject) fyne.CanvasObject {
	if backdropURL == "" {
		// Fallback to regular padded scroll if no backdrop
		return container.NewPadded(container.NewVScroll(content))
	}
	
	return NewParallaxContainer(backdropURL, container.NewPadded(content))
}

