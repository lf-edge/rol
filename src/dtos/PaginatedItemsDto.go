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
type PaginatedItemsDto[DataType any] struct {
	//Pagination info about pagination
	Pagination PaginationInfoDto
	//Items slice of items
	Items []DataType
}
