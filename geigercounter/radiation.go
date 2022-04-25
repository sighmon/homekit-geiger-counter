package geigercounter

import (
	"github.com/brutella/hc/characteristic"
)

const (
	TypeRadiationCpm = "8DBA11DD-991F-4265-9B0F-4985045222E1"
	TypeRadiationNsv = "665E5B09-1B6D-4296-BB98-398A61388C90"
	TypeRadiationUsv = "43C282AF-45C3-47E1-9BE2-79D054DE00E8"
)

type RadiationLevel struct {
	*characteristic.Float
}

func NewRadiationLevel(val float64) *RadiationLevel {
	p := RadiationLevel{characteristic.NewFloat("")}
	p.Value = val
	p.Format = characteristic.FormatFloat
	p.Perms = characteristic.PermsRead()
	p.SetMinValue(0)
	p.SetMaxValue(181600)
	p.SetValue(0)

	return &p
}
