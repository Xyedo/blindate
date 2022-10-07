package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func addAutho(t *testing.T, req *http.Request, tokennizer service.Jwt, id, typeAuth string) {
	token, err := tokennizer.GenerateAccessToken(id)
	assert.NoError(t, err)
	authoHeader := fmt.Sprintf("%s %s", typeAuth, token)
	req.Header.Set("Authorization", authoHeader)

}

func Test_AuthMiddleware(t *testing.T) {
	var (
		validId = "e590666c-3ea8-4fda-958c-c2dc6c2599b5"
	)
	jwt := service.NewJwt("test-access-secret", "test-refresh-secret", "1s", "720h")
	tests := []struct {
		name         string
		id           string
		setupAuth    func(t *testing.T, req *http.Request, tokenizer service.Jwt)
		checkRespose func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "valid Authorization",
			id:   validId,
			setupAuth: func(t *testing.T, req *http.Request, tokenizer service.Jwt) {
				addAutho(t, req, tokenizer, validId, "Bearer")
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "No Authorization",
			id:        validId,
			setupAuth: func(t *testing.T, req *http.Request, tokenizer service.Jwt) {},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupported Authorization",
			id:   validId,
			setupAuth: func(t *testing.T, req *http.Request, tokenizer service.Jwt) {
				addAutho(t, req, tokenizer, validId, "Basic")
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Authorization format",
			id:   validId,
			setupAuth: func(t *testing.T, req *http.Request, tokenizer service.Jwt) {
				addAutho(t, req, tokenizer, validId, "")
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired Token",
			id:   validId,
			setupAuth: func(t *testing.T, req *http.Request, tokenizer service.Jwt) {
				addAutho(t, req, tokenizer, validId, "Bearer")
				time.Sleep(2 * time.Second)
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Id",
			id:   "1",
			setupAuth: func(t *testing.T, req *http.Request, tokenizer service.Jwt) {
				addAutho(t, req, tokenizer, validId, "Bearer")

			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Valid Id but uri not the same with token",
			id:   "d3aa0883-4a29-4a39-8f0e-2413c169bd9d",
			setupAuth: func(t *testing.T, req *http.Request, tokenizer service.Jwt) {
				addAutho(t, req, tokenizer, validId, "Bearer")

			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, r := gin.CreateTestContext(rr)
			r.GET("/api/v1/users/:id", validateUser(jwt), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, nil)
			})
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", tt.id), nil)
			assert.NoError(t, err)
			c.Request = req
			tt.setupAuth(t, c.Request, jwt)
			r.ServeHTTP(rr, c.Request)
			tt.checkRespose(t, rr)
		})
	}
}

func Test_InterestMidleware(t *testing.T) {
	tests := []struct {
		name         string
		interestId   string
		expectedCode int
	}{
		{
			name:         "valid interestId",
			interestId:   util.RandomUUID(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "interestId not uuid",
			interestId:   util.RandomString(12),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "interestId nil",
			interestId:   "",
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, r := gin.CreateTestContext(rr)
			r.GET("/interests/:interestId", validateInterest(), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, nil)
			})
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/interests/%s", tt.interestId), nil)
			assert.NoError(t, err)
			c.Request = req
			r.ServeHTTP(rr, c.Request)
			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}

}
