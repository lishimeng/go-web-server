package web

import (
	"net/http"
	"path"
	"strings"
)

type multiHandlerServer struct {
	schemas map[string]map[string]string // schema: { pathRegex: handlerIndex(schema+pathRegex) }

	delegates map[string]http.Handler // handlerIndex:handler
}

func newServer() *multiHandlerServer {
	return &multiHandlerServer{
		schemas:   make(map[string]map[string]string),
		delegates: make(map[string]http.Handler),
	}
}

func (m *multiHandlerServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	findProxy(m, req).ServeHTTP(resp, req)
}

func (m *multiHandlerServer) RegisterHandler(schema string, requestPath string, handler http.Handler) *multiHandlerServer {

	handlerIndex := handlerGroupKey(schema, requestPath)
	var group map[string]string
	var hasGroup bool
	if group, hasGroup = m.schemas[schema]; !hasGroup {
		group = make(map[string]string)
		m.schemas[schema] = group
	}
	group[requestPath] = handlerIndex

	m.delegates[handlerIndex] = handler
	return m
}

func findProxy(m *multiHandlerServer, req *http.Request) (handler http.Handler) {
	p := req.URL.Path
	p = cleanPath(p)
	schema := req.URL.Scheme
	if schema == "" {
		schema = DefaultSchema
	}

	if handlerGroup, hasSchema := m.schemas[schema]; hasSchema {
		handlerKey, exist := matchWithRegex(handlerGroup, p)
		if exist {
			if hg, hasHandler := m.delegates[handlerKey]; hasHandler {
				handler = hg
			}
		}
	}
	if handler == nil {
		handler = http.NotFoundHandler()
	}

	return handler
}

func handlerGroupKey(schema string, path string) string {
	return schema + "_" + path
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}

func matchWithRegex(toCheck map[string]string, toMatch string) (value string, exist bool) {
	exist = false
	matchLen := 0
	for k, v := range toCheck {
		// Check if key exists.

		if strings.HasPrefix(toMatch, k) {
			exist = true
			if len(k) > matchLen { // 匹配相对精确的http#handler
				value = v
				matchLen = len(k)
			}
		}
	}
	return value, exist
}
