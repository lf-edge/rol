// Package domain stores the main structures of the program
package domain

import "github.com/google/uuid"

//TFTPPathRatio TFTP path ratio entity
type TFTPPathRatio struct {
	EntityUUID
	//TFTPConfigID TFTP config ID
	TFTPConfigID uuid.UUID
	//ActualPath actual file path
	ActualPath string
	//VirtualPath virtual file path
	VirtualPath string
}
