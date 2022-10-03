package domain

import (
	"github.com/google/uuid"
	"time"
)

//DHCP4Lease information about DHCP v4 lease
type DHCP4Lease struct {
	Entity
	IP            string `gorm:"type:varchar(15);index"`
	MAC           string `gorm:"type:varchar(17);index"`
	Expires       time.Time
	DHCP4ConfigID uuid.UUID `gorm:"type:varchar(36);index"`
}
