package tui

import (
	"grafikart/lazyradio/radio"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var radios = []radioItem{
	{
		name:    "Radio Nova",
		fetcher: radio.FetcherNova(910),
	},
	{
		name:    "Nova Hip-Hop",
		fetcher: radio.FetcherNova(196018),
	},
	{
		name:    "Nova Reggae",
		fetcher: radio.FetcherNova(220781),
	},
	{
		name:    "Nova Soul",
		fetcher: radio.FetcherNova(216011),
	},
	{
		name:    "Nouvo Nova",
		fetcher: radio.FetcherNova(79676),
	},
	{
		name:    "Nova la nuit",
		fetcher: radio.FetcherNova(916),
	},
	{
		name:    "Nova Classics",
		fetcher: radio.FetcherNova(913),
	},
	{
		name:    "Nova Danse",
		fetcher: radio.FetcherNova(560),
	},
	{
		name:    "Nova la plage",
		fetcher: radio.FetcherNova(268797),
	},
	{name: "Fip ", fetcher: radio.FetcherFipTracks("fip")},
	{name: "Fip Sacré français", fetcher: radio.FetcherFipTracks("radio-sacre-francais")},
	{name: "Fip Rock", fetcher: radio.FetcherFipTracks("radio-rock")},
	{name: "Fip Jazz", fetcher: radio.FetcherFipTracks("radio-jazz")},
	{name: "Fip Groove", fetcher: radio.FetcherFipTracks("radio-groove")},
	{name: "Fip Reggae", fetcher: radio.FetcherFipTracks("radio-reggae")},
	{name: "Fip Pop", fetcher: radio.FetcherFipTracks("radio-pop")},
	{name: "Fip Electro", fetcher: radio.FetcherFipTracks("radio-electro")},
	{name: "Fip Monde", fetcher: radio.FetcherFipTracks("radio-monde")},
	{name: "Fip Nouveautés", fetcher: radio.FetcherFipTracks("radio-nouveautes")},
	{name: "Fip Metal", fetcher: radio.FetcherFipTracks("radio-metal")},
	{name: "Fip Hip Hop", fetcher: radio.FetcherFipTracks("radio-hip-hop")},
}

type radioItem struct {
	name    string
	fetcher func() ([]radio.TrackItem, error)
}

type radioSelectedMsg radioItem

func (s radioItem) FilterValue() string {
	return s.name
}

type radioListModel struct {
	list   list.Model
	width  int
	height int
}

func newSidebar() radioListModel {
	l := list.New(sidebarItemsToListItems(radios), newListDelegate(), 0, 0)
	l.Title = "Radios"
	l.DisableQuitKeybindings()
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	return radioListModel{
		list: l,
	}
}

func (m radioListModel) Init() tea.Cmd {
	return func() tea.Msg {
		return radioSelectedMsg(radioItem(m.list.Items()[0].(radioItem)))
	}
}

func (m radioListModel) Update(msg tea.Msg) (radioListModel, tea.Cmd) {
	var cmdList tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			item := m.list.SelectedItem().(radioItem)
			return m, func() tea.Msg { return radioSelectedMsg(item) }
		}
	}
	m.list, cmdList = m.list.Update(msg)
	return m, cmdList
}

func (m radioListModel) View(focused bool) string {
	style := panelStyle
	if !focused {
		style = mutedPanelStyle
	}
	return style.
		Width(m.width).
		Height(m.height).
		Render(m.list.View())
}

func (m *radioListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

// Converts a list of radioItem into a list compatible items
func sidebarItemsToListItems(sitems []radioItem) []list.Item {
	items := make([]list.Item, len(sitems))
	for i, item := range sitems {
		items[i] = list.Item(item)
	}
	return items
}
