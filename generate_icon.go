// +build ignore

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func main() {
	img := createIcon()
	f, err := os.Create("Icon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func createIcon() image.Image {
	size := 512
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	
	// Background - simple dark background
	bgColor := color.RGBA{R: 30, G: 30, B: 40, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
	
	// Draw film strip frame
	drawFilmFrame(img, size)
	
	// Draw play button in center
	drawPlayButton(img, size)
	
	return img
}

func drawFilmFrame(img *image.RGBA, size int) {
	frameColor := color.RGBA{R: 220, G: 220, B: 230, A: 255}
	innerColor := color.RGBA{R: 50, G: 50, B: 60, A: 255}
	
	thickness := size / 12
	margin := size / 6
	
	// Draw outer frame
	// Top
	drawRect(img, margin, margin, size-2*margin, thickness, frameColor)
	// Bottom
	drawRect(img, margin, size-margin-thickness, size-2*margin, thickness, frameColor)
	// Left
	drawRect(img, margin, margin, thickness, size-2*margin, frameColor)
	// Right
	drawRect(img, size-margin-thickness, margin, thickness, size-2*margin, frameColor)
	
	// Draw film perforations on sides
	perfSize := thickness / 2
	numPerfs := 6
	spacing := (size - 2*margin) / (numPerfs + 1)
	
	for i := 1; i <= numPerfs; i++ {
		y := margin + i*spacing - perfSize/2
		// Left perforations
		drawRect(img, margin+thickness/4, y, perfSize, perfSize, innerColor)
		// Right perforations
		drawRect(img, size-margin-thickness/4-perfSize, y, perfSize, perfSize, innerColor)
	}
}

func drawPlayButton(img *image.RGBA, size int) {
	centerX, centerY := size/2, size/2
	
	// Outer circle
	outerRadius := size / 4
	circleColor := color.RGBA{R: 100, G: 150, B: 255, A: 255}
	drawCircle(img, centerX, centerY, outerRadius, circleColor)
	
	// Inner circle (lighter)
	innerRadius := outerRadius - 10
	innerColor := color.RGBA{R: 120, G: 170, B: 255, A: 255}
	drawCircle(img, centerX, centerY, innerRadius, innerColor)
	
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
	for y := -size; y <= size; y++ {
		for x := 0; x <= size; x++ {
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

