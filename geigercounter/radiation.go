package geigercounter

import (
	"github.com/brutella/hc/characteristic"
)

const (
	TypeRadiation = "8DBA11DD-991F-4265-9B0F-4985045222E1"
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