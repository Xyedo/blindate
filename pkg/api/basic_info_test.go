package api

import (
	"context"
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
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	mockrepo "github.com/xyedo/blindate/pkg/repository/mock"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_postBasicInfoHandler(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		reqBody   string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) basicinfo
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid Body fullReq",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"fromLoc":"Jakarta Indonesia",
				"height":173,
				"educationLevel":"Bachelor",
				"drinking":"Ocassionally",
				"smoking":"Never",
				"relationshipPref":"Serious",
				"lookingFor":"Female",
				"zodiac":"Virgo",
				"kids": 0,
				"Work":"Software Engineer"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender: "Male",
					FromLoc: sql.NullString{
						Valid:  true,
						String: "Jakarta Indonesia",
					},
					Height: sql.NullInt16{
						Valid: true,
						Int16: 173,
					},
					EducationLevel: sql.NullString{
						Valid:  true,
						String: "Bachelor",
					},
					Drinking: sql.NullString{
						Valid:  true,
						String: "Ocassionally",
					},
					Smoking: sql.NullString{
						Valid:  true,
						String: "Never",
					},
					RelationshipPref: sql.NullString{
						Valid:  true,
						String: "Serious",
					},
					LookingFor: "Female",
					Zodiac: sql.NullString{
						Valid:  true,
						String: "Virgo",
					},
					Kids: sql.NullInt16{
						Valid: true,
						Int16: 0,
					},
					Work: sql.NullString{
						Valid:  true,
						String: "Software Engineer",
					},
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(1), nil)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status":  "success",
				"message": "basic info created!",
			},
		},
		{
			name: "Valid Body but Only Required",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId:     "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender:     "Male",
					LookingFor: "Female",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(1), nil)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status":  "success",
				"message": "basic info created!",
			},
		},
		{
			name: "Duplicate UserId",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId:     "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender:     "Male",
					LookingFor: "Female",
				}
				pqErr := &pq.Error{
					Code:       "23505",
					Constraint: "basic_info_pkey",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "already created",
			},
		},
		{
			name: "Invalid Body",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":0,
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Any()).Times(0)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"gender\"",
			},
		},
		{
			name: "Valid Body but Invalid Required Gender",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Any()).Times(0)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"gender": "required and must have an gender enums",
				},
			},
		},
		{
			name: "Invalid Field on Gender enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Non-Binary",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId:     "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender:     "Non-Binary",
					LookingFor: "Female",
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_gender",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"gender": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on EducationLevel enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"educationLevel":"Bachelor",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender: "Male",
					EducationLevel: sql.NullString{
						Valid:  true,
						String: "Bachelor",
					},
					LookingFor: "Female",
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_education_level",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"educationLevel": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on Drinking enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"drinking":"Sukakkulah",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender: "Male",
					Drinking: sql.NullString{
						Valid:  true,
						String: "Sukakkulah",
					},
					LookingFor: "Female",
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_drinking",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"drinking": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on Smoking enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"smoking":"Sukakkulah",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender: "Male",
					Smoking: sql.NullString{
						Valid:  true,
						String: "Sukakkulah",
					},
					LookingFor: "Female",
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_smoking",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"smoking": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on RelationshipPref enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"relationshipPref":"Sukakkulah",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender: "Male",
					RelationshipPref: sql.NullString{
						Valid:  true,
						String: "Sukakkulah",
					},
					LookingFor: "Female",
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_relationship_pref",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"relationshipPref": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on LookingFor enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"lookingFor":"Non-Binary"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId:     "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender:     "Male",
					LookingFor: "Non-Binary",
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_looking_for",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"lookingFor": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on Zodiac enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"lookingFor":"Non-Binary",
				"zodiac":"Propagandalf"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId:     "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender:     "Male",
					LookingFor: "Non-Binary",
					Zodiac: sql.NullString{
						Valid:  true,
						String: "Propagandalf",
					},
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_zodiac",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"zodiac": "not valid enums",
				},
			},
		},
		{
			name: "Invalid UserId",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Non-Binary",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId:     "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Gender:     "Non-Binary",
					LookingFor: "Female",
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_user_id",
				}
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Eq(basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"user_id": "not found",
				},
			},
		},
		{
			name: "Invalid All Validation",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbafasdfasfdasfasdfasfasdfs;gbasjdgasdgba",
				"fromLoc":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbasafdasdfasfasfasfdfasfasfasdfasas;gbasjdgasdgba",
				"height":325,
				"educationLevel":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgsdfasfasfdbasdgbas;gbasjdgasdgba",
				"drinking":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdfsdfaasdfasdfasdfadsfasfagbas;gbasjdgasdgba",
				"smoking":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbfasdsdasdfasdffasdfasfafasdfasfasas;gbasjdgasdgba",
				"relationShipPref":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgasdfsdfasdfasfdasdfasdfabasdgbas;gbasjdgasdgba",				
				"lookingFor":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbas;gbasjasdfasdfasdfasfasdfasdfadgasdgba",
				"zodiac":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbas;gbasjdgasdgbaasdfasdfasdfasdfasdfasdfasfasdfasf",
				"kids":125,
				"work":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbas;gbasjdgasdgbaasdfasdfasdfasdfasdfasdfasfasdfasf"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfoRepo.EXPECT().InsertBasicInfo(gomock.Any()).Times(0)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"drinking":         "maximal character is 50",
					"educationLevel":   "maximal character is 50",
					"fromLoc":          "maximal character is 100",
					"gender":           "required and must have an gender enums",
					"height":           "must have valid height in cm",
					"kids":             "minimum is 0 and maximal number is 100",
					"lookingFor":       "required and must have an gender enums",
					"relationshipPref": "maximal character is 50",
					"smoking":          "maximal character is 50",
					"work":             "maximal character is 50",
					"zodiac":           "must have zodiac enums",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			basicInfoH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/users/%s/basic-info", tt.id), strings.NewReader(tt.reqBody))
			c.Request = req
			basicInfoH.postBasicInfoHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			expResBody, err := json.Marshal(tt.wantResp)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}

