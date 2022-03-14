package geigercounter

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type Accessory struct {
	*accessory.Accessory

	GeigerCounter *Service
}

func NewAccessory(info accessory.Info) *Accessory {
	a := accessory.New(info, accessory.TypeSensor)
	svc := NewService(info.Name)

	a.AddService(svc.Service)

	return &Accessory{a, svc}
}

const TypeAirQualitySensor = "8D"

type Service struct {
	*service.Service

	Name	*characteristic.Name
	Cpm		*RadiationLevel
	Nsvh	*RadiationLevel
	Usvh	*RadiationLevel
}

func NewService(name string) *Service {
	nameChar := characteristic.NewName()
	nameChar.SetValue(name)

	countsPerMinute := NewRadiationLevel(0)
	countsPerMinute.Type = characteristic.UnitPPM
	countsPerMinute.Unit = "CPM"
	countsPerMinute.Description = "Counts per minute"

	nanoSievert := NewRadiationLevel(0)
	nanoSievert.Type = characteristic.UnitPPM
	nanoSievert.Unit = "nSv/h"
	nanoSievert.Description = "Nanosieverts per hour"

	microSievert := NewRadiationLevel(0)
	microSievert.Type = characteristic.UnitPPM
	microSievert.Unit = "µSv/h"
	microSievert.Description = "Microsieverts per hour"

	svc := service.New(TypeAirQualitySensor)
	svc.AddCharacteristic(countsPerMinute.Characteristic)
	svc.AddCharacteristic(nanoSievert.Characteristic)
	svc.AddCharacteristic(microSievert.Characteristic)

	return &Service{svc, nameChar, countsPerMinute, nanoSievert, microSievert}
}
