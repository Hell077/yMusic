package main

import (
	"fmt"
	"os"

	"ymusic/internal/api"
	"ymusic/internal/config"
	"ymusic/internal/player"
	"ymusic/internal/theme"
	"ymusic/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Handle --logout
	for _, arg := range os.Args[1:] {
		if arg == "--logout" || arg == "-logout" {
			cfg, _ := config.Load()
			if cfg != nil {
				cfg.Token = ""
				cfg.Save()
			}
			fmt.Println("Logged out. Token cleared.")
			return
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	theme.SetTheme(cfg.Theme)

	var client *api.Client
	if cfg.Token != "" {
		client = api.NewClient(cfg.Token)
	}

	ctrl := player.NewController(float64(cfg.Volume))

	root := ui.NewRoot(cfg, client, ctrl)

	p := tea.NewProgram(root,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ctrl.Quit()
}
