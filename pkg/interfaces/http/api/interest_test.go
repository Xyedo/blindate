package api

import (
	"bytes"
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
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/applications/service"
	apiError "github.com/xyedo/blindate/pkg/common/error"
	"github.com/xyedo/blindate/pkg/common/util"
	interestEntity "github.com/xyedo/blindate/pkg/domain/interest/entities"
	mockrepo "github.com/xyedo/blindate/pkg/infra/repository/mock"
)

func Test_postInterestBioHandler(t *testing.T) {

	tests := []struct {
		name      string
		id        string
		reqBody   string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid post interest",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &interestEntity.BioDTO{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).Return(nil)
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(1).Return(nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},

			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status":  "success",
				"message": "interest bio created",
				"data": map[string]any{
					"interestId": "",
				},
			},
		},
		{
			name: "valid but duplicate bio",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &interestEntity.BioDTO{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				pqErr := &pq.Error{Code: "23505", Constraint: "user_id_unique"}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).Return(apiError.WrapWithMsg(pqErr, apiError.ErrUniqueConstraint23505, "interest already created"))
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interest already created",
			},
		},
		{
			name: "invalid body interest",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			reqBody: `{
			}`,
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"bio": "at least an empty string and maximal character length is less than 300"},
			},
		},
		{
			name: "valid but empty",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &interestEntity.BioDTO{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).Return(nil)
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(1).Return(nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			reqBody: `{
				"bio":""
			}`,
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status":  "success",
				"message": "interest bio created",
				"data":    map[string]interface{}{"interestId": ""},
			},
		},
		{
			name: "Valid but userId not found",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &interestEntity.BioDTO{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).
					Return(apiError.WrapWithMsg(&pq.Error{Code: "23505", Constraint: "user_id"}, apiError.ErrRefNotFound23503, "userId is not invalid"))
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "userId is not invalid",
			},
		},
		{
			name: "Conflict",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &interestEntity.BioDTO{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).
					Return(apiError.ErrTooLongAccessingDB)
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
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
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/users/%s", tt.id), strings.NewReader(tt.reqBody))
			c.Request = req

			interestH.postInterestBioHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			require.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			expResBody, err := json.Marshal(tt.wantResp)
			require.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}

func Test_getInterestBioHandler(t *testing.T) {
	validId := util.RandomUUID()
	validBio := interestEntity.FullDTO{
		BioDTO: interestEntity.BioDTO{
			Id:     util.RandomUUID(),
			UserId: validId,
			Bio:    "apa sih loe",
		},
	}
	tests := []struct {
		name      string
		id        string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid getter with bio",
			id:   validId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)

				interestRepo.EXPECT().GetInterest(gomock.Eq(validId)).Times(1).Return(validBio, nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": gin.H{
					"interest": validBio,
				},
			},
		},
		{
			name: "Invalid",
			id:   validId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)

				interestRepo.EXPECT().GetInterest(gomock.Eq(validId)).Times(1).Return(interestEntity.FullDTO{}, apiError.ErrResourceNotFound)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
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
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/users/%s", tt.id), nil)
			c.Request = req

			interestH.getInterestHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			require.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			expResBody, err := json.Marshal(tt.wantResp)
			require.NoError(t, err)
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}

