package models

import (
	"context"
	"github.com/bunkr-cli/bunkr/cmd/tui/delegate"
	"github.com/bunkr-cli/bunkr/cmd/tui/messages"
	"github.com/bunkr-cli/bunkr/cmd/tui/styles"
	"github.com/bunkr-cli/bunkr/internal/scrape"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"time"
)

type Albums struct {
	list          list.Model
	keys          *listKeyMap
	albumsLoading uint
	delegateKeys  *delegate.DelegateKeyMap
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

	delegate := delegate.NewItemDelegate(delegateKeys)
	albumList := list.New([]list.Item{}, delegate, 0, 0)
	albumList.Title = "Fetching Bunkr Albums..."
	albumList.Styles.Title = styles.TitleStyle
	albumList.Styles.PaginationStyle = albumList.Styles.StatusBar
	albumList.Paginator.Type = paginator.Arabic
	albumList.SetStatusBarItemName("album", "albums")
	albumList.StatusMessageLifetime = 5 * time.Second
	albumList.StartSpinner()
	albumList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.downloadAlbum,
		}
	}
	albumList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.reloadAlbums,
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
	return tea.Batch(tea.EnterAltScreen, m.list.StartSpinner(), messages.ListAlbums(context.Background(), false))
}

func (m Albums) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

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
			if m.albumsLoading == 0 {
				cmd = m.list.StartSpinner()
				cmds = append(cmds, cmd)
			}
			m.albumsLoading += 1

			i, ok := m.list.SelectedItem().(*scrape.Album)
			if !ok {
				return m, nil
			}

			cmds = append(cmds, messages.HydrateAlbum(context.Background(), i))
			return m, tea.Batch(cmds...)

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

		case key.Matches(msg, m.keys.reloadAlbums):
			m.list.Title = "Fetching Bunkr Albums..."
			cmds = append(cmds, m.list.StartSpinner())
			cmds = append(cmds, messages.ListAlbums(context.Background(), true))
			return m, tea.Batch(cmds...)
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
				if zone.Get(item.Identifier).InBounds(msg) {
					if m.list.SelectedItem() == listItem {
						if err := item.Open(); err != nil {
							return m, func() tea.Msg {
								return messages.NewErrMsg("Unable to open item", err)
							}
						}
					} else {
						m.list.Select(i)
					}
					break
				}
			}
		}

	case messages.AlbumsReadyMsg:
		items := make([]list.Item, 0, len(msg.Albums))
		for i := range msg.Albums {
			items = append(items, msg.Albums[i])
		}
		cmds = append(cmds, m.list.SetItems(items))
		m.list.StopSpinner()
		m.list.Title = "Bunkr Albums"

	case messages.AlbumHydratedMsg:
		m.albumsLoading -= 1
		if m.albumsLoading == 0 {
			m.list.StopSpinner()
		}
		return m, nil
	}

	m.list, cmd = m.list.Update(msg)
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
	reloadAlbums     key.Binding
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
		reloadAlbums: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "reload")),
		downloadAlbum: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "download")),
	}
}
