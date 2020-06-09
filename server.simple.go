package web

import (
	"context"
	"github.com/kataras/iris"
	"net/http"
	"time"
)

type SimpleServer struct {
	config        ServerConfig
	delegate  *iris.Application
}

func (s *SimpleServer) GetDelegate() *iris.Application {
	return s.delegate
}

func (s *SimpleServer) Start(ctx context.Context) error {
	if err := s.delegate.Configure(iris.WithCharset("UTF-8")).Build(); err != nil {
		return err
	}
	srv := http.Server{
		Addr:    s.config.Listen,
		Handler: s.delegate,
	}
	go s.shutdownFuture(&srv, ctx)

	return srv.ListenAndServe()
}

func (s *SimpleServer) shutdownFuture(srv *http.Server, ctx context.Context) {
	if ctx == nil {
		return
	}
	var c context.Context
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			c, cancel = context.WithTimeout(context.Background(), 3*time.Second)
			if err := srv.Shutdown(c); nil != err {
			}
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func NewSimpleServer(config ServerConfig) (handler *Server) {

	s := Server{
		config:        config,
		primaryProxy:  iris.New(),
		delegate:      newServer(),
		primaryPath:   "/",
		primarySchema: DefaultSchema,
	}
	return &s
}