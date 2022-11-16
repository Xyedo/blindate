package repository

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
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

func Test_CreateProfilePicture(t *testing.T) {
	repo := NewUser(testQuery)
	t.Run("create valid pp", func(t *testing.T) {
		usr := createNewAccount(t)
		id, err := repo.CreateProfilePicture(usr.ID, util.RandomUUID()+".png", true)
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})
	t.Run("create valid pp but not false selected", func(t *testing.T) {
		usr := createNewAccount(t)
		id, err := repo.CreateProfilePicture(usr.ID, util.RandomUUID()+".png", false)
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})
	t.Run("create multiple profpic", func(t *testing.T) {
		usr := createNewAccount(t)
		for i := 0; i < 3; i++ {
			id, err := repo.CreateProfilePicture(usr.ID, util.RandomUUID()+".png", false)
			require.NoError(t, err)
			assert.NotEmpty(t, id)
		}
	})

}

func Test_SelectProfilePic(t *testing.T) {
	repo := NewUser(testQuery)
	setupFunc := func(t *testing.T) string {
		usr := createNewAccount(t)
		for i := 0; i < 3; i++ {
			id, err := repo.CreateProfilePicture(usr.ID, util.RandomUUID()+".png", false)
			require.NoError(t, err)
			assert.NotEmpty(t, id)
		}
		return usr.ID
	}

	t.Run("valid Select Profile Pic", func(t *testing.T) {
		userId := setupFunc(t)
		profpics, err := repo.SelectProfilePicture(userId, nil)
		require.NoError(t, err)
		assert.Equal(t, len(profpics), 3)

	})
	t.Run("valid Select with Params > Return 0", func(t *testing.T) {
		userId := setupFunc(t)
		selected := true
		profpics, err := repo.SelectProfilePicture(userId, &entity.ProfilePicQuery{
			Selected: &selected,
		})
		require.NoError(t, err)
		assert.Equal(t, len(profpics), 0)

	})
	t.Run("valid Select with Params > Return 1", func(t *testing.T) {
		userId := setupFunc(t)
		id, err := repo.CreateProfilePicture(userId, util.RandomUUID()+".png", true)
		require.NoError(t, err)
		require.NotEmpty(t, id)
		selected := true
		profpics, err := repo.SelectProfilePicture(userId, &entity.ProfilePicQuery{
			Selected: &selected,
		})
		require.NoError(t, err)
		assert.Equal(t, len(profpics), 1)
	})
}

func Test_ProfilePicSelectedToFalse(t *testing.T) {
	repo := NewUser(testQuery)
	usr := createNewAccount(t)
	for i := 0; i < 3; i++ {
		id, err := repo.CreateProfilePicture(usr.ID, util.RandomUUID()+".png", true)
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	}
	selected := true
	profpics, err := repo.SelectProfilePicture(usr.ID, &entity.ProfilePicQuery{
		Selected: &selected,
	})
	require.NoError(t, err)
	require.Equal(t, len(profpics), 3)
	require.Equal(t, profpics[0].UserId, usr.ID)

	row, err := repo.ProfilePicSelectedToFalse(usr.ID)
	require.NoError(t, err)
	require.Equal(t, row, int64(3))

	actualProfPic, err := repo.SelectProfilePicture(usr.ID, &entity.ProfilePicQuery{
		Selected: &selected,
	})
	require.NoError(t, err)
	require.Equal(t, len(actualProfPic), 0)

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
