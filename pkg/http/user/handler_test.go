package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/mock"
	"github.com/xyedo/blindate/pkg/http/user"
	"github.com/xyedo/blindate/pkg/validation"
)

func Test_PostUserHandler(t *testing.T) {

	userH := user.New(mock.UserService{})

	gin.SetMode(gin.TestMode)
	v, ok := binding.Validator.Engine().(*validator.Validate)
	assert.True(t, ok)
	err := v.RegisterValidation("validdob", validation.ValidDob)
	assert.NoError(t, err)
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
				"status": "success",
				"data": map[string]string{
					"id": "1",
				},
			},
		},
		{
			name: "Invalid Email Body",
			body: map[string]any{
				"fullName": validFullName,
				"email":    "hahahahha",
				"password": validPassword,
				"dob":      validDOB,
			},
			wantCode: http.StatusBadRequest,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Invalid Password Body",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": "012",
				"dob":      validDOB,
			},
			wantCode: http.StatusBadRequest,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
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
			wantCode: http.StatusBadRequest,
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
			name: "Password to long",
			body: map[string]any{
				"fullName": validFullName,
				"email":    validEmail,
				"password": "adamdkmadnianafnafnafnkanflafnkafnalfknanfoasnfnkaongao[jgbaobg[oawubgawgbasdasdadadawga",
				"dob":      validDOB,
			},
			wantCode: http.StatusInternalServerError,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:     "No Body",
			body:     nil,
			wantCode: http.StatusBadRequest,
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
			wantCode: http.StatusBadRequest,
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
			wantCode: http.StatusBadRequest,
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
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(bodyStr))
			c.Request = req
			userH.PostUserHandler(c)
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
