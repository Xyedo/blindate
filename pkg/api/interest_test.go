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
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/domain"
	mockrepo "github.com/xyedo/blindate/pkg/repository/mock"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_postInterestBioHandler(t *testing.T) {

	tests := []struct {
		name      string
		id        string
		reqBody   string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *interest
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid post interest",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &domain.Bio{
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
			},
		},
		{
			name: "valid but duplicate bio",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &domain.Bio{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).Return(&pq.Error{Code: "23505", Constraint: "user_id_unique"})
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interest with this user id is already created",
			},
		},
		{
			name: "invalid body interest",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
					"Bio": "at least an empty string and maximal character length is less than 300"},
			},
		},
		{
			name: "valid but empty",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &domain.Bio{
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
			},
		},
		{
			name: "Valid but userId not found",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &domain.Bio{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).Return(service.ErrRefUserIdField)
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "userId is not match with our resource",
			},
		},
		{
			name: "Valid but userId not found",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &domain.Bio{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).Return(service.ErrRefUserIdField)
				interestRepo.EXPECT().InsertNewStats(gomock.Eq(interest.Id)).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, "alah lo"),
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "userId is not match with our resource",
			},
		},
		{
			name: "Conflict",
			id:   "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				interest := &domain.Bio{
					UserId: "8c540e20-75d1-4513-a8e3-72dc4bc68619",
					Bio:    "alah lo",
				}
				interestRepo.EXPECT().InsertInterestBio(gomock.Eq(interest)).Times(1).Return(domain.ErrTooLongAccesingDB)
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
	validBio := &domain.Interest{
		Bio: domain.Bio{
			Id:     util.RandomUUID(),
			UserId: validId,
			Bio:    "apa sih loe",
		},
	}
	tests := []struct {
		name      string
		id        string
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *interest
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid getter with bio",
			id:   validId,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
		setupFunc func(t *testing.T, ctrl *gomock.Controller) *interest
		wantCode  int
		wantResp  map[string]any
	}{
		{
			name: "Valid put",
			id:   validId,
			reqBody: `{
				"bio":"im not that good with bio"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				validBio := domain.Bio{
					Id:     util.RandomUUID(),
					UserId: validId,
					Bio:    "old bio",
				}
				interestRepo.EXPECT().SelectInterestBio(gomock.Eq(validId)).Times(1).Return(&validBio, nil)
				updatedBio := validBio
				updatedBio.Bio = "im not that good with bio"
				interestRepo.EXPECT().UpdateInterestBio(gomock.Eq(&updatedBio)).Times(1).Return(nil)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				validBio := domain.Bio{
					Id:     util.RandomUUID(),
					UserId: validId,
					Bio:    "old bio",
				}
				interestRepo.EXPECT().SelectInterestBio(gomock.Eq(validId)).Times(1).Return(&validBio, nil)
				interestRepo.EXPECT().UpdateInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"Bio": "required, maximal character length is less than 300",
				},
			},
		},
		{
			name: "Valid but Not Changed",
			id:   validId,
			reqBody: `{
				"bio":"old bio"
			}`,
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				validBio := domain.Bio{
					Id:     util.RandomUUID(),
					UserId: validId,
					Bio:    "old bio",
				}
				interestRepo.EXPECT().SelectInterestBio(gomock.Eq(validId)).Times(1).Return(&validBio, nil)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				validBio := domain.Bio{
					Id:     util.RandomUUID(),
					UserId: validId,
					Bio:    "old bio",
				}
				interestRepo.EXPECT().SelectInterestBio(gomock.Eq(validId)).Times(1).Return(&validBio, nil)
				interestRepo.EXPECT().UpdateInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors": map[string]any{
					"Bio": "required, maximal character length is less than 300",
				},
			},
		},
		{
			name: "Invalid userId",
			id:   validId,
			reqBody: fmt.Sprintf(`{
				"bio":"%s"
			}`, util.RandomString(500)),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)

				interestRepo.EXPECT().SelectInterestBio(gomock.Any()).Times(1).Return(nil, sql.ErrNoRows)
				interestRepo.EXPECT().UpdateInterestBio(gomock.Any()).Times(0)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "userId is not match with our resource",
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
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *interest

		wantCode int
		wantResp map[string]any
	}{
		{
			name:       "Valid post interest",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, validHobbies),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]domain.Hobbie, 0)
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "main",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, domain.Hobbie{
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]domain.Hobbie, 0, 11)
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "main",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "coding",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "gatau",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "pengen",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "lebih",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "dari",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "sepuluh",
				})
				pqErr := pq.Error{Code: "24514", Constraint: "interests_statistics_hobbie_count_chk"}
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).Return(&pqErr)

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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
					"Hobbies": "each hobbies must be unique and has more than 2 and less than 50 character",
				},
			},
		},
		{
			name:       "Unique but less than 2 Character",
			interestId: validId,
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, `["m", "a"]`),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]domain.Hobbie, 0)
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "main",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "coding",
				})
				pqErr := pq.Error{Code: "23503", Constraint: "interest_id_ref"}
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Eq(hobbies)).Times(1).Return(&pqErr)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]any{
				"status":  "fail",
				"message": "interestId is not found",
			},
		},
		{
			name:       "valid but violates hobbies uniqueness",
			interestId: util.RandomUUID(),
			reqBody: fmt.Sprintf(`{
				"hobbies":%s
			}`, validHobbies),
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := make([]domain.Hobbie, 0)
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "main",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "mendaki",
				})
				hobbies = append(hobbies, domain.Hobbie{
					Hobbie: "coding",
				})
				pqErr := pq.Error{Code: "23505", Constraint: "interest_id_ref"}
				interestRepo.EXPECT().InsertInterestHobbies(gomock.Any(), gomock.Eq(hobbies)).Times(1).Return(&pqErr)
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
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
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
		setupFunc  func(t *testing.T, ctrl *gomock.Controller) *interest
		wantCode   int
		wantResp   map[string]any
	}{
		{
			name:       "valid Body",
			interestId: validId,
			reqBody: map[string]any{
				"hobbies": []domain.Hobbie{
					{
						Hobbie: "playing",
					},
				},
			},
			setupFunc: func(t *testing.T, ctrl *gomock.Controller) *interest {
				interestRepo := mockrepo.NewMockInterest(ctrl)
				hobbies := []domain.Hobbie{
					{
						Hobbie: "playing",
					},
				}
				interestRepo.EXPECT().UpdateInterestHobbies(gomock.Eq(validId), gomock.Eq(hobbies)).Times(1).Return(int64(1), nil)
				interestSvc := service.NewInterest(interestRepo)
				return NewInterest(interestSvc)
			},
			wantCode: http.StatusOK,
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
