package dtos

//PaginatedListDto paginated list DTO structure
type PaginatedListDto[DtoType any] struct {
	//	Page - page number
	Page int
	//	Size - page size
	Size int
	//	Total - total number of items
	Total int64
	//	Items - array of DTO items
	Items *[]DtoType
}
