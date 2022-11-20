package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/domain"
	mockrepo "github.com/xyedo/blindate/pkg/repository/mock"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_AddRefreshToken(t *testing.T) {
	validToken := util.RandomUUID()
	tests := []struct {
		name string

		token   string
		stub    func(t *testing.T, ctrl *gomock.Controller) *auth
		wantRet error
	}{
		{
			name:  "valid refresh token",
			token: validToken,
			stub: func(t *testing.T, ctrl *gomock.Controller) *auth {
				authRepo := mockrepo.NewMockAuth(ctrl)
				authRepo.EXPECT().AddRefreshToken(gomock.Eq(validToken)).Return(int64(1), nil).Times(1)
				return NewAuth(authRepo)
			},
			wantRet: nil,
		},
		{
			name:  "return context canceled",
			token: util.RandomUUID(),
			stub: func(t *testing.T, ctrl *gomock.Controller) *auth {
				authRepo := mockrepo.NewMockAuth(ctrl)
				authRepo.EXPECT().AddRefreshToken(gomock.Any()).Return(int64(0), context.Canceled).Times(1)
				return NewAuth(authRepo)
			},
			wantRet: domain.ErrTooLongAccesingDB,
		},
		{
			name:  "err pqErr",
			token: util.RandomUUID(),
			stub: func(t *testing.T, ctrl *gomock.Controller) *auth {
				authRepo := mockrepo.NewMockAuth(ctrl)
				authRepo.EXPECT().AddRefreshToken(gomock.Any()).Return(int64(0), &pq.Error{Code: pq.ErrorCode("23505"), Constraint: "token_primary"}).Times(1)
				return NewAuth(authRepo)
			},
			wantRet: domain.ErrUniqueConstraint23505,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authSvc := tt.stub(t, ctrl)
			err := authSvc.AddRefreshToken(tt.token)
			assert.EqualValues(t, err, tt.wantRet)
		})
	}
}
