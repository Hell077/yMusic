# ymusic

Terminal UI client for Yandex Music.

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-lightgrey)
![License](https://img.shields.io/badge/License-MIT-green)

## Features

- Browse home feed, playlists, albums, artists
- Full-text search with tabbed results (tracks, albums, artists)
- Collection: liked tracks, playlists, liked albums
- My Wave radio with auto-advancement
- Playback via mpv: play/pause, seek, next/prev, volume
- Queue with shuffle and repeat modes
- 4 color themes: Dark, Light, Solarized, Nord
- Keyboard-driven navigation with vim-style keys
- ESC as back navigation + overlay menu

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Hell077/yMusic/main/install.sh | sh
```

This will download the latest release for your platform and install the `ymusic` binary.

### Requirements

- **mpv** (audio backend)

```bash
# Fedora
sudo dnf install mpv

# Ubuntu / Debian
sudo apt install mpv

# macOS
brew install mpv

# Arch
sudo pacman -S mpv
```

### Build from source

```bash
git clone https://github.com/Hell077/yMusic.git
cd yMusic
go build -o ymusic .
```

## Usage

```bash
./ymusic
```

On first launch you'll be prompted to authenticate with Yandex. A browser window will open — log in and paste the authorization code back into the terminal.

```bash
# Clear saved token
./ymusic --logout
```

## Keyboard Shortcuts

| Key | Action |
|---|---|
| `space` | Play / Pause |
| `n` / `p` | Next / Previous track |
| `+` / `-` | Volume up / down |
| `>` / `<` | Seek forward / back 10s |
| `up` / `down` / `j` / `k` | Navigate lists |
| `left` / `right` / `h` / `l` | Switch tabs |
| `enter` | Select / Play |
| `tab` | Switch focus (sidebar / content) |
| `/` | Search |
| `L` | Like track |
| `s` | Toggle shuffle |
| `r` | Cycle repeat (off → all → one) |
| `esc` | Back / Menu |
| `q` | Quit |

## Config

Stored at `~/.config/ymusic/config.json`:

```json
{
  "token": "...",
  "theme": "dark",
  "volume": 70
}
```

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [mpv](https://mpv.io) — audio playback via JSON IPC
