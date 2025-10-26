package gui

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

var (
	imageCache = make(map[string]*canvas.Image)
	cacheMutex sync.RWMutex
)

// LoadImageFromURL loads an image from a URL and caches it
func LoadImageFromURL(url string, width, height float32) *canvas.Image {
	if url == "" {
		return createPlaceholderImage(width, height)
	}

	// Check cache first
	cacheMutex.RLock()
	if cachedImg, ok := imageCache[url]; ok {
		cacheMutex.RUnlock()
		// Create a new image instance with the same resource
		img := canvas.NewImageFromImage(cachedImg.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(width, height))
		return img
	}
	cacheMutex.RUnlock()

	// Create placeholder
	placeholder := createPlaceholderImage(width, height)

	// Load image asynchronously
	go func() {
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		
		resp, err := client.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		imgData, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		img, _, err := image.Decode(io.NopCloser(io.Reader(newBytesReader(imgData))))
		if err != nil {
			return
		}

		// Cache the image
		cacheMutex.Lock()
		canvasImg := canvas.NewImageFromImage(img)
		canvasImg.FillMode = canvas.ImageFillContain
		imageCache[url] = canvasImg
		cacheMutex.Unlock()

		// Update placeholder
		fyne.Do(func() {
			placeholder.Image = img
			placeholder.Refresh()
		})
	}()

	return placeholder
}

// bytesReader wraps []byte to implement io.Reader
type bytesReader struct {
	data []byte
	pos  int
}

func newBytesReader(data []byte) *bytesReader {
	return &bytesReader{data: data}
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// createPlaceholderImage creates a placeholder image while loading
func createPlaceholderImage(width, height float32) *canvas.Image {
	img := canvas.NewImageFromImage(createPlaceholderRect())
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}

// createPlaceholderRect creates a simple colored rectangle for placeholder
func createPlaceholderRect() image.Image {
	width, height := 300, 450
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill with card background color
	bgColor := GetCardColor()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}
	
	return img
}

// ClearImageCache clears the image cache (useful for theme changes)
func ClearImageCache() {
	cacheMutex.Lock()
	imageCache = make(map[string]*canvas.Image)
	cacheMutex.Unlock()
}

// LoadBackdropImage loads a backdrop image with blur and dimming effect
func LoadBackdropImage(url string, width, height float32) *canvas.Image {
	if url == "" {
		return createPlaceholderBackdrop(width, height)
	}

	// Check cache first (with special key for backdrop)
	cacheKey := "backdrop_" + url
	cacheMutex.RLock()
	if cachedImg, ok := imageCache[cacheKey]; ok {
		cacheMutex.RUnlock()
		img := canvas.NewImageFromImage(cachedImg.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(width, height))
		return img
	}
	cacheMutex.RUnlock()

	// Create placeholder
	placeholder := createPlaceholderBackdrop(width, height)

	// Load image asynchronously
	go func() {
		client := &http.Client{
			Timeout: 15 * time.Second,
		}
		
		resp, err := client.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		imgData, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		img, _, err := image.Decode(io.NopCloser(io.Reader(newBytesReader(imgData))))
		if err != nil {
			return
		}

		// Apply blur and dim effect
		processedImg := applyBackdropEffect(img)

		// Cache the processed image
		cacheMutex.Lock()
		canvasImg := canvas.NewImageFromImage(processedImg)
		canvasImg.FillMode = canvas.ImageFillContain
		imageCache[cacheKey] = canvasImg
		cacheMutex.Unlock()

		// Update placeholder
		fyne.Do(func() {
			placeholder.Image = processedImg
			placeholder.Refresh()
		})
	}()

	return placeholder
}

// applyBackdropEffect applies blur and dimming to create a backdrop effect
func applyBackdropEffect(img image.Image) image.Image {
	bounds := img.Bounds()
	dimmed := image.NewRGBA(bounds)
	
	// Apply dimming effect (reduce brightness by overlaying semi-transparent black)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			
			// Reduce brightness by 60% and desaturate slightly
			r = uint32(float64(r) * 0.4)
			g = uint32(float64(g) * 0.4)
			b = uint32(float64(b) * 0.4)
			
			dimmed.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}
	
	// Simple box blur for performance
	blurred := boxBlur(dimmed, 8)
	
	return blurred
}

// boxBlur applies a simple box blur to the image
func boxBlur(img *image.RGBA, radius int) *image.RGBA {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a uint32
			count := 0
			
			// Sample pixels in the blur radius
			for dy := -radius; dy <= radius; dy += 2 {
				for dx := -radius; dx <= radius; dx += 2 {
					nx, ny := x+dx, y+dy
					if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
						c := img.At(nx, ny)
						pr, pg, pb, pa := c.RGBA()
						r += pr
						g += pg
						b += pb
						a += pa
						count++
					}
				}
			}
			
			if count > 0 {
				result.Set(x, y, color.RGBA{
					R: uint8((r / uint32(count)) >> 8),
					G: uint8((g / uint32(count)) >> 8),
					B: uint8((b / uint32(count)) >> 8),
					A: uint8((a / uint32(count)) >> 8),
				})
			}
		}
	}
	
	return result
}

// createPlaceholderBackdrop creates a dark placeholder for backdrop
func createPlaceholderBackdrop(width, height float32) *canvas.Image {
	bgImg := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	bgColor := GetBackgroundColor()
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
	
	img := canvas.NewImageFromImage(bgImg)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}

