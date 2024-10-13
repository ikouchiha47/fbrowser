package render

import (
	"log"
	"os"

	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os/exec"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type ImageRender struct{}

func (ir *ImageRender) Render(fontPath string) {
	img := image.NewRGBA(image.Rect(0, 0, 800, 200))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Load the font
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("could not read font file: %v", err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("could not parse font: %v", err)
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.Point26_6{fixed.I(50), fixed.I(100)},
	}
	drawer.DrawString("The quick brown fox jumps over the lazy dog")

	// Save the image as PNG
	file, err := os.Create("/tmp/font_preview.png")
	if err != nil {
		log.Fatalf("could not create image: %v", err)
	}
	defer file.Close()
	png.Encode(file, img)

	// Open the image using default viewer
	// exec.Command("open", "/tmp/font_preview.png").Run()     // For macOS
	exec.Command("xdg-open", "/tmp/font_preview.png").Run() // For Linux
}
