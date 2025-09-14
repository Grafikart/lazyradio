package tui

import (
	"fmt"
	"grafikart/lazyradio/radio"
	"grafikart/lazyradio/utils"
	"net/url"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
)

var (
	sidebarWidth = 25
	footerHeight = 2
)

type Model struct {
	footer    footer
	trackList trackList

	// State
	tracks []radio.RadioItem
	width  int
	height int
	err    error
}

func NewModel() Model {
	p := radio.NewPlayer()
	return Model{
		trackList: newTrackList(p),
		footer:    newFooter(p),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.trackList.Init(),
		m.footer.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Keypress
	case tea.KeyMsg:
		switch {
		// Quit when pressing CTRL C or Esc
		case msg.Type == tea.KeyCtrlC, msg.Type == tea.KeyEsc:
			return m, tea.Quit
		}

	// Resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
	}

	var cmdList, cmdFooter tea.Cmd
	m.trackList, cmdList = m.trackList.Update(msg)
	m.footer, cmdFooter = m.footer.Update(msg)
	return m, tea.Batch(cmdList, cmdFooter)
}

func (m Model) View() string {
	contentHeight := m.height - footerHeight - 4
	sidebar := panelStyle.
		Width(sidebarWidth).
		Height(contentHeight).
		Render("Sidebar")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			sidebar,
			m.trackList.View(),
		),
		m.footer.View(),
	)
}

func (m *Model) updateSizes() {
	contentWidth := m.width - sidebarWidth - 4
	contentHeight := m.height - footerHeight - 4
	m.footer.SetSize(m.width-2, footerHeight)
	m.trackList.SetSize(contentWidth, contentHeight)
}

func openYtMusic(item list.Item) {
	u := fmt.Sprintf(
		"https://music.youtube.com/search?q=%s", url.QueryEscape(item.FilterValue()),
	)
	utils.OpenBrowser(u)
}
