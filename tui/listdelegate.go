package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type simpleListDelegate struct {
	normalStyle   lipgloss.Style
	selectedStyle lipgloss.Style
}

func newSimpleListDelegate(focused bool) simpleListDelegate {
	d := newDefaultListDelegate(focused)
	return simpleListDelegate{
		normalStyle:   d.Styles.NormalTitle,
		selectedStyle: d.Styles.SelectedTitle,
	}
}

func newDefaultListDelegate(focused bool) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	if focused {
		return d
	}
	d.Styles.SelectedTitle = d.Styles.NormalTitle
	d.Styles.SelectedDesc = d.Styles.NormalTitle
	d.Styles.NormalTitle = d.Styles.NormalDesc
	return d
}

func (l simpleListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var content string
	if m.Index() == index {
		content = l.selectedStyle.Render(item.FilterValue())
	} else {
		content = l.normalStyle.Render(item.FilterValue())
	}
	fmt.Fprint(w, content)
}

func (l simpleListDelegate) Height() int {
	return 1
}

func (l simpleListDelegate) Spacing() int {
	return 1
}

func (l simpleListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
