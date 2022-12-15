// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//GormDevicePowerStateRepository repository for domain.DevicePowerState entity
type GormDevicePowerStateRepository struct {
	*GormGenericRepository[uuid.UUID, domain.DevicePowerState]
}

//NewGormDevicePowerStateRepository constructor for domain.DevicePowerState GORM generic repository
//Params
//	db - gorm database
//	log - logrus logger
//Return
//	interfaces.IGenericRepository[uuid.UUID, domain.DevicePowerState] - new device power state repository
func NewGormDevicePowerStateRepository(db *gorm.DB, log *logrus.Logger) interfaces.IGenericRepository[uuid.UUID, domain.DevicePowerState] {
	genericRepository := NewGormGenericRepository[uuid.UUID, domain.DevicePowerState](db, log)
	return &GormDevicePowerStateRepository{
		genericRepository,
	}
}

//Insert entity to the repository
//
//Params
//	ctx - context
//	entity - entity to save
//Return
//	EntityType - created entity
//	error - if an error occurs, otherwise nil
func (r *GormDevicePowerStateRepository) Insert(ctx context.Context, entity domain.DevicePowerState) (domain.DevicePowerState, error) {
	r.log(ctx, "debug", fmt.Sprintf("Insert: entity=%+v", entity))
	queryBuilder := r.NewQueryBuilder(ctx)
	queryBuilder.Where("DeviceID", "==", entity.DeviceID)
	stateForRemove, err := r.GetList(ctx, "CreatedAt", "desc", 51, 1, queryBuilder)
	if err != nil {
		return domain.DevicePowerState{}, errors.Internal.Wrapf(err, "power states count failed for device %s", entity.DeviceID.String())
	}
	if len(stateForRemove) > 0 {
		back := r.Db
		r.Db = r.Db.Unscoped()
		err := r.Delete(ctx, stateForRemove[0].ID)
		if err != nil {
			r.log(ctx, "debug", fmt.Sprintf("failed to remove oldest(51) power state for device %s", entity.DeviceID.String()))
		}
		r.Db = back
	}
	if err = r.Db.Create(&entity).Error; err != nil {
		return domain.DevicePowerState{}, errors.Internal.Wrap(err, "gorm failed create entity")
	}
	r.log(ctx, "debug", fmt.Sprintf("Insert: newID=%+v", entity.GetID()))
	return entity, nil
}
