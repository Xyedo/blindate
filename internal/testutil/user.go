package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/user"
	userDTOs "github.com/xyedo/blindate/pkg/user/dtos"
	userEntity "github.com/xyedo/blindate/pkg/user/entities"
	"golang.org/x/crypto/bcrypt"
)

func CreateNewAccount(t *testing.T, userRepo user.Repository) userEntity.User {
	hashed, err := bcrypt.GenerateFromPassword([]byte(RandomString(12)), 12)
	assert.NoError(t, err)
	user := userDTOs.RegisterUser{
		FullName: "Hafid Mahdi",
		Email:    RandomEmail(23),
		Password: string(hashed),
		Dob:      RandDOB(1980, 2000),
	}
	userId, err := userRepo.InsertUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, userId)
	return userEntity.User{
		ID:       userId,
		FullName: user.FullName,
		Email:    user.Email,
		Password: user.Password,
		Dob:      user.Dob,
	}
}
