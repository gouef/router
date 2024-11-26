package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Router struct {
	routes []Group
}

type Group struct {
	relativePath string
	handlers     gin.HandlerFunc
}

func (r *Router) Run() {
	myRouter := gin.Default()

	for rg := range r.routes {
		myRouter.Group(rg.relativePath, rg.handlers)
	}
	log.Fatal(http.ListenAndServe(":8000", myRouter))

}
