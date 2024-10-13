package main

import (
	"context"
	"fmt"
	"font-browser/finders"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Paths to font directories for Linux and macOS
var fontDirs = []string{
	"~/.local/share/fonts",
	"/usr/local/share/fonts",
	"/usr/share/fonts",
	"/Library/Fonts",
	"~/Library/Fonts",
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// Model represents the Bubble Tea model
type model struct {
	list     list.Model
	selected string
	cancel   context.CancelFunc
}

// Initialize the model
func setup() *model {
	ctx, cancel := context.WithCancel(context.Background())

	fontfinder := finders.NewFontFinder(fontDirs)
	fontItems := fontfinder.JayWalk(ctx)

	var items []list.Item

	for _, font := range fontItems {
		items = append(items, font)
	}

	// Bubble Tea list configuration
	l := list.New(items, list.NewDefaultDelegate(), 20, 32)
	l.Title = "Available Fonts"

	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)

	// l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &model{
		list:   l,
		cancel: cancel,
	}
}

// Update handles the incoming messages for Bubble Tea
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selectedItem, ok := m.list.SelectedItem().(finders.FontItem); ok {
				m.selected = selectedItem.Path()
				fmt.Println(m.selected)
			}
		case "q", "ctrl+c":
			m.cancel()
			return m, tea.Quit
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) Init() tea.Cmd {
	// No initial commands, so we return nil here
	return nil
}

// View renders the Bubble Tea view
func (m *model) View() string {
	return fmt.Sprint("\n", m.list.View())
}

type Game struct {
	model *model
}

// func (g *Game) Update() error {
// 	// Call the Update method from the Bubble Tea model
// 	_, cmd := g.model.Update(tea.KeyMsg{Type: tea.KeyRunes})
// 	if cmd != nil {
// 		return cmd()
// 	}
// 	return nil
// }
//
// func (g *Game) Draw(screen *ebiten.Image) {
// 	// Call the Draw method from the Bubble Tea model
// 	g.model.Draw(screen)
// }
//
// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return 800, 600 // Set the desired window size
// }

func main() {
	p := tea.NewProgram(setup())
	if _, err := p.Run(); err != nil {
		log.Fatalf("could not start program: %v", err)
	}
}
