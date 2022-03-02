package mappers

import (
	"rol/domain/entities"
)

func GetEmptyEntityFromArrayType(entityArr interface{}) interface{} {
	switch entityArr.(type) {
	case *[]*entities.EthernetSwitch:
		return &entities.EthernetSwitch{}
	case *[]*entities.EthernetSwitchPort:
		return &entities.EthernetSwitchPort{}
	}
	return nil
}
