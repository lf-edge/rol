package tests

import (
	"github.com/google/uuid"
	"rol/domain"
	"rol/dtos"
	"testing"
	"time"
)

func Test_EthernetSwitchCreateDtoToEntity(t *testing.T) {
	tester := NewGenericMapperToEntity[dtos.EthernetSwitchCreateDto, domain.EthernetSwitch]()
	dto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "Mapper",
			Serial:      "serial",
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "Map",
		},
		Password: "per",
	}
	entity := new(domain.EthernetSwitch)
	err := tester.MapToEntity(dto, entity)
	if err != nil {
		t.Error(err)
	}

	if err := compareDtoAndEntity(dto, entity); err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchEntityToDto(t *testing.T) {
	tester := NewGenericMapperToEntity[dtos.EthernetSwitchDto, domain.EthernetSwitch]()
	entity := domain.EthernetSwitch{
		Entity: domain.Entity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Now(),
		},
		Name:        "ASD",
		Serial:      "asd",
		SwitchModel: 0,
		Address:     "asd",
		Username:    "asd",
		Password:    "asd",
		Ports:       nil,
	}
	dto := new(dtos.EthernetSwitchDto)
	err := tester.MapToDto(entity, dto)
	if err != nil {
		t.Error(err)
	}

	if err := compareDtoAndEntity(entity, dto); err != nil {
		t.Error(err)
	}
}
