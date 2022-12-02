package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/domain"
	mockrepo "github.com/xyedo/blindate/pkg/repository/mock"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

type onlineUserMatcher struct {
	arg *domain.Online
}

func (e onlineUserMatcher) Matches(x any) bool {
	arg, ok := x.(*domain.Online)
	if !ok {
		return false
	}
	return e.arg.UserId == arg.UserId && e.arg.IsOnline == arg.IsOnline
}

func (e onlineUserMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqOnlineUser(arg *domain.Online) gomock.Matcher {
	return onlineUserMatcher{
		arg: arg,
	}
}
func Test_postUserOnlineHandler(t *testing.T) {
	validUserId := util.RandomUUID()

	tests := []struct {
		name      string
		userId    string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *online
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name:   "valid online post",
			userId: validUserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				online := &domain.Online{
					UserId:     validUserId,
					LastOnline: time.Now(),
					IsOnline:   false,
				}
				onlineRepo.EXPECT().InsertNewOnline(EqOnlineUser(online)).Times(1).Return(nil)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user-online created",
			},
		},
		{
			name:   "userId not found",
			userId: validUserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				online := &domain.Online{
					UserId:     validUserId,
					LastOnline: time.Now(),
					IsOnline:   false,
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "user_id",
				}
				onlineRepo.EXPECT().InsertNewOnline(EqOnlineUser(online)).Times(1).Return(&pqErr)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "provided userId is not match with our resource",
			},
		},
		{
			name:   "duplicate userId",
			userId: validUserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				online := &domain.Online{
					UserId:     validUserId,
					LastOnline: time.Now(),
					IsOnline:   false,
				}
				pqErr := pq.Error{
					Code:       "23505",
					Constraint: "user_id",
				}
				onlineRepo.EXPECT().InsertNewOnline(EqOnlineUser(online)).Times(1).Return(&pqErr)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "already created",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			onlineH := tt.setupFunc(t, ctrl)
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.userId)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/users/%s/online", tt.userId), nil)
			c.Request = req

			onlineH.postUserOnlineHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			require.Contains(t, rr.Header().Get("Content-Type"), "application/json")

			if tt.wantResp != nil {
				expResBody, err := json.Marshal(tt.wantResp)
				require.NoError(t, err)
				assert.JSONEq(t, string(expResBody), rr.Body.String())
			}
		})
	}
}

func Test_getUserOnlineHandler(t *testing.T) {
	validUserId := util.RandomUUID()

	tests := []struct {
		name      string
		userId    string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *online
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name:   "valid online get",
			userId: validUserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				online := &domain.Online{
					UserId:     validUserId,
					LastOnline: time.Now(),
					IsOnline:   false,
				}
				onlineRepo.EXPECT().SelectOnline(gomock.Eq(validUserId)).Times(1).Return(online, nil)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "userId not found",
			userId: validUserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				onlineRepo.EXPECT().SelectOnline(gomock.Eq(validUserId)).Times(1).Return(nil, sql.ErrNoRows)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "resource not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			onlineH := tt.setupFunc(t, ctrl)
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.userId)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/online", tt.userId), nil)
			c.Request = req

			onlineH.getUserOnlineHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			require.Contains(t, rr.Header().Get("Content-Type"), "application/json")

			if tt.wantResp != nil {
				expResBody, err := json.Marshal(tt.wantResp)
				require.NoError(t, err)
				assert.JSONEq(t, string(expResBody), rr.Body.String())
			}
		})
	}
}
func Test_putsUserOnlineHandler(t *testing.T) {
	validUserId := util.RandomUUID()

	tests := []struct {
		name      string
		userId    string
		reqBody   map[string]bool
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *online
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name:   "valid online put true",
			userId: validUserId,
			reqBody: map[string]bool{
				"online": true,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				onlineRepo.EXPECT().UpdateOnline(gomock.Eq(validUserId), gomock.Eq(true)).Times(1).Return(nil)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user-online updated",
			},
		},
		{
			name:   "valid online put false",
			userId: validUserId,
			reqBody: map[string]bool{
				"online": false,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				onlineRepo.EXPECT().UpdateOnline(gomock.Eq(validUserId), gomock.Eq(false)).Times(1).Return(nil)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "user-online updated",
			},
		},
		{
			name:   "userId not Found",
			userId: validUserId,
			reqBody: map[string]bool{
				"online": true,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *online {
				onlineRepo := mockrepo.NewMockOnline(ctrl)
				onlineRepo.EXPECT().UpdateOnline(gomock.Eq(validUserId), gomock.Eq(true)).Times(1).Return(sql.ErrNoRows)
				onlineSvc := service.NewOnline(onlineRepo)
				return NewOnline(onlineSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "resource not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			onlineH := tt.setupFunc(t, ctrl)
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.userId)

			reqJsonBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/online", tt.userId), bytes.NewReader(reqJsonBody))
			c.Request = req

			onlineH.putuserOnlineHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			require.Contains(t, rr.Header().Get("Content-Type"), "application/json")

			if tt.wantResp != nil {
				expResBody, err := json.Marshal(tt.wantResp)
				require.NoError(t, err)
				assert.JSONEq(t, string(expResBody), rr.Body.String())
			}
		})
	}
}
