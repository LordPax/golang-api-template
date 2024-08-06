package models

import (
	"time"
)

type VerificationCode struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int       `gorm:"not null"`
	Email     string    `gorm:"not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VerifyUserDto struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type RequestCodeDto struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordDto struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (vc *VerificationCode) Save() error {
	return DB.Save(vc).Error
}

func (vc *VerificationCode) IsExpired() bool {
	return time.Now().After(vc.ExpiresAt)
}

func (vc *VerificationCode) FindOneByCodeAndEmail(code, email string) error {
	return DB.Where("code = ? AND email = ?", code, email).First(vc).Error
}

func (vc *VerificationCode) Delete() error {
	return DB.Delete(vc).Error
}

func DeleteExpiredCodes() error {
	return DB.Where("expires_at < ?", time.Now()).Delete(&VerificationCode{}).Error
}
