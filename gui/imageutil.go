package gui

import (
	"image"
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

