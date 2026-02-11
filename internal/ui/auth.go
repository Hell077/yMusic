package ui

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"ymusic/internal/config"
	"ymusic/internal/theme"
)

type AuthModel struct {
	input      textinput.Model
	spinner    spinner.Model
	err        error
	exchanging bool
	width      int
	height     int
}

func NewAuth() AuthModel {
	ti := textinput.New()
	ti.Placeholder = "Вставьте URL или код..."
	ti.CharLimit = 2000
	ti.Width = 60
	ti.Focus()
	s := spinner.New()
	s.Spinner = spinner.Dot
	return AuthModel{input: ti, spinner: s}
}

func (m AuthModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick, openBrowserCmd)
}

func openBrowserCmd() tea.Msg {
	openBrowser(config.AuthURL)
	return browserOpenedMsg{}
}

type browserOpenedMsg struct{}
type codeExchangeMsg struct {
	token string
	err   error
}

func (m AuthModel) Update(msg tea.Msg) (AuthModel, tea.Cmd) {
	switch msg := msg.(type) {
	case browserOpenedMsg:
		return m, nil

	case codeExchangeMsg:
		m.exchanging = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		return m, func() tea.Msg {
			return AuthCompleteMsg{Token: msg.token}
		}

	case AuthErrorMsg:
		m.err = msg.Err
		m.exchanging = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.exchanging {
			return m, nil
		}
		switch msg.String() {
		case "enter":
			raw := strings.TrimSpace(m.input.Value())
			if raw == "" {
				return m, nil
			}
			code := extractCode(raw)
			if code == "" {
				m.err = fmt.Errorf("не найден код авторизации")
				return m, nil
			}
			m.err = nil
			m.exchanging = true
			return m, exchangeCode(code)
		case "o":
			if m.input.Value() == "" {
				openBrowser(config.AuthURL)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func exchangeCode(code string) tea.Cmd {
	return func() tea.Msg {
		token, err := config.ExchangeCode(code)
		if err != nil {
			return codeExchangeMsg{err: err}
		}
		return codeExchangeMsg{token: token}
	}
}

// extractCode pulls the authorization code from a URL like
// https://music.yandex.ru/?code=XXXXXX or accepts a raw code string.
func extractCode(raw string) string {
	raw = strings.TrimSpace(raw)

	// Try parsing as URL with ?code= parameter
	if strings.Contains(raw, "code=") {
		if u, err := url.Parse(raw); err == nil {
			if c := u.Query().Get("code"); c != "" {
				return c
			}
		}
		// Brute-force
		for _, part := range strings.Split(raw, "&") {
			part = strings.TrimLeft(part, "?")
			if strings.HasPrefix(part, "code=") {
				return strings.TrimPrefix(part, "code=")
			}
		}
	}

	// Raw code — typically 7 digits
	if len(raw) >= 6 && len(raw) <= 20 && !strings.Contains(raw, " ") && !strings.Contains(raw, "/") {
		return raw
	}

	return ""
}

func (m AuthModel) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(theme.S.Title.Render("  ymusic — авторизация") + "\n\n")

	if m.err != nil {
		b.WriteString(theme.S.Error.Render("  Ошибка: "+m.err.Error()) + "\n\n")
	}

	if m.exchanging {
		b.WriteString("  " + m.spinner.View() + " Получаем токен...\n")
		return b.String()
	}

	b.WriteString("  1. Браузер уже открыт — залогиньтесь в Яндекс\n")
	b.WriteString("  2. После авторизации вас перекинет на страницу\n")
	b.WriteString("     с " + theme.S.Primary.Render("?code=XXXXXXX") + " в адресной строке\n")
	b.WriteString("  3. Скопируйте " + theme.S.Primary.Render("код") +
		" или весь URL и вставьте сюда:\n\n")
	b.WriteString("     " + m.input.View() + "\n\n")
	b.WriteString(theme.S.Muted.Render("  [o] открыть браузер  [Enter] подтвердить  [q] выход") + "\n")

	return b.String()
}

func (m *AuthModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.input.Width = w - 12
	if m.input.Width > 80 {
		m.input.Width = 80
	}
}

func openBrowser(rawURL string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	default:
		return
	}
	cmd.Start()
}
