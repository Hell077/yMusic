package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name       string
	Background lipgloss.Color
	Foreground lipgloss.Color
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Muted      lipgloss.Color
	Error      lipgloss.Color
	Border     lipgloss.Color
	Highlight  lipgloss.Color
	Surface    lipgloss.Color
}

var (
	Dark = Theme{
		Name:       "dark",
		Background: lipgloss.Color("#1a1a2e"),
		Foreground: lipgloss.Color("#e0e0e0"),
		Primary:    lipgloss.Color("#ffcc00"),
		Secondary:  lipgloss.Color("#a0a0ff"),
		Accent:     lipgloss.Color("#ff6b6b"),
		Muted:      lipgloss.Color("#666680"),
		Error:      lipgloss.Color("#ff4444"),
		Border:     lipgloss.Color("#333355"),
		Highlight:  lipgloss.Color("#2a2a4a"),
		Surface:    lipgloss.Color("#16213e"),
	}

	Light = Theme{
		Name:       "light",
		Background: lipgloss.Color("#fafafa"),
		Foreground: lipgloss.Color("#1a1a1a"),
		Primary:    lipgloss.Color("#d4a017"),
		Secondary:  lipgloss.Color("#4040c0"),
		Accent:     lipgloss.Color("#cc3333"),
		Muted:      lipgloss.Color("#999999"),
		Error:      lipgloss.Color("#cc0000"),
		Border:     lipgloss.Color("#cccccc"),
		Highlight:  lipgloss.Color("#e8e8f0"),
		Surface:    lipgloss.Color("#f0f0f5"),
	}

	Solarized = Theme{
		Name:       "solarized",
		Background: lipgloss.Color("#002b36"),
		Foreground: lipgloss.Color("#839496"),
		Primary:    lipgloss.Color("#b58900"),
		Secondary:  lipgloss.Color("#268bd2"),
		Accent:     lipgloss.Color("#dc322f"),
		Muted:      lipgloss.Color("#586e75"),
		Error:      lipgloss.Color("#dc322f"),
		Border:     lipgloss.Color("#073642"),
		Highlight:  lipgloss.Color("#073642"),
		Surface:    lipgloss.Color("#073642"),
	}

	Nord = Theme{
		Name:       "nord",
		Background: lipgloss.Color("#2e3440"),
		Foreground: lipgloss.Color("#d8dee9"),
		Primary:    lipgloss.Color("#ebcb8b"),
		Secondary:  lipgloss.Color("#81a1c1"),
		Accent:     lipgloss.Color("#bf616a"),
		Muted:      lipgloss.Color("#4c566a"),
		Error:      lipgloss.Color("#bf616a"),
		Border:     lipgloss.Color("#3b4252"),
		Highlight:  lipgloss.Color("#3b4252"),
		Surface:    lipgloss.Color("#3b4252"),
	}

	Themes  = []Theme{Dark, Light, Solarized, Nord}
	Current = Dark
)

func SetTheme(name string) {
	for _, t := range Themes {
		if t.Name == name {
			Current = t
			RefreshStyles()
			return
		}
	}
	Current = Dark
	RefreshStyles()
}

func SetThemeByIndex(i int) {
	if i >= 0 && i < len(Themes) {
		Current = Themes[i]
		RefreshStyles()
	}
}
