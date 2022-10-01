package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/interfaces"
	"rol/domain"
	"rol/infrastructure"
	"runtime"
	"testing"
)

type switchVLANRepoTester struct {
	repo       interfaces.IGenericRepository[domain.EthernetSwitchVLAN]
	dbPath     string
	insertedID uuid.UUID
}

var vlanRepoTester *switchVLANRepoTester

func Test_EthernetSwitchVLANRepository_Prepare(t *testing.T) {
	vlanRepoTester = &switchVLANRepoTester{}
	vlanRepoTester.dbPath = "ethernetSwitchVlanRepo_test.db"
	dbConnection := sqlite.Open(vlanRepoTester.dbPath)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.EthernetSwitchVLAN),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	vlanRepoTester.repo = infrastructure.NewEthernetSwitchVLANRepository(testGenDb, logger)

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), vlanRepoTester.dbPath)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(vlanRepoTester.dbPath)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchVLANRepository_Insert(t *testing.T) {
	vlan := domain.EthernetSwitchVLAN{
		VlanID:           10,
		EthernetSwitchID: uuid.New(),
		UntaggedPorts:    "test;test",
		TaggedPorts:      "test;test",
	}

	insertedVlan, err := vlanRepoTester.repo.Insert(context.Background(), vlan)
	if err != nil {
		t.Errorf("failed to insert vlan: %q", err)
	}
	if insertedVlan.ID == uuid.Nil {
		t.Error("failed to insert vlan: nil uuid")
	}
	vlanRepoTester.insertedID = insertedVlan.ID
}

func Test_EthernetSwitchVLANRepository_GetByID(t *testing.T) {
	vlan, err := vlanRepoTester.repo.GetByID(context.Background(), vlanRepoTester.insertedID)
	if err != nil {
		t.Errorf("failed to get vlan by id: %q", err)
	}
	if vlan.VlanID != 10 {
		t.Error("failed to get vlan by id: wrong vlanID")
	}
}

func Test_EthernetSwitchVLANRepository_GetList(t *testing.T) {
	vlans, err := vlanRepoTester.repo.GetList(context.Background(), "", "", 1, 1, nil)
	if err != nil {
		t.Errorf("failed to get vlan list: %q", err)
	}
	if len(vlans) != 1 {
		t.Error("failed to get vlan list: wrong list length")
	}
}

func Test_EthernetSwitchVLANRepository_Delete(t *testing.T) {
	err := vlanRepoTester.repo.Delete(context.Background(), vlanRepoTester.insertedID)
	if err != nil {
		t.Errorf("failed to delete vlan: %q", err)
	}
	_, err = vlanRepoTester.repo.GetByID(context.Background(), vlanRepoTester.insertedID)
	if err == nil {
		t.Error("failed to delete vlan: received deleted vlan")
	}
}

func Test_EthernetSwitchVLANRepository_Insert20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		vlan := domain.EthernetSwitchVLAN{
			VlanID:           i,
			EthernetSwitchID: uuid.New(),
			UntaggedPorts:    "test;test",
			TaggedPorts:      fmt.Sprintf("tag_%d", i),
		}
		_, err := vlanRepoTester.repo.Insert(context.Background(), vlan)
		if err != nil {
			t.Errorf("failed to insert vlan: %q", err)
		}
	}
}

func Test_EthernetSwitchVLANRepository_Pagination(t *testing.T) {
	vlans, err := vlanRepoTester.repo.GetList(context.Background(), "VlanID", "desc", 1, 20, nil)
	if err != nil {
		t.Errorf("failed to get vlan list: %q", err)
	}
	if len(vlans) != 20 {
		t.Error("failed to get vlan list: wrong list length")
	}
	if vlans[5].VlanID != 15 {
		t.Error("pagination failed")
	}
}

func Test_EthernetSwitchVLANRepository_Filter(t *testing.T) {
	queryBuilder := vlanRepoTester.repo.NewQueryBuilder(context.Background())
	queryGroupBuilder := vlanRepoTester.repo.NewQueryBuilder(context.Background())
	queryBuilder.Where("VlanID", ">=", 5).Where("VlanID", "<=", 10).
		WhereQuery(queryGroupBuilder.Where("TaggedPorts", "==", "tag_6").
			Or("TaggedPorts", "==", "tag_8"))

	vlans, err := vlanRepoTester.repo.GetList(context.Background(), "", "", 1, 20, queryBuilder)
	if err != nil {
		t.Errorf("failed to get vlan list: %q", err)
	}
	if len(vlans) != 2 {
		t.Error("failed to get vlan list: wrong list length")
	}
}

func Test_EthernetSwitchVLANRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	err := vlanRepoTester.repo.CloseDb()
	if err != nil {
		t.Errorf("failed to close db: %q", err)
	}
	if err = os.Remove(vlanRepoTester.dbPath); err != nil {
		t.Errorf("failed to remove db: %q", err)
	}
}
