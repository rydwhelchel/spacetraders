package view

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rydwhelchel/spacetraders/api"
)

type fleetview struct {
	list list.Model
	// initialized tracks whether we have initialized the contained list with width&height
	initialized   bool
	traderService *api.TraderService
}

func newFleetView() *fleetview {
	return &fleetview{}
}

// initFleetList is in charge of creating the list object with the correct width & height, populating may happen elsewhere
func (fv *fleetview) initFleetList(width, height int) {
	// TODO: Subtract height of fleetview header
	fv.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	fv.list.Title = "Fleet"
}

func (fv *fleetview) Init() tea.Cmd {
	return nil
}

func (fv *fleetview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return fv, nil
}

func (fv *fleetview) View() string {
	return ""
}
