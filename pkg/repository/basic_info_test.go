package repository_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

func Test_InsertBasicInfo(t *testing.T) {
	repo := repository.NewBasicInfo(testQuery)
	tests := []struct {
		name         string
		setupFunc    func() entity.BasicInfo
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name: "Valid BasicInfo",
			setupFunc: func() entity.BasicInfo {
				return createBasicInfo(t)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Valid BasicInfo But Twice",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
			},
		},
		{
			name: "Invalid Gender Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Gender = "Non-binary"
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid User_Id Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b5"
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid Education_level Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.EducationLevel = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid Drinking Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Drinking = sql.NullString{
					Valid:  true,
					String: "IDK Man",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid Smoking Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Smoking = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid relationship_pref Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.RelationshipPref = sql.NullString{
					Valid:  true,
					String: "IDK MANS",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid looking_for Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.LookingFor = "Non-Binary"
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid Zodiac Columns",
			setupFunc: func() entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Zodiac = sql.NullString{
					Valid:  true,
					String: "Non-zodiac",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basicInfo := tt.setupFunc()
			err := repo.InsertBasicInfo(basicInfo)
			tt.expectedFunc(t, err)
		})
	}
}

func Test_GetBasicInfoByUserId(t *testing.T) {
	repo := repository.NewBasicInfo(testQuery)
	t.Run("Valid BasicInfo", func(t *testing.T) {
		expected := createBasicInfo(t)
		err := repo.InsertBasicInfo(expected)
		require.NoError(t, err)
		basicInfo, err := repo.GetBasicInfoByUserId(expected.UserId)
		require.NoError(t, err)
		assert.Equal(t, expected.Gender, basicInfo.Gender)
		assert.Equal(t, expected.FromLoc, basicInfo.FromLoc)
		assert.Equal(t, expected.Height, basicInfo.Height)
		assert.Equal(t, expected.EducationLevel, basicInfo.EducationLevel)
		assert.Equal(t, expected.Drinking, basicInfo.Drinking)
		assert.Equal(t, expected.Smoking, basicInfo.Smoking)
		assert.Equal(t, expected.RelationshipPref, basicInfo.RelationshipPref)
		assert.Equal(t, expected.Zodiac, basicInfo.Zodiac)
		assert.Equal(t, expected.Kids, basicInfo.Kids)
		assert.Equal(t, expected.Work, basicInfo.Work)
	})
	t.Run("Invalid UseriD", func(t *testing.T) {
		resc, err := repo.GetBasicInfoByUserId("e590666c-3ea8-4fda-958c-c2dc6c2599b5")
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)
		assert.Zero(t, resc)
	})
}

func Test_UpdateBasicInfo(t *testing.T) {
	repo := repository.NewBasicInfo(testQuery)
	tests := []struct {
		name         string
		setupFunc    func() entity.BasicInfo
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name: "Valid but Not Change Basic Info",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Valid BasicInfo",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Invalid Gender Columns",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Gender = "Non-binary"
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid User_Id Columns",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b6"
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrResourceNotFound)
			},
		},
		{
			name: "Invalid Education_level Columns",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.EducationLevel = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid Drinking Columns",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Drinking = sql.NullString{
					Valid:  true,
					String: "IDK Man",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid Smoking Columns",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Smoking = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid relationship_pref Columns",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.RelationshipPref = sql.NullString{
					Valid:  true,
					String: "IDK MANS",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
		{
			name: "Invalid zodiac Columns",
			setupFunc: func() entity.BasicInfo {
				repo := repository.NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Zodiac = sql.NullString{
					Valid:  true,
					String: "Non-zodiac",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrRefNotFound23503)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedBasicInfo := tt.setupFunc()
			err := repo.UpdateBasicInfo(updatedBasicInfo)
			tt.expectedFunc(t, err)
		})
	}
}

func createBasicInfo(t *testing.T) entity.BasicInfo {
	user := createNewAccount(t)
	basicInfo := entity.BasicInfo{
		UserId: user.ID,
		Gender: "Male",
		FromLoc: sql.NullString{
			Valid:  true,
			String: "Jakarta, Indonesia",
		},
		Height: sql.NullInt16{
			Valid: true,
			Int16: 173,
		},
		EducationLevel:   sql.NullString{},
		Drinking:         sql.NullString{},
		Smoking:          sql.NullString{},
		RelationshipPref: sql.NullString{},
		LookingFor:       "Female",
		Zodiac:           sql.NullString{},
		Kids:             sql.NullInt16{},
		Work: sql.NullString{
			Valid:  true,
			String: "Software Developer",
		},
	}
	return basicInfo
}
