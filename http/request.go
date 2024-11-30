package http

import (
	"errors"
	"net/url"
	"strings"
)

type Request struct {
	url           Url
	post          map[string]interface{}
	files         []interface{}
	cookies       map[string]interface{}
	headers       map[string]string
	method        string
	remoteAddress *string
	remoteHost    *string
}

func NewRequest(
	url Url,
	post map[string]interface{},
	files []interface{},
	cookies map[string]interface{},
	headers map[string]string,
	method string,
	remoteAddress *string,
	remoteHost *string) *Request {
	return &Request{
		url:           url,
		post:          post,
		files:         files,
		cookies:       cookies,
		headers:       headers,
		method:        method,
		remoteAddress: remoteAddress,
		remoteHost:    remoteHost,
	}
}

func (r *Request) GetUrl() Url {
	return r.url
}

func (r *Request) GetQuery(key string) interface{} {
	return r.url.GetQueryParameter(key)
}

func (r *Request) GetQueries() url.Values {
	return r.url.GetQueries()
}

func (r *Request) GetPost(key string) any {
	return r.post[key]
}

func (r *Request) GetPosts() map[string]interface{} {
	return r.post
}

func (r *Request) GetFile(key string) {

}

func (r *Request) GetFiles() []any {

}

func (r *Request) GetCookie(key string) interface{} {
	return r.cookies[key]
}

func (r *Request) GetCookies() map[string]interface{} {
	return r.cookies
}

func (r *Request) GetMethod() string {
	return r.method
}

func (r *Request) IsMethod(method string) bool {
	return strings.ToLower(r.method) == strings.ToLower(method)
}

func (r *Request) GetHeader(header string) (string, error) {
	_, ok := r.headers[header]
	if !ok {
		return "", errors.New("header does not exists")
	}

	return r.headers[header], nil
}

func (r *Request) GetHeaders() map[string]string {
	return r.headers
}

func (r *Request) IsSecured() bool {
	return r.url.scheme == "https"
}

func (r *Request) IsAjax() bool {
	header, err := r.GetHeader("X-Requested-With")

	if err != nil {
		return false
	}

	return header == "XMLHttpRequest"
}

func (r *Request) GetRemoteAddress() *string {
	return r.remoteAddress
}

func (r *Request) GetRemoteHost() *string {
	return r.remoteHost
}

func (r *Request) GetRawBody() *string {

}
