package signinapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jslee/JiRigo/pkg/infra/response"
	"github.com/jslee/JiRigo/pkg/services/signin/signinmodel"
)

// GET: /user/detail
// 사용자 정보 상세 조회
func (sapi *SigninAPI) getUserByUID(w http.ResponseWriter, r *http.Request) {
	// Gin Context 생성
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	var query signinmodel.GetUserByUIDQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "잘못된 요청 데이터", err)
		return
	}

	// 사용자 목록 조회 서비스 호출
	users, err := sapi.signinService.GetUserByUID(c, query.UID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "사용자 조회 실패", err)
		return
	}

	response.Success(c, http.StatusOK, "사용자 조회 성공", users)
}

// POST: /user/sign-up
// 사용자 생성
func (sapi *SigninAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Gin Context 생성
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	var cmd signinmodel.CreateUserCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, "bad request data", err)
		return
	}

	// 사용자 생성 서비스 호출
	if err := sapi.signinService.CreateUser(c, cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, "사용자 생성 실패", err)
		return
	}

	response.Success(c, http.StatusOK, "사용자 생성 성공", nil)
}

// PATCH: /user/edit
// 사용자 정보 수정
// func (sapi *SigninAPI) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	// Gin Context 생성
// 	c, _ := gin.CreateTestContext(w)
// 	c.Request = r

// 	var cmd signinmodel.UpdateUserCmd
// 	if err := c.ShouldBindJSON(&cmd); err != nil {
// 		response.Error(c, http.StatusBadRequest, "bad request data", err)
// 		return
// 	}

// 	// 사용자 정보 수정 서비스 호출
// 	if err := sapi.signinService.UpdateUser(c, cmd); err != nil {
// 		response.Error(c, http.StatusInternalServerError, "사용자 정보 수정 실패", err)
// 		return
// 	}

// 	response.Success(c, http.StatusOK, "사용자 정보 수정 성공", nil)
// }

// DELETE: /user/delete(USER)
// 사용자 삭제
func (sapi *SigninAPI) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Gin Context 생성
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	var query signinmodel.GetUserByUIDQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "잘못된 요청 데이터", err)
		return
	}

	// 사용자 삭제 서비스 호출
	if err := sapi.signinService.DeleteUser(c, query.UID); err != nil {
		response.Error(c, http.StatusInternalServerError, "사용자 삭제 실패", err)
		return
	}

	response.Success(c, http.StatusOK, "사용자 삭제 성공", nil)
}
