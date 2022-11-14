package tests

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pin/tftp/v3"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/interfaces"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"testing"
)

type tftpServiceTester struct {
	configRepo      interfaces.IGenericRepository[uuid.UUID, domain.TFTPConfig]
	pathsRepo       interfaces.IGenericRepository[uuid.UUID, domain.TFTPPathRatio]
	service         *services.TFTPServerService
	tftpAddress     string
	createdServerID uuid.UUID
	createdPathID   uuid.UUID
	deleteServerID  uuid.UUID
	dbFileName      string
}

var tftpTester *tftpServiceTester

func Test_TFTPServerService_Prepare(t *testing.T) {
	tftpTester = &tftpServiceTester{}
	tftpTester.dbFileName = "tftpServerService_test.db"
	//remove old test db file
	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), tftpTester.dbFileName)); err == nil {
		err = os.Remove(tftpTester.dbFileName)
		if err != nil {
			t.Errorf("remove db failed:  %q", err)
		}
	}
	dbConnection := sqlite.Open(tftpTester.dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.TFTPConfig),
		new(domain.TFTPPathRatio),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	tftpTester.configRepo = infrastructure.NewGormGenericRepository[uuid.UUID, domain.TFTPConfig](testGenDb, logger)
	tftpTester.pathsRepo = infrastructure.NewGormGenericRepository[uuid.UUID, domain.TFTPPathRatio](testGenDb, logger)
	factory, _ := infrastructure.NewPinTFTPServerFactory()
	tftpTester.service = services.NewTFTPServerService(tftpTester.configRepo, tftpTester.pathsRepo, factory, logger)
	if err != nil {
		t.Errorf("create new service failed: %q", err)
	}

	err = os.MkdirAll("files", os.ModePerm)
	if err != nil {
		t.Errorf("error when creating a directory: %q", err)
	}
	f, err := os.Create("files/test.txt")
	if err != nil {
		t.Errorf("error when creating a file: %q", err)
	}
	defer f.Close()

	_, err = f.WriteString("test string")

	if err != nil {
		t.Errorf("write string failed: %q", err)
	}

	links, err := netlink.LinkList()
	if err != nil {
		t.Errorf("error getting a list of link devices: %q", err)
	}

	for _, link := range links {
		if link.Attrs().Name != "lo" && link.Type() != "vlan" {
			addr, err := netlink.AddrList(link, netlink.FAMILY_V4)
			if err != nil {
				t.Errorf("error getting a list of link addresses: %q", err)
			}
			for _, a := range addr {
				tftpTester.tftpAddress = a.IP.String()
			}
			break
		}
	}
}

func Test_TFTPServerService_InitializeTest(t *testing.T) {
	cfg1 := domain.TFTPConfig{
		Address: tftpTester.tftpAddress,
		Port:    "6969",
		Enabled: true,
	}

	a, _ := tftpTester.configRepo.Insert(context.Background(), cfg1)

	paths := []domain.TFTPPathRatio{{
		TFTPConfigID: a.ID,
		ActualPath:   "files/test.txt",
		VirtualPath:  "MAC123123/test.txt",
	}, {
		TFTPConfigID: a.ID,
		ActualPath:   "files/test.txt",
		VirtualPath:  "MAC322/test.txt",
	}}
	for _, pathRatio := range paths {
		_, err := tftpTester.pathsRepo.Insert(context.Background(), pathRatio)
		if err != nil {
			t.Errorf("insert path ratio failed: %q", err)
		}
	}
	err := services.TFTPServerServiceInit(tftpTester.service)
	if err != nil {
		t.Errorf("initialize tftp server failed: %q", err)
	}

	p := "files/testReceive.txt"
	c, err := tftp.NewClient(fmt.Sprintf("%s:6969", tftpTester.tftpAddress))
	if err != nil {
		t.Errorf("create tftp client failed: %q", err)
	}
	wt, err := c.Receive("MAC322/test.txt", "octet")
	if err != nil {
		t.Errorf("receive incoming file transmission failed: %q", err)
	}
	file, err := os.Create(p)
	if err != nil {
		t.Errorf("Create creates named file failed: %q", err)
	}
	// Optionally obtain transfer size before actual data.
	if n, ok := wt.(tftp.IncomingTransfer).Size(); ok {
		fmt.Printf("Transfer size: %d\n", n)
	}
	n, err := wt.WriteTo(file)
	if n < 10 {
		t.Errorf("write to file failed: %q", err)
	}
}

