package routing

import "net/http"

// HTTP 라우트를 등록하기 위한 인터페이스
type RouteRegister interface {
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Put(path string, handler http.HandlerFunc)
	Delete(path string, handler http.HandlerFunc)
	Group(path string, fn func(RouteRegister))
	UseMiddleware(middleware ...func(http.Handler) http.Handler)
}
