package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rydwhelchel/spacetraders/api"
)

type systemview struct {
	traderService *api.TraderService
}

func newSystemView(ts *api.TraderService) *systemview {
	return &systemview{ts}
}

func (sv *systemview) Init() tea.Cmd {
	return nil
}

func (sv *systemview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return sv, nil
}

func (sv *systemview) View() string {
	return "System view!"
}
