package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type viewFleetMsg struct{}

func viewFleetCmd() tea.Msg {
	return viewFleetMsg{}
}

type viewSystemsMsg struct{}

func viewSystemsCmd() tea.Msg {
	return viewSystemsMsg{}
}

// TODO: How to make this a tea.Cmd?
// type viewWaypointsMsg struct {
// 	system string
// }
//
// func viewWaypointsCmd(system string) tea.Msg {
// 	return viewWaypointsMsg{system}
// }
