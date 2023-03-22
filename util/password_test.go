package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NotEmpty(t, hashedPassword)
	require.NoError(t, err)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(5)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