func Test_getBasicInfoHandler(t *testing.T) {
	validBasicInfo := domain.BasicInfo{
		UserId:     util.RandomUUID(),
		Gender:     "Female",
		LookingFor: "Male",
	}
	tests := []struct {
		name      string
		id        string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) basicinfo
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid Getter",
			id:   validBasicInfo.UserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := &entity.BasicInfo{
					UserId:     validBasicInfo.UserId,
					Gender:     validBasicInfo.Gender,
					LookingFor: validBasicInfo.LookingFor,
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq(validBasicInfo.UserId)).Times(1).Return(basicInfo, nil)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"basicInfo": validBasicInfo,
				},
			},
		},
		{
			name: "invalid Id",
			id:   validBasicInfo.UserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq(validBasicInfo.UserId)).Times(1).Return(nil, sql.ErrNoRows)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "resource not found",
			},
		},
		{
			name: "Accessing Too Long ",
			id:   validBasicInfo.UserId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq(validBasicInfo.UserId)).Times(1).Return(nil, context.Canceled)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusConflict,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "request conflicted, please try again",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			basicInfoH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/basic-info", tt.id), nil)
			c.Request = req
			basicInfoH.getBasicInfoHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			expResBody, err := json.Marshal(tt.wantResp)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}

func Test_patchBasicInfoHandler(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		reqBody   string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) basicinfo
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid Body fullReq",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
			"gender":"Male",
			"fromLoc":"Jakarta Indonesia",
			"height":173,
			"educationLevel":"Bachelor",
			"drinking":"Ocassionally",
			"smoking":"Never",
			"relationshipPref":"Serious",
			"lookingFor":"Female",
			"zodiac":"Virgo",
			"kids": 0,
			"Work":"Software Engineer"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfo := createValidBasicInfo()
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&basicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicInfo)).Times(1).Return(int64(1), nil)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status":  "success",
				"message": "basic info updated!",
			},
		},
		{
			name: "Accessing Too Long",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(nil, context.Canceled)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Any()).Times(0)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusConflict,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "request conflicted, please try again",
			},
		},
		{
			name: "Invalid Body",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":0,
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicinfo := createValidBasicInfo()
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&basicinfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Any()).Times(0)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"gender\"",
			},
		},
		{
			name: "Invalid All Validation",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbafasdfasfdasfasdfasfasdfs;gbasjdgasdgba",
				"fromLoc":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbasafdasdfasfasfasfdfasfasfasdfasas;gbasjdgasdgba",
				"height":325,
				"educationLevel":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgsdfasfasfdbasdgbas;gbasjdgasdgba",
				"drinking":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdfsdfaasdfasdfasdfadsfasfagbas;gbasjdgasdgba",
				"smoking":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbfasdsdasdfasdffasdfasfafasdfasfasas;gbasjdgasdgba",
				"relationShipPref":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgasdfsdfasdfasfdasdfasdfabasdgbas;gbasjdgasdgba",				
				"lookingFor":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbas;gbasjasdfasdfasdfasfasdfasdfadgasdgba",
				"zodiac":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbas;gbasjdgasdgbaasdfasdfasdfasdfasdfasdfasfasdfasf",
				"kids":125,
				"work":"MaleSDfhSOIFHoshdfasofhhaosdfhaojghaosjgbasodgjbasjgbasdgbas;gbasjdgasdgbaasdfasdfasdfasdfasdfasdfasfasdfasf"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				basicinfo := createValidBasicInfo()
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&basicinfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Any()).Times(0)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"drinking":         "maximal character is 50",
					"educationLevel":   "maximal character is 50",
					"fromLoc":          "maximal character is 100",
					"gender":           "maximal character is 25",
					"height":           "must have valid height in cm",
					"kids":             "minimum is 0 and maximal number is 100",
					"lookingFor":       "maximal character is 25",
					"relationshipPref": "maximal character is 50",
					"smoking":          "maximal character is 50",
					"work":             "maximal character is 50",
					"zodiac":           "maximal character is 50",
				},
			},
		},
		{
			name: "Invalid Field on Gender enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Non-Binary",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				validbasicInfo := createValidBasicInfo()
				basicinfo := createValidBasicInfo()
				basicinfo.Gender = "Non-Binary"
				basicinfo.LookingFor = "Female"

				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_gender",
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&validbasicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicinfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"gender": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on EducationLevel enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"educationLevel":"Bachelor",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				validBasicInfo := createValidBasicInfo()
				basicInfo := createValidBasicInfo()
				basicInfo.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68619"
				basicInfo.Gender = "Male"
				basicInfo.EducationLevel = sql.NullString{
					Valid:  true,
					String: "Bachelor",
				}
				basicInfo.LookingFor = "Female"
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_education_level",
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&validBasicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"educationLevel": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on Drinking enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"drinking":"Sukakkulah",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				validBasicInfo := createValidBasicInfo()
				basicInfo := createValidBasicInfo()
				basicInfo.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68619"
				basicInfo.Gender = "Male"
				basicInfo.Drinking = sql.NullString{
					Valid:  true,
					String: "Sukakkulah",
				}
				basicInfo.LookingFor = "Female"
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_drinking",
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&validBasicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"drinking": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on Smoking enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"smoking":"Sukakkulah",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				validBasicInfo := createValidBasicInfo()
				basicInfo := createValidBasicInfo()
				basicInfo.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68619"
				basicInfo.Gender = "Male"
				basicInfo.Smoking = sql.NullString{
					Valid:  true,
					String: "Sukakkulah",
				}
				basicInfo.LookingFor = "Female"
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_smoking",
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&validBasicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"smoking": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on RelationshipPref enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"relationshipPref":"Sukakkulah",
				"lookingFor":"Female"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				validBasicInfo := createValidBasicInfo()
				basicInfo := createValidBasicInfo()
				basicInfo.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68619"
				basicInfo.Gender = "Male"
				basicInfo.RelationshipPref = sql.NullString{
					Valid:  true,
					String: "Sukakkulah",
				}
				basicInfo.LookingFor = "Female"

				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_relationship_pref",
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&validBasicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"relationshipPref": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on LookingFor enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"lookingFor":"Non-Binary"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				validBasicInfo := createValidBasicInfo()
				basicInfo := createValidBasicInfo()
				basicInfo.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68619"
				basicInfo.Gender = "Male"
				basicInfo.LookingFor = "Non-Binary"
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_looking_for",
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&validBasicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"lookingFor": "not valid enums",
				},
			},
		},
		{
			name: "Invalid Field on Zodiac enums",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: `{
				"gender":"Male",
				"lookingFor":"Female",
				"zodiac":"Propagandalf"
				}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) basicinfo {
				basicInfoRepo := mockrepo.NewMockBasicInfo(ctrl)
				validBasicInfo := createValidBasicInfo()
				basicInfo := createValidBasicInfo()
				basicInfo.UserId = "8c540e20-75d1-4513-a8e3-72dc4bc68619"
				basicInfo.Gender = "Male"
				basicInfo.Zodiac = sql.NullString{
					Valid:  true,
					String: "Propagandalf",
				}
				basicInfo.LookingFor = "Female"
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "basic_info_zodiac",
				}
				basicInfoRepo.EXPECT().GetBasicInfoByUserId(gomock.Eq("8c540e20-75d1-4513-a8e3-72dc4bc68619")).Times(1).Return(&validBasicInfo, nil)
				basicInfoRepo.EXPECT().UpdateBasicInfo(gomock.Eq(&basicInfo)).Times(1).Return(int64(0), &pqErr)
				basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
				return NewBasicInfo(basicInfoSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"zodiac": "not valid enums",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			basicInfoH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%s/basic-info", tt.id), strings.NewReader(tt.reqBody))
			c.Request = req
			basicInfoH.patchBasicInfoHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			expResBody, err := json.Marshal(tt.wantResp)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}

func createValidBasicInfo() entity.BasicInfo {
	validBasicInfo := entity.BasicInfo{
		UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
		Gender: "Male",
		FromLoc: sql.NullString{
			Valid:  true,
			String: "Jakarta Indonesia",
		},
		Height: sql.NullInt16{
			Valid: true,
			Int16: 173,
		},
		EducationLevel: sql.NullString{
			Valid:  true,
			String: "Bachelor",
		},
		Drinking: sql.NullString{
			Valid:  true,
			String: "Ocassionally",
		},
		Smoking: sql.NullString{
			Valid:  true,
			String: "Never",
		},
		RelationshipPref: sql.NullString{
			Valid:  true,
			String: "Serious",
		},
		LookingFor: "Female",
		Zodiac: sql.NullString{
			Valid:  true,
			String: "Virgo",
		},
		Kids: sql.NullInt16{
			Valid: true,
			Int16: 0,
		},
		Work: sql.NullString{
			Valid:  true,
			String: "Software Engineer",
		},
	}
	return validBasicInfo
}
