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
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "Map",
		},
		//  pragma: allowlist nextline secret
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
	now := time.Now()
	entity := domain.EthernetSwitch{
		Entity: domain.Entity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: &now,
		},
		Name:        "ASD",
		Serial:      "asd",
		SwitchModel: "unifi_switch_us-24-250w",
		Address:     "asd",
		Username:    "asd",
		//  pragma: allowlist nextline secret
		Password: "asd",
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
