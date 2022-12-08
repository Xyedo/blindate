package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/domain/entity"
	mockrepo "github.com/xyedo/blindate/pkg/repository/mock"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_PostLocationByUserIdHandler(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		reqBody   map[string]any
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *Location
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid Body && Id",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: map[string]any{
				"lat": "80.23231",
				"lng": "170.12112",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := &entity.Location{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Geog:   "Point(80.23231 170.12112)",
				}
				locationRepo.EXPECT().InsertNewLocation(gomock.Eq(location)).Times(1).Return(int64(1), nil)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status":  "success",
				"message": "location created",
			},
		},
		{
			name: "Valid Body but Invalid Id",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "80.23231",
				"lng": "170.12112",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := &entity.Location{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68618",
					Geog:   "Point(80.23231 170.12112)",
				}
				pqErr := pq.Error{
					Code: "23503",
				}
				locationRepo.EXPECT().InsertNewLocation(gomock.Eq(location)).Times(1).Return(int64(0), &pqErr)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "provided userId is not match with our resource",
			},
		},
		{
			name: "Valid Body but userId already created",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "80.23231",
				"lng": "170.12112",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := &entity.Location{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68618",
					Geog:   "Point(80.23231 170.12112)",
				}
				pqErr := pq.Error{
					Code: "23505",
				}
				locationRepo.EXPECT().InsertNewLocation(gomock.Eq(location)).Times(1).Return(int64(0), &pqErr)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "already created",
			},
		},
		{
			name: "Invalid Body",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": 80.1391,
				"lng": "170.12112",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				locationRepo.EXPECT().InsertNewLocation(gomock.Any()).Times(0)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"lat\"",
			},
		},
		{
			name: "Invalid Lat validation",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "-91.2323232",
				"lng": "170.12112",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				locationRepo.EXPECT().InsertNewLocation(gomock.Any()).Times(0)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"lat": "required and must be valid lat geometry",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			locationH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)
			reqstring, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/users/%s/location", tt.id), bytes.NewBuffer(reqstring))
			c.Request = req
			locationH.postLocationByUserIdHandler(c)
			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			expResBody, err := json.Marshal(tt.wantResp)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}

func Test_getLocationByUserIdHandler(t *testing.T) {
	validLoc := createNewLocation(t)
	geog := strings.TrimPrefix(validLoc.Geog, "Point(")
	geog = strings.TrimSuffix(geog, ")")
	latlng := strings.Fields(geog)
	tests := []struct {
		name      string
		id        string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *Location
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "valid Id",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := validLoc
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Eq(location.UserId)).Times(1).Return(location, nil)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"location": map[string]string{
						"lat": latlng[0],
						"lng": latlng[1],
					},
				},
			},
		},
		{
			name: "Invalid Id",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Eq(location.UserId)).Times(1).Return(nil, sql.ErrNoRows)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			locationH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/location", tt.id), nil)
			c.Request = req
			locationH.getLocationByUserIdHandler(c)
			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			if tt.wantResp != nil {
				expResBody, err := json.Marshal(tt.wantResp)
				assert.NoError(t, err)
				assert.JSONEq(t, string(expResBody), rr.Body.String())

			}
		})
	}
}

func Test_patchLocationByUserIdHandler(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		reqBody   map[string]any
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *Location
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid Patching",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "70.1891",
				"lng": "80.1291",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Eq(location.UserId)).Times(1).Return(location, nil)
				location.Geog = fmt.Sprintf("Point(%0.4f %0.4f)", 70.1891, 80.1291)
				locationRepo.EXPECT().UpdateLocation(gomock.Eq(location)).Times(1).Return(int64(1), nil)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "location updated",
			},
		},
		{
			name: "Valid But only lng",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lng": "80.1291",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Eq(location.UserId)).Times(1).Return(location, nil)
				location.Geog = fmt.Sprintf("Point(%0.4f %0.4f)", 70.1891, 80.1291)
				locationRepo.EXPECT().UpdateLocation(gomock.Eq(location)).Times(1).Return(int64(1), nil)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "location updated",
			},
		},
		{
			name: "Valid But only lat",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "70.1891",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Eq(location.UserId)).Times(1).Return(location, nil)
				location.Geog = fmt.Sprintf("Point(%0.4f %0.4f)", 70.1891, 80.1291)
				locationRepo.EXPECT().UpdateLocation(gomock.Eq(location)).Times(1).Return(int64(1), nil)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "location updated",
			},
		},
		{
			name: "Invalid Patching",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "70.1891asdad",
				"lng": "80.1291",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Any()).Times(0)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"lat": "must be valid lat geometry",
				},
			},
		},
		{
			name: "Invalid Patching On Lng",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "70.1891",
				"lng": "80.1291asda",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Any()).Times(0)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"lng": "must be valid lng geometry",
				},
			},
		},
		{
			name: "Invalid Body",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": 121221,
				"lng": "80.1291",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Any()).Times(0)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"lat\"",
			},
		},
		{
			name: "Invalid Id",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68618",
			reqBody: map[string]any{
				"lat": "70.1891",
				"lng": "80.1291",
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Location {
				locationRepo := mockrepo.NewMockLocation(ctrl)
				location := createNewLocation(t)
				location.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68618"
				locationRepo.EXPECT().GetLocationByUserId(gomock.Eq(location.UserId)).Times(1).Return(nil, sql.ErrNoRows)
				locationSvc := service.NewLocation(locationRepo)
				return NewLocation(locationSvc)
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
			locationH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)

			reqBody, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%s/location", tt.id), bytes.NewBuffer(reqBody))
			c.Request = req
			locationH.patchLocationByUserIdHandler(c)
			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			if tt.wantResp != nil {
				expResBody, err := json.Marshal(tt.wantResp)
				assert.NoError(t, err)
				assert.JSONEq(t, string(expResBody), rr.Body.String())

			}

		})
	}
}

func createNewLocation(t *testing.T) *entity.Location {
	return &entity.Location{
		UserId: util.RandomUUID(),
		Geog:   util.RandomPoint(5),
	}
}
