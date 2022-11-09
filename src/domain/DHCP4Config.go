package domain

//DHCP4Config configuration for dhcp v4 server
type DHCP4Config struct {
	//EntityUUID - nested base entity where ID type is uuid.UUID
	EntityUUID
	Interface string `gorm:"type:varchar(64);index"`
	Gateway   string `gorm:"type:varchar(15)"`
	NTP       string `gorm:"type:varchar(15)"`
	//ServerID server id DHCP option
	ServerID   string `gorm:"type:varchar(15)"`
	NextServer string `gorm:"type:varchar(15)"`
	Mask       string `gorm:"type:varchar(15)"`
	DNS        string
	Range      string
	Enabled    bool
	Port       int
	LeaseTime  int
}
