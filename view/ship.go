package view

import (
	"fmt"
	"time"

	"github.com/rydwhelchel/spacetraders/api"
	"github.com/rydwhelchel/spacetraders/api/openapi"
)

// separate from ship
type shipview struct {
	ship
	traderService *api.TraderService
}

type ship struct {
	openapi.Ship
}

func (s *ship) Title() string {
	return s.Symbol + " " + string(s.Registration.Role)
}

// Description taken from oldview for now
func (s *ship) Description() string {
	shipString := fmt.Sprintf("Ship %v", s.Registration.GetName())
	// Ship details
	shipString += fmt.Sprintf(" ~ %v ~ F%v/%v", s.Registration.Role, s.Fuel.Current, s.Frame.FuelCapacity)
	if s.Cargo.Capacity > 0 {
		shipString += fmt.Sprintf(" ~ C%v/%v", s.Cargo.Units, s.Cargo.Capacity)
	}
	shipString += s.getCooldownText()
	return shipString
}

// FilterValue allows you to filter on ship type for now
func (s *ship) FilterValue() string {
	return string(s.Registration.Role)
}

// TODO: Progress bar?
func (s *ship) getCooldownText() (cdText string) {
	// Ship cooldown details
	exp := s.Cooldown.Expiration
	if exp != nil {
		cd := s.Cooldown.Expiration.Sub(time.Now())
		if cd > 0 {
			seconds := int(cd.Seconds()) % 60
			minutes := int(cd.Seconds()) / 60
			cdText += fmt.Sprintf(" on cooldown for ")
			if minutes > 0 {
				cdText += fmt.Sprintf("%v minutes and ", minutes)
			}
			cdText += fmt.Sprintf("%v seconds", seconds)
		}
	}

	return cdText
}
