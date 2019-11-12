package test_test

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUser(t *testing.T) *model.User {
	t.Helper()
	return &model.User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := testUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptPassword)
}

func TestUser_ComparePassword(t *testing.T) {
	u := testUser(t)
	assert.NoError(t, u.BeforeCreate())
	if !u.ComparePassword("password") || u.Password != "" {
		t.Error("fail compare same passwords")
	}
}

func TestUserInput_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				return testUser(t)
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *model.User {
				u := testUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				u := testUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				u := testUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *model.User {
				u := testUser(t)
				u.Password = "short"

				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}
