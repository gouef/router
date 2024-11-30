package http

import "net/url"

type IRequest interface {
	GetUrl() Url
	GetQuery(key string) interface{}
	GetQueries() url.Values
	GetPost(key string) any
	GetPosts() map[string]interface{}
	GetFile(key string)
	GetFiles() []any
	GetCookie(key string) interface{}
	GetCookies() map[string]interface{}
	GetMethod() string
	IsMethod(method string) bool
	GetHeader(header string) *string
	GetHeaders() map[string]string
	IsSecured() bool
	IsAjax() bool
	GetRemoteAddress() *string
	GetRemoteHost() *string
	GetRawBody() *string
}

const (
	Get     = "GET"
	Post    = "POST"
	Head    = "HEAD"
	Put     = "PUT"
	Delete  = "DELETE"
	Patch   = "PATCH"
	Options = "OPTIONS"
	GET     = Get
	POST    = Post
	HEAD    = Head
	PUT     = Put
	DELETE  = Delete
	PATCH   = Patch
	OPTIONS = Options
)
