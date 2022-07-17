package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"rol/app/errors"
	"strconv"
)

type getListQueryParams struct {
	orderBy        string
	orderDirection string
	search         string
	page           int
	pageSize       int
}

//restoreBody restore body in gin.Context for logging it later in middleware
func restoreBody(body interface{}, ctx *gin.Context) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return nil
}

func parseGetListQueryParams(ctx *gin.Context, orderByDef string) getListQueryParams {
	orderBy := ctx.DefaultQuery("orderBy", orderByDef)
	orderDirection := ctx.DefaultQuery("orderDirection", "asc")
	search := ctx.DefaultQuery("search", "")
	page := ctx.DefaultQuery("page", "1")
	pageInt64, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		pageInt64 = 1
	}
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageSizeInt64, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		pageSizeInt64 = 10
	}
	return getListQueryParams{
		orderBy:        orderBy,
		orderDirection: orderDirection,
		search:         search,
		page:           int(pageInt64),
		pageSize:       int(pageSizeInt64),
	}
}

func parseUUIDFromUrl(ctx *gin.Context, paramName string) (uuid.UUID, error) {
	strID := ctx.Param(paramName)
	return uuid.Parse(strID)
}

func abortByTypedError(ctx *gin.Context, err error) {
	if errors.As(err, errors.Internal) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if errors.As(err, errors.Validation) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if errors.As(err, errors.NotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if errors.As(err, errors.NoType) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
