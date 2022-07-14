package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"rol/app/errors"
)

//RestoreBody restore body in gin.Context for logging it later in middleware
func RestoreBody(body interface{}, ctx *gin.Context) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return nil
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
