package dtos

//ResponseStatusDto DTO structure for response status
type ResponseStatusDto struct {
	//	Code - 0 if OK, otherwise 1
	Code int `json:"code"`
	//	Message - response message
	Message string `json:"message"`
}

//ResponseDataDto STO structure for response with DATA
type ResponseDataDto struct {
	//	Status - response status structure
	Status ResponseStatusDto `json:"status"`
	//	Data - any data to send
	Data interface{} `json:"data"`
}

//ResponseDto DTO structure for response without DATA
type ResponseDto struct {
	//	Status - response status structure
	Status ResponseStatusDto `json:"status"`
}
