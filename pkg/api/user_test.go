package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/domain"
	mockrepo "github.com/xyedo/blindate/pkg/repository/mock"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_PostUserHandler(t *testing.T) {
	var (
		validFullName = "Uncle Bob"
		validEmail    = "bob23@gmail.com"
		validPassword = "validPa$$word"
		validDOB      = "2012-04-23T18:25:43.511Z"
	)
	tests := []struct {
		name       string
		body       map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *user
		wantCode   int
		wantHeader map[string]string
		wantResp   map[string]any
	}{
		{
			name: "Valid Submission",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusCreated,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "success",
				"message": "confirmation email sent, check your email!",
				"data": map[string]string{
					"id": "",
				},
			},
		},
		{
			name: "Invalid Body on Type FullName",
			body: map[string]any{
				"fullName": 1223232,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusBadRequest,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"fullName\"",
			},
		},
		{
			name: "Invalid Body on Type email",
			body: map[string]any{
				"fullName": validFullName,
				"email":    124151,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusBadRequest,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"email\"",
			},
		},
		{
			name: "Invalid Body on Type Password",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": 124541,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusBadRequest,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"password\"",
			},
		},
		{
			name: "Invalid Req Body",
			body: map[string]any{
				"fullName": validFullName,
				"email":    "hahahahha",
				"password": "pass",
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]string{
					"Email":    "must have an valid email",
					"Password": "must have more than 8 character",
				},
			},
		},
		{
			name: "Invalid Body Fields",
			body: map[string]any{
				"fasdullName": validFullName,
				"emailasd":    validEmail,
				"passasdword": "012",
				"dobasd":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Valid but have Unknown Fields",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
				"foo":      "bar",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusCreated,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Valid but have Duplicate Email",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				pqErr := pq.Error{
					Code:       "23505",
					Constraint: "users_email_unique",
				}
				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).Return(&pqErr)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "No Body",
			body: nil,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},

		{

			name: "Invalid Greater Dob",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": validPassword,
				"dob":      time.Now().AddDate(0, 1, 0),
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{

			name: "Invalid Less Dob",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": validPassword,
				"dob":      time.Date(1900, time.January, 0, 0, 0, 0, 0, time.UTC),
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Valid but have Sql Blocking",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).Return(context.Canceled)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusConflict,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			bodyStr, err := json.Marshal(tt.body)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(bodyStr))
			c.Request = req
			userH.postUserHandler(c)
			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), tt.wantHeader["Content-Type"])
			if tt.wantResp != nil {
				expResBody, err := json.Marshal(tt.wantResp)
				assert.NoError(t, err)
				assert.JSONEq(t, string(expResBody), rr.Body.String())
			}
		})
	}
}
func Test_GetUserByIdHandler(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) (userSvc, *domain.User)
		respFunc  func(t *testing.T, user *domain.User, resp *httptest.ResponseRecorder)
	}{

		{
			name: "Valid Submission",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (userSvc, *domain.User) {
				userRepo := mockrepo.NewMockUser(ctrl)
				users := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(users, nil)
				userService := service.NewUser(userRepo)
				return userService, users
			},
			respFunc: func(t *testing.T, user *domain.User, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Contains(t, resp.Header().Get("Content-Type"), "application/json")

				expResBody, err := json.Marshal(map[string]any{
					"status": "success",
					"data": map[string]any{
						"user": user,
					},
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(expResBody), resp.Body.String())
			},
		},
		{
			name: "Valid URL Params but User Not Found",
			id:   "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (userSvc, *domain.User) {
				userRepo := mockrepo.NewMockUser(ctrl)
				users := createNewUser(t)
				userRepo.EXPECT().GetUserById("d3aa0883-4a29-4a39-8f0e-2413c169bd9d").Times(1).Return(nil, sql.ErrNoRows)
				userService := service.NewUser(userRepo)
				return userService, users
			},
			respFunc: func(t *testing.T, user *domain.User, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
				assert.Contains(t, resp.Header().Get("Content-Type"), "application/json")

				expResBody, err := json.Marshal(map[string]any{
					"status":  "fail",
					"message": "id not found",
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(expResBody), resp.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, user := tt.setupFunc(t, ctrl)
			userH := NewUser(userSvc)
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", tt.id), nil)
			c.Request = req
			userH.getUserByIdHandler(c)
			tt.respFunc(t, user, rr)
		})
	}

}

func Test_PatchUserById(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		reqBody   string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *user
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name:    "Valid Submission On FullName",
			id:      "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{"fullName":"Bob Martin"}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				user.FullName = "Bob Martin"
				userRepo.EXPECT().UpdateUser(gomock.Eq(user)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user updated",
			},
		},
		{
			name:    "Valid Submission On Email",
			id:      "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{"email":"bob@martin.com"}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				user.Email = "bob@martin.com"
				userRepo.EXPECT().UpdateUser(gomock.Eq(user)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user updated",
			},
		},
		{
			name:     "Valid Submission On NewPassword",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"oldPassword":"pa55word","newPassword":"newPa55word"}`,
			wantCode: http.StatusOK,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				hashed, err := bcrypt.GenerateFromPassword([]byte("pa55word"), 12)
				assert.NoError(t, err)
				user.HashedPassword = string(hashed)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				userRepo.EXPECT().GetUserByEmail(gomock.Eq(user.Email)).Times(1).Return(user, nil)
				user.Password = "newPa55word"
				userRepo.EXPECT().UpdateUser(gomock.Not(nil)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantResp: map[string]any{
				"status":  "success",
				"message": "user updated",
			},
		},
		{
			name:     "Valid Submission On NewPassword But Invalid OldPassword",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"oldPassword":"pa55word","newPassword":"newPa55word"}`,
			wantCode: http.StatusUnauthorized,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				userRepo.EXPECT().GetUserByEmail(gomock.Eq(user.Email)).Times(1).Return(user, nil)
				userRepo.EXPECT().UpdateUser(gomock.Not(nil)).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "email or password do not match",
			},
		},
		{
			name:    "Invalid Submission On NewPassword",
			id:      "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{"newPassword":"newPa55word"}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				userRepo.EXPECT().GetUserByEmail(gomock.Any()).Times(0)
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]string{
					"OldPassword": "must be more than 8 character and pairing with NewPassword",
				},
			},
		},
		{
			name:     "Invalid Submission On OldPassword",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"oldPassword":"pa55word"}`,
			wantCode: http.StatusUnprocessableEntity,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				userRepo.EXPECT().GetUserByEmail(gomock.Any()).Times(0)
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]string{
					"NewPassword": "must be more than 8 character and pairing with OldPassword",
				},
			},
		},
		{
			name:    "Valid URL Params but User Not Found",
			id:      "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			reqBody: `{"fullName":"Bob Martin"}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().GetUserById(gomock.Eq("d3aa0883-4a29-4a39-8f0e-2413c169bd9d")).Times(1).Return(nil, sql.ErrNoRows)
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "id not found",
			},
		},
		{
			name:    "Invalid Body",
			id:      "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			reqBody: `{"1", 1, 102}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *user {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("d3aa0883-4a29-4a39-8f0e-2413c169bd9d")).Times(1).Return(user, nil)
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				return NewUser(userService)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains badly-formed JSON (at character 5)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%s", tt.id), strings.NewReader(tt.reqBody))
			c.Request = req
			userH.patchUserByIdHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

			expResBody, err := json.Marshal(tt.wantResp)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}

func createNewUser(t *testing.T) *domain.User {
	pass := util.RandomString(10)
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	assert.NoError(t, err)
	return &domain.User{
		ID:             util.RandomUUID(),
		FullName:       util.RandomString(12),
		Email:          util.RandomEmail(12),
		Password:       pass,
		HashedPassword: string(hashed),
		Dob:            util.RandDOB(1980, 2000),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
