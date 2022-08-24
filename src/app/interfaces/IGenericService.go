package interfaces

import (
	"context"
	"rol/dtos"

	"github.com/google/uuid"
)

//IGenericService generic service interface for IEntityModel and it's DTOs
type IGenericService[DtoType interface{},
	CreateDtoType interface{},
	UpdateDtoType interface{},
	EntityType IEntityModel] interface {
	//GetList Get list of elements with search and pagination.
	//Params
	//	search - string value to search
	//	orderBy - order field name
	//	orderDirection - Order direction, desc/asc
	//	page - number of the page
	//	pageSize - size of the page
	//Return
	//	dtos.PaginatedItemsDto[DtoType] - pointer to the struct with pagination info and entities
	//	error - if an error occurred, otherwise nil
	GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[DtoType], error)
	//GetByID Get entity by ID from service.
	//Params
	//	id - entity id
	//Return
	//	*DtoType - pointer to the entity DTO.
	//	error - if an error occurred, otherwise nil
	GetByID(ctx context.Context, id uuid.UUID) (DtoType, error)
	//Update Save the changes to the existing entity in the service.
	//Params
	//  updateDto - pointer to the entity DTO.
	//	id - entity id
	//Return
	// 	DtoType - updated entity dto
	//	error - if an error occurred, otherwise nil
	Update(ctx context.Context, updateDto UpdateDtoType, id uuid.UUID) (DtoType, error)
	//Create Add entity to the service.
	//Params
	//  createDto - pointer to the entity DTO.
	//Return
	// 	DtoType - created entity dto
	//	error - if an error occurred, otherwise nil
	Create(ctx context.Context, createDto CreateDtoType) (DtoType, error)
	//Delete entity from the service.
	//Params
	//	id - entity id
	//Return
	//	error - if an error occurred, otherwise nil
	Delete(ctx context.Context, id uuid.UUID) error
}
