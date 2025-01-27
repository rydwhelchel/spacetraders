package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rydwhelchel/spacetraders/api"
	"github.com/rydwhelchel/spacetraders/view/styles"
)

// NOTE: Can create a goroutine that "ticks" every second, sending an updateMsg to model.Update(updateMsg), this will allow kicking off requests based on time and automatically updating cooldown text every second
//			not convinced by this method, there is a "Timer" component in Bubbles which may serve better for this sort of issue.

type page int

const (
	main page = iota
	fleet
	systems
	waypoints
)

type Model struct {
	traderService *api.TraderService

	width  int
	height int

	menu       list.Model
	fleetview  *fleetview
	systemview *systemview

	// Pages
	page page
}

type menuOpt struct {
	title       string
	description string
	initFunc    func(int, int)
	cmd         tea.Cmd
}

func (mo menuOpt) FilterValue() string {
	return mo.title
}

func (mo menuOpt) Title() string {
	return mo.title
}

func (mo menuOpt) Description() string {
	return mo.description
}

func NewModel(traderService *api.TraderService) *Model {
	return &Model{
		traderService: traderService,
		fleetview:     newFleetView(traderService),
		systemview:    newSystemView(traderService),
		page:          main,
	}
}

// NOTE: Not convinced by this approach -- revisit
func (m *Model) initMenu() {
	menuOptions := []list.Item{
		menuOpt{
			title:       "Fleet",
			description: "View ship fleet",
			initFunc:    m.fleetview.initFleetList,
			cmd:         viewFleetCmd,
		},
		menuOpt{
			title:       "Systems",
			description: "View all systems",
			cmd:         viewSystemsCmd,
		},
		// TODO: How to make this a tea.Cmd? (Make it accept an argument.. Maybe currying?)
		// TODO: This command should take you to a system selection screen probably
		// TODO: We should have another option on ship menu which takes systemSymbol as an argument at menu creation
		// 	{
		// 		title:       "Waypoints",
		// 		description: "View waypoints in given system TBD",
		// 		action:      viewWaypointsCmd,
		// 	},
	}

	m.menu = list.New([]list.Item{}, list.NewDefaultDelegate(), m.width, m.height)
	m.menu.SetShowTitle(false)
	m.menu.SetFilteringEnabled(false)
	m.menu.SetShowStatusBar(false)
	m.menu.SetItems(menuOptions)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		// Subtract 1 line from height for header, and for every newline character subtract 1 more
		m.height = msg.Height - strings.Count(m.getHeader(), "\n") - 1

		// Init lists with width/height
		m.initMenu()
	// TODO: Return early to prevent calling update functions of child pages; or maybe we only need to return early in children's Updates
	// NOTE: think it only needs to happen in children

	case tea.KeyMsg:
		switch msg.String() {
		// TODO: Prob not necessary, think it's handled already
		case "ctrl+c", "q":
			return m, tea.Quit
		case tea.KeyBackspace.String(): // Go back to main menu TODO: Add helper text for bksp
			m.page = main
			return m, tea.ClearScreen
		case tea.KeyEnter.String(): // Choose a menu option and return its command
			mo := m.menu.SelectedItem().(menuOpt)
			mo.initFunc(m.width, m.height)
			return m, mo.cmd
		}

	case viewFleetMsg:
		m.page = fleet

	case viewSystemsMsg:
		m.page = systems
	}

	switch m.page {
	// TODO: Should main menu be a child, or part of main model?
	case main:
		menu, cmd := m.menu.Update(msg)
		m.menu = menu
		return m, cmd
	case fleet:
		fv, cmd := m.fleetview.Update(msg)
		m.fleetview = fv.(*fleetview)
		return m, cmd
	case systems:
		sv, cmd := m.systemview.Update(msg)
		m.systemview = sv.(*systemview)
		return m, cmd
	}
	return m, nil
}

func (m *Model) View() string {
	screen := m.getHeader() + "\n"
	switch m.page {
	case main:
		screen += m.menu.View()
	case fleet:
		screen += m.fleetview.View()
	case systems:
		screen += m.systemview.View()
	}

	return screen
}

// getHeader returns the header for the entire app, this will always show
func (m *Model) getHeader() string {
	// TODO: This header is styled super ugly right now, clean it up
	// Style header
	var padLeftStyle = lipgloss.NewStyle().
		PaddingLeft(2)
	var baseStyle = lipgloss.NewStyle().
		Background(styles.MANTLE).
		Bold(true)
	var agentStyle = baseStyle.
		Foreground(styles.RED)
	var hqStyle = baseStyle.
		Foreground(styles.BLUE).
		AlignHorizontal(lipgloss.Right)
	var credStyle = baseStyle.
		Foreground(styles.FLAMINGO)
	var shipStyle = baseStyle.
		Foreground(styles.FLAMINGO)

	// TODO: Make this look less crap :)
	header := fmt.Sprintf("%v - %v - Ships: %v ~~~~~~~~~~~~ %v",
		agentStyle.Render(m.traderService.Data.Agent.Symbol),
		credStyle.Render("$"+strconv.Itoa(int(m.traderService.Data.Agent.Credits))),
		shipStyle.Render(strconv.Itoa(int(m.traderService.Data.Agent.ShipCount))),
		hqStyle.Render("Headquartered in: "+m.traderService.Data.Agent.Headquarters))
	return padLeftStyle.Render(header)

}
