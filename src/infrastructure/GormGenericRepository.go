package infrastructure

import (
	"context"
	"fmt"
	"reflect"
	"rol/app/interfaces"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormGenericRepository[EntityType interfaces.IEntityModel] struct {
	Db            *gorm.DB
	logger        *logrus.Logger
	logSourceName string
}

//NewGormGenericRepository GORM generic repository constructor
//Params
//	*gorm.DB - gorm database
//	*logrus.Logger - logrus logger
//Return
//	*GormGenericRepository[EntityType] - repository for instantiated entity
func NewGormGenericRepository[EntityType interfaces.IEntityModel](db *gorm.DB, log *logrus.Logger) *GormGenericRepository[EntityType] {
	model := new(EntityType)
	return &GormGenericRepository[EntityType]{
		Db:            db,
		logger:        log,
		logSourceName: fmt.Sprintf("GormGenericRepository<%s>", reflect.TypeOf(*model).Name()),
	}
}

func (ger *GormGenericRepository[EntityType]) log(ctx context.Context, level, message string) {
	if ctx != nil {
		actionId := uuid.UUID{}
		if ctx.Value("requestId") != nil {
			actionId = ctx.Value("requestId").(uuid.UUID)
		}

		entry := ger.logger.WithFields(logrus.Fields{
			"actionId": actionId,
			"source":   ger.logSourceName,
		})
		switch level {
		case "err":
			entry.Error(message)
		case "info":
			entry.Info(message)
		case "warn":
			entry.Warn(message)
		case "debug":
			entry.Debug(message)
		}
	}
}

//NewQueryBuilder
//	Gets new query builder
//Params
//	ctx - context is used only for logging
//Return
//	interfaces.IQueryBuilder - new query builder
func (ger *GormGenericRepository[EntityType]) NewQueryBuilder(ctx context.Context) interfaces.IQueryBuilder {
	ger.log(ctx, "debug", "Call method NewQueryBuilder")
	return NewGormQueryBuilder()
}

func (ger *GormGenericRepository[EntityType]) addQueryToGorm(gormQuery *gorm.DB, queryBuilder interfaces.IQueryBuilder) error {
	if queryBuilder != nil {
		query, err := queryBuilder.Build()
		if err != nil {
			return err
		}
		arrQuery := query.([]interface{})
		// TODO: We need more checks here
		switch arrQuery[0].(type) {
		case string:
			queryString := arrQuery[0].(string)
			queryArgs := make([]interface{}, 0)
			for i := 1; i < len(arrQuery); i++ {
				queryArgs = append(queryArgs, arrQuery[i])
			}
			gormQuery.Where(queryString, queryArgs...)
		}
	}
	return nil
}

func generateOrderString(orderBy string, orderDirection string) string {
	order := ""
	if len(orderBy) > 0 {
		order = orderBy
		if len(orderDirection) > 0 {
			order = order + " " + orderDirection
		}
	}
	if len(order) < 1 {
		order = "created_at"
	}
	return order
}

//GetList
//	Get list of elements with filtering and pagination
//Params
//	ctx - context is used only for logging
//	orderBy - order by string parameter
//	orderDirection - ascending or descending order
//	page - page number
//	size - page size
//	queryBuilder - query builder for filtering
//Return
//	*[]EntityType - pointer to array of entities
//	error - if an error occurs, otherwise nil
func (ger *GormGenericRepository[EntityType]) GetList(ctx context.Context, orderBy string, orderDirection string, page int, size int, queryBuilder interfaces.IQueryBuilder) (*[]EntityType, error) {
	ger.log(ctx, "debug", fmt.Sprintf("GetList: IN: orderBy=%s, orderDirection=%s, page=%d, size=%d, queryBuilder=%s", orderBy, orderDirection, page, size, queryBuilder))
	model := new(EntityType)
	entities := &[]EntityType{}
	offset := int64((page - 1) * size)
	if len(orderBy) > 1 {
		orderBy = ToSnakeCase(orderBy)
	}
	orderString := generateOrderString(orderBy, orderDirection)
	gormQuery := ger.Db.Model(&model).Order(orderString)
	err := ger.addQueryToGorm(gormQuery, queryBuilder)
	if err != nil {
		return nil, err
	}
	err = gormQuery.Offset(int(offset)).Limit(size).Find(entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

//Count
//	Gets total count of entities with current query
//Params
//	ctx - context is used only for logging
//	queryBuilder - query for entities to count
//Return
//	int64 - number of entities
//	error - if an error occurs, otherwise nil
func (ger *GormGenericRepository[EntityType]) Count(ctx context.Context, queryBuilder interfaces.IQueryBuilder) (int64, error) {
	ger.log(ctx, "debug", fmt.Sprintf("Count: IN: queryBuilder=%+v", queryBuilder))
	count := int64(0)
	model := new(EntityType)
	gormQuery := ger.Db.Model(&model)
	err := ger.addQueryToGorm(gormQuery, queryBuilder)
	if err != nil {
		return 0, err
	}
	err = gormQuery.Count(&count).Error
	if err != nil {
		return 0, err
	}
	ger.log(ctx, "debug", fmt.Sprintf("Count: OUT: count=%d", count))
	return count, nil
}

//GetById
//	Get entity by ID from repository
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	*EntityType - point to entity
//	error - if an error occurs, otherwise nil
func (ger *GormGenericRepository[EntityType]) GetById(ctx context.Context, id uuid.UUID) (*EntityType, error) {
	ger.log(ctx, "debug", fmt.Sprintf("GetByID: id=%d", id))
	entity := new(EntityType)
	err := ger.Db.First(entity, id).Error
	if err != nil {
		return nil, err
	}
	ger.log(ctx, "debug", fmt.Sprintf("GetByID: entity=%+v", entity))
	return entity, nil
}

//Update
//	Save the changes to the existing entity in the repository
//Params
//	ctx - context is used only for logging
//	entity - updated entity to save
//Return
//	error - if an error occurs, otherwise nil
func (ger *GormGenericRepository[EntityType]) Update(ctx context.Context, entity *EntityType) error {
	ger.log(ctx, "debug", fmt.Sprintf("Update: entity=%+v", entity))
	return ger.Db.Save(entity).Error
}

//Insert
//	Add entity to the repository
//Params
//	ctx - context is used only for logging
//	entity - entity to save
//Return
//	uuid.UUID - new entity id
//	error - if an error occurs, otherwise nil
func (ger *GormGenericRepository[EntityType]) Insert(ctx context.Context, entity EntityType) (uuid.UUID, error) {
	ger.log(ctx, "debug", fmt.Sprintf("Insert: entity=%+v", entity))
	if err := ger.Db.Create(&entity).Error; err != nil {
		return uuid.UUID{}, err
	}
	ger.log(ctx, "debug", fmt.Sprintf("Insert: newId=%d", entity.GetID()))
	return entity.GetID(), nil
}

//Delete entity from the database
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	error - if an error occurs, otherwise nil
func (ger *GormGenericRepository[EntityType]) Delete(ctx context.Context, id uuid.UUID) error {
	ger.log(ctx, "debug", fmt.Sprintf("Delete: id=%d", id))
	entity := new(EntityType)
	gormQuery := ger.Db.Model(entity).Select(clause.Associations)
	return gormQuery.Delete(entity, id).Error
}

//CloseDb Closes current database connection
//Return
//	error - if an error occurs, otherwise nil
func (ger *GormGenericRepository[EntityType]) CloseDb() error {
	sqlDb, err := ger.Db.DB()
	if err != nil {
		return err
	}
	if err := sqlDb.Close(); err != nil {
		return err
	}
	return nil
}
