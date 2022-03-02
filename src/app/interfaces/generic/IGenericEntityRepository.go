package generic

import "rol/app/interfaces"

type IGenericEntityRepository interface {
	//GetList
	//	Get list of elements with filtering and pagination.
	//Params
	//  entities - pointer to array of pointers to the entities.
	//  orderBy - order field name
	//	orderDirection - Order direction, desc/asc
	//	page - number of the page
	//	size - size of the page
	//	query - query string
	//  args - args for query string
	//Return
	//	error - if an error occurred, otherwise nil
	GetList(entities interface{}, orderBy string, orderDirection string, page int, size int, query string, args ...interface{}) (int64, error)
	//GetAll
	//	Get all entities from repository.
	//Params
	//  entities - pointer to array of pointers to the entities.
	//Return
	//	error - if an error occurred, otherwise nil
	GetAll(entities interface{}) error
	//GetById
	//	Get entity by ID from repository.
	//Params
	//  entity - pointer to the entity.
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	GetById(entity interfaces.IEntityModel, id uint) error
	//Update
	//	Save the changes to the existing entity in the repository.
	//Params
	//  entity - pointer to the entity.
	//Return
	//	error - if an error occurred, otherwise nil
	Update(entity interfaces.IEntityModel) error
	//Insert
	//	Add entity to the repository.
	//Params
	//  entity - pointer to the entity.
	//Return
	//	uint - entity id
	//	error - if an error occurred, otherwise nil
	Insert(entity interfaces.IEntityModel) (uint, error)
	//Delete
	//	Mark an entity as deleted in the repository.
	//Params
	//  entity - pointer to the entity.
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(entity interfaces.IEntityModel) error
}
