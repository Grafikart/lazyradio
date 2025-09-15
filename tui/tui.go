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

const (
	sidebarPanel = iota
	tracksPanel  = iota
)

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	mutedPanelStyle = panelStyle.
			BorderForeground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
	normalText = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})
	mutedText = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
)

var (
	sidebarWidth = 25
	footerHeight = 2
)

type Model struct {
	footer    footer
	trackList trackListModel
	sidebar   radioListModel

	// State
	player      *radio.Player
	tracks      []radio.TrackItem
	width       int
	height      int
	err         error
	renderCount int
	panel       int
}

func NewModel() Model {
	p := radio.NewPlayer()
	return Model{
		player:    p,
		trackList: newTrackList(p),
		footer:    newFooter(p),
		sidebar:   newSidebar(),
		panel:     sidebarPanel,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.trackList.Init(),
		m.footer.Init(),
		m.sidebar.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.renderCount++
	var cmdList, cmdFooter tea.Cmd
	switch msg := msg.(type) {

	// Keypress
	case tea.KeyMsg:
		switch msg.Type {
		// Quit when pressing CTRL C or Esc
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyLeft, tea.KeyRight:
			m.togglePanel()
		}
		switch msg.String() {
		case "p":
			m.player.Pause()
		}

	// Resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
	case radioSelectedMsg:
		if m.panel == sidebarPanel {
			m.togglePanel()
		}
		m.trackList, cmdList = m.trackList.Update(msg)
		return m, cmdList
	}
	m.footer, cmdFooter = m.footer.Update(msg)
	if m.panel == tracksPanel {
		m.trackList, cmdList = m.trackList.Update(msg)
	} else {
		m.sidebar, cmdList = m.sidebar.Update(msg)
	}
	return m, tea.Batch(cmdList, cmdFooter)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.sidebar.View(m.panel == sidebarPanel),
			m.trackList.View(m.panel == tracksPanel),
		),
		m.footer.View(),
	)
}

func (m *Model) updateSizes() {
	contentWidth := m.width - sidebarWidth - 4
	contentHeight := m.height - footerHeight - 4
	m.footer.SetSize(m.width-2, footerHeight)
	m.sidebar.SetSize(sidebarWidth, contentHeight)
	m.trackList.SetSize(contentWidth, contentHeight)
}

func (m *Model) togglePanel() {
	if m.panel == tracksPanel {
		m.panel = sidebarPanel
	} else {
		m.panel = tracksPanel
	}
}

func openYtMusic(item list.Item) {
	u := fmt.Sprintf(
		"https://music.youtube.com/search?q=%s", url.QueryEscape(item.FilterValue()),
	)
	utils.OpenBrowser(u)
}
