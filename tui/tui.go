package tui

import (
	"fmt"
	"grafikart/lazyradio/radio"
	"grafikart/lazyradio/utils"
	"log"
	"net/url"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samber/lo"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type listKeyMap struct {
	openBrowser key.Binding
	play        key.Binding
	refresh     key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		openBrowser: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open yt music"),
		),
		play: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "play music"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "play music"),
		),
	}
}

type Model struct {
	tracks []ListItem
	list   list.Model
	keys   *listKeyMap
}

type ListItem struct {
	title, desc string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }

func NewModel() Model {
	listKeys := newListKeyMap()
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.openBrowser,
		}
	}
	l.DisableQuitKeybindings()
	l.Title = "Last tracks"
	return Model{
		list: l,
		keys: listKeys,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchTracks
}

type TracksMsg []radio.RadioItem
type ErrMsg error

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Keypress
	case tea.KeyMsg:
		switch {
		// Quit when pressing CTRL C or Esc
		case msg.Type == tea.KeyCtrlC, msg.Type == tea.KeyEsc:
			return m, tea.Quit
		// Open a track in the browser
		case key.Matches(msg, m.keys.openBrowser), msg.Type == tea.KeyEnter:
			openYtMusic(m.list.SelectedItem())
		// Start playing a track
		case key.Matches(msg, m.keys.play):
			play(m.list.SelectedItem())
		// Start playing a track
		case key.Matches(msg, m.keys.refresh):
			return m, fetchTracks
		}

	// Resize
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	// New tracks are loaded
	case TracksMsg:
		items := lo.Map(msg, func(v radio.RadioItem, i int) list.Item {
			return ListItem{title: v.Name, desc: v.Artist}
		})
		m.list.SetItems(items)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func fetchTracks() tea.Msg {
	tracks, err := radio.FetchLastNovaTracks()
	if err != nil {
		return ErrMsg(err)
	}
	return TracksMsg(tracks)
}

func openYtMusic(v list.Item) {
	switch vv := v.(type) {
	case ListItem:
		utils.OpenBrowser(fmt.Sprintf("https://music.youtube.com/search?q=%s", url.QueryEscape(vv.title+" "+vv.desc)))
	default:
		log.Fatalln("List element should be an ListItem")
	}
}

func play(v list.Item) {
	switch vv := v.(type) {
	case ListItem:
		cmd := exec.Command(
			"mpv",
			"--no-video",
			"--ytdl-format=bestaudio",
			fmt.Sprintf("ytdl://ytsearch1:%s", vv.title+" "+vv.desc),
		)
		err := cmd.Start()
		if err != nil {
			log.Fatal(fmt.Errorf("could not start mpv: %w", err))
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatal(fmt.Errorf("could not wait for mpv: %w", err))
		}
	default:
		log.Fatalln("List element should be an ListItem")
	}
}
