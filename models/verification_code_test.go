package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVerificationCode_Save(t *testing.T) {
	assert := assert.New(t)

	if err := ConnectDB(true); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := CloseDB(); err != nil {
			t.Fatal("Failed to close database:", err)
		}
	}()

	vc := &VerificationCode{
		UserID:    1,
		Email:     "test@example.com",
		Code:      "12345",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := vc.Save(); err != nil {
		t.Fatalf("failed to save verification code: %v", err)
	}

	var savedVC VerificationCode
	err := DB.Where("user_id = ? AND email = ?", vc.UserID, vc.Email).First(&savedVC).Error
	assert.NoError(err)
	assert.Equal(vc.Code, savedVC.Code)
	assert.Equal(vc.ExpiresAt.Unix(), savedVC.ExpiresAt.Unix())
}

func TestIsExpired(t *testing.T) {
	assert := assert.New(t)

	vc := &VerificationCode{
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	assert.True(vc.IsExpired(), "expected verification code to be expired")
}

func TestDeleteExpiredCodes(t *testing.T) {
	assert := assert.New(t)

	if err := ConnectDB(true); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := CloseDB(); err != nil {
			t.Fatal("Failed to close database:", err)
		}
	}()

	vc1 := &VerificationCode{
		UserID:    1,
		Email:     "expired1@example.com",
		Code:      "11111",
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	err := vc1.Save()
	assert.NoError(err)

	vc2 := &VerificationCode{
		UserID:    2,
		Email:     "expired2@example.com",
		Code:      "22222",
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	err = vc2.Save()
	assert.NoError(err)

	vc3 := &VerificationCode{
		UserID:    3,
		Email:     "valid@example.com",
		Code:      "33333",
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	err = vc3.Save()
	assert.NoError(err)

	err = DeleteExpiredCodes()
	assert.NoError(err)

	var count int64
	err = DB.Model(&VerificationCode{}).Where("expires_at < ?", time.Now()).Count(&count).Error
	assert.NoError(err)
	assert.Equal(int64(0), count)

	err = DB.Model(&VerificationCode{}).Where("id = ?", vc3.ID).Count(&count).Error
	assert.NoError(err)
	assert.Equal(int64(1), count)
}
