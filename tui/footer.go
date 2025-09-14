package tui

import (
	"fmt"
	"grafikart/lazyradio/radio"

	"github.com/charmbracelet/bubbles/list"
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
		m.spinner.Tick,
		m.listenCmd,
	)
}

func (m footer) Update(msg tea.Msg) (footer, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case radio.PlayerProgressMsg:
		if msg.Progress == 100 {
			m.info = "Ended"
			return m, tea.Batch(
				m.listenCmd,
				m.endCmd,
			)
		}
		return m, m.listenCmd
	case radio.PlayerMsg:
		return m, m.listenCmd
	}

	return m, nil
}

func (m footer) View() string {
	style := panelStyle.
		Width(m.width).
		Height(m.height)
	var content string

	if m.player.State() == radio.Loading {
		content += m.spinner.View()
	}

	content += fmt.Sprintf("%s / %s - %d", m.player.Info().Current, m.player.Info().Duration, m.player.State())

	return style.
		Render(
			content,
			m.prog.ViewAs(float64(m.player.Info().Progress)/100),
		)
}

func (m *footer) listenCmd() tea.Msg {
	return <-m.player.Ch()
}

func (m *footer) endCmd() tea.Msg {
	return endMsg{}
}

func (m *footer) play(item list.Item) {
	m.player.Play(item.FilterValue())
}

func (m *footer) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.prog.Width = w
}
