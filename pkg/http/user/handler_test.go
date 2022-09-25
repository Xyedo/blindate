package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/mock"
	"github.com/xyedo/blindate/pkg/http/user"
)

func Test_CreateUser(t *testing.T) {

	userH := user.New(mock.UserService{})
	gin.SetMode(gin.TestMode)
	var (
		validFullName = "Uncle Bob"
		validEmail    = "bob23@gmail.com"
		validPassword = "validPa$$word"
		validDOB      = "2012-04-23T18:25:43.511Z"
	)
	tests := []struct {
		name     string
		body     map[string]any
		wantCode int
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
		})
	}
}
