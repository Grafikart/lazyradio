package tui

import (
	"fmt"
	"grafikart/lazyradio/radio"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type trackList struct {
	// components
	list list.Model

	// state
	item   radio.RadioItem
	width  int
	height int
	msg    string
	player *radio.Player
}

type tracklistKeyMap struct {
	openBrowser key.Binding
	refresh     key.Binding
}

var keys = tracklistKeyMap{
	openBrowser: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open yt music"),
	),
	refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
}

func newTrackList(p *radio.Player) trackList {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Track List"
	l.SetShowStatusBar(false)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.openBrowser,
		}
	}
	return trackList{
		list:   l,
		player: p,
	}
}

func (m trackList) Init() tea.Cmd {
	return fetchTrackCmd
}

type tracksLoadedMsg []list.Item
type nextTrackMsg struct{}

func (m trackList) Update(msg tea.Msg) (trackList, tea.Cmd) {
	switch msg := msg.(type) {

	// Keypress
	case tea.KeyMsg:
		switch {
		// Open a track in the browser
		case key.Matches(msg, keys.openBrowser):
			openYtMusic(m.list.SelectedItem())
			return m, nil
		// Play the track
		case msg.Type == tea.KeyEnter:
			item := m.list.SelectedItem()
			m.player.Play(item.FilterValue())
			return m, nil
		// Start playing a track
		case key.Matches(msg, keys.refresh):
			return m, fetchTrackCmd
		}
	case tracksLoadedMsg:
		item := m.list.SelectedItem()
		if item == nil {
			return m, m.list.SetItems(msg)
		}
		_, index, _ := lo.FindIndexOf(msg, func(l list.Item) bool {
			return l.FilterValue() == item.FilterValue()
		})
		if index > -1 {
			m.list.Select(index)
		}
		return m, m.list.SetItems(msg)
	case endMsg:
		return m, tea.Sequence(
			fetchTrack2Cmd,
			nextCmd,
		)
	case nextTrackMsg:
		m.list.CursorUp()
		return m, func() tea.Msg { return tea.KeyMsg{Type: tea.KeyEnter} }
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m trackList) View() string {
	return panelStyle.
		Width(m.width).
		Height(m.height).
		Render(m.msg, m.list.View())
}

func (m *trackList) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.list.SetSize(w, h)
}

func fetchTrackCmd() tea.Msg {
	tracks, err := radio.FetchLastNovaTracks(0)
	if err != nil {
		fmt.Printf("Error fetching nova tracks %s", err)
	}
	return tracksLoadedMsg(tracksToListItems(tracks))
}

func nextCmd() tea.Msg {
	return nextTrackMsg{}
}

func fetchTrack2Cmd() tea.Msg {
	tracks, err := radio.FetchLastNovaTracks(1)
	if err != nil {
		fmt.Printf("Error fetching nova tracks %s", err)
	}
	return tracksLoadedMsg(tracksToListItems(tracks))
}

// Converts a list of radioItem into a list compatible items
func tracksToListItems(tracks []radio.RadioItem) []list.Item {
	items := make([]list.Item, len(tracks))
	for i, item := range tracks {
		items[i] = list.Item(item)
	}
	return items
}
