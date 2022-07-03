package interfaces

import (
	"context"
	"github.com/google/uuid"
)

//IGenericRepository generic repository interface for IEntityModel
type IGenericRepository[EntityType IEntityModel] interface {
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
	//	 *[]EntityType - pointer to array of the entities.
	//	error - if an error occurred, otherwise nil
	GetList(ctx context.Context, orderBy string, orderDirection string, page int, size int, queryBuilder IQueryBuilder) (*[]EntityType, error)
	//Count
	// Get count of entities with filtering
	//Params
	//	ctx - context
	//	queryBuilder - QueryBuilder is repo.NewQueryBuilder()
	//Return
	//	int64 - count of entities
	//	error - if an error occurred, otherwise nil
	Count(ctx context.Context, queryBuilder IQueryBuilder) (int64, error)
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
	//  *EntityType - pointer to the entity.
	//	error - if an error occurred, otherwise nil
	GetByID(ctx context.Context, id uuid.UUID) (*EntityType, error)
	//GetByIDExtended Get entity by ID and query from repository
	//Params
	//	ctx - context is used only for logging
	//	id - entity id
	//	queryBuilder - extended query conditions
	//Return
	//	*EntityType - point to entity if existed, otherwise nil
	//	error - if an error occurs, otherwise nil
	GetByIDExtended(ctx context.Context, id uuid.UUID, queryBuilder IQueryBuilder) (*EntityType, error)
	//Update
	//	Save the changes to the existing entity in the repository.
	//Params
	//	ctx - context
	//  entity - pointer to the entity.
	//Return
	//	error - if an error occurred, otherwise nil
	Update(ctx context.Context, entity *EntityType) error
	//Insert
	//	Add entity to the repository.
	//Params
	//	ctx - context
	//  entity - pointer to the entity for update.
	//Return
	//	uuid.UUID - entity id
	//	error - if an error occurred, otherwise nil
	Insert(ctx context.Context, entity EntityType) (uuid.UUID, error)
	//Delete
	//	Mark an entity as deleted in the repository.
	//Params
	//	ctx - context
	//  id - id of the entity for delete.
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(ctx context.Context, id uuid.UUID) error
	//CloseDb
	//	Closes current database connection
	//Return
	//	error - if an error occurred, otherwise nil
	CloseDb() error
}
