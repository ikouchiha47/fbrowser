package finders

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FontFinder struct {
	locations []string
}

type FontItem struct {
	name string
	path string
}

type FontItems []FontItem

// Len is the number of elements in the collection.
func (fontitems FontItems) Len() int {
	return len(fontitems)
}

func (fontitems FontItems) Less(i int, j int) bool {
	return fontitems[i].name < fontitems[j].name
}

func (fontitems FontItems) Swap(i int, j int) {
	fontitems[i], fontitems[j] = fontitems[j], fontitems[i]
}

func (f FontItem) Title() string { return f.name }

func (f FontItem) Path() string { return f.path }

func (f FontItem) Description() string { return f.path }

func (f FontItem) FilterValue() string { return f.name }

func NewFontFinder(locations []string) *FontFinder {
	return &FontFinder{locations: locations}
}

func (ff *FontFinder) JayWalk(ctx context.Context) []FontItem {
	fontChan := make(chan FontItem)
	doneChan := make(chan struct{})
	fontDirs := ff.locations

	var fonts FontItems = []FontItem{}

	go func() {
		for font := range fontChan {
			fonts = append(fonts, font)
		}
		doneChan <- struct{}{}
	}()

	for _, dir := range fontDirs {
		go func(d string) {
			searchFontsInDir(ctx, expandHome(d), fontChan)
			doneChan <- struct{}{}
		}(dir)
	}

	for range fontDirs {
		<-doneChan
	}

	close(fontChan)
	<-doneChan

	sort.Sort(fonts)
	return fonts
}

func searchFontsInDir(ctx context.Context, dir string, fontChan chan<- FontItem) {
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			// Continue
		}

		if err != nil {
			return nil // Ignore errors and continue
		}

		if !d.IsDir() && (strings.HasSuffix(path, ".ttf") || strings.HasSuffix(path, ".otf")) {
			fontChan <- FontItem{
				name: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
				path: path,
			}
		}
		return nil
	})
}

// Expands the tilde (~) to the user home directory
func expandHome(path string) string {
	if path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