func Test_putInterestBioHandler(t *testing.T) {
	validId := util.RandomUUID()

	tests := []struct {
		name      string
		id        string
		reqBody   string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid put",
			id:   validId,
			reqBody: `{
				"bio":"im not that good with bio"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				validBio := interestEntity.BioDTO{
					Id:     util.RandomUUID(),
					UserId: validId,
					Bio:    "old bio",
				}
				interestRepo.EXPECT().SelectInterestBio(gomock.Eq(validId)).Times(1).Return(validBio, nil)
				updatedBio := validBio
				updatedBio.Bio = "im not that good with bio"
				interestRepo.EXPECT().UpdateInterestBio(gomock.Eq(updatedBio)).Times(1).Return(nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "bio successfully changed",
			},
		},
		{
			name: "invalid Bio",
			id:   validId,
			reqBody: `{
				"bio": null
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().SelectInterestBio(gomock.Any()).Times(0)
				interestRepo.EXPECT().UpdateInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"bio": "required, maximal character length is less than 300",
				},
			},
		},
		{
			name: "Valid but Not Changed",
			id:   validId,
			reqBody: `{
				"bio":"old bio"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				validBio := interestEntity.BioDTO{
					Id:     util.RandomUUID(),
					UserId: validId,
					Bio:    "old bio",
				}
				interestRepo.EXPECT().SelectInterestBio(gomock.Eq(validId)).Times(1).Return(validBio, nil)
				interestRepo.EXPECT().UpdateInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status":  "success",
				"message": "nothing changed",
			},
		},
		{
			name: "Invalid, Too long Bio",
			id:   validId,
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, util.RandomString(500)),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().SelectInterestBio(gomock.Eq(validId)).Times(0)
				interestRepo.EXPECT().UpdateInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"bio": "required, maximal character length is less than 300",
				},
			},
		},
		{
			name: "Invalid userId",
			id:   validId,
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, util.RandomString(12)),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)

				interestRepo.EXPECT().SelectInterestBio(gomock.Any()).Times(1).
					Return(interestEntity.BioDTO{}, apiError.WrapWithMsg(
						&pq.Error{Code: "23503", Constraint: "user_id"},
						apiError.ErrRefNotFound23503,
						"userId is not invalid"),
					)
				interestRepo.EXPECT().UpdateInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "userId is not invalid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("userId", tt.id)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", tt.id), strings.NewReader(tt.reqBody))
			c.Request = req

			interestH.putInterestBioHandler(c)

			assert.Equal(t, tt.wantCode, rr.Code)
			require.Contains(t, rr.Header().Get("Content-Type"), "application/json")
			expResBody, err := json.Marshal(tt.wantResp)
			require.NoError(t, err)
			require.NotZero(t, rr.Body.String())
			assert.JSONEq(t, string(expResBody), rr.Body.String())
		})
	}
}
func Test_postInterestHobbiesHandler(t *testing.T) {
	validId := util.RandomUUID()
	validHobbies := `["main", "mendaki", "coding"]`
	tests := []struct {
		name       string
		interestId string
		reqBody    string
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest

		wantCode int
		wantResp map[string]any
	}{
		{
			name:       "Valid post interest",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, validHobbies),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]interestEntity.HobbieDTO, 0)
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "main",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "coding",
				})
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).Return(nil)

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"hobbies": []map[string]any{
						{
							"hobbie": "main",
						},
						{
							"hobbie": "mendaki",
						},
						{
							"hobbie": "coding",
						},
					},
				},
			},
		},
		{
			name:       "Valid hobbies but too much over 10",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, `["main", "mendaki","coding", "gatau", "pengen", "lebih", "dari", "sepuluh"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]interestEntity.HobbieDTO, 0, 11)
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "main",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "coding",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "gatau",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "pengen",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "lebih",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "dari",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "sepuluh",
				})
				pqErr := pq.Error{Code: "23514", Constraint: "interests_statistics_hobbie_count_chk"}
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "hobbies must less than 10"))

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "hobbies must less than 10",
			},
		},
		{
			name:       "Non Unique Hobbie",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, `["main", "main"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"hobbies": "each hobbies must be unique, less than 10 and has more than 2 and less than 50 character",
				},
			},
		},
		{
			name:       "Unique but less than 2 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, `["m", "a"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Invalid InterestId",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, validHobbies),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]interestEntity.HobbieDTO, 0)
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "main",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "coding",
				})
				pqErr := pq.Error{Code: "23503", Constraint: "interest_id_ref"}
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Eq(hobbies)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interestId is invalid",
			},
		},
		{
			name:       "valid but violates hobbies uniqueness",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, validHobbies),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]interestEntity.HobbieDTO, 0)
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "main",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, interestEntity.HobbieDTO{
					Hobbie: "coding",
				})
				pqErr := pq.Error{Code: "23505", Constraint: "hobbie_unique"}
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Eq(hobbies)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every hobbies must be unique"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every hobbies must be unique",
			},
		},
		{
			name:       "invalid Body",
			interestId: util.RandomUUID(),
			reqBody: `{
				"hobbies":"not valid array"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"hobbies\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/interests/%s", tt.interestId), strings.NewReader(tt.reqBody))
			c.Request = req

			interestH.postInterestHobbiesHandler(c)

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

func Test_putInterestHobbiesHandler(t *testing.T) {
	validId := util.RandomUUID()

	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"hobbies": []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).Return(nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"hobbies": []any{
						map[string]any{
							"hobbie": "playing",
						},
					},
				},
			},
		},
		{
			name:       "duplicate hobbie in db level",
			interestId: validId,
			reqBody: map[string]any{
				"hobbies": []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23505",
					Constraint: "hobbie_unique",
				}
				interestRepo.EXPECT().UpdateInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every hobbies must be unique"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every hobbies must be unique",
			},
		},
		{
			name:       "hobbies more than 10",
			interestId: validId,
			reqBody: map[string]any{
				"hobbies": []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23514",
					Constraint: "hobbie_count",
				}
				interestRepo.EXPECT().UpdateInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "hobbies must less than 10"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "hobbies must less than 10",
			},
		},
		{
			name:       "invalid Ref Id",
			interestId: validId,
			reqBody: map[string]any{
				"hobbies": []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "interest_id",
				}
				interestRepo.EXPECT().UpdateInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "non-unique hobbies",
			interestId: validId,
			reqBody: map[string]any{
				"hobbies": []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
					{
						Hobbie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := []interestEntity.HobbieDTO{
					{
						Hobbie: "playing",
					},
					{
						Hobbie: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"hobbies": "each hobbies must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new hobbies",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.putInterestHobbiesHandler(c)

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

func Test_deleteInterestHobbiesHandler(t *testing.T) {
	validId := util.RandomUUID()
	hobbiesId := util.RandomUUID()
	insertedIds := []string{util.RandomUUID(), util.RandomUUID(), util.RandomUUID()}
	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"ids": []string{
					util.RandomUUID(), util.RandomUUID(), util.RandomUUID(),
				},
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestHobbies(gomock.Eq(validId), gomock.Any()).Times(1).
					Return([]string{util.RandomUUID(), util.RandomUUID(), util.RandomUUID()}, nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
		},
		{
			name:       "one of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": insertedIds,
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestHobbies(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(insertedIds[:len(insertedIds)-1], nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"deletedIds": insertedIds[:len(insertedIds)-1],
				},
			},
		},
		{
			name:       "all of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": []string{
					util.RandomUUID(), util.RandomUUID(), util.RandomUUID(),
				},
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestHobbies(gomock.Eq(validId), gomock.Any()).Times(1).Return(nil, apiError.ErrResourceNotFound)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "resource not found",
			},
		},
		{
			name:       "Non unique id",
			interestId: validId,
			reqBody: map[string]any{
				"ids": []string{
					hobbiesId, hobbiesId,
				},
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestHobbies(gomock.Eq(validId), gomock.Eq([]string{hobbiesId, hobbiesId})).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"ids": "each ids must be unique and uuid character",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.deleteInterestHobbiesHandler(c)

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

func Test_postInterestMovieSeriesHandler(t *testing.T) {
	validId := util.RandomUUID()
	validMovieSeries := `["main", "mendaki", "coding"]`
	tests := []struct {
		name       string
		interestId string
		reqBody    string
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest

		wantCode int
		wantResp map[string]any
	}{
		{
			name:       "Valid post interest",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, validMovieSeries),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := make([]interestEntity.MovieSerieDTO, 0)
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "main",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "mendaki",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "coding",
				})
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Eq(validId), gomock.Eq(movieSeries)).Times(1).Return(nil)

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"movieSeries": []map[string]any{
						{
							"movieSerie": "main",
						},
						{
							"movieSerie": "mendaki",
						},
						{
							"movieSerie": "coding",
						},
					},
				},
			},
		},
		{
			name:       "Valid movieSeries but too much over 10",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, `["main", "mendaki","coding", "gatau", "pengen", "lebih", "dari", "sepuluh"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := make([]interestEntity.MovieSerieDTO, 0, 11)
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "main",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "mendaki",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "coding",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "gatau",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "pengen",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "lebih",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "dari",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "sepuluh",
				})
				pqErr := pq.Error{Code: "23514", Constraint: "interests_statistics_movie_serie_count_chk"}
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Eq(validId), gomock.Eq(movieSeries)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "movieSeries must less than 10"))

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "movieSeries must less than 10",
			},
		},
		{
			name:       "Non Unique movieSeries",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, `["main", "main"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"movieSeries": "each movieSeries must be unique, less than 10 and has more than 2 and less than 50 character",
				},
			},
		},
		{
			name:       "Unique but less than 2 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, `["m", "a"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Invalid InterestId",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, validMovieSeries),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := make([]interestEntity.MovieSerieDTO, 0)
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "main",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "mendaki",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "coding",
				})
				pqErr := pq.Error{Code: "23503", Constraint: "interest_id_ref"}
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Any(), gomock.Eq(movieSeries)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interestId is invalid",
			},
		},
		{
			name:       "valid but violates movieSeries uniqueness",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"movieSeries":%s
			}`, validMovieSeries),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := make([]interestEntity.MovieSerieDTO, 0)
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "main",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "mendaki",
				})
				movieSeries = append(movieSeries, interestEntity.MovieSerieDTO{
					MovieSerie: "coding",
				})
				pqErr := pq.Error{Code: "23505", Constraint: "movie_serie_unique"}
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Any(), gomock.Eq(movieSeries)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every moviesSeries must be unique"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every moviesSeries must be unique",
			},
		},
		{
			name:       "invalid Body",
			interestId: util.RandomUUID(),
			reqBody: `{
				"movieSeries":"not valid array"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestMovieSeries(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"movieSeries\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/interests/%s", tt.interestId), strings.NewReader(tt.reqBody))
			c.Request = req

			interestH.postInterestMovieSeriesHandler(c)

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

