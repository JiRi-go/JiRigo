package signinapi

import (
	"github.com/jslee/JiRigo/pkg/infra/routing"
	"github.com/jslee/JiRigo/pkg/services/signin"
)

type SigninAPI struct {
	signinService signin.Service
}

func ProvideSigninAPI(
	routeRegister routing.RouteRegister,
	signinService signin.Service,
) *SigninAPI {
	sapi := &SigninAPI{
		signinService: signinService,
	}

	sapi.registerRoutes(routeRegister)
	return sapi
}

// 라우트 등록 메서드 수정
func (sapi *SigninAPI) registerRoutes(router routing.RouteRegister) {
	router.Group("/user", func(signinRoute routing.RouteRegister) {
		// 단일 사용자 조회
		signinRoute.Post("/detail", sapi.getUserByUID)

		// 사용자 생성
		signinRoute.Post("/", sapi.CreateUser)

		// 사용자 정보 수정
		// signinRoute.Patch("/edit", sapi.UpdateUser)

		// 사용자 삭제
		signinRoute.Post("/delete", sapi.DeleteUser)
	})
}
