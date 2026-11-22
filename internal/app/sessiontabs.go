package app

import (
    "finallyshell/internal/icons"
    "finallyshell/internal/widgets"
    "finallyshell/locales"

    "github.com/egoist/mygo/ui"
)

// sessionTabs is the compact top bar used by the main terminal area.
// Saved connections are opened from the folder button; open sessions stay
// visible as horizontal tabs so the server list does not consume the sidebar.
func (a *App) sessionTabs(c *ui.Context) {
    type tabItem struct {
		session *Session
		index   int
	}
	var items []tabItem
	selected := 0
	active := a.current()
	if active != nil && active.isPanelTab() {
		active = active.panelOf
	}
	for i, s := range a.sessions {
		if s.isPanelTab() {
			continue
		}
		if s == active {
			selected = len(items)
		}
		items = append(items, tabItem{session: s, index: i})
	}
	if len(items) == 0 {
		selected = 0
	}

	var closing *Session
	ui.Row(c).Height(34).PaddingX(6).Gap(4).AlignItems(ui.Center).
		Background(theme.Background).BorderWidth(0, 0, 1, 0).BorderColor(theme.Border).Children(func() {
			if widgets.IconButton(c, "folder", locales.ConnectionsPage()).Clicked() {
				a.showManager(false)
			}
			ui.ScrollHorizontal(c).Grow(1).MinWidth(0).Scrollbars(false).Children(func() {
				chosen := selected
				tabs := ui.TabsBase(c, &chosen, len(items))
				tabs.List.Gap(2).Children(func() {
					for i, item := range items {
						tab := tabs.Tab(i).Row().Gap(5).AlignItems(ui.Center).Padding(4, 8).
							MinWidth(72).MaxWidth(220).Shrink(0).Radius(4)
						if i == chosen {
							tab.Background(theme.Surface)
						}
						tab.Children(func() {
							ui.Icon(c, icons.Get("server")).FontSize(13).Shrink(0)
							ui.Text(c, item.session.label()).SingleLine().Shrink(1).MinWidth(0)
							close := widgets.IconButton(c, "x", locales.CloseSession()).Size(18, 18)
							if close.Clicked() {
								closing = item.session
							}
						})
					}
				})
				if chosen != selected && chosen >= 0 && chosen < len(items) {
					a.active = items[chosen].index
				}
			})
			if widgets.IconButton(c, "plus", locales.NewTab()).Clicked() {
				a.NewTab()
			}
		})
	if closing != nil {
		a.closeSession(closing)
	}
}

