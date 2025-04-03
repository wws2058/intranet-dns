package models

import "github.com/gin-gonic/gin"

// common request interface
type IRequest interface {
	ParseRequest() error
	DBOperation() error
	ExtraOperation() error
	GetResponse() (interface{}, error)
}

// handle common request interface
func ProcessRequest(req IRequest) (interface{}, error) {
	if errs := req.ParseRequest(); errs != nil {
		return nil, errs
	}
	if errs := req.DBOperation(); errs != nil {
		return nil, errs
	}
	if errs := req.ExtraOperation(); errs != nil {
		return nil, errs
	}
	return req.GetResponse()
}

// handle request params
func GenericBodyBinding(c *gin.Context, target interface{}) (err error) {
	c.ShouldBind(target)
	return
}
