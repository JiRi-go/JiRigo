package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinRouteRegister struct {
	engine   *gin.Engine
	group    *gin.RouterGroup
	basePath string
}

func NewGinRouteRegister(engine *gin.Engine) *GinRouteRegister {
	return &GinRouteRegister{
		engine: engine,
		group:  &engine.RouterGroup,
	}
}

func (r *GinRouteRegister) Get(path string, handler http.HandlerFunc) {
	r.group.GET(path, wrapHandler(handler))
}

func (r *GinRouteRegister) Post(path string, handler http.HandlerFunc) {
	r.group.POST(path, wrapHandler(handler))
}

func (r *GinRouteRegister) Put(path string, handler http.HandlerFunc) {
	r.group.PUT(path, wrapHandler(handler))
}

func (r *GinRouteRegister) Delete(path string, handler http.HandlerFunc) {
	r.group.DELETE(path, wrapHandler(handler))
}

// 지정된 경로 접두사를 가진 새 라우트 그룹을 생성
func (r *GinRouteRegister) Group(path string, fn func(RouteRegister)) {
	group := r.group.Group(path)
	subRouter := &GinRouteRegister{
		engine:   r.engine,
		group:    group,
		basePath: r.basePath + path,
	}
	fn(subRouter)
}

// 라우터에 미들웨어를 추가
func (r *GinRouteRegister) UseMiddleware(middleware ...func(http.Handler) http.Handler) {
	for _, m := range middleware {
		r.group.Use(func(c *gin.Context) {
			m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.Next()
			})).ServeHTTP(c.Writer, c.Request)
		})
	}
}

// http.HandlerFunc를 gin.HandlerFunc으로 변환
func wrapHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}
