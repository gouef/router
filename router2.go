package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Router2 struct {
	groups []Group
}

func (r *Router2) Groups() []Group {
	return r.groups
}

type Group struct {
	relativePath string
	handlers     gin.HandlerFunc
}

func (g Group) RelativePath() string {
	return g.relativePath
}

func (g Group) Handlers() gin.HandlerFunc {
	return g.handlers
}

func (r *Router2) Run() {
	myRouter := gin.Default()

	for _, group := range r.Groups() {
		myRouter.Group(group.RelativePath(), group.Handlers())
	}

	log.Fatal(http.ListenAndServe(":8000", myRouter))

}
