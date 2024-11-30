package http

import (
	"errors"
	"fmt"
	"github.com/Gouef/utils"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Url struct {
	defaultPorts map[string]int
	scheme       string
	user         string
	password     string
	host         string
	port         *int
	path         string
	query        url.Values
	fragment     string
	scriptPath   string
	basePath     string
}

var DefaultPorts = map[string]int{
	"http":  80,
	"https": 443,
	"ftp":   21,
}

func NewUrl(urlStr interface{}) (*Url, error) {
	switch v := urlStr.(type) {
	case string:
		return NewUrlString(v)
	case UrlImmutable:
		return NewUrlFromImmutable(v)
	case Url:
		return NewUrlFromSelf(v)
	}

	return &Url{}, nil
}

func NewUrlFromSelf(url Url) (*Url, error) {
	return &Url{
		defaultPorts: DefaultPorts,
		scheme:       url.scheme,
		user:         url.user,
		password:     url.password,
		host:         url.host,
		port:         url.port,
		path:         url.path,
		query:        url.query,
		fragment:     url.fragment,
	}, nil
}

func NewUrlFromImmutable(urlImmutable UrlImmutable) (*Url, error) {
	return &Url{
		defaultPorts: DefaultPorts,
		scheme:       urlImmutable.scheme,
		user:         urlImmutable.user,
		password:     urlImmutable.password,
		host:         urlImmutable.host,
		port:         urlImmutable.port,
		path:         urlImmutable.path,
		query:        urlImmutable.query,
		fragment:     urlImmutable.fragment,
	}, nil
}

func NewUrlString(urlStr string) (*Url, error) {
	p, err := url.Parse(urlStr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Malformed or unsupported URI '%s'", urlStr))
	}

	portStr := p.Port()

	var port *int

	if portStr == "" {
		port = nil
	} else {

		po, err := strconv.Atoi(portStr)
		port = &po
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error converting port to int: %s", err))
		}
	}
	pass, _ := p.User.Password()

	return &Url{
		defaultPorts: DefaultPorts,
		scheme:       p.Scheme,
		port:         port,
		host:         p.Hostname(),
		user:         p.User.Username(),
		password:     pass,
		path:         p.Path,
		query:        p.Query(),
		fragment:     p.Fragment,
	}, nil
}

func (u *Url) GetDomain(level int) string {
	var parts []string
	_, err := utils.Ip2long(u.host)

	if err == nil {
		parts = append(parts, u.host)
	} else {
		parts = utils.Explode(".", u.host)
	}

	if level >= 0 {
		parts = parts[:-level]
	} else {
		parts = parts[0:level]
	}

	return utils.Implode(".", parts)
}

func (u *Url) SetScheme(scheme string) *Url {
	u.scheme = scheme
	return u
}

func (u *Url) GetScheme() string {
	return u.scheme
}

func (u *Url) SetPath(path string) *Url {
	u.path = path

	if u.host != "" && strings.HasPrefix(u.path, "/") {
		u.path = "/" + u.path
	}
	return u
}

func (u *Url) GetPath() string {
	return u.path
}

func (u *Url) SetQuery(query url.Values) *Url {
	u.query = query
	return u
}

func (u *Url) SetQueryString(query string) *Url {
	u.query = ParseQuery(query)
	return u
}

func (u *Url) AppendQuery(query url.Values) *Url {
	for k, v := range query {
		u.query[k] = append(u.query[k], v...)
	}
	return u
}

func (u *Url) AppendQueryString(query string) *Url {
	u.query = ParseQuery(u.GetQuery() + "&" + query)
	return u
}

func (u *Url) GetQuery() string {
	return u.query.Encode()
}

func (u *Url) GetQueries() url.Values {
	return u.query
}

func (u *Url) GetQueryParameter(name string) string {
	return u.query.Get(name)
}

func (u *Url) SetQueryParameter(name string, value string) *Url {
	u.query.Set(name, value)
	return u
}

func (u *Url) GetDefaultPort() *int {
	defaultPort := u.defaultPorts[u.scheme]
	return &defaultPort
}

func (u *Url) GetAuthority() string {
	if u.host == "" {
		return ""
	}

	result := ""

	if u.user != "" {
		result = url.QueryEscape(u.user)

		if u.password != "" {
			result += ":" + url.QueryEscape(u.password)
		}

		result += "@"
	}

	result += u.host

	if u.port != nil && u.port != u.GetDefaultPort() {
		result += fmt.Sprintf(":%i", u.port)
	}

	return result
}

func (u *Url) GetHostUrl() string {
	hostUrl := u.scheme
	authority := u.GetAuthority()

	if authority != "" {
		hostUrl += "//" + authority
	}

	return hostUrl
}

func (u *Url) GetAbsoluteUrl() string {
	absoluteUrl := u.GetHostUrl() + u.path
	query := u.GetQuery()

	if query != "" {
		absoluteUrl += "?"
	}
	absoluteUrl += query

	if u.fragment != "" {
		absoluteUrl += "#"
	}

	absoluteUrl += u.fragment
	return absoluteUrl
}

func (u *Url) SetScriptPath(scriptPath string) error {
	path := u.path

	if scriptPath == "" {
		scriptPath = path
	}

	pos := strings.LastIndex(scriptPath, "/")

	if pos == -1 || utils.StringCompare(scriptPath, path, pos+1) == 1 {
		return errors.New(fmt.Sprintf("ScriptPath '%s' doesn't match path '%s'", scriptPath, path))
	}

	u.scriptPath = scriptPath
	pos += 1
	u.basePath = utils.Substr(scriptPath, 0, &pos)
	return nil
}

func (u *Url) GetScriptPath() string {
	return u.scriptPath
}

func (u *Url) GetBasePath() string {
	return u.basePath
}

func (u *Url) GetRelativePath() string {
	return utils.Substr(u.path, len(u.basePath), nil)
}

func (u *Url) GetBaseUrl() string {
	return u.GetHostUrl() + u.basePath
}

func (u *Url) GetRelativeUrl() string {
	return utils.Substr(u.GetAbsoluteUrl(), len(u.GetBaseUrl()), nil)
}

func (u *Url) GetPathInfo() string {
	return utils.Substr(u.path, len(u.scriptPath), nil)
}

func (u *Url) mergePath(path string) string {
	return u.basePath + path
}

func ParseQuery(s string) map[string][]string {
	s = strings.ReplaceAll(s, "%5B", "[")
	s = strings.ReplaceAll(s, "%5b", "[")

	separator := "&"

	re := regexp.MustCompile(fmt.Sprintf(`([%s])([^%s=]+)([^%s]*)`, regexp.QuoteMeta(separator), regexp.QuoteMeta(separator), regexp.QuoteMeta(separator)))
	s = re.ReplaceAllString("&0[$2]$3", "&"+s)

	values, err := url.ParseQuery(s)
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return nil
	}

	return values
}
