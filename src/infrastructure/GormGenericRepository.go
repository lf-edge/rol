package infrastructure

import (
	"context"
	def_err "errors"
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
type GormGenericRepository[EntityIDType comparable, EntityType interfaces.IEntityModel[EntityIDType]] struct {
	//Db - gorm database
	Db *gorm.DB
	//logger - logrus logger
	logger *logrus.Logger
	//logSourceName - logger recording source
	logSourceName string
}

//NewGormGenericRepository GORM generic repository constructor
//
//Params
//	*gorm.DB - gorm database
//	*logrus.Logger - logrus logger
//Return
//	*GormGenericRepository[EntityType] - repository for instantiated entity
func NewGormGenericRepository[EntityIDType comparable, EntityType interfaces.IEntityModel[EntityIDType]](db *gorm.DB, log *logrus.Logger) *GormGenericRepository[EntityIDType, EntityType] {
	model := new(EntityType)
	repo := &GormGenericRepository[EntityIDType, EntityType]{
		Db:            db,
		logger:        log,
		logSourceName: fmt.Sprintf("GormGenericRepository<%s>", reflect.TypeOf(*model).Name()),
	}
	err := repo.Db.Callback().Create().Before("gorm:before_create").Replace("rol:before_create", repo.beforeCreate)
	if err != nil {
		panic(err)
	}
	return repo
}

// BeforeCreate hook with uuid.UUID id set
func (g *GormGenericRepository[EntityIDType, EntityType]) beforeCreate(db *gorm.DB) {
	if db.Statement.Schema != nil {
		idField := db.Statement.Schema.PrioritizedPrimaryField
		switch idField.FieldType {
		case reflect.TypeOf(uuid.UUID{}):
			idAlreadySet := false
			if fieldValue, isZero := idField.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
				if id, ok := fieldValue.(uuid.UUID); ok {
					if id != uuid.Nil {
						idAlreadySet = true
					}
				}
			}
			if !idAlreadySet {
				err := idField.Set(db.Statement.Context, db.Statement.ReflectValue, uuid.New())
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (g *GormGenericRepository[EntityIDType, EntityType]) log(ctx context.Context, level, message string) {
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

//NewQueryBuilder gets new query builder
//
//Params
//	ctx - context is used only for logging
//Return
//	interfaces.IQueryBuilder - new query builder
func (g *GormGenericRepository[EntityIDType, EntityType]) NewQueryBuilder(ctx context.Context) interfaces.IQueryBuilder {
	g.log(ctx, "debug", "Call method NewQueryBuilder")
	return NewGormQueryBuilder()
}

func (g *GormGenericRepository[EntityIDType, EntityType]) addQueryToGorm(gormQuery *gorm.DB, queryBuilder interfaces.IQueryBuilder) error {
	if queryBuilder != nil {
		query, err := queryBuilder.Build()
		if err != nil {
			return errors.Internal.Wrap(err, "error building a query")
		}
		arrQuery := query.([]interface{})
		if len(arrQuery) > 0 {
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
		order = "created_at desc"
	}
	return order
}

//GetList of elements with filtering and pagination
//
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
func (g *GormGenericRepository[EntityIDType, EntityType]) GetList(ctx context.Context, orderBy string, orderDirection string, page int, size int, queryBuilder interfaces.IQueryBuilder) ([]EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("GetList: IN: orderBy=%s, orderDirection=%s, page=%d, size=%d, queryBuilder=%s", orderBy, orderDirection, page, size, queryBuilder))
	model := new(EntityType)
	entities := []EntityType{}
	offset := int64((page - 1) * size)
	if len(orderBy) > 1 {
		orderBy = ToSnakeCase(orderBy)
	}
	orderString := generateOrderString(orderBy, orderDirection)
	gormQuery := g.Db.Model(model).Order(orderString)
	err := g.addQueryToGorm(gormQuery, queryBuilder)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "adding query to gorm failed")
	}
	err = gormQuery.Offset(int(offset)).Limit(size).Find(&entities).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error finding entities with gorm query")
	}
	return entities, nil
}

//IsExist checks that entity is existed in repository
//
//Params
//	ctx - context is used only for logging
//	id - id of the entity
//	queryBuilder - query builder with addition conditions, can be nil
//Return
//	bool - true if existed, otherwise false
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) IsExist(ctx context.Context, id EntityIDType, queryBuilder interfaces.IQueryBuilder) (bool, error) {
	entity := new(EntityType)
	gormQuery := g.Db.Model(entity)
	rootQueryBuilder := g.NewQueryBuilder(ctx)
	rootQueryBuilder.Where("id", "==", id)
	if queryBuilder != nil {
		rootQueryBuilder.WhereQuery(queryBuilder)
	}
	err := g.addQueryToGorm(gormQuery, rootQueryBuilder)
	if err != nil {
		return false, errors.Internal.Wrap(err, "adding query to gorm failed")
	}
	res := gormQuery.First(entity)
	if def_err.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if res.Error != nil {
		return false, errors.Internal.Wrap(res.Error, "failed to check entity existence")
	}
	return true, nil
}

//Count gets total count of entities with current query
//
//Params
//	ctx - context
//	queryBuilder - query builder with conditions
//Return
//	int64 - number of entities
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) Count(ctx context.Context, queryBuilder interfaces.IQueryBuilder) (int, error) {
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
	return int(count), nil
}

