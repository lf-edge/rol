package dtos

import "github.com/google/uuid"

//HTTPLogDto dto for http log
type HTTPLogDto struct {
	BaseDto[uuid.UUID]
	//	HTTPMethod - http method
	HTTPMethod string
	//	Domain - domain that processed the request
	Domain string
	//	RelativePath - path to the endpoint
	RelativePath string
	//	QueryParams - query params passed
	QueryParams string
	//	ClientIP - client IP address
	ClientIP string
	//	Latency - latency in milliseconds
	Latency int
	//	RequestBody - body of the request
	RequestBody string
	//	ResponseBody - body of the response
	ResponseBody string
	//	RequestHeaders - headers of the request
	RequestHeaders string
	//	CustomRequestHeaders - custom headers of the request
	CustomRequestHeaders string
}
