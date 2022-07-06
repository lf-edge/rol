package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
