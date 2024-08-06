package models

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupToken() {
	setup()
	os.Setenv("JWT_SECRET_KEY", "secret")
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func createSampleTokenData() (User, Token) {
	user := User{
		ID:        1,
		Firstname: fake.Person().FirstName(),
		Lastname:  fake.Person().LastName(),
		Username:  fake.Person().Name(),
		Email:     fake.Internet().Email(),
		Password:  fake.Internet().Password(),
		Roles:     []string{"user"},
	}
	token := Token{
		ID:        1,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	DB.Create(&user)
	DB.Create(&token)

	return user, token
}

func TestNewToken(t *testing.T) {
	setupToken()
	defer close()
	user, _ := createSampleTokenData()

	token, err := NewToken("access", user.ID)

	assert.Nil(t, err)
	assert.Equal(t, user.ID, token.UserID)
}

func TestFindTokensByUserID(t *testing.T) {
	setupToken()
	defer close()
	user, _ := createSampleTokenData()

	tokens, err := FindTokensByUserID(user.ID)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(tokens))
}

func TestDeleteTokensByUserID(t *testing.T) {
	setupToken()
	defer close()
	user, _ := createSampleTokenData()

	err := DeleteTokensByUserID(user.ID)

	assert.Nil(t, err)
	var count int64
	DB.Model(&Token{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestParseJWTToken(t *testing.T) {
	setupToken()
	defer close()
	user, _ := createSampleTokenData()

	token, _ := NewToken("access", user.ID)
	_ = token.GenerateAccessToken()

	claims, err := ParseJWTToken(token.AccessToken)

	assert.Nil(t, err)
	assert.Equal(t, user.ID, claims.UserID)
}

func TestGenerateAccessToken(t *testing.T) {
	setupToken()
	defer close()
	_, token := createSampleTokenData()

	err := token.GenerateAccessToken()

	assert.Nil(t, err)
	assert.NotEmpty(t, token.AccessToken)
}

func TestGenerateRefreshToken(t *testing.T) {
	setupToken()
	defer close()
	_, token := createSampleTokenData()

	err := token.GenerateRefreshToken()

	assert.Nil(t, err)
	assert.NotEmpty(t, token.RefreshToken)
}

func TestGenerateTokens(t *testing.T) {
	setupToken()
	defer close()
	_, token := createSampleTokenData()

	err := token.GenerateTokens()

	assert.Nil(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
}

func TestFindOneTokenByID(t *testing.T) {
	setupToken()
	defer close()
	_, token := createSampleTokenData()

	var foundToken Token
	err := foundToken.FindOneById(token.ID)

	assert.Nil(t, err)
	assert.Equal(t, token.ID, foundToken.ID)
}

func TestFindOneToken(t *testing.T) {
	setupToken()
	defer close()
	_, token := createSampleTokenData()

	var foundToken Token
	err := foundToken.FindOne("access_token", token.AccessToken)

	assert.Nil(t, err)
	assert.Equal(t, token.AccessToken, foundToken.AccessToken)
}

func TestSaveToken(t *testing.T) {
	setupToken()
	defer close()
	_, token := createSampleTokenData()
	token.AccessToken = "new_access_token"

	err := token.Save()

	assert.Nil(t, err)
	assert.Equal(t, "new_access_token", token.AccessToken)
}

func TestDeleteToken(t *testing.T) {
	setupToken()
	defer close()
	_, token := createSampleTokenData()

	err := token.Delete()

	assert.Nil(t, err)
	var count int64
	DB.Model(&Token{}).Count(&count)
	assert.Equal(t, int64(0), count)
}