//GetByID gets entity by ID from repository
//
//Params
//	ctx - context
//	id - entity id
//Return
//	EntityType - point to entity
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) GetByID(ctx context.Context, id EntityIDType) (EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("GetByID: id=%+v", id))
	entity := new(EntityType)
	err := g.Db.First(entity, id).Error
	if err != nil {
		if def_err.Is(err, gorm.ErrRecordNotFound) {
			return *new(EntityType), errors.NotFound.New("entity not found in database")
		}
		return *new(EntityType), errors.Internal.Wrap(err, "error finding first record")
	}
	g.log(ctx, "debug", fmt.Sprintf("GetByID: entity=%+v", entity))
	return *entity, nil
}

//GetByIDExtended Get entity by ID and query from repository
//
//Params
//	ctx - context
//	id - entity id
//	queryBuilder - extended query conditions
//Return
//	*EntityType - point to entity
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) GetByIDExtended(ctx context.Context, id EntityIDType, queryBuilder interfaces.IQueryBuilder) (EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("GetByIDExtended: id=%+v, query builder: %s", id, queryBuilder))
	model := new(EntityType)
	gormQuery := g.Db.Model(model)
	fullQueryBuilder := g.NewQueryBuilder(ctx)
	fullQueryBuilder.Where("ID", "==", id)
	if queryBuilder != nil {
		fullQueryBuilder.WhereQuery(queryBuilder)
	}
	err := g.addQueryToGorm(gormQuery, fullQueryBuilder)
	if err != nil {
		return *new(EntityType), errors.Internal.Wrap(err, "failed add query to gorm query")
	}
	entities := &[]EntityType{}
	err = gormQuery.Find(entities).Error
	if err != nil {
		return *new(EntityType), errors.Internal.Wrap(err, "error finding entities with gorm query")
	}
	if len(*entities) < 1 {
		return *new(EntityType), errors.NotFound.New("entity is not exist in repository")
	}
	g.log(ctx, "debug", fmt.Sprintf("GetByIDExtended: entity=%+v", (*entities)[0]))
	return (*entities)[0], nil
}

//Update save the changes to the existing entity in the repository
//
//Params
//	ctx - context
//	entity - updated entity to save
//Return
//	EntityType - updated entity
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) Update(ctx context.Context, entity EntityType) (EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("Update: entity=%+v", entity))
	err := g.Db.Save(&entity).Error
	if err != nil {
		if def_err.As(err, &gorm.ErrRecordNotFound) {
			return *new(EntityType), errors.NotFound.New("entity not found in database")
		}
		return *new(EntityType), errors.Internal.Wrap(err, "error finding first record")
	}
	return entity, nil
}

//Insert entity to the repository
//
//Params
//	ctx - context
//	entity - entity to save
//Return
//	EntityType - created entity
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) Insert(ctx context.Context, entity EntityType) (EntityType, error) {
	g.log(ctx, "debug", fmt.Sprintf("Insert: entity=%+v", entity))
	if err := g.Db.Create(&entity).Error; err != nil {
		return *new(EntityType), errors.Internal.Wrap(err, "gorm failed create entity")
	}
	g.log(ctx, "debug", fmt.Sprintf("Insert: newID=%+v", entity.GetID()))
	return entity, nil
}

//Delete entity from the repository
//
//Params
//	ctx - context
//	id - entity id
//Return
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) Delete(ctx context.Context, id EntityIDType) error {
	g.log(ctx, "debug", fmt.Sprintf("Delete: id=%+v", id))
	exist, err := g.IsExist(ctx, id, nil)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to check existence of the entity")
	}
	if !exist {
		return errors.NotFound.New("entity not found in database")
	}
	entity := new(EntityType)
	gormQuery := g.Db.Model(entity).Select(clause.Associations)
	err = gormQuery.Delete(entity, id).Error
	if err != nil {
		return errors.Internal.Wrap(err, "gorm failed delete entity")
	}
	return nil
}

//DeleteAll entities matching the condition
//
//Params
//	ctx - context
//	queryBuilder - query builder with conditions
//Return
//	error - if an error occurs, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) DeleteAll(ctx context.Context, queryBuilder interfaces.IQueryBuilder) error {
	g.log(ctx, "debug", fmt.Sprintf("DeleteAll: IN: queryBuilder=%+v", queryBuilder))
	entity := new(EntityType)
	gormQuery := g.Db.Model(entity).Select(clause.Associations)
	if queryBuilder != nil {
		err := g.addQueryToGorm(gormQuery, queryBuilder)
		if err != nil {
			return errors.Internal.Wrap(err, "failed add query to gorm query")
		}
	}
	res := gormQuery.Delete(entity)
	if res.Error != nil {
		return errors.Internal.Wrap(res.Error, "failed to delete entities by conditions")
	}
	return nil
}

//Dispose releases all resources
//
//Return
//	error - if an error occurred, otherwise nil
func (g *GormGenericRepository[EntityIDType, EntityType]) Dispose() error {
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
