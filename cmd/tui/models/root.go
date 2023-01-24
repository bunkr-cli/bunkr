package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Root struct {
	albums tea.Model
}

func NewRoot() (tea.Model, error) {
	var root Root
	var err error

	root.albums, err = NewAlbums()
	if err != nil {
		return root, err
	}

	return root, nil
}

func (m Root) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.albums, cmd = m.albums.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Root) View() string {
	return m.albums.View()
}
