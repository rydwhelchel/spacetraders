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

func newFleetView(ts *api.TraderService) *fleetview {
	return &fleetview{
		traderService: ts,
	}
}

// initFleetList is in charge of creating the list object with the correct width & height, populating may happen elsewhere
func (fv *fleetview) initFleetList(width, height int) {
	// TODO: Subtract height of fleetview header
	fv.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	fv.list.Title = "Fleet"

	if !fv.initialized {
		fleetMap := fv.traderService.GetFleetData()
		var items []list.Item
		for _, v := range fleetMap {
			items = append(items, ship{v})
		}

		fv.list.SetItems(items)
	}

	fv.initialized = true
}

func (fv *fleetview) Init() tea.Cmd {
	return nil
}

func (fv *fleetview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !fv.initialized {
		return fv, nil
	}

	list, cmd := fv.list.Update(msg)
	fv.list = list
	return fv, cmd
}

func (fv *fleetview) View() string {
	if !fv.initialized {
		return "Loading..."
	}
	return fv.list.View()
}
