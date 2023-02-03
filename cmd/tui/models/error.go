package models

import (
	"github.com/bunkr-cli/bunkr/cmd/tui/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

func NewErr(title string, err error) Err {
	return Err{
		title: title,
		err:   err,
		keys:  keys,
		help:  help.New(),
	}
}

type Err struct {
	title string
	err   error
	keys  keyMap
	help  help.Model
	w, h  int
}

func (m Err) Init() tea.Cmd {
	return nil
}

func (m Err) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Err) View() string {
	s := styles.Err(wordwrap.String(m.title+".\n\n"+m.err.Error(), m.w))
	s += "\n\n"
	s += m.help.View(m.keys)
	return s
}

type keyMap struct {
	quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = keyMap{
	quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
