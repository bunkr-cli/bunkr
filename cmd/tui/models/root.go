package models

import (
	"github.com/bunkr-cli/bunkr/cmd/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

type Root struct {
	main tea.Model
	w, h int
}

func NewRoot() (tea.Model, error) {
	var root Root
	var err error

	root.main, err = NewAlbums()
	if err != nil {
		return root, err
	}

	return root, nil
}

func (m Root) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.main.Init())
}

func (m Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height

	case messages.ErrMsg:
		m.main = NewErr(msg.Title(), msg)
		cmds = append(cmds, m.main.Init(), messages.TriggerSizeMsg(m.w, m.h))
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.main, cmd = m.main.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Root) View() string {
	return m.main.View()
}
