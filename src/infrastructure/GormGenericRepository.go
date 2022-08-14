package infrastructure

import (
	"context"
	"fmt"
	"reflect"
	"rol/app/errors"
	"rol/app/interfaces"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//GormGenericRepository is implementation of interfaces.IGenericRepository
type GormGenericRepository[EntityType interfaces.IEntityModel] struct {
	//Db - gorm database
	Db *gorm.DB
	//logger - logrus logger
	logger *logrus.Logger
	//logSourceName - logger recording source
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

func (g *GormGenericRepository[EntityType]) log(ctx context.Context, level, message string) {
	if ctx != nil {
		actionID := uuid.UUID{}
		if ctx.Value("requestID") != nil {
			actionID = ctx.Value("requestID").(uuid.UUID)
		}

		entry := g.logger.WithFields(logrus.Fields{
			"actionID": actionID,
			"source":   g.logSourceName,
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
func (g *GormGenericRepository[EntityType]) NewQueryBuilder(ctx context.Context) interfaces.IQueryBuilder {
	g.log(ctx, "debug", "Call method NewQueryBuilder")
	return NewGormQueryBuilder()
}

func (g *GormGenericRepository[EntityType]) addQueryToGorm(gormQuery *gorm.DB, queryBuilder interfaces.IQueryBuilder) error {
	if queryBuilder != nil {
		query, err := queryBuilder.Build()
		if err != nil {
			return errors.Internal.Wrap(err, "error building a query")
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
func (g *GormGenericRepository[EntityType]) GetList(ctx context.Context, orderBy string, orderDirection string, page int, size int, queryBuilder interfaces.IQueryBuilder) (*[]EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("GetList: IN: orderBy=%s, orderDirection=%s, page=%d, size=%d, queryBuilder=%s", orderBy, orderDirection, page, size, queryBuilder))
	model := new(EntityType)
	entities := &[]EntityType{}
	offset := int64((page - 1) * size)
	if len(orderBy) > 1 {
		orderBy = ToSnakeCase(orderBy)
	}
	orderString := generateOrderString(orderBy, orderDirection)
	gormQuery := g.Db.Model(&model).Order(orderString)
	err := g.addQueryToGorm(gormQuery, queryBuilder)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "adding query to gorm failed")
	}
	err = gormQuery.Offset(int(offset)).Limit(size).Find(entities).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error finding entities with gorm query")
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
func (g *GormGenericRepository[EntityType]) Count(ctx context.Context, queryBuilder interfaces.IQueryBuilder) (int64, error) {
	g.log(ctx, "debug", fmt.Sprintf("Count: IN: queryBuilder=%+v", queryBuilder))
	count := int64(0)
	model := new(EntityType)
	gormQuery := g.Db.Model(&model)
	err := g.addQueryToGorm(gormQuery, queryBuilder)
	if err != nil {
		return 0, errors.Internal.Wrap(err, "adding query to gorm failed")
	}
	err = gormQuery.Count(&count).Error
	if err != nil {
		return 0, errors.Internal.Wrap(err, "gorm failed counting entities")
	}
	g.log(ctx, "debug", fmt.Sprintf("Count: OUT: count=%d", count))
	return count, nil
}

//GetByID
//	Get entity by ID from repository
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	*EntityType - point to entity
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityType]) GetByID(ctx context.Context, id uuid.UUID) (*EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("GetByID: id=%d", id))
	entity := new(EntityType)
	err := g.Db.First(entity, id).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error finding first record")
	}
	g.log(ctx, "debug", fmt.Sprintf("GetByID: entity=%+v", entity))
	return entity, nil
}

//GetByIDExtended Get entity by ID and query from repository
//Params
//	ctx - context is used only for logging
//	id - entity id
//	queryBuilder - extended query conditions
//Return
//	*EntityType - point to entity
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityType]) GetByIDExtended(ctx context.Context, id uuid.UUID, queryBuilder interfaces.IQueryBuilder) (*EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("GetByIDExtended: id=%s, query builder: %s", id, queryBuilder))
	model := new(EntityType)
	gormQuery := g.Db.Model(model)
	fullQueryBuilder := g.NewQueryBuilder(ctx)
	fullQueryBuilder.Where("ID", "==", id)
	if queryBuilder != nil {
		fullQueryBuilder.WhereQuery(queryBuilder)
	}
	err := g.addQueryToGorm(gormQuery, fullQueryBuilder)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed add query to gorm query")
	}
	entities := &[]EntityType{}
	var entity *EntityType
	err = gormQuery.Find(entities).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error finding entities with gorm query")
	}
	if len(*entities) < 1 {
		return nil, nil
	}
	entity = &(*entities)[0]
	g.log(ctx, "debug", fmt.Sprintf("GetByID: entity=%+v", entity))
	return entity, nil
}

//Update
//	Save the changes to the existing entity in the repository
//Params
//	ctx - context is used only for logging
//	entity - updated entity to save
//Return
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityType]) Update(ctx context.Context, entity *EntityType) error {
	g.log(ctx, "debug", fmt.Sprintf("Update: entity=%+v", entity))
	err := g.Db.Save(entity).Error
	if err != nil {
		return errors.Internal.Wrap(err, "error saving entity")
	}
	return nil
}

//Insert
//	Add entity to the repository
//Params
//	ctx - context is used only for logging
//	entity - entity to save
//Return
//	uuid.UUID - new entity id
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityType]) Insert(ctx context.Context, entity EntityType) (uuid.UUID, error) {
	g.log(ctx, "debug", fmt.Sprintf("Insert: entity=%+v", entity))
	if err := g.Db.Create(&entity).Error; err != nil {
		return uuid.UUID{}, errors.Internal.Wrap(err, "gorm failed create entity")
	}
	g.log(ctx, "debug", fmt.Sprintf("Insert: newID=%d", entity.GetID()))
	return entity.GetID(), nil
}

//Delete entity from the database
//Params
//	ctx - context is used only for logging
//	id - entity id
//Return
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityType]) Delete(ctx context.Context, id uuid.UUID) error {
	g.log(ctx, "debug", fmt.Sprintf("Delete: id=%d", id))
	entity := new(EntityType)
	gormQuery := g.Db.Model(entity).Select(clause.Associations)
	err := gormQuery.Delete(entity, id).Error
	if err != nil {
		return errors.Internal.Wrap(err, "gorm failed delete entity")
	}
	return nil
}

//CloseDb Closes current database connection
//Return
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityType]) CloseDb() error {
	sqlDb, err := g.Db.DB()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get db connection")
	}
	err = sqlDb.Close()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to close db connection")
	}
	return nil
}
