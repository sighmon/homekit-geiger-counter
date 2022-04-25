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

const TypeGeigerCounter = "AC6273AD-E687-4D03-9907-21991198E9E4"

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
	countsPerMinute.Type = TypeRadiationCpm
	countsPerMinute.Unit = "CPM"
	countsPerMinute.Description = "Counts per minute"

	nanoSievert := NewRadiationLevel(0)
	nanoSievert.Type = TypeRadiationNsv
	nanoSievert.Unit = "nSv/h"
	nanoSievert.Description = "Nanosieverts per hour"

	microSievert := NewRadiationLevel(0)
	microSievert.Type = TypeRadiationUsv
	microSievert.Unit = "µSv/h"
	microSievert.Description = "Microsieverts per hour"

	svc := service.New(TypeGeigerCounter)
	svc.AddCharacteristic(countsPerMinute.Characteristic)
	svc.AddCharacteristic(nanoSievert.Characteristic)
	svc.AddCharacteristic(microSievert.Characteristic)

	return &Service{svc, nameChar, countsPerMinute, nanoSievert, microSievert}
}
