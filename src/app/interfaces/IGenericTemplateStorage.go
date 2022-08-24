package interfaces

import (
	"context"
)

//IGenericTemplateStorage is the interface is needed to get templates
type IGenericTemplateStorage[TemplateType interface{}] interface {
	//GetByName get template by name
	//Params
	//	ctx - usually http context gin.Context
	//	name - template name
	//Return
	//	*TemplateType - pointer to the struct
	//	error - if an error occurred, otherwise nil
	GetByName(ctx context.Context, templateName string) (TemplateType, error)
	//GetList get list of templates with pagination.
	//Params
	//	ctx - usually http context gin.Context
	//	orderBy - order field name
	//	orderDirection - order direction, desc/asc
	//	page - number of the page
	//	pageSize - size of the page
	//	queryBuilder - is an instance of IQueryBuilder
	//Return
	//	*[]TemplateType - pointer to the struct with pagination info and entities
	//	error - if an error occurred, otherwise nil
	GetList(ctx context.Context, orderBy, orderDirection string, page, pageSize int, queryBuilder IQueryBuilder) ([]TemplateType, error)
	//Count
	// Get count of entities with filtering
	//Params
	//	ctx - usually http context gin.Context
	//	queryBuilder - is an instance of IQueryBuilder
	//Return
	//	int64 - count of templates
	//	error - if an error occurred, otherwise nil
	Count(ctx context.Context, queryBuilder IQueryBuilder) (int64, error)
	//NewQueryBuilder
	//	Get QueryBuilder
	//Return
	//	IQueryBuilder pointer to object that implements IQueryBuilder interface
	NewQueryBuilder(ctx context.Context) IQueryBuilder
}
