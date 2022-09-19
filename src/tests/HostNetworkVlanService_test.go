//go:build linux

package tests

import (
	"os"
	"path/filepath"
	"rol/app/errors"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"strings"
	"testing"
)

var (
	vlanService            *services.HostNetworkVlanService
	configServiceFilePath  string
	serviceMasterInterface string
	createdVlanName        string
)

func Test_HostNetworkVlanService_Prepare(t *testing.T) {
	_, filePath, _, _ := runtime.Caller(0)
	configServiceFilePath = filepath.Join(filepath.Dir(filePath), "hostNetworkConfig.yaml")
	configStorage := infrastructure.NewYamlHostNetworkConfigStorage(domain.GlobalDIParameters{RootPath: filepath.Dir(configServiceFilePath)})
	networkManager, err := infrastructure.NewHostNetworkManager(configStorage)
	if err != nil {
		t.Error("error to create host network manager")
	}
	vlanService = services.NewHostNetworkVlanService(networkManager)

	links, err := networkManager.GetList()
	if err != nil {
		t.Errorf("error getting list: %s", err.Error())
	}
	for _, link := range links {
		if link.GetName() != "lo" && link.GetType() != "vlan" {
			serviceMasterInterface = link.GetName()
			break
		}
	}
}

func Test_HostNetworkVlanService_CreateVlan(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID: 132,
		Master: serviceMasterInterface,
		Addresses: []string{
			"123.123.123.123/24",
			"123.123.124.124/24",
		},
	}
	dto, err := vlanService.Create(createDto)
	if err != nil {
		t.Errorf("error creating vlan: %s", err.Error())
	}
	createdVlanName = dto.Name
	if !strings.Contains(createdVlanName, "rol.") {
		t.Errorf("wrong vlan name: %s, expect rol.{%d}.{%s}", dto.Name, dto.VlanID, serviceMasterInterface)
	}
}

func Test_HostNetworkVlanService_CreateVlanWithIncorrectID(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID:    5000,
		Parent:    serviceMasterInterface,
		Addresses: []string{},
	}
	dto, err := vlanService.Create(createDto)
	if err == nil {
		_ = vlanService.Delete(dto.Name)
		t.Error("successfully created vlan with incorrect vlan id")
	}
}

func Test_HostNetworkVlanService_CreateVlanWithIncorrectMasterName(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID:    133,
		Parent:    " incorrect",
		Addresses: []string{},
	}
	dto, err := vlanService.Create(createDto)
	if err != nil {
		if !errors.As(err, errors.Validation) {
			t.Errorf("expected error is not Validation error: %s", err.Error())
		}
	} else {
		_ = vlanService.Delete(dto.Name)
		t.Error("successfully created vlan with incorrect master interface name")
	}
}

func Test_HostNetworkVlanService_CreateVlanWithNotExistedMasterInterface(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID:    133,
		Parent:    "notexisted",
		Addresses: []string{},
	}
	dto, err := vlanService.Create(createDto)
	if err != nil {
		if !errors.As(err, errors.Validation) {
			t.Errorf("expected error is not Validation error: %s", err.Error())
		}
	} else {
		_ = vlanService.Delete(dto.Name)
		t.Error("successfully created vlan with incorrect master interface name")
	}
}

func Test_HostNetworkVlanService_UpdateVlan(t *testing.T) {
	updateDto := dtos.HostNetworkVlanUpdateDto{
		Addresses: []string{
			"123.123.125.125/24",
		},
	}
	dto, err := vlanService.Update(createdVlanName, updateDto)
	if err != nil {
		t.Errorf("error creating vlan: %s", err.Error())
	}
	if len(dto.Addresses) != 1 {
		t.Error("failed to update vlan addresses")
	}
	for _, addressStr := range dto.Addresses {
		if addressStr != "123.123.125.125/24" {
			t.Error("failed to update vlan addresses")
			return
		}
	}
}

func Test_HostNetworkVlanService_UpdateIncorrectAddress(t *testing.T) {
	updateDto := dtos.HostNetworkVlanUpdateDto{
		Addresses: []string{
			"123.123.125.1252/24",
		},
	}
	_, err := vlanService.Update(createdVlanName, updateDto)
	if err != nil {
		if !errors.As(err, errors.Validation) {
			t.Errorf("expected error is not Validation error: %s", err.Error())
		}
		return
	}
	t.Errorf("error is nil")
}

func Test_HostNetworkVlanService_GetByNameLo(t *testing.T) {
	_, err := vlanService.GetByName("lo")
	if !errors.As(err, errors.NotFound) {
		t.Errorf("unexpected behavior, expect not found error")
	}
}

func Test_HostNetworkVlanService_GetByNameVlan(t *testing.T) {
	vlan, err := vlanService.GetByName(createdVlanName)
	if err != nil {
		t.Errorf("get vlan by name failed: %s", err.Error())
	}
	if vlan.Name != createdVlanName {
		t.Errorf("wrong vlan name: %s, expect: %s", vlan.Name, createdVlanName)
	}
}

func Test_HostNetworkVlanService_GetList(t *testing.T) {
	vlans, err := vlanService.GetList()
	if err != nil {
		t.Errorf("get list failed: %s", err.Error())
	}
	vlanFound := false
	for _, vlan := range vlans {
		if vlan.Name == "lo" {
			t.Error("got lo interface through vlan service")
		}
		if vlan.Name == createdVlanName {
			vlanFound = true
		}
	}
	if !vlanFound {
		t.Error("created vlan was not found")
	}
}

func Test_HostNetworkVlanService_Delete(t *testing.T) {
	err := vlanService.Delete(createdVlanName)
	if err != nil {
		t.Errorf("delete vlan failed: %s", err.Error())
	}
	_, err = vlanService.GetByName(createdVlanName)
	if err == nil {
		t.Error("deleted vlan was received")
	}
}

func Test_HostNetworkVlanService_CleaningAfterTests(t *testing.T) {
	err := os.Remove(configServiceFilePath)
	if err != nil {
		t.Errorf("remove network config file failed:  %q", err)
	}
}
