package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chatMsg struct {
	text    string
	sent    bool   // true = user, false = bot
	time    string
	actions []BotAction
}

type model struct {
	baseURL  string
	from     string
	channel  string
	channels []string
	chanIdx  int

	messages []chatMsg
	input    textinput.Model
	viewport viewport.Model

	// action selection
	actions     []BotAction
	actionIdx   int
	actionMode  bool

	loading bool
	status  string
	width   int
	height  int
	ready   bool
}

// ── Tea messages ───────────────────────────────────────────────────────────

type botRespMsg struct{ r *BotReply }
type errMsg struct{ err error }
type initMsg struct{ r *BotReply }

// ── Init ───────────────────────────────────────────────────────────────────

func newModel(baseURL, from string) model {
	ti := textinput.New()
	ti.Placeholder = "Message…"
	ti.Focus()
	ti.CharLimit = 500

	return model{
		baseURL:  baseURL,
		from:     from,
		channel:  "whatsapp",
		channels: []string{"whatsapp", "imessage", "sms"},
		chanIdx:  0,
		input:    ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		fetchInit(m.baseURL, m.from, m.channel),
	)
}

func fetchInit(baseURL, from, channel string) tea.Cmd {
	return func() tea.Msg {
		r, err := sendMessage(baseURL, from, "", channel)
		if err != nil {
			return errMsg{err}
		}
		return initMsg{r}
	}
}

// ── Update ─────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerH := 3
		statusH := 1
		inputH := 3
		vpHeight := m.height - headerH - statusH - inputH
		if vpHeight < 4 {
			vpHeight = 4
		}
		if !m.ready {
			m.viewport = viewport.New(m.width, vpHeight)
			m.ready = true
		} else {
			m.viewport.Width = m.width
			m.viewport.Height = vpHeight
		}
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case initMsg:
		m.loading = false
		m.addBotMsg(msg.r)

	case botRespMsg:
		m.loading = false
		m.status = ""
		m.addBotMsg(msg.r)

	case errMsg:
		m.loading = false
		m.status = "Error: " + msg.err.Error()
		m.messages = append(m.messages, chatMsg{
			text: "⚠️  " + msg.err.Error(),
			sent: false,
			time: nowStr(),
		})
		m.refreshViewport()

	case tea.KeyMsg:
		// Action mode: arrow keys + enter navigate chips.
		if m.actionMode && len(m.actions) > 0 {
			switch msg.String() {
			case "up", "left", "shift+tab":
				if m.actionIdx > 0 {
					m.actionIdx--
				}
			case "down", "right", "tab":
				if m.actionIdx < len(m.actions)-1 {
					m.actionIdx++
				}
			case "enter":
				a := m.actions[m.actionIdx]
				return m, m.sendAction(a)
			case "esc":
				m.actionMode = false
				m.input.Focus()
			// Number shortcuts 1–9
			default:
				for i := range m.actions {
					if msg.String() == fmt.Sprintf("%d", i+1) {
						return m, m.sendAction(m.actions[i])
					}
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		case "tab":
			// Switch to action mode if actions are available.
			if len(m.actions) > 0 {
				m.actionMode = true
				m.actionIdx = 0
				m.input.Blur()
			}

		case "ctrl+left", "ctrl+h":
			// Cycle channel left.
			m.chanIdx = (m.chanIdx - 1 + len(m.channels)) % len(m.channels)
			m.channel = m.channels[m.chanIdx]

		case "ctrl+right", "ctrl+l":
			// Cycle channel right.
			m.chanIdx = (m.chanIdx + 1) % len(m.channels)
			m.channel = m.channels[m.chanIdx]

		case "enter":
			if m.loading {
				return m, nil
			}
			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				// Enter with empty input and actions → enter action mode.
				if len(m.actions) > 0 {
					m.actionMode = true
					m.actionIdx = 0
					m.input.Blur()
				}
				return m, nil
			}
			m.input.SetValue("")
			return m, m.sendText(text)
		}
	}

	if m.ready {
		var vpCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		cmds = append(cmds, vpCmd)
	}

	var tiCmd tea.Cmd
	m.input, tiCmd = m.input.Update(msg)
	cmds = append(cmds, tiCmd)

	return m, tea.Batch(cmds...)
}

func (m *model) sendText(text string) tea.Cmd {
	m.messages = append(m.messages, chatMsg{text: text, sent: true, time: nowStr()})
	m.actions = nil
	m.actionMode = false
	m.loading = true
	m.status = "Sending…"
	m.refreshViewport()
	from, baseURL, channel := m.from, m.baseURL, m.channel
	return func() tea.Msg {
		r, err := sendMessage(baseURL, from, text, channel)
		if err != nil {
			return errMsg{err}
		}
		return botRespMsg{r}
	}
}

