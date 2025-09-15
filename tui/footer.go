package tui

import (
	"fmt"
	"grafikart/lazyradio/radio"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type endMsg struct{}

type footer struct {
	// gui
	spinner spinner.Model
	prog    progress.Model

	// state
	player *radio.Player
	info   string
	err    error
	width  int
	height int
}

func newFooter(player *radio.Player) footer {
	spin := spinner.New()
	spin.Spinner = spinner.MiniDot
	prog := progress.New(progress.WithDefaultGradient())
	prog.ShowPercentage = false
	return footer{
		player:  player,
		spinner: spin,
		prog:    prog,
	}
}

func (m footer) Init() tea.Cmd {
	return tea.Batch(
		m.listenCmd,
	)
}

func (m footer) Update(msg tea.Msg) (footer, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		if m.player.State() == radio.Loading {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	case radio.PlayerOutputMsg:
		m.info = string(msg)
		return m, m.listenCmd
	case radio.PlayerErrorMsg:
		m.info = fmt.Sprintf("%v", msg)
		m.info = fmt.Sprintf("%v", msg)
		m.info = fmt.Sprintf("%v", msg)
		m.info = fmt.Sprintf("%v", msg)
		fmt.Printf("Error: %v\n", msg)
		return m, m.listenCmd
	case radio.PlayerProgressMsg:
		m.info = ""
		if msg.Progress == 100 {
			return m, tea.Batch(
				m.listenCmd,
				endCmd,
			)
		}
		return m, m.listenCmd
	case radio.PlayerStateChangedMsg:
		if msg == radio.Loading {
			return m, m.spinner.Tick
		}
	case radio.PlayerMsg:
		return m, m.listenCmd
	}

	return m, nil
}

func (m footer) View() string {
	style := panelStyle.
		Padding(0, 1).
		Width(m.width).
		Height(m.height)
	var content string

	if m.player.State() == radio.Loading {
		content += m.spinner.View() + " "
	}

	if m.player.Info().Current != "" {
		content += fmt.Sprintf("%s/%s", m.player.Info().Current, m.player.Info().Duration)
	}

	if m.info != "" {
		content += " - " + m.info
	}

	return style.
		Render(
			mutedText.Render(content),
			m.prog.ViewAs(float64(m.player.Info().Progress)/100),
		)
}

func (m *footer) listenCmd() tea.Msg {
	return <-m.player.Ch()
}

func endCmd() tea.Msg {
	return endMsg{}
}

func (m *footer) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.prog.Width = w - 2
}
