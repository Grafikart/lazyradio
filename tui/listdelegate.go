package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	NormalTitle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			Padding(0, 0, 0, 2) //nolint:mnd

	SelectedTitle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Padding(0, 0, 0, 1)
)

type listDelegate struct{}

func newListDelegate() listDelegate {
	return listDelegate{}
}

func (l listDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var content string
	if m.Index() == index {
		content = SelectedTitle.Render(item.FilterValue())
	} else {
		content = NormalTitle.Render(item.FilterValue())
	}
	fmt.Fprint(w, content)
}

func (l listDelegate) Height() int {
	return 1
}

func (l listDelegate) Spacing() int {
	return 1
}

func (l listDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