func Test_TFTPServerService_GetList(t *testing.T) {
	servers, err := tftpTester.service.GetServerList(context.Background(), "", "", "", 1, 1)
	if err != nil {
		t.Errorf("get tftp servers list failed: %q", err)
	}
	if len(servers.Items) != 1 {
		t.Error("wrong number of tftp servers")
	}
	if servers.Items[0].Address != tftpTester.tftpAddress {
		t.Error("wrong address of tftp server received")
	}
	tftpTester.createdServerID = servers.Items[0].ID
}

func Test_TFTPServerService_GetByID(t *testing.T) {
	server, err := tftpTester.service.GetServerByID(context.Background(), tftpTester.createdServerID)
	if err != nil {
		t.Errorf("get tftp server by id failed: %q", err)
	}
	if server.Address != tftpTester.tftpAddress {
		t.Error("wrong address of tftp server received")
	}
}

func Test_TFTPServerService_Create(t *testing.T) {
	dto := dtos.TFTPServerCreateDto{
		TFTPServerBaseDto: dtos.TFTPServerBaseDto{
			Address: "111.111.111.111",
			Port:    "6969",
			Enabled: false,
		},
	}

	server, err := tftpTester.service.CreateServer(context.Background(), dto)
	if err != nil {
		t.Errorf("create tftp server failed: %q", err)
	}
	if server.Address != dto.Address && server.Port != dto.Port && server.Enabled != dto.Enabled {
		t.Error("received wrong fields on create server")
	}
	tftpTester.deleteServerID = server.ID
}

func Test_TFTPServerService_Update(t *testing.T) {
	dto := dtos.TFTPServerUpdateDto{
		TFTPServerBaseDto: dtos.TFTPServerBaseDto{
			Address: "222.222.222.222",
			Port:    "6969",
			Enabled: false,
		},
	}

	server, err := tftpTester.service.UpdateServer(context.Background(), dto, tftpTester.deleteServerID)
	if err != nil {
		t.Errorf("update tftp server failed: %q", err)
	}
	if server.Address != dto.Address && server.Port != dto.Port && server.Enabled != dto.Enabled {
		t.Error("received wrong fields on update server")
	}
}

func Test_TFTPServerService_Delete(t *testing.T) {
	err := tftpTester.service.DeleteServer(context.Background(), tftpTester.deleteServerID)
	if err != nil {
		t.Errorf("delete tftp server failed: %q", err)
	}
	_, err = tftpTester.service.GetServerByID(context.Background(), tftpTester.deleteServerID)
	if err == nil {
		t.Error("received deleted tftp server")
	}
}

func Test_TFTPServerService_GetPaths(t *testing.T) {
	paths, err := tftpTester.service.GetPathsList(context.Background(), tftpTester.createdServerID, "", "", 1, 10)
	if err != nil {
		t.Errorf("get server paths failed: %q", err)
	}
	if len(paths.Items) != 2 {
		t.Error("wrong paths number")
	}
}

func Test_TFTPServerService_CreatePath(t *testing.T) {
	dto := dtos.TFTPPathCreateDto{
		TFTPPathBaseDto: dtos.TFTPPathBaseDto{
			ActualPath:  "test/path",
			VirtualPath: "my/test/path",
		}}

	createdPath, err := tftpTester.service.CreatePath(context.Background(), tftpTester.createdServerID, dto)
	if err != nil {
		t.Errorf("create server paths failed: %q", err)
	}
	if createdPath.ActualPath != dto.ActualPath && createdPath.VirtualPath != dto.VirtualPath {
		t.Error("received wrong fields on created path")
	}
	tftpTester.createdPathID = createdPath.ID
}

func Test_TFTPServerService_DeletePath(t *testing.T) {
	err := tftpTester.service.DeletePath(context.Background(), tftpTester.createdServerID, tftpTester.createdPathID)
	if err != nil {
		t.Errorf("delete server paths failed: %q", err)
	}
	paths, err := tftpTester.service.GetPathsList(context.Background(), tftpTester.createdServerID, "", "", 1, 10)
	if err != nil {
		t.Errorf("get server paths failed: %q", err)
	}

	for _, p := range paths.Items {
		if p.ActualPath == "test/path" {
			t.Error("wrong paths number")
		}
	}
}

func Test_TFTPServerService_CleaningAfterTests(t *testing.T) {
	err := os.Remove(tftpTester.dbFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
	err = os.Remove("files/test.txt")
	if err != nil {
		t.Errorf("remove test file failed:  %q", err)
	}
	err = os.Remove("files/testReceive.txt")
	if err != nil {
		t.Errorf("remove test file failed:  %q", err)
	}
}
