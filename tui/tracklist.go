package tui

import (
	"fmt"
	"grafikart/lazyradio/radio"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type trackListModel struct {
	// components
	list list.Model

	// state
	item   radio.TrackItem
	width  int
	height int
	msg    string
	player *radio.Player
	radio  radioItem
}

type tracklistKeyMap struct {
	openBrowser key.Binding
	refresh     key.Binding
	pause       key.Binding
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
	pause: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "pause/play"),
	),
}

func newTrackList(p *radio.Player) trackListModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Track List"
	l.DisableQuitKeybindings()
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.openBrowser,
			keys.refresh,
			keys.pause,
		}
	}
	return trackListModel{
		list:   l,
		player: p,
	}
}

func (m trackListModel) Init() tea.Cmd {
	return nil
}

type tracksLoadedMsg []list.Item
type nextTrackMsg struct{}

func (m trackListModel) Update(msg tea.Msg) (trackListModel, tea.Cmd) {
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
		// Refresh track list
		case key.Matches(msg, keys.refresh):
			return m, func() tea.Msg { return radioSelectedMsg(m.radio) }
		}
	case radioSelectedMsg:
		m.list.ToggleSpinner()
		m.radio = radioItem(msg)
		return m, func() tea.Msg {
			items, err := msg.fetcher()
			if err != nil {
				fmt.Printf("Error fetching track items: %s\n", err)
			}
			return tracksLoadedMsg(tracksToListItems(items))
		}
	case tracksLoadedMsg:
		m.list.ToggleSpinner()
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
			func() tea.Msg { return radioSelectedMsg(m.radio) },
			func() tea.Msg { return nextTrackMsg{} },
		)
	case nextTrackMsg:
		m.list.CursorUp()
		return m, func() tea.Msg { return tea.KeyMsg{Type: tea.KeyEnter} }
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m trackListModel) View(focused bool) string {
	style := panelStyle
	if !focused {
		style = mutedPanelStyle
	}
	return style.
		Width(m.width).
		Height(m.height).
		Render(m.msg, m.list.View())
}

func (m *trackListModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.list.SetSize(w, h)
}

// Converts a list of radioItem into a list compatible items
func tracksToListItems(tracks []radio.TrackItem) []list.Item {
	items := make([]list.Item, len(tracks))
	for i, item := range tracks {
		items[i] = list.Item(item)
	}
	return items
}
