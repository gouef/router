package router

import (
	"fmt"
	"github.com/gouef/router/http"
	"net/url"
)

type SimpleRouter struct {
	defaults url.Values
}

func NewSimpleRouter(defaults url.Values) *SimpleRouter {
	return &SimpleRouter{defaults: defaults}
}

func (r *SimpleRouter) match(httpRequest http.IRequest) *url.Values {
	httpUrl := httpRequest.GetUrl()
	if httpUrl.GetPathInfo() == "" {
		result := httpRequest.GetQueries()
		for k, v := range r.defaults {
			for _, value := range v {
				result.Add(k, value)
			}
		}

		return &result
	}

	return nil
}

func (r *SimpleRouter) constructUrl(params map[string]interface{}, refUrl *http.Url) *string {
	values := url.Values{}

	for key, value := range params {
		if val, exists := r.defaults[key]; !exists || fmt.Sprintf("%v", value) != fmt.Sprintf("%v", val) {
			values.Add(fmt.Sprintf("%v", key), fmt.Sprintf("%v", value))
		}
	}

	u := url.URL{
		Scheme: refUrl.GetScheme(),
		Host:   refUrl.GetHostUrl(),
		Path:   refUrl.GetPath(),
	}

	u.RawQuery = values.Encode()
	result := u.String()
	return &result
}

func (r *SimpleRouter) GetDefaults() url.Values {
	return r.defaults
}
