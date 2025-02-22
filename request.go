package router

import "github.com/gin-gonic/gin"

// Request struct to work with input data
type Request struct {
	Context *gin.Context
}

// GetParam return value of GET param
func (req *Request) GetParam(name string) string {
	return req.Context.Query(name)
}

// GetBody read JSON body to struct
func (req *Request) GetBody(dest interface{}) error {
	return req.Context.BindJSON(dest)
}
