package router

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type Response struct {
	Context *gin.Context
}

// SendTemplate render HTML template
func (r *Response) SendTemplate(name string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + name + ".html")
	if err != nil {
		r.Context.String(http.StatusInternalServerError, "Template error: %v", err)
		return
	}
	r.Context.Header("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(r.Context.Writer, data)
	if err != nil {
		r.Context.String(http.StatusInternalServerError, "Render error: %v", err)
	}
}

// SendJSON send JSON response
func (r *Response) SendJSON(data interface{}) {
	r.Context.JSON(http.StatusOK, data)
}

// SendXML send XML response
func (r *Response) SendXML(data interface{}) {
	r.Context.XML(http.StatusOK, data)
}
