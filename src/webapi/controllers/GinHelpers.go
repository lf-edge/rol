package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"rol/app/errors"
	"rol/dtos"
)

func parseUUIDParam(ctx *gin.Context, paramName string) (uuid.UUID, error) {
	uuidStr := ctx.Param(paramName)
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return [16]byte{}, errors.NotFound.New("incorrect format for entity uuid")
	}
	return uuid, nil
}

//getRequestDtoAndRestoreBody parse json body to dto object and restore body in context
//for logging it later in middleware
func getRequestDtoAndRestoreBody[reqDtoType any](ctx *gin.Context) (reqDtoType, error) {
	reqDto := new(reqDtoType)
	err := ctx.ShouldBindJSON(reqDto)
	if err != nil {
		return *reqDto, errors.Validation.New("incorrect json")
	}
	buf, err := json.Marshal(reqDto)
	if err != nil {
		return *reqDto, errors.Internal.New("failed to marshal object back tp json for logging")
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return *reqDto, nil
}

func validationErrorToValidationErrorDto(err error) dtos.ValidationErrorDto {
	validationErrorDto := dtos.ValidationErrorDto{
		Message: err.Error(),
		Errors:  []dtos.ValidationErrorElemDto{},
	}
	//For now, we can wrap validation error in the top level to another typed error
	//and current error will be not validation error
	//We need fix this later
	if errors.GetType(err) == errors.Validation {
		context := errors.GetErrorContext(err)
		for key, value := range context {
			validationErrorDto.Errors = append(
				validationErrorDto.Errors,
				dtos.ValidationErrorElemDto{
					Source: "body",
					Field:  key,
					Error:  value,
				})
		}
	}
	return validationErrorDto
}

//abortWithStatusByErrorType call ctx.AbortWithStatus by error type,
//if error = nil, then abort with http.StatusInternalServerError
func abortWithStatusByErrorType(ctx *gin.Context, err error) {
	if errors.As(err, errors.Internal) || errors.As(err, errors.NoType) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	} else if errors.As(err, errors.Validation) {
		ctx.JSON(http.StatusBadRequest, validationErrorToValidationErrorDto(err))
		return
	} else if errors.As(err, errors.NotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

//handle current request by error type, if error = nil, set http status to 204 (No Content)
func handle(ctx *gin.Context, err error) {
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

//handleWithData same as handle, but
//if error = nil, set response body to data as json
//and http status code to 200
func handleWithData(ctx *gin.Context, err error, data interface{}) {
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
