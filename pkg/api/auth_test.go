package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockrepo "github.com/xyedo/blindate/pkg/repository/mock"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

func Test_postAuthHandler(t *testing.T) {
	jwt := service.NewJwt("test-access-secret", "test-refresh-secret", "1s", "720h")

	tests := []struct {
		name      string
		reqBody   string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) auth
		wantCode  int
		respFunc  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Submission",
			reqBody: `{
				"email":"uncleBob23@cool.com",
				"password":"pa55word"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) auth {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)

				validUser := createNewUser(t)
				hashed, err := bcrypt.GenerateFromPassword([]byte("pa55word"), 12)
				assert.NoError(t, err)
				validUser.HashedPassword = string(hashed)

				userRepo.EXPECT().GetUserByEmail(gomock.Eq("uncleBob23@cool.com")).Times(1).Return(validUser, nil)
				authRepo.EXPECT().AddRefreshToken(gomock.Any()).Times(1).Return(int64(1), nil)

				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt)
			},
			wantCode: http.StatusCreated,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.NotZero(t, rr.Result().Cookies()[0].Value)
				fmt.Println(rr.Result().Cookies()[0].Value)
				var result map[string]any
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				assert.NoError(t, err)

				assert.Equal(t, "success", result["status"])
				data, ok := result["data"].(map[string]any)
				assert.True(t, ok)
				assert.NotZero(t, data["accessToken"])
			},
		},
		{
			name: "Invalid Type Req Body",
			reqBody: `{
				"email":0,
				"password":"pa55word"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) auth {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().GetUserByEmail(gomock.Any()).Times(0)
				authRepo.EXPECT().AddRefreshToken(gomock.Any()).Times(0)

				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt)
			},
			wantCode: http.StatusBadRequest,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				respBody, err := json.Marshal(map[string]any{
					"status":  "fail",
					"message": "body contains incorrect JSON type for field \"email\"",
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(respBody), rr.Body.String())
			},
		},
		{
			name: "Invalid Validation All",
			reqBody: `{
				"email":"hummm",
				"password":"pa55w"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) auth {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)

				userRepo.EXPECT().GetUserByEmail(gomock.Any()).Times(0)
				authRepo.EXPECT().AddRefreshToken(gomock.Any()).Times(0)

				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt)
			},
			wantCode: http.StatusUnprocessableEntity,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				respBody, err := json.Marshal(map[string]any{
					"status":  "fail",
					"message": "please refer to the documentation",
					"errors": map[string]any{
						"email":    "required and must be valid email",
						"password": "required and must be over 8 character",
					},
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(respBody), rr.Body.String())
			},
		},
		{
			name: "Valid Submission but Invalid Credentials",
			reqBody: `{
				"email":"uncleBob23@cool.com",
				"password":"pa55word"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) auth {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)

				validUser := createNewUser(t)
				userRepo.EXPECT().GetUserByEmail(gomock.Eq("uncleBob23@cool.com")).Times(1).Return(validUser, nil)
				authRepo.EXPECT().AddRefreshToken(gomock.Any()).Times(0)

				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt)
			},
			wantCode: http.StatusUnauthorized,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				respBody, err := json.Marshal(map[string]any{
					"status":  "fail",
					"message": "invalid credentials",
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(respBody), rr.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth", strings.NewReader(tt.reqBody))
			c.Request = req
			authH.postAuthHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			tt.respFunc(t, rr)
		})
	}
}

func Test_putAuthHandler(t *testing.T) {
	jwt := service.NewJwt("test-access-secret", "test-refresh-secret", "1s", "720h")

	tests := []struct {
		name      string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) (auth, string)
		wantCode  int
		respFunc  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "Valid PutAuthHandler",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (auth, string) {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)
				id := util.RandomUUID()
				token, err := jwt.GenerateRefreshToken(id)
				assert.NoError(t, err)
				authRepo.EXPECT().VerifyRefreshToken(gomock.Eq(token)).Times(1).Return(nil)

				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt), token
			},
			wantCode: http.StatusOK,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result map[string]any
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				assert.NoError(t, err)

				assert.Equal(t, "success", result["status"])
				data, ok := result["data"].(map[string]any)
				assert.True(t, ok)
				assert.NotZero(t, data["accessToken"])
			},
		},
		{
			name: "Invalid RefreshToken",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (auth, string) {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)
				id := util.RandomUUID()
				token, err := jwt.GenerateRefreshToken(id)
				assert.NoError(t, err)
				authRepo.EXPECT().VerifyRefreshToken(gomock.Eq(token)).Times(1).Return(sql.ErrNoRows)

				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt), token
			},
			wantCode: http.StatusUnauthorized,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				respBody, err := json.Marshal(map[string]any{
					"status":  "fail",
					"message": "invalid credentials",
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(respBody), rr.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authH, token := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/auth", nil)
			req.AddCookie(&http.Cookie{
				Name:     "refreshToken",
				Value:    token,
				Path:     "/api/v1",
				Domain:   "localhost",
				MaxAge:   2592000,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			})
			c.Request = req

			authH.putAuthHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			tt.respFunc(t, rr)
		})
	}
}

func Test_deleteAuthHandler(t *testing.T) {
	jwt := service.NewJwt("test-access-secret", "test-refresh-secret", "1s", "720h")
	tests := []struct {
		name      string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) (auth, string)
		wantCode  int
		respFunc  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Log out",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (auth, string) {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)
				id := util.RandomUUID()
				token, err := jwt.GenerateRefreshToken(id)
				assert.NoError(t, err)
				authRepo.EXPECT().VerifyRefreshToken(gomock.Eq(token)).Times(1).Return(nil)
				authRepo.EXPECT().DeleteRefreshToken(gomock.Eq(token)).Times(1).Return(int64(1), nil)
				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt), token
			},
			wantCode: http.StatusOK,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Zero(t, rr.Result().Cookies()[0].Value)
				respBody, err := json.Marshal(map[string]any{
					"status":  "success",
					"message": "log out success",
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(respBody), rr.Body.String())
			},
		},
		{
			name: "InValid Log out",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) (auth, string) {
				authRepo := mockrepo.NewMockAuth(ctrl)
				userRepo := mockrepo.NewMockUser(ctrl)
				id := util.RandomUUID()
				token, err := jwt.GenerateRefreshToken(id)
				assert.NoError(t, err)
				authRepo.EXPECT().VerifyRefreshToken(gomock.Eq(token)).Times(1).Return(sql.ErrNoRows)
				authRepo.EXPECT().DeleteRefreshToken(gomock.Any()).Times(0)
				authSvc := service.NewAuth(authRepo)
				userSvc := service.NewUser(userRepo)
				return NewAuth(authSvc, userSvc, jwt), token
			},
			wantCode: http.StatusUnauthorized,
			respFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				respBody, err := json.Marshal(map[string]any{
					"status":  "fail",
					"message": "invalid credentials",
				})
				assert.NoError(t, err)
				assert.JSONEq(t, string(respBody), rr.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authH, token := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/auth", nil)
			req.AddCookie(&http.Cookie{
				Name:     "refreshToken",
				Value:    token,
				Path:     "/api/v1",
				Domain:   "localhost",
				MaxAge:   2592000,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			})
			c.Request = req

			authH.deleteAuthHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			tt.respFunc(t, rr)
		})
	}
}
