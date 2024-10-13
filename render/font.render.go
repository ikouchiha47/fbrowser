package render

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type GameRender struct {
	ttfFont font.Face
}

func (gr *GameRender) Render(fontPath string) font.Face {
	ttfData, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the font
	ttfParsed, err := opentype.Parse(ttfData)
	if err != nil {
		log.Fatal(err)
	}

	// Set the font size
	const dpi = 72
	ttfFont, err := opentype.NewFace(ttfParsed, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return ttfFont
}

// Implement Ebiten's Update function (empty for now)
func (gr *GameRender) Update() error {
	return nil
}

// Implement Ebiten's Draw function to render the text
func (gr *GameRender) Draw(screen *ebiten.Image) {
	// Set the color for the text
	textColor := colornames.White
	ttfFont := gr.ttfFont

	// Draw the "The quick brown fox" with the loaded font
	text.Draw(screen, "The quick brown fox jumps over the lazy dog", ttfFont, 20, 100, textColor)
}

// func main() {
// 	ebiten.SetWindowSize(640, 480)
// 	ebiten.SetWindowTitle("Font Preview")
//
// 	// Start the game loop
// 	if err := ebiten.RunGame(&Game{}); err != nil {
// 		log.Fatal(err)
// 	}
// }
