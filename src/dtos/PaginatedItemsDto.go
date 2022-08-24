package dtos

//PaginationInfoDto struct for paginated result response
type PaginationInfoDto struct {
	//Page - page number
	Page int
	//Size - page size
	Size int
	//TotalCount - total number of items
	TotalCount int
	//TotalPages - total number of pages
	TotalPages int
}

//PaginatedItemsDto struct with slice of items and pagination info
type PaginatedItemsDto[T any] struct {
	//Pagination info about pagination
	Pagination PaginationInfoDto
	//Items slice of items
	Items []T
}

func pageCount(totalItems, pageSize int) int {
	pages := totalItems / pageSize

	if totalItems%pageSize > 0 {
		pages++
	}
	if pages < 1 {
		return 1
	}
	return pages
}

//NewEmptyPaginatedItemsDto constructor for PaginatedItemsDto[DataType] without data, used for errors
func NewEmptyPaginatedItemsDto[T any]() PaginatedItemsDto[T] {
	return PaginatedItemsDto[T]{
		Pagination: PaginationInfoDto{
			Page:       1,
			Size:       10,
			TotalCount: 0,
			TotalPages: 1,
		},
		Items: []T{},
	}
}

//NewPaginatedItemsDto constructor for PaginatedItemsDto[DataType]
func NewPaginatedItemsDto[T any](page, size, totalCount int, items []T) PaginatedItemsDto[T] {
	return PaginatedItemsDto[T]{
		Pagination: PaginationInfoDto{
			Page:       page,
			Size:       size,
			TotalCount: totalCount,
			TotalPages: pageCount(totalCount, size),
		},
		Items: items,
	}
}
