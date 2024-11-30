package http

import "net/url"

type UrlImmutable struct {
	scheme    string
	user      string
	password  string
	host      string
	port      *int
	path      string
	query     url.Values
	fragment  string
	authority *string
}

func NewUrlImmutable() *UrlImmutable {
	return &UrlImmutable{}
}
