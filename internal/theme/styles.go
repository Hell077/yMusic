package theme

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Base         lipgloss.Style
	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	Secondary    lipgloss.Style
	Muted        lipgloss.Style
	Primary      lipgloss.Style
	Accent       lipgloss.Style
	Error        lipgloss.Style
	Border       lipgloss.Style
	Selected     lipgloss.Style
	ListItem     lipgloss.Style
	ListActive   lipgloss.Style
	ListPlaying  lipgloss.Style
	SidebarItem  lipgloss.Style
	SidebarActive lipgloss.Style
	PlayerBar    lipgloss.Style
	ProgressFull lipgloss.Style
	ProgressEmpty lipgloss.Style
	Overlay      lipgloss.Style
	OverlayItem  lipgloss.Style
	OverlayActive lipgloss.Style
	HelpKey      lipgloss.Style
	HelpDesc     lipgloss.Style
	Tab          lipgloss.Style
	ActiveTab    lipgloss.Style
	StatusBar    lipgloss.Style
	SearchInput  lipgloss.Style
}

var S Styles

func init() {
	RefreshStyles()
}

func RefreshStyles() {
	t := Current
	S = Styles{
		Base: lipgloss.NewStyle().
			Foreground(t.Foreground).
			Background(t.Background),
		Title: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(t.Secondary),
		Secondary: lipgloss.NewStyle().
			Foreground(t.Secondary),
		Muted: lipgloss.NewStyle().
			Foreground(t.Muted),
		Primary: lipgloss.NewStyle().
			Foreground(t.Primary),
		Accent: lipgloss.NewStyle().
			Foreground(t.Accent),
		Error: lipgloss.NewStyle().
			Foreground(t.Error),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border),
		Selected: lipgloss.NewStyle().
			Background(t.Highlight).
			Foreground(t.Primary),
		ListItem: lipgloss.NewStyle().
			Foreground(t.Foreground).
			PaddingLeft(2),
		ListActive: lipgloss.NewStyle().
			Foreground(t.Primary).
			Background(t.Highlight).
			PaddingLeft(1).
			Bold(true),
		ListPlaying: lipgloss.NewStyle().
			Foreground(t.Primary).
			Background(t.Surface).
			PaddingLeft(2).
			Bold(true),
		SidebarItem: lipgloss.NewStyle().
			Foreground(t.Foreground).
			PaddingLeft(2).
			PaddingRight(2),
		SidebarActive: lipgloss.NewStyle().
			Foreground(t.Primary).
			Background(t.Highlight).
			PaddingLeft(1).
			PaddingRight(2).
			Bold(true),
		PlayerBar: lipgloss.NewStyle().
			Background(t.Surface).
			Foreground(t.Foreground).
			Padding(0, 1),
		ProgressFull: lipgloss.NewStyle().
			Foreground(t.Primary),
		ProgressEmpty: lipgloss.NewStyle().
			Foreground(t.Muted),
		Overlay: lipgloss.NewStyle().
			Background(t.Surface).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Padding(1, 2),
		OverlayItem: lipgloss.NewStyle().
			Foreground(t.Foreground).
			PaddingLeft(2),
		OverlayActive: lipgloss.NewStyle().
			Foreground(t.Primary).
			Background(t.Highlight).
			PaddingLeft(1).
			Bold(true),
		HelpKey: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),
		HelpDesc: lipgloss.NewStyle().
			Foreground(t.Muted),
		Tab: lipgloss.NewStyle().
			Foreground(t.Muted).
			Padding(0, 2),
		ActiveTab: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			Padding(0, 2).
			Underline(true),
		StatusBar: lipgloss.NewStyle().
			Background(t.Surface).
			Foreground(t.Muted).
			Padding(0, 1),
		SearchInput: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Padding(0, 1),
	}
}
