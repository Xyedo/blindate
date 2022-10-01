package repository

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

func Test_InsertUser(t *testing.T) {
	t.Run("Valid NewAcc", func(t *testing.T) {
		createNewAccount(t)
	})
	t.Run("Duplicate Email", func(t *testing.T) {
		user := createNewAccount(t)
		repo := NewUser(testQuery)
		err := repo.InsertUser(user)
		var pqErr *pq.Error
		assert.Error(t, err)
		assert.ErrorAs(t, err, &pqErr)
		assert.Equal(t, pq.ErrorCode("23505"), pqErr.Code)
		assert.True(t, strings.Contains(pqErr.Constraint, "users_email"))
	})

}
func Test_UpdateUser(t *testing.T) {
	repo := NewUser(testQuery)

	t.Run("Not Found UserId", func(t *testing.T) {
		user := createNewAccount(t)
		user.ID = "e590666c-3ea8-4fda-958c-c2dc6c2599b5"
		user.FullName = util.RandomString(12)
		user.Email = util.RandomEmail(12)
		user.Active = true
		err := repo.UpdateUser(user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
	t.Run("Success Updating", func(t *testing.T) {
		user := createNewAccount(t)
		user.FullName = util.RandomString(12)
		user.Email = util.RandomEmail(12)
		user.Active = true
		err := repo.UpdateUser(user)
		assert.NoError(t, err)
	})
}

func Test_GetUserById(t *testing.T) {
	repo := NewUser(testQuery)
	t.Run("Valid UserId", func(t *testing.T) {
		expectedUser := createNewAccount(t)
		user, err := repo.GetUserById(expectedUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.FullName, user.FullName)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.HashedPassword, user.HashedPassword)
		assert.Equal(t, expectedUser.Dob.Year(), user.Dob.Year())
		assert.Equal(t, expectedUser.Dob.Month(), user.Dob.Month())
		assert.Equal(t, expectedUser.Dob.Day(), user.Dob.Day())
	})
	t.Run("Invalid Id", func(t *testing.T) {
		_, err := repo.GetUserById("e590666c-3ea8-4fda-958c-c2dc6c2599b5")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

}
func Test_GetUserByEmail(t *testing.T) {
	repo := NewUser(testQuery)
	t.Run("Valid UserId", func(t *testing.T) {
		expectedUser := createNewAccount(t)
		user, err := repo.GetUserByEmail(expectedUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.HashedPassword, user.HashedPassword)
	})
	t.Run("Invalid Id", func(t *testing.T) {
		_, err := repo.GetUserByEmail(util.RandomEmail(12))
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

}
func createNewAccount(t *testing.T) *domain.User {
	repo := NewUser(testQuery)
	hashed, err := bcrypt.GenerateFromPassword([]byte(util.RandomString(12)), 12)
	assert.NoError(t, err)
	user := &domain.User{
		FullName:       "Hafid Mahdi",
		Email:          util.RandomEmail(23),
		HashedPassword: string(hashed),
		Dob:            util.RandDOB(1980, 2000),
	}
	repo.InsertUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	return user
}
