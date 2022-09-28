package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/mock"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

func Test_PostUserHandler(t *testing.T) {
	userService := service.NewUser(mock.UserRepository{})
	userH := NewUser(userService)

	gin.SetMode(gin.TestMode)
	registerValidDObValidator()
	var (
		validFullName = "Uncle Bob"
		validEmail    = "bob23@gmail.com"
		validPassword = "validPa$$word"
		validDOB      = "2012-04-23T18:25:43.511Z"
	)
	tests := []struct {
		name       string
		body       map[string]any
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
			wantCode: http.StatusCreated,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "success",
				"message": "confirmation email sent, check your email!",
				"data": map[string]string{
					"id": "1",
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
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"error": map[string]string{
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
			wantCode: http.StatusCreated,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Duplicate Email",
			body: map[string]any{
				"fullName": validFullName,
				"email":    "dupli23@gmail.com",
				"password": validPassword,
				"dob":      validDOB,
			},
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:     "No Body",
			body:     nil,
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
			wantCode: http.StatusUnprocessableEntity,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
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
	userService := service.NewUser(mock.UserRepository{})
	userH := NewUser(userService)

	gin.SetMode(gin.TestMode)
	var (
		validFullName = "Uncle Bob"
		validEmail    = "bob@example.com"
		validDOB      = time.Date(2000, time.August, 23, 0, 0, 0, 0, time.UTC)
	)
	var realuser = domain.User{
		ID:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
		FullName: validFullName,
		Email:    validEmail,
		Dob:      validDOB,
	}

	tests := []struct {
		name     string
		id       string
		wantCode int
		wantResp map[string]any
	}{

		{
			name:     "Valid Submission",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"user": realuser,
				},
			},
		},
		{
			name:     "Valid URL Params but User Not Found",
			id:       "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "id not found!",
			},
		},
		{
			name:     "Invalid URL Params",
			id:       "1",
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "must have uuid in uri!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, r := gin.CreateTestContext(rr)
			r.GET("/api/v1/users/:id", userH.getUserByIdHandler)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", tt.id), nil)
			c.Request = req
			r.ServeHTTP(rr, c.Request)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

			expResBody, err := json.Marshal(tt.wantResp)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}

}

func Test_PatchUserById(t *testing.T) {
	userService := service.NewUser(mock.UserRepository{})
	userH := NewUser(userService)

	gin.SetMode(gin.TestMode)
	// var (
	// 	validFullName    = "Bob Martin"
	// 	validEmail       = "bob@martin.com"
	// 	validOldPassword = "pa55word"
	// 	// validNewPassword = "newPa55word"
	// 	// validDOB         = time.Date(2000, time.August, 23, 0, 0, 0, 0, time.UTC)
	// )

	tests := []struct {
		name     string
		id       string
		reqBody  string
		wantCode int
		wantResp map[string]any
	}{
		{
			name:     "Valid Submission On FullName",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"fullName":"Bob Martin"}`,
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user updated",
			},
		},
		{
			name:     "Valid Submission On Email",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"email":"bob@martin.com"}`,
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
			wantResp: map[string]any{
				"status":  "success",
				"message": "user updated",
			},
		},
		{
			name:     "Invalid Submission On NewPassword",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"newPassword":"newPa55word"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"error": map[string]string{
					"OldPassword": "must be more than 8 character and pairing with NewPassword",
				},
			},
		},
		{
			name:     "Invalid Submission On OldPassword",
			id:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody:  `{"oldPassword":"pa55word"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"error": map[string]string{
					"NewPassword": "must be more than 8 character and pairing with OldPassword",
				},
			},
		},
		{
			name:     "Valid URL Params but User Not Found",
			id:       "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			reqBody:  `{"fullName":"Bob Martin"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "id not found!",
			},
		},
		{
			name:     "Invalid URL Params",
			id:       "1",
			wantCode: http.StatusBadRequest,
			reqBody:  `{"fullName":"Bob Martin"}`,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "must have uuid in uri!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, r := gin.CreateTestContext(rr)
			r.PATCH("/api/v1/users/:id", userH.patchUserByIdHandler)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%s", tt.id), strings.NewReader(tt.reqBody))
			c.Request = req
			r.ServeHTTP(rr, c.Request)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

			expResBody, err := json.Marshal(tt.wantResp)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}
