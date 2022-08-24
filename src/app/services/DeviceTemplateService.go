package services

import (
	"context"
	"github.com/sirupsen/logrus"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/domain"
	"rol/dtos"
	"strings"
)

//DeviceTemplateService device template service structure for domain.DeviceTemplate
type DeviceTemplateService struct {
	storage interfaces.IGenericTemplateStorage[domain.DeviceTemplate]
	logger  *logrus.Logger
}

//NewDeviceTemplateService constructor for DeviceTemplateService
//Params
//	storage - generic storage for domain.DeviceTemplate
//	log - logrus logger
func NewDeviceTemplateService(storage interfaces.IGenericTemplateStorage[domain.DeviceTemplate], log *logrus.Logger) (*DeviceTemplateService, error) {
	return &DeviceTemplateService{
		storage: storage,
		logger:  log,
	}, nil
}

//GetList get list of domain.DeviceTemplate with filtering and pagination
//Params
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	dtos.PaginatedItemsDto[dtos.DeviceTemplateDto] - pointer to paginated list of device templates
//	error - if an error occurs, otherwise nil
func (d *DeviceTemplateService) GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.DeviceTemplateDto], error) {
	paginatedItemsDto := dtos.NewEmptyPaginatedItemsDto[dtos.DeviceTemplateDto]()
	pageFinal := page
	pageSizeFinal := pageSize
	if page < 1 {
		pageFinal = 1
	}
	if pageSize < 1 {
		pageSizeFinal = 10
	}
	queryBuilder := d.storage.NewQueryBuilder(ctx)
	if len(search) > 3 {
		queryBuilder = d.addSearchInAllFields(search, queryBuilder)
	}
	templatesArr, err := d.storage.GetList(ctx, orderBy, orderDirection, pageFinal, pageSizeFinal, queryBuilder)
	if err != nil {
		return paginatedItemsDto, err
	}
	count, err := d.storage.Count(ctx, queryBuilder)
	if err != nil {
		return paginatedItemsDto, errors.Internal.Wrap(err, "failed to count device templates")
	}
	dtosArr := make([]dtos.DeviceTemplateDto, len(templatesArr))
	for i := range templatesArr {
		err := mappers.MapEntityToDto((templatesArr)[i], &dtosArr[i])
		if err != nil {
			return paginatedItemsDto, errors.Internal.Wrap(err, "failed to map template to dto")
		}
	}
	paginatedItemsDto = dtos.NewPaginatedItemsDto[dtos.DeviceTemplateDto](pageFinal, pageSizeFinal, int(count), dtosArr)
	return paginatedItemsDto, nil
}

//GetByName Get ethernet switch by name
//Params
//	ctx - context is used only for logging
//	name - device template name
//Return
//	*dtos.DeviceTemplateDto - point to ethernet switch dto
//	error - if an error occurs, otherwise nil
func (d *DeviceTemplateService) GetByName(ctx context.Context, templateName string) (dtos.DeviceTemplateDto, error) {
	dto := *new(dtos.DeviceTemplateDto)
	template, err := d.storage.GetByName(ctx, templateName)
	if err != nil {
		return dto, err
	}
	mappers.MapDeviceTemplateToDto(template, &dto)
	return dto, nil
}

func (d *DeviceTemplateService) addSearchInAllFields(search string, queryBuilder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	template := new(domain.DeviceTemplate)
	stringFieldNames := &[]string{}
	utils.GetStringFieldsNames(template, stringFieldNames)
	queryGroup := d.storage.NewQueryBuilder(nil)
	for i := 0; i < len(*stringFieldNames); i++ {
		fieldName := (*stringFieldNames)[i]
		containPass := strings.Contains(strings.ToLower(fieldName), "pass")
		containKey := strings.Contains(strings.ToLower(fieldName), "key")
		if containPass || containKey {
			continue
		}
		queryGroup.Or(fieldName, "LIKE", search)
	}
	return queryBuilder.WhereQuery(queryGroup)
}
