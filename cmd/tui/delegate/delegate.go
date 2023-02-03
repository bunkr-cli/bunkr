package delegate

import (
	"fmt"
	"github.com/bunkr-cli/bunkr/cmd/tui/styles"
	"github.com/bunkr-cli/bunkr/internal/scrape"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func NewItemDelegate(keys *DelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		i, ok := m.SelectedItem().(*scrape.Album)
		if !ok {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.open):
				if err := i.Open(); err != nil {
					return m.NewStatusMessage(err.Error())
				}
				return m.NewStatusMessage(styles.StatusMessageStyle(fmt.Sprintf(`Opening "%s" at %s`, i.Name, i.URL().String())))
			}
		}

		return nil
	}

	d.ShortHelpFunc = keys.ShortHelp
	d.FullHelpFunc = keys.FullHelp
	return d
}

type DelegateKeyMap struct {
	open key.Binding
}

func (d DelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.open,
	}
}

func (d DelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.open,
		},
	}
}

func NewDelegateKeyMap() *DelegateKeyMap {
	return &DelegateKeyMap{
		open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open"),
		),
	}
}
