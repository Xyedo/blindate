package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/entity"
)

func Test_InsertBasicInfo(t *testing.T) {
	repo := NewBasicInfo(testQuery)
	tests := []struct {
		name         string
		setupFunc    func() *entity.BasicInfo
		expectedFunc func(t *testing.T, affectedRow int64, err error)
	}{
		{
			name: "Valid BasicInfo",
			setupFunc: func() *entity.BasicInfo {
				return createBasicInfo(t)
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, int(row))
			},
		},
		{
			name: "Valid BasicInfo But Twice",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				assert.Error(t, err)
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23505"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "basic_info_pkey"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid Gender Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Gender = "Non-binary"
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "gender"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid User_Id Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b5"
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "user_id"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid Education_level Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.EducationLevel = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "education_level"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid Drinking Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Drinking = sql.NullString{
					Valid:  true,
					String: "IDK Man",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "drinking"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid Smoking Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Smoking = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "smoking"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid relationship_pref Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.RelationshipPref = sql.NullString{
					Valid:  true,
					String: "IDK MANS",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "relationship_pref"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid looking_for Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.LookingFor = "Non-Binary"
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "looking_for"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid Zodiac Columns",
			setupFunc: func() *entity.BasicInfo {
				validBasicInfo := createBasicInfo(t)
				validBasicInfo.Zodiac = sql.NullString{
					Valid:  true,
					String: "Non-zodiac",
				}
				return validBasicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				if assert.ErrorAs(t, err, &pqErr) {
					assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
					fmt.Println(pqErr.Constraint)
					assert.True(t, strings.Contains(pqErr.Constraint, "zodiac"))
					assert.Zero(t, row)
				}

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basicInfo := tt.setupFunc()
			row, err := repo.InsertBasicInfo(basicInfo)
			tt.expectedFunc(t, row, err)
		})
	}
}

func Test_GetBasicInfoByUserId(t *testing.T) {
	repo := NewBasicInfo(testQuery)
	t.Run("Valid BasicInfo", func(t *testing.T) {
		expected := createBasicInfo(t)
		row, err := repo.InsertBasicInfo(expected)
		assert.NoError(t, err)
		assert.NotZero(t, row)
		basicInfo, err := repo.GetBasicInfoByUserId(expected.UserId)
		assert.NoError(t, err)
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
		_, err := repo.GetBasicInfoByUserId("e590666c-3ea8-4fda-958c-c2dc6c2599b5")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func Test_UpdateBasicInfo(t *testing.T) {
	repo := NewBasicInfo(testQuery)
	tests := []struct {
		name         string
		setupFunc    func() *entity.BasicInfo
		expectedFunc func(t *testing.T, affectedRow int64, err error)
	}{
		{
			name: "Valid but Not Change Basic Info",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				return basicInfo
			},
			expectedFunc: func(t *testing.T, affectedRow int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int(affectedRow), 1)
			},
		},
		{
			name: "Valid BasicInfo",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, int(row))
			},
		},
		{
			name: "Invalid Gender Columns",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Gender = "Non-binary"
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "gender"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid User_Id Columns",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b6"
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				assert.Equal(t, int(row), 0)
			},
		},
		{
			name: "Invalid Education_level Columns",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.EducationLevel = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "education_level"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid Drinking Columns",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Drinking = sql.NullString{
					Valid:  true,
					String: "IDK Man",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "drinking"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid Smoking Columns",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Smoking = sql.NullString{
					Valid:  true,
					String: "IDK man",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "smoking"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid relationship_pref Columns",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.RelationshipPref = sql.NullString{
					Valid:  true,
					String: "IDK MANS",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "relationship_pref"))
				assert.Zero(t, row)
			},
		},
		{
			name: "Invalid zodiac Columns",
			setupFunc: func() *entity.BasicInfo {
				repo := NewBasicInfo(testQuery)
				basicInfo := createBasicInfo(t)
				repo.InsertBasicInfo(basicInfo)
				basicInfo.Zodiac = sql.NullString{
					Valid:  true,
					String: "Non-zodiac",
				}
				return basicInfo
			},
			expectedFunc: func(t *testing.T, row int64, err error) {
				var pqErr *pq.Error
				assert.ErrorAs(t, err, &pqErr)
				assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
				assert.True(t, strings.Contains(pqErr.Constraint, "zodiac"))
				assert.Zero(t, row)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedBasicInfo := tt.setupFunc()
			row, err := repo.UpdateBasicInfo(updatedBasicInfo)
			tt.expectedFunc(t, row, err)
		})
	}
}

func createBasicInfo(t *testing.T) *entity.BasicInfo {
	user := createNewAccount(t)
	basicInfo := &entity.BasicInfo{
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