func Test_putInterestMovieSeriesHandler(t *testing.T) {
	validId := util.RandomUUID()

	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"movieSeries": []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestMovieSeries(gomock.Eq(validId), gomock.Eq(movieSeries)).Times(1).Return(nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"movieSeries": []any{
						map[string]any{
							"movieSerie": "playing",
						},
					},
				},
			},
		},
		{
			name:       "duplicate movieSeries in db level",
			interestId: validId,
			reqBody: map[string]any{
				"movieSeries": []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23505",
					Constraint: "movie_serie_unique",
				}
				interestRepo.EXPECT().UpdateInterestMovieSeries(gomock.Eq(validId), gomock.Eq(movieSeries)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every moviesSeries must be unique"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every moviesSeries must be unique",
			},
		},
		{
			name:       "movieSeries more than 10",
			interestId: validId,
			reqBody: map[string]any{
				"movieSeries": []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "24514",
					Constraint: "movie_serie_count",
				}
				interestRepo.EXPECT().UpdateInterestMovieSeries(gomock.Eq(validId), gomock.Eq(movieSeries)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "movieSeries must less than 10"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "movieSeries must less than 10",
			},
		},
		{
			name:       "invalid Ref Id",
			interestId: validId,
			reqBody: map[string]any{
				"movieSeries": []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "interest_id",
				}
				interestRepo.EXPECT().UpdateInterestMovieSeries(gomock.Eq(validId), gomock.Eq(movieSeries)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "non-unique hobbies",
			interestId: validId,
			reqBody: map[string]any{
				"movieSeries": []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
					{
						MovieSerie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				movieSeries := []interestEntity.MovieSerieDTO{
					{
						MovieSerie: "playing",
					},
					{
						MovieSerie: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestMovieSeries(gomock.Eq(validId), gomock.Eq(movieSeries)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"movieSeries": "each movieSeries must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new movieSeries",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.putInterestMovieSeriesHandler(c)

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

func Test_deleteInterestMovieSeriesHandler(t *testing.T) {
	validId := util.RandomUUID()
	movieSeriesId := util.RandomUUID()
	movieSeriesIds := []string{util.RandomUUID(), util.RandomUUID(), util.RandomUUID()}
	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"ids": movieSeriesIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestMovieSeries(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(movieSeriesIds, nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"deletedIds": movieSeriesIds,
				},
			},
		},
		{
			name:       "One of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": movieSeriesIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestMovieSeries(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(movieSeriesIds[:len(movieSeriesIds)-1], nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"deletedIds": movieSeriesIds[:len(movieSeriesIds)-1],
				},
			},
		},
		{
			name:       "all of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": movieSeriesIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestMovieSeries(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(nil, apiError.ErrResourceNotFound)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:       "Non unique id",
			interestId: validId,
			reqBody: map[string]any{
				"ids": []string{
					movieSeriesId, movieSeriesId,
				},
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestMovieSeries(gomock.Eq(validId), gomock.Eq([]string{movieSeriesId, movieSeriesId})).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"ids": "each ids must be unique and uuid character",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.deleteInterestMovieSeriesHandler(c)

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

func Test_postInterestTravelingHandler(t *testing.T) {
	validId := util.RandomUUID()
	validTravels := `["main", "mendaki", "coding"]`
	tests := []struct {
		name       string
		interestId string
		reqBody    string
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest

		wantCode int
		wantResp map[string]any
	}{
		{
			name:       "Valid post interest",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, validTravels),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := make([]interestEntity.TravelDTO, 0)
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "main",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "mendaki",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "coding",
				})
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Eq(validId), gomock.Eq(travels)).Times(1).Return(nil)

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"travels": []map[string]any{
						{
							"travel": "main",
						},
						{
							"travel": "mendaki",
						},
						{
							"travel": "coding",
						},
					},
				},
			},
		},
		{
			name:       "Valid movieSeries but too much over 10",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, `["main", "mendaki","coding", "gatau", "pengen", "lebih", "dari", "sepuluh"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := make([]interestEntity.TravelDTO, 0, 11)
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "main",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "mendaki",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "coding",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "gatau",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "pengen",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "lebih",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "dari",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "sepuluh",
				})
				pqErr := pq.Error{Code: "23514", Constraint: "interests_statistics_traveling_count_chk"}
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Eq(validId), gomock.Eq(travels)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "travels must less than 10"))

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "travels must less than 10",
			},
		},
		{
			name:       "Non Unique travels",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, `["main", "main"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"travels": "each travels must be unique, less than 10 and has more than 2 and less than 50 character",
				},
			},
		},
		{
			name:       "Unique but less than 2 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, `["m", "a"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Invalid InterestId",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, validTravels),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := make([]interestEntity.TravelDTO, 0)
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "main",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "mendaki",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "coding",
				})
				pqErr := pq.Error{Code: "23503", Constraint: "interest_id_ref"}
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Any(), gomock.Eq(travels)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interestId is invalid",
			},
		},
		{
			name:       "valid but violates travels uniqueness",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"travels":%s
			}`, validTravels),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := make([]interestEntity.TravelDTO, 0)
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "main",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "mendaki",
				})
				travels = append(travels, interestEntity.TravelDTO{
					Travel: "coding",
				})
				pqErr := pq.Error{Code: "23505", Constraint: "travel_unique"}
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Any(), gomock.Eq(travels)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every travels must be unique"))

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every travels must be unique",
			},
		},
		{
			name:       "invalid Body",
			interestId: util.RandomUUID(),
			reqBody: `{
				"travels":"not valid array"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestTraveling(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"travels\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/interests/%s", tt.interestId), strings.NewReader(tt.reqBody))
			c.Request = req

			interestH.postInterestTravelingHandler(c)

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

func Test_putInterestTravelingHandler(t *testing.T) {
	validId := util.RandomUUID()

	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"travels": []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestTraveling(gomock.Eq(validId), gomock.Eq(travels)).Times(1).Return(nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"travels": []any{
						map[string]any{
							"travel": "playing",
						},
					},
				},
			},
		},
		{
			name:       "duplicate travels in db level",
			interestId: validId,
			reqBody: map[string]any{
				"travels": []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23505",
					Constraint: "travel_unique",
				}
				interestRepo.EXPECT().UpdateInterestTraveling(gomock.Eq(validId), gomock.Eq(travels)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every travels must be unique"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every travels must be unique",
			},
		},
		{
			name:       "travels more than 10",
			interestId: validId,
			reqBody: map[string]any{
				"travels": []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23514",
					Constraint: "traveling_count",
				}
				interestRepo.EXPECT().UpdateInterestTraveling(gomock.Eq(validId), gomock.Eq(travels)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "travels must less than 10"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "travels must less than 10",
			},
		},
		{
			name:       "invalid Ref Id",
			interestId: validId,
			reqBody: map[string]any{
				"travels": []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "interest_id",
				}
				interestRepo.EXPECT().UpdateInterestTraveling(gomock.Eq(validId), gomock.Eq(travels)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interestId is invalid",
			},
		},
		{
			name:       "non-unique hobbies",
			interestId: validId,
			reqBody: map[string]any{
				"travels": []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
					{
						Travel: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				travels := []interestEntity.TravelDTO{
					{
						Travel: "playing",
					},
					{
						Travel: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestTraveling(gomock.Eq(validId), gomock.Eq(travels)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"travels": "each travels must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new travel.",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.putInterestTravelingHandler(c)

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

func Test_deleteInterestTravelingHandler(t *testing.T) {
	validId := util.RandomUUID()
	travelId := util.RandomUUID()
	travelIds := []string{util.RandomUUID(), util.RandomUUID(), util.RandomUUID()}
	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"ids": []string{
					util.RandomUUID(), util.RandomUUID(), util.RandomUUID(),
				},
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestTraveling(gomock.Eq(validId), gomock.Any()).Times(1).Return(travelIds, nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"deletedIds": travelIds,
				},
			},
		},
		{
			name:       "One of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": travelIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestTraveling(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(travelIds[:len(travelIds)-1], nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"deletedIds": travelIds[:len(travelIds)-1],
				},
			},
		},
		{
			name:       "all of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": travelIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestTraveling(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(nil, apiError.ErrResourceNotFound)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:       "Non unique id",
			interestId: validId,
			reqBody: map[string]any{
				"ids": []string{
					travelId, travelId,
				},
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestTraveling(gomock.Eq(validId), gomock.Eq([]string{travelId, travelId})).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"ids": "each ids must be unique and uuid character",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.deleteInterestTravelingHandler(c)

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
func Test_postInterestSportsHandler(t *testing.T) {
	validId := util.RandomUUID()
	validSports := `["main", "mendaki", "coding"]`
	tests := []struct {
		name       string
		interestId string
		reqBody    string
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest

		wantCode int
		wantResp map[string]any
	}{
		{
			name:       "Valid post *Interest",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, validSports),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := make([]interestEntity.SportDTO, 0)
				sports = append(sports, interestEntity.SportDTO{
					Sport: "main",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "mendaki",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "coding",
				})
				interestRepo.EXPECT().InsertInterestSports(gomock.Eq(validId), gomock.Eq(sports)).Times(1).Return(nil)

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusCreated,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"sports": []map[string]any{
						{
							"sport": "main",
						},
						{
							"sport": "mendaki",
						},
						{
							"sport": "coding",
						},
					},
				},
			},
		},
		{
			name:       "Valid sports but too much over 10",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, `["main", "mendaki","coding", "gatau", "pengen", "lebih", "dari", "sepuluh"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := make([]interestEntity.SportDTO, 0, 11)
				sports = append(sports, interestEntity.SportDTO{
					Sport: "main",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "mendaki",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "coding",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "gatau",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "pengen",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "lebih",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "dari",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "sepuluh",
				})
				pqErr := pq.Error{Code: "23514", Constraint: "interests_statistics_sport_count_chk"}
				interestRepo.EXPECT().InsertInterestSports(gomock.Eq(validId), gomock.Eq(sports)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "sports must less than 10"))

				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "sports must less than 10",
			},
		},
		{
			name:       "Non Unique sports",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, `["main", "main"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestSports(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"sports": "each sports must be unique, less than 10 and has more than 2 and less than 50 character",
				},
			},
		},
		{
			name:       "Unique but less than 2 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, `["m", "a"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestSports(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestSports(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Unique but max than 50 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, fmt.Sprintf(`["%s", "%s"]`, util.RandomString(60), util.RandomString(60))),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestSports(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Invalid InterestId",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, validSports),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := make([]interestEntity.SportDTO, 0)
				sports = append(sports, interestEntity.SportDTO{
					Sport: "main",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "mendaki",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "coding",
				})
				pqErr := pq.Error{Code: "23503", Constraint: "interest_id_ref"}
				interestRepo.EXPECT().InsertInterestSports(gomock.Any(), gomock.Eq(sports)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interestId is invalid",
			},
		},
		{
			name:       "valid but violates sports uniqueness",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"sports":%s
			}`, validSports),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := make([]interestEntity.SportDTO, 0)
				sports = append(sports, interestEntity.SportDTO{
					Sport: "main",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "mendaki",
				})
				sports = append(sports, interestEntity.SportDTO{
					Sport: "coding",
				})
				pqErr := pq.Error{Code: "23505", Constraint: "sport_unique"}
				interestRepo.EXPECT().InsertInterestSports(gomock.Any(), gomock.Eq(sports)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every sports must be unique"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every sports must be unique",
			},
		},
		{
			name:       "invalid Body",
			interestId: util.RandomUUID(),
			reqBody: `{
				"sports":"not valid array"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().InsertInterestSports(gomock.Any(), gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "body contains incorrect JSON type for field \"sports\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/interests/%s", tt.interestId), strings.NewReader(tt.reqBody))
			c.Request = req

			interestH.postInterestSportHandler(c)

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

func Test_putInterestSportsHandler(t *testing.T) {
	validId := util.RandomUUID()

	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"sports": []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestSport(gomock.Eq(validId), gomock.Eq(sports)).Times(1).Return(nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"sports": []any{
						map[string]any{
							"sport": "playing",
						},
					},
				},
			},
		},
		{
			name:       "duplicate sports in db level",
			interestId: validId,
			reqBody: map[string]any{
				"sports": []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23505",
					Constraint: "sport_unique",
				}
				interestRepo.EXPECT().UpdateInterestSport(gomock.Eq(validId), gomock.Eq(sports)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrUniqueConstraint23505, "every sports must be unique"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "every sports must be unique",
			},
		},
		{
			name:       "sports more than 10",
			interestId: validId,
			reqBody: map[string]any{
				"sports": []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23514",
					Constraint: "sport_count",
				}
				interestRepo.EXPECT().UpdateInterestSport(gomock.Eq(validId), gomock.Eq(sports)).Times(1).
					Return(apiError.WrapWithNewSentinel(&pqErr, http.StatusUnprocessableEntity, "sports must less than 10"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "sports must less than 10",
			},
		},
		{
			name:       "invalid Ref Id",
			interestId: validId,
			reqBody: map[string]any{
				"sports": []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
				}
				pqErr := pq.Error{
					Code:       "23503",
					Constraint: "interest_id",
				}
				interestRepo.EXPECT().UpdateInterestSport(gomock.Eq(validId), gomock.Eq(sports)).Times(1).
					Return(apiError.WrapWithMsg(&pqErr, apiError.ErrRefNotFound23503, "interestId is invalid"))
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interestId is invalid",
			},
		},
		{
			name:       "non-unique sports",
			interestId: validId,
			reqBody: map[string]any{
				"sports": []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
					{
						Sport: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				sports := []interestEntity.SportDTO{
					{
						Sport: "playing",
					},
					{
						Sport: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestSport(gomock.Eq(validId), gomock.Eq(sports)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"sports": "each sports must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new sports.",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.putInterestSportHandler(c)

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

func Test_deleteInterestSportsHandler(t *testing.T) {
	validId := util.RandomUUID()
	sportId := util.RandomUUID()
	sportIds := []string{util.RandomUUID(), util.RandomUUID(), util.RandomUUID()}
	tests := []struct {
		name       string
		interestId string
		reqBody    map[string]any
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *Interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"ids": sportIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestSports(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(sportIds, nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"deletedIds": sportIds,
				},
			},
		},
		{
			name:       "One of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": sportIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestSports(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(sportIds[:len(sportIds)-1], nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
			wantResp: map[string]any{
				"status": "success",
				"data": map[string]any{
					"deletedIds": sportIds[:len(sportIds)-1],
				},
			},
		},
		{
			name:       "all of the id is not found",
			interestId: validId,
			reqBody: map[string]any{
				"ids": sportIds,
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestSports(gomock.Eq(validId), gomock.Any()).Times(1).
					Return(nil, apiError.ErrResourceNotFound)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:       "Non unique id",
			interestId: validId,
			reqBody: map[string]any{
				"ids": []string{
					sportId, sportId,
				},
			},

			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *Interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interestRepo.EXPECT().DeleteInterestSports(gomock.Eq(validId), gomock.Eq([]string{sportId, sportId})).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"ids": "each ids must be unique and uuid character",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interestH := tt.setupFunc(t, ctrl)

			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Set("interestId", tt.interestId)
			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/interests/%s", tt.interestId), bytes.NewReader(reqBody))
			c.Request = req

			interestH.deleteInterestSportHandler(c)

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
