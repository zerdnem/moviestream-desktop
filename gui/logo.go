package gui

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// CreateAppIcon generates the MovieStream application icon
func CreateAppIcon() image.Image {
	size := 512
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	
	// Background - deep purple/blue gradient
	bgColor := color.RGBA{R: 20, G: 20, B: 30, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
	
	// Draw film strip background
	drawFilmStrip(img, size)
	
	// Draw play button in center
	drawPlayButton(img, size)
	
	return img
}

// CreateAppLogo creates a smaller logo for the UI header
func CreateAppLogo(width, height float32) *canvas.Image {
	img := createLogoImage(int(width), int(height))
	logo := canvas.NewImageFromImage(img)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(width, height))
	return logo
}

func createLogoImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Transparent background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Draw simplified film strip
	filmColor := GeistGray10
	centerY := height / 2
	stripHeight := height / 3
	
	// Main film strip bar
	for y := centerY - stripHeight/2; y < centerY+stripHeight/2; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, filmColor)
		}
	}
	
	// Film holes
	holeSize := stripHeight / 4
	numHoles := 5
	spacing := width / (numHoles + 1)
	
	for i := 1; i <= numHoles; i++ {
		holeX := i * spacing
		// Top holes
		drawRect(img, holeX-holeSize/2, centerY-stripHeight/2+holeSize/2, holeSize, holeSize/2, color.RGBA{R: 0, G: 0, B: 0, A: 255})
		// Bottom holes
		drawRect(img, holeX-holeSize/2, centerY+stripHeight/2-holeSize, holeSize, holeSize/2, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	}
	
	// Draw play triangle in center
	playColor := GeistGray10
	triSize := stripHeight / 2
	drawTriangle(img, width/2, centerY, triSize, playColor)
	
	return img
}

func drawFilmStrip(img *image.RGBA, size int) {
	stripColor := color.RGBA{R: 60, G: 60, B: 80, A: 255}
	holeColor := color.RGBA{R: 20, G: 20, B: 30, A: 255}
	
	// Top and bottom strips
	stripHeight := size / 8
	
	// Top strip
	for y := 0; y < stripHeight; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, stripColor)
		}
	}
	
	// Bottom strip
	for y := size - stripHeight; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, stripColor)
		}
	}
	
	// Film holes
	holeSize := stripHeight / 2
	numHoles := 8
	spacing := size / (numHoles + 1)
	
	for i := 1; i <= numHoles; i++ {
		holeX := i * spacing
		// Top strip holes
		drawCircle(img, holeX, stripHeight/2, holeSize/2, holeColor)
		// Bottom strip holes
		drawCircle(img, holeX, size-stripHeight/2, holeSize/2, holeColor)
	}
}

func drawPlayButton(img *image.RGBA, size int) {
	centerX, centerY := size/2, size/2
	
	// Outer circle (button background)
	outerRadius := size / 3
	circleColor := color.RGBA{R: 120, G: 100, B: 255, A: 255}
	drawCircle(img, centerX, centerY, outerRadius, circleColor)
	
	// Inner glow
	glowColor := color.RGBA{R: 140, G: 120, B: 255, A: 200}
	drawCircle(img, centerX, centerY, outerRadius-10, glowColor)
	
	// Play triangle
	triangleSize := outerRadius / 2
	triangleColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	drawTriangle(img, centerX+5, centerY, triangleSize, triangleColor)
}

func drawCircle(img *image.RGBA, centerX, centerY, radius int, col color.Color) {
	for y := centerY - radius; y <= centerY+radius; y++ {
		for x := centerX - radius; x <= centerX+radius; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx+dy*dy <= radius*radius {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, col)
				}
			}
		}
	}
}

func drawTriangle(img *image.RGBA, centerX, centerY, size int, col color.Color) {
	// Right-pointing triangle
	for y := -size; y <= size; y++ {
		for x := 0; x <= size; x++ {
			// Triangle formula: check if point is inside
			if x <= size && math.Abs(float64(y)) <= float64(size)-float64(x) {
				px := centerX + x
				py := centerY + y
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, col)
				}
			}
		}
	}
}

func drawRect(img *image.RGBA, x, y, width, height int, col color.Color) {
	for py := y; py < y+height; py++ {
		for px := x; px < x+width; px++ {
			if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
				img.Set(px, py, col)
			}
		}
	}
}

