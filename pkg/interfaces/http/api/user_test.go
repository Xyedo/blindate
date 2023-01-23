package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/applications/service"
	mocksvc "github.com/xyedo/blindate/pkg/applications/service/mock"
	apiError "github.com/xyedo/blindate/pkg/common/error"
	"github.com/xyedo/blindate/pkg/common/util"
	userEntity "github.com/xyedo/blindate/pkg/domain/user/entities"
	mockrepo "github.com/xyedo/blindate/pkg/infra/repository/mock"
)

func Test_PostUserHandler(t *testing.T) {
	var (
		validUUID     = util.RandomUUID()
		validFullName = "Uncle Bob"
		validAlias    = "bobbies"
		validEmail    = "bob23@gmail.com"
		validPassword = "validPa$$word"
		validDOB      = "2012-04-23T18:25:43.511Z"
	)
	tests := []struct {
		name       string
		body       map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *User
		wantCode   int
		wantHeader map[string]string
		wantResp   map[string]any
	}{
		{
			name: "Valid Submission",
			body: map[string]any{
				"fullName": validFullName,
				"alias":    util.RandomString(5),
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).Return(validUUID, nil)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantCode: http.StatusCreated,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "success",
				"message": "confirmation email sent, check your email!",
				"data": map[string]string{
					"id": validUUID,
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]string{
					"email":    "must be required and have an valid email",
					"password": "must be required and have more than 8 character",
					"alias":    "must be required and between 1-15 characters",
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
				"alias":    validAlias,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
				"foo":      "bar",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).Return(validUUID, nil)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantCode: http.StatusCreated,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "success",
				"message": "confirmation email sent, check your email!",
				"data": map[string]string{
					"id": validUUID,
				},
			},
		},
		{
			name: "Valid but have Duplicate Email",
			body: map[string]any{
				"fullName": validFullName,
				"alias":    validAlias,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				pqErr := pq.Error{
					Code:       "23505",
					Constraint: "users_email_unique",
				}
				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).
					Return("", apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "email already taken"))
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "No Body",
			body: nil,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().InsertUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
				"alias":    validAlias,
				"email":    validEmail,
				"password": validPassword,
				"dob":      validDOB,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().InsertUser(gomock.Not(nil)).Times(1).
					Return("", apiError.Wrap(context.Canceled, apiError.ErrTooLongAccessingDB))
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
		setupFunc func(t *testing.T, ctrl *gomock.Controller) (userSvc, service.Attachment, userEntity.FullDTO)
		respFunc  func(t *testing.T, user userEntity.FullDTO, resp *httptest.ResponseRecorder)
	}{

		{
			name: "Valid Submission",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (userSvc, service.Attachment, userEntity.FullDTO) {
				userRepo := mockrepo.NewMockUser(ctrl)
				users := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(users, nil)
				attachSvc := mocksvc.NewMockAttachment(ctrl)

				userService := service.NewUser(userRepo)
				return userService, attachSvc, users
			},
			respFunc: func(t *testing.T, user userEntity.FullDTO, resp *httptest.ResponseRecorder) {
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (userSvc, service.Attachment, userEntity.FullDTO) {
				userRepo := mockrepo.NewMockUser(ctrl)
				users := createNewUser(t)
				userRepo.EXPECT().GetUserById("d3aa0883-4a29-4a39-8f0e-2413c169bd9d").Times(1).
					Return(userEntity.FullDTO{}, apiError.Wrap(sql.ErrNoRows, apiError.ErrResourceNotFound))
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				userService := service.NewUser(userRepo)
				return userService, attachSvc, users
			},
			respFunc: func(t *testing.T, user userEntity.FullDTO, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
				assert.Contains(t, resp.Header().Get("Content-Type"), "application/json")

				expResBody, err := json.Marshal(map[string]any{
					"status":  "fail",
					"message": "resource not found",
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
			userSvc, attachSvc, user := tt.setupFunc(t, ctrl)

			userH := NewUser(userSvc, attachSvc)
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
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *User
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name:    "Valid Submission On FullName",
			id:      "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{"fullName":"Bob Martin"}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				user.FullName = "Bob Martin"
				userRepo.EXPECT().UpdateUser(gomock.Eq(user)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				user.Email = "bob@martin.com"
				userRepo.EXPECT().UpdateUser(gomock.Eq(user)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				hashed, err := bcrypt.GenerateFromPassword([]byte("pa55word"), 12)
				assert.NoError(t, err)
				user.Password = string(hashed)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				user.Password = "newPa55word"
				userRepo.EXPECT().UpdateUser(gomock.Not(nil)).Times(1).Return(nil)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				user := createNewUser(t)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(user, nil)
				userRepo.EXPECT().UpdateUser(gomock.Not(nil)).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "invalid credentials",
			},
		},
		{
			name:    "Invalid Submission On NewPassword",
			id:      "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{"newPassword":"newPa55word"}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(0)
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]string{
					"oldPassword": "must be more than 8 character and pairing with newPassword",
				},
			},
		},
		{
			name:     "Invalid Submission On OldPassword",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"oldPassword":"pa55word"}`,
			wantCode: http.StatusUnprocessableEntity,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().GetUserById(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(0)
				userRepo.EXPECT().GetUserByEmail(gomock.Any()).Times(0)
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]string{
					"newPassword": "must be more than 8 character and pairing with oldPassword",
				},
			},
		},
		{
			name:    "Valid URL Params but User Not Found",
			id:      "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			reqBody: `{"fullName":"Bob Martin"}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().GetUserById(gomock.Eq("d3aa0883-4a29-4a39-8f0e-2413c169bd9d")).Times(1).
					Return(userEntity.FullDTO{}, apiError.Wrap(sql.ErrNoRows, apiError.ErrResourceNotFound))
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "resource not found",
			},
		},
		{
			name:    "Invalid Body",
			id:      "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			reqBody: `{"1", 1, 102}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				userRepo.EXPECT().GetUserById(gomock.Eq("d3aa0883-4a29-4a39-8f0e-2413c169bd9d")).Times(0)
				userRepo.EXPECT().UpdateUser(gomock.Any()).Times(0)
				userService := service.NewUser(userRepo)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				return NewUser(userService, attachSvc)
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

type uploadTestMatcher struct {
	actualV any
}

func (u *uploadTestMatcher) Matches(x any) bool {
	u.actualV = x
	_, ok := x.(io.Reader)
	return ok
}
func (u *uploadTestMatcher) String() string {
	return fmt.Sprintln(u.actualV)
}
func uploadValid() *uploadTestMatcher {
	return &uploadTestMatcher{}
}
func Test_PutUserImageProfile(t *testing.T) {
	validUserId := util.RandomUUID()
	validProfPicId := util.RandomUUID()
	writeToPng := func(writer *multipart.Writer) {
		defer writer.Close()
		part, err := writer.CreateFormFile("file", "img-test.png")
		require.NoError(t, err)
		img := util.CreateDefaultImage(200, 200)
		err = png.Encode(part, img)
		require.NoError(t, err)
	}
	tests := []struct {
		name      string
		id        string
		params    map[string]string
		writoMime func(writer *multipart.Writer)
		stubFunc  func(t *testing.T, ctrl *gomock.Controller, pr *io.PipeReader) *User
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "valid with selected profile-pictures",
			id:   validUserId,
			params: map[string]string{
				"selected": "true",
			},
			writoMime: writeToPng,
			stubFunc: func(t *testing.T, ctrl *gomock.Controller, pr *io.PipeReader) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				validKey := "profile-picture/" + util.RandomUUID() + ".png"
				user := createNewUser(t)
				user.ID = validUserId
				attachSvc.EXPECT().UploadBlob(uploadValid(), gomock.Any()).Return(validKey, nil).Times(1)

				profPic := make([]userEntity.ProfilePic, 0, 4)
				for i := 0; i < 3; i++ {
					profPic = append(profPic, createRandomProfPic(user.ID))
				}

				userRepo.EXPECT().SelectProfilePicture(gomock.Eq(user.ID), gomock.Nil()).Return(profPic, nil).Times(1)
				userRepo.EXPECT().ProfilePicSelectedToFalse(gomock.Eq(user.ID)).Return(int64(3), nil).Times(1)
				userRepo.
					EXPECT().
					CreateProfilePicture(
						gomock.Eq(user.ID),
						gomock.Eq(filepath.Base(validKey)),
						gomock.Eq(true),
					).
					Return(validProfPicId, nil).
					Times(1)
				userSvc := service.NewUser(userRepo)
				return NewUser(userSvc, attachSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user profile-picture uploaded",
				"data": gin.H{
					"profilePicture": gin.H{
						"id": validProfPicId,
					},
				},
			},
		},
		{
			name:      "valid but unselected",
			id:        validUserId,
			writoMime: writeToPng,
			stubFunc: func(t *testing.T, ctrl *gomock.Controller, pr *io.PipeReader) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				validKey := "profile-picture/" + util.RandomUUID() + ".png"
				user := createNewUser(t)
				user.ID = validUserId
				attachSvc.EXPECT().UploadBlob(uploadValid(), gomock.Any()).Return(validKey, nil).Times(1)

				profPic := make([]userEntity.ProfilePic, 0, 4)
				for i := 0; i < 3; i++ {
					profPic = append(profPic, createRandomProfPic(user.ID))
				}

				userRepo.EXPECT().SelectProfilePicture(gomock.Eq(user.ID), gomock.Nil()).Return(profPic, nil).Times(1)
				userRepo.EXPECT().ProfilePicSelectedToFalse(gomock.Any()).Times(0)
				userRepo.
					EXPECT().
					CreateProfilePicture(
						gomock.Eq(user.ID),
						gomock.Eq(filepath.Base(validKey)),
						gomock.Eq(false),
					).
					Return(validProfPicId, nil).
					Times(1)
				userSvc := service.NewUser(userRepo)
				return NewUser(userSvc, attachSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user profile-picture uploaded",
				"data": gin.H{
					"profilePicture": gin.H{
						"id": validProfPicId,
					},
				},
			},
		},
		{
			name: "invalid MIME types",
			id:   validUserId,
			params: map[string]string{
				"selected": "true",
			},
			writoMime: func(writer *multipart.Writer) {
				defer writer.Close()
				part, err := writer.CreateFormFile("file", "text.txt")
				require.NoError(t, err)
				str := util.RandomString(128)
				read := strings.NewReader(str)
				_, err = io.Copy(part, read)
				require.NoError(t, err)
			},
			stubFunc: func(t *testing.T, ctrl *gomock.Controller, pr *io.PipeReader) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				validKey := "profile-picture/" + util.RandomUUID() + ".png"
				user := createNewUser(t)
				user.ID = validUserId
				attachSvc.EXPECT().UploadBlob(gomock.Any(), gomock.Any()).Times(0)

				userRepo.EXPECT().SelectProfilePicture(gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().ProfilePicSelectedToFalse(gomock.Any()).Times(0)
				userRepo.
					EXPECT().
					CreateProfilePicture(
						gomock.Eq(user.ID),
						gomock.Eq(filepath.Base(validKey)),
						gomock.Eq(false),
					).
					Times(0)
				userSvc := service.NewUser(userRepo)
				return NewUser(userSvc, attachSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "not valid mime-type",
			},
		},
		{
			name: "not having a file",
			id:   validUserId,
			params: map[string]string{
				"selected": "true",
			},
			writoMime: func(writer *multipart.Writer) {
				defer writer.Close()
				part, err := writer.CreateFormField("file")
				require.NoError(t, err)
				str := util.RandomString(128)
				read := strings.NewReader(str)
				_, err = io.Copy(part, read)
				require.NoError(t, err)
			},
			stubFunc: func(t *testing.T, ctrl *gomock.Controller, pr *io.PipeReader) *User {
				userRepo := mockrepo.NewMockUser(ctrl)
				attachSvc := mocksvc.NewMockAttachment(ctrl)
				validKey := "profile-picture/" + util.RandomUUID() + ".png"
				user := createNewUser(t)
				user.ID = validUserId
				attachSvc.EXPECT().UploadBlob(gomock.Any(), gomock.Any()).Times(0)

				userRepo.EXPECT().SelectProfilePicture(gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().ProfilePicSelectedToFalse(gomock.Any()).Times(0)
				userRepo.
					EXPECT().
					CreateProfilePicture(
						gomock.Eq(user.ID),
						gomock.Eq(filepath.Base(validKey)),
						gomock.Eq(false),
					).
					Times(0)
				userSvc := service.NewUser(userRepo)
				return NewUser(userSvc, attachSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "request did not contain a file",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			pr, pw := io.Pipe()
			writer := multipart.NewWriter(pw)
			go tt.writoMime(writer)

			userH := tt.stubFunc(t, ctrl, pr)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)
			req := httptest.NewRequest("PUT", "/users/"+tt.id+"/profile-picture", pr)
			req.Header.Add("Content-Type", writer.FormDataContentType())

			if tt.params != nil {
				queryParams := req.URL.Query()
				for k, v := range tt.params {
					queryParams.Add(k, v)
				}
				req.URL.RawQuery = queryParams.Encode()
			}

			c.Request = req
			userH.putUserImageProfileHandler(c)
			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			if tt.wantResp != nil {
				expResBody, err := json.Marshal(tt.wantResp)
				require.NoError(t, err)
				assert.JSONEq(t, string(expResBody), rr.Body.String())
			}
		})
	}
	t.Run("not multipart", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mockrepo.NewMockUser(ctrl)
		attachSvc := mocksvc.NewMockAttachment(ctrl)
		validKey := "profile-picture/" + util.RandomUUID() + ".png"
		user := createNewUser(t)
		user.ID = validUserId
		attachSvc.EXPECT().UploadBlob(gomock.Any(), gomock.Any()).Times(0)

		userRepo.EXPECT().SelectProfilePicture(gomock.Any(), gomock.Any()).Times(0)
		userRepo.EXPECT().ProfilePicSelectedToFalse(gomock.Any()).Times(0)
		userRepo.
			EXPECT().
			CreateProfilePicture(
				gomock.Eq(user.ID),
				gomock.Eq(filepath.Base(validKey)),
				gomock.Eq(false),
			).
			Times(0)
		userSvc := service.NewUser(userRepo)
		userApi := NewUser(userSvc, attachSvc)
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		c.Set("userId", validUserId)
		req := httptest.NewRequest("PUT", "/users/"+validUserId+"/profile-picture", strings.NewReader(`{"fullName":"Bob Martin"}`))
		c.Request = req
		userApi.putUserImageProfileHandler(c)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		require.Contains(t, rr.Header().Get("Content-Type"), "application/json")
		expResBody, err := json.Marshal(map[string]any{"status": "fail", "message": "content-Type header is not valid"})
		assert.NoError(t, err)
		assert.JSONEq(t, string(expResBody), rr.Body.String())
	})
}

func createNewUser(t *testing.T) userEntity.FullDTO {
	pass := util.RandomString(10)
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	assert.NoError(t, err)
	return userEntity.FullDTO{
		ID:        util.RandomUUID(),
		FullName:  util.RandomString(12),
		Email:     util.RandomEmail(12),
		Password:  string(hashed),
		Dob:       util.RandDOB(1980, 2000),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createRandomProfPic(userId string) userEntity.ProfilePic {
	return userEntity.ProfilePic{
		Id:          util.RandomUUID(),
		UserId:      userId,
		PictureLink: util.RandomUUID() + ".png",
		Selected:    util.RandomBool(),
	}
}
