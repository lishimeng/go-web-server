package web

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

type stubHttpHandler struct {
	Name string
}

func (stub *stubHttpHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(stub.Name)
}

func TestMultiHandlerServer_RegisterHandler(t *testing.T) {

	var stub1 = stubHttpHandler{Name: "HandlerA"}
	var stub2 = stubHttpHandler{Name: "HandlerB"}

	var s = newServer()
	s.RegisterHandler("ws", "/index/p", &stub1)
	s.RegisterHandler("https", "/", &stub2)
	h := findProxy(s, &http.Request{URL: &url.URL{Scheme: "ws", Path: "/index/p/something"}})
	stub, ok := h.(*stubHttpHandler)
	if !ok {
		t.Fatalf("not stubHttpHandler type")
		return
	}
	if stub.Name != "HandlerA" {
		t.Fatalf("handler not match")
	}
}
