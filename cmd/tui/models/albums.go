package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bunkr-cli/bunkr/cmd/tui/delegate"
	"github.com/bunkr-cli/bunkr/cmd/tui/styles"
	"github.com/bunkr-cli/bunkr/internal/scrape"
	zone "github.com/lrstanley/bubblezone"
	"time"
)

type Albums struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegate.DelegateKeyMap
}

func NewAlbums() (tea.Model, error) {
	var (
		delegateKeys = delegate.NewDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	m := Albums{
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}

	albums, err := scrape.DefaultScraper.Albums(false)
	items := make([]list.Item, 0, len(albums))
	if err != nil {
		return m, err
	}
	for i := range albums {
		items = append(items, albums[i])
	}

	delegate := delegate.NewItemDelegate(delegateKeys)
	albumList := list.New(items, delegate, 0, 0)
	albumList.Title = "Bunkr Albums"
	albumList.Styles.Title = styles.TitleStyle
	albumList.Styles.PaginationStyle = albumList.Styles.StatusBar
	albumList.Paginator.Type = paginator.Arabic
	albumList.SetStatusBarItemName("album", "albums")
	albumList.StatusMessageLifetime = 5 * time.Second
	albumList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.downloadAlbum,
		}
	}
	albumList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.downloadAlbum,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}
	m.list = albumList

	return m, nil
}

func (m Albums) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Albums) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.AppStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.downloadAlbum):
			i, ok := m.list.SelectedItem().(*scrape.Album)
			if !ok {
				return m, nil
			}

			if err := scrape.DefaultScraper.HydrateAlbum(i); err != nil {
				return m, nil
			}

			fmt.Println(i)
			return m, nil

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.list.CursorUp()
			return m, nil
		case tea.MouseWheelDown:
			m.list.CursorDown()
			return m, nil
		case tea.MouseLeft:
			for i, listItem := range m.list.VisibleItems() {
				item, _ := listItem.(*scrape.Album)
				if zone.Get(item.Name).InBounds(msg) {
					m.list.Select(i)
					break
				}
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Albums) View() string {
	return zone.Scan(styles.AppStyle.Render(m.list.View()))
}

type listKeyMap struct {
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	downloadAlbum    key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		downloadAlbum: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "download")),
	}
}
