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

type vlanServiceTester struct {
	service         *services.HostNetworkService
	configFilePath  string
	masterInterface string
	createdVlanName string
}

var vlanTester *vlanServiceTester

func Test_HostNetworkVlanService_Prepare(t *testing.T) {
	vlanTester = &vlanServiceTester{}
	_, filePath, _, _ := runtime.Caller(0)
	vlanTester.configFilePath = filepath.Join(filepath.Dir(filePath), "hostNetworkConfig.yaml")
	configStorage := infrastructure.NewYamlHostNetworkConfigStorage(domain.GlobalDIParameters{RootPath: filepath.Dir(vlanTester.configFilePath)})
	networkManager, err := infrastructure.NewHostNetworkManager(configStorage)
	if err != nil {
		t.Error("error to create host network manager")
	}
	vlanTester.service = services.NewHostNetworkService(networkManager)

	links, err := networkManager.GetList()
	if err != nil {
		t.Errorf("error getting list: %s", err.Error())
	}
	for _, link := range links {
		if link.GetName() != "lo" && link.GetType() != "vlan" {
			vlanTester.masterInterface = link.GetName()
			break
		}
	}
}

func Test_HostNetworkVlanService_CreateVlan(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID: 132,
		Parent: vlanTester.masterInterface,
		Addresses: []string{
			"123.123.123.123/24",
			"123.123.124.124/24",
		},
	}
	dto, err := vlanTester.service.CreateVlan(createDto)
	if err != nil {
		t.Errorf("error creating vlan: %s", err.Error())
	}
	vlanTester.createdVlanName = dto.Name
	if !strings.Contains(vlanTester.createdVlanName, "rol.") {
		t.Errorf("wrong vlan name: %s, expect rol.{%d}.{%s}", dto.Name, dto.VlanID, vlanTester.masterInterface)
	}
}

func Test_HostNetworkVlanService_CreateVlanWithIncorrectID(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID:    5000,
		Parent:    vlanTester.masterInterface,
		Addresses: []string{},
	}
	dto, err := vlanTester.service.CreateVlan(createDto)
	if err == nil {
		_ = vlanTester.service.DeleteVlan(dto.Name)
		t.Error("successfully created vlan with incorrect vlan id")
	}
}

func Test_HostNetworkVlanService_CreateVlanWithIncorrectMasterName(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID:    133,
		Parent:    " incorrect",
		Addresses: []string{},
	}
	dto, err := vlanTester.service.CreateVlan(createDto)
	if err != nil {
		if !errors.As(err, errors.Validation) {
			t.Errorf("expected error is not Validation error: %s", err.Error())
		}
	} else {
		_ = vlanTester.service.DeleteVlan(dto.Name)
		t.Error("successfully created vlan with incorrect master interface name")
	}
}

func Test_HostNetworkVlanService_CreateVlanWithNotExistedMasterInterface(t *testing.T) {
	createDto := dtos.HostNetworkVlanCreateDto{
		VlanID:    133,
		Parent:    "notexisted",
		Addresses: []string{},
	}
	dto, err := vlanTester.service.CreateVlan(createDto)
	if err != nil {
		if !errors.As(err, errors.Validation) {
			t.Errorf("expected error is not Validation error: %s", err.Error())
		}
	} else {
		_ = vlanTester.service.DeleteVlan(dto.Name)
		t.Error("successfully created vlan with incorrect master interface name")
	}
}

func Test_HostNetworkVlanService_UpdateVlan(t *testing.T) {
	updateDto := dtos.HostNetworkVlanUpdateDto{
		Addresses: []string{
			"123.123.125.125/24",
		},
	}
	dto, err := vlanTester.service.UpdateVlan(vlanTester.createdVlanName, updateDto)
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
	_, err := vlanTester.service.UpdateVlan(vlanTester.createdVlanName, updateDto)
	if err != nil {
		if !errors.As(err, errors.Validation) {
			t.Errorf("expected error is not Validation error: %s", err.Error())
		}
		return
	}
	t.Errorf("error is nil")
}

func Test_HostNetworkVlanService_GetByNameLo(t *testing.T) {
	_, err := vlanTester.service.GetVlanByName("lo")
	if !errors.As(err, errors.NotFound) {
		t.Errorf("unexpected behavior, expect not found error")
	}
}

func Test_HostNetworkVlanService_GetByNameVlan(t *testing.T) {
	vlan, err := vlanTester.service.GetVlanByName(vlanTester.createdVlanName)
	if err != nil {
		t.Errorf("get vlan by name failed: %s", err.Error())
	}
	if vlan.Name != vlanTester.createdVlanName {
		t.Errorf("wrong vlan name: %s, expect: %s", vlan.Name, vlanTester.createdVlanName)
	}
}

func Test_HostNetworkVlanService_GetList(t *testing.T) {
	vlans, err := vlanTester.service.GetVlanList()
	if err != nil {
		t.Errorf("get list failed: %s", err.Error())
	}
	vlanFound := false
	for _, vlan := range vlans {
		if vlan.Name == "lo" {
			t.Error("got lo interface through vlan service")
		}
		if vlan.Name == vlanTester.createdVlanName {
			vlanFound = true
		}
	}
	if !vlanFound {
		t.Error("created vlan was not found")
	}
}

func Test_HostNetworkVlanService_Delete(t *testing.T) {
	err := vlanTester.service.DeleteVlan(vlanTester.createdVlanName)
	if err != nil {
		t.Errorf("delete vlan failed: %s", err.Error())
	}
	_, err = vlanTester.service.GetVlanByName(vlanTester.createdVlanName)
	if err == nil {
		t.Error("deleted vlan was received")
	}
}

func Test_HostNetworkVlanService_CleaningAfterTests(t *testing.T) {
	err := os.Remove(vlanTester.configFilePath)
	if err != nil {
		t.Errorf("remove network config file failed:  %q", err)
	}
}