func (m *model) sendAction(a BotAction) tea.Cmd {
	m.messages = append(m.messages, chatMsg{text: a.Label, sent: true, time: nowStr()})
	m.actions = nil
	m.actionMode = false
	m.loading = true
	m.status = "Sending…"
	m.refreshViewport()
	from, baseURL, channel, val := m.from, m.baseURL, m.channel, a.Value
	return func() tea.Msg {
		r, err := sendMessage(baseURL, from, val, channel)
		if err != nil {
			return errMsg{err}
		}
		return botRespMsg{r}
	}
}

func (m *model) addBotMsg(r *BotReply) {
	if r.Reply != "" {
		m.messages = append(m.messages, chatMsg{
			text:    r.Reply,
			sent:    false,
			time:    nowStr(),
			actions: r.Actions,
		})
	}
	m.actions = r.Actions
	m.actionIdx = 0
	m.actionMode = len(r.Actions) > 0
	if m.actionMode {
		m.input.Blur()
	} else {
		m.input.Focus()
	}
	m.refreshViewport()
}

func (m *model) refreshViewport() {
	if !m.ready {
		return
	}
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom()
}

// ── View ───────────────────────────────────────────────────────────────────

func (m model) View() string {
	if !m.ready {
		return "Loading…\n"
	}

	// Header
	channels := make([]string, len(m.channels))
	for i, ch := range m.channels {
		s := strings.ToUpper(ch[:1]) + ch[1:]
		if i == m.chanIdx {
			channels[i] = channelActiveStyle.Render(s)
		} else {
			channels[i] = channelStyle.Render(s)
		}
	}
	header := headerStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			lipgloss.NewStyle().Foreground(colAccent).Bold(true).Render("🏢 TowerShare"),
			"  ",
			strings.Join(channels, " "),
		),
	)

	// Messages viewport
	vp := m.viewport.View()

	// Action chips
	actionBar := ""
	if len(m.actions) > 0 {
		chips := make([]string, len(m.actions))
		for i, a := range m.actions {
			label := fmt.Sprintf("%d. %s", i+1, a.Label)
			if m.actionMode && i == m.actionIdx {
				chips[i] = chipSelectedStyle.Render(label)
			} else {
				chips[i] = chipStyle.Render(label)
			}
		}
		hint := dotStyle.Render("  ← Tab to select, Enter to confirm, Esc to type")
		actionBar = "\n" + lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().PaddingLeft(1).Render(strings.Join(chips, " ")),
		) + hint
	}

	// Input bar
	var inputContent string
	if m.actionMode {
		inputContent = statusStyle.Render("Use ↑↓ or 1-9 to pick an option, Enter to send, Esc to type")
	} else {
		inputContent = m.input.View()
	}
	inputBar := inputBarStyle.Width(m.width).Render(inputContent)

	// Status
	status := ""
	if m.loading {
		status = statusStyle.Render("● thinking…")
	} else if m.status != "" {
		status = statusStyle.Render(m.status)
	} else {
		status = statusStyle.Render("Ctrl+C to quit  •  Tab for actions  •  Ctrl+←/→ to change channel")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		vp,
		actionBar,
		inputBar,
		status,
	)
}

func (m model) renderMessages() string {
	if len(m.messages) == 0 {
		return statusStyle.Render("Connecting…")
	}

	width := m.width
	if width < 20 {
		width = 80
	}
	maxBubble := width * 2 / 3

	var lines []string
	for _, msg := range m.messages {
		text := wordWrap(msg.text, maxBubble-4)
		if msg.sent {
			bubble := sentStyle.MaxWidth(maxBubble).Render(text)
			// Right-align
			bw := lipgloss.Width(bubble)
			pad := width - bw - 2
			if pad < 0 {
				pad = 0
			}
			lines = append(lines, strings.Repeat(" ", pad)+bubble)
		} else {
			label := senderStyle.Render("Bot · " + msg.time)
			bubble := receivedStyle.MaxWidth(maxBubble).Render(text)
			lines = append(lines, label)
			lines = append(lines, "  "+bubble)
		}
		lines = append(lines, "") // spacer
	}
	return strings.Join(lines, "\n")
}

// ── Helpers ────────────────────────────────────────────────────────────────

func nowStr() string {
	return time.Now().Format("15:04")
}

func wordWrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var result strings.Builder
	for _, para := range strings.Split(s, "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			result.WriteString("\n")
			continue
		}
		line := words[0]
		for _, w := range words[1:] {
			if len(line)+1+len(w) > width {
				result.WriteString(line + "\n")
				line = w
			} else {
				line += " " + w
			}
		}
		result.WriteString(line + "\n")
	}
	return strings.TrimRight(result.String(), "\n")
}
