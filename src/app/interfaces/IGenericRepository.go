package interfaces

import (
	"context"
)

//IGenericRepository generic repository interface for IEntityModel
type IGenericRepository[EntityIDType comparable, EntityType IEntityModel[EntityIDType]] interface {
	//GetList
	//	Get list of elements with filtering and pagination.
	//Params
	//	ctx - context
	//  orderBy - order field name
	//	orderDirection - Order direction, desc/asc
	//	page - number of the page
	//	size - size of the page
	//	queryBuilder - QueryBuilder is repo.NewQueryBuilder()
	//Return
	//	 []EntityType - array of the entities.
	//	error - if an error occurred, otherwise nil
	GetList(ctx context.Context, orderBy string, orderDirection string, page int, size int, queryBuilder IQueryBuilder) ([]EntityType, error)
	//Count
	// Get count of entities with filtering
	//Params
	//	ctx - context
	//	queryBuilder - QueryBuilder is repo.NewQueryBuilder()
	//Return
	//	int64 - count of entities
	//	error - if an error occurred, otherwise nil
	Count(ctx context.Context, queryBuilder IQueryBuilder) (int, error)
	//IsExist checks that entity is existed in repository
	//
	//Params
	//	ctx - context is used only for logging
	//	id - id of the entity
	//	queryBuilder - query builder with addition conditions, can be nil
	//Return
	//	bool - true if existed, otherwise false
	//	error - if an error occurs, otherwise nil
	IsExist(ctx context.Context, id EntityIDType, queryBuilder IQueryBuilder) (bool, error)
	//NewQueryBuilder
	//	Get QueryBuilder
	//Params
	//	ctx - context
	//Return
	//	IQueryBuilder pointer to object that implements IQueryBuilder interface for this repository
	NewQueryBuilder(ctx context.Context) IQueryBuilder
	//GetByID
	//	Get entity by ID from repository.
	//Params
	//	ctx - context
	//	id - entity id
	//Return
	//  EntityType - entity from repository
	//	error - if an error occurred, otherwise nil
	GetByID(ctx context.Context, id EntityIDType) (EntityType, error)
	//GetByIDExtended Get entity by ID and query from repository
	//Params
	//	ctx - context is used only for logging
	//	id - entity id
	//	queryBuilder - extended query conditions
	//Return
	//  EntityType - entity from repository
	//	error - if an error occurs, otherwise nil
	GetByIDExtended(ctx context.Context, id EntityIDType, queryBuilder IQueryBuilder) (EntityType, error)
	//Update
	//	Save the changes to the existing entity in the repository.
	//Params
	//	ctx - context
	//  entity - pointer to the entity.
	//Return
	//	EntityType - updated entity
	//	error - if an error occurred, otherwise nil
	Update(ctx context.Context, entity EntityType) (EntityType, error)
	//Insert
	//	Add entity to the repository.
	//Params
	//	ctx - context
	//  entity - pointer to the entity for update.
	//Return
	//	uuid.UUID - entity id
	//	error - if an error occurred, otherwise nil
	Insert(ctx context.Context, entity EntityType) (EntityType, error)
	//Delete
	//	Mark an entity as deleted in the repository.
	//Params
	//	ctx - context
	//  id - id of the entity for delete.
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(ctx context.Context, id EntityIDType) error
	//DeleteAll entities matching the condition
	//
	//Params
	//	ctx - context is used only for logging
	//	queryBuilder - query builder with conditions
	//Return
	//	error - if an error occurs, otherwise nil
	DeleteAll(ctx context.Context, queryBuilder IQueryBuilder) error
	//Dispose releases all resources
	//
	//Return
	//	error - if an error occurred, otherwise nil
	Dispose() error
}
