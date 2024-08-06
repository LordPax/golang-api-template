package controllers

import (
	"fmt"
	"golang-api/models"
	"golang-api/services"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Login godoc
//
//	@Summary		Login user
//	@Description	Login user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body		models.LoginUserDto	true	"Login user"
//	@Success		200		{object}	models.LoginSuccessResponse
//	@Failure		400		{object}	utils.HttpError
//	@Router			/auth/login [post]
func Login(c *gin.Context) {
	var user models.User
	body, _ := c.MustGet("body").(models.LoginUserDto)

	if err := user.FindOne("email", body.Email); err != nil {
		models.ErrorLogf([]string{"auth", "login"}, "Failed to find user with email %s", body.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.Verified {
		models.ErrorLogf([]string{"auth", "login"}, "Account not verified")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not verified"})
		return
	}

	if !user.ComparePassword(body.Password) {
		models.ErrorLogf([]string{"auth", "login"}, "Invalid password for user %s", user.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := models.DeleteTokensByUserID(user.ID); err != nil {
		models.ErrorLogf([]string{"auth", "login"}, "Failed to delete tokens for user %s", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tokens"})
		return
	}

	token, err := models.NewToken("access_token", user.ID)
	if err != nil {
		models.ErrorLogf([]string{"auth", "login"}, "Failed to create token for user %s", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	if token.GenerateTokens() != nil {
		models.ErrorLogf([]string{"auth", "login"}, "Failed to generate tokens for user %s", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	if token.Save() != nil {
		models.ErrorLogf([]string{"auth", "login"}, "Failed to save token for user %s", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	setAuthCookies(c, token)
	models.PrintLogf([]string{"auth", "login"}, "User %s logged in", user.Email)
	c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}

// Refresh godoc
//
//	@Summary		Refresh token
//	@Description	Refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	models.RefreshSuccessResponse
//	@Failure		400		{object}	utils.HttpError
//	@Router			/auth/refresh [post]
func Refresh(c *gin.Context) {
	var token models.Token

	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token not found"})
		return
	}

	if token.FindOne("refresh_token", refreshTokenString) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token not found"})
		return
	}

	if token.GenerateAccessToken() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	if token.Save() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate or save token"})
		return
	}

	setAuthCookies(c, &token)
	c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken})
}

// Register godoc
//
//	@Summary		Register user
//	@Description	Register user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			register	body	models.CreateUserDto	true	"Register user"
//	@Success		201
//	@Failure		400	{object}	utils.HttpError
//	@Router			/auth/register [post]
func Register(c *gin.Context) {
	body, _ := c.MustGet("body").(models.CreateUserDto)

	if isUserExists(body) {
		models.ErrorLogf([]string{"auth", "register"}, "Email or Username already exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Username already exists"})
		return
	}

	defaultAvatarUrl := fmt.Sprintf("https://api.dicebear.com/9.x/initials/svg?seed=%s", body.Firstname+body.Lastname)

	user := models.User{
		Firstname:       body.Firstname,
		Lastname:        body.Lastname,
		Username:        body.Username,
		Email:           body.Email,
		ProfileImageUrl: defaultAvatarUrl,
		Roles:           []string{models.ROLE_USER},
	}

	if err := user.HashPassword(body.Password); err != nil {
		models.ErrorLogf([]string{"auth", "register"}, "Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := user.Save(); err != nil {
		models.ErrorLogf([]string{"auth", "register"}, "Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	models.PrintLogf([]string{"auth", "register"}, "User %s registered", user.Email)
	sendWelcomeAndVerificationEmails(user, c)
}

// Verify godoc
//
//	@Summary		Verify account
//	@Description	Verify user account with code
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			verification	body		models.VerifyUserDto	true	"Verify user"
//	@Success		200
//	@Failure		400	{object}	utils.HttpError
//	@Router			/auth/verify [post]
func Verify(c *gin.Context) {
	var user models.User
	var verificationCode models.VerificationCode

	body, _ := c.MustGet("body").(models.VerifyUserDto)

	if verificationCode.FindOneByCodeAndEmail(body.Code, body.Email) != nil {
		models.ErrorLogf([]string{"auth", "verify"}, "Verification code not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Verification code not found"})
		return
	}

	if verificationCode.IsExpired() {
		models.ErrorLogf([]string{"auth", "verify"}, "Invalid or expired verification code")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification code"})
		return
	}

	if user.FindOne("email", verificationCode.Email) != nil {
		models.ErrorLogf([]string{"auth", "verify"}, "User %s not found", verificationCode.Email)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Verified = true

	if user.Save() != nil {
		models.ErrorLogf([]string{"auth", "verify"}, "Failed to verify account")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify account"})
		return
	}

	if verificationCode.Delete() != nil {
		models.ErrorLogf([]string{"auth", "verify"}, "Failed to delete verification code")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete verification code"})
		return
	}

	models.PrintLogf([]string{"auth", "verify"}, "User %s verified", user.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}

// RequestVerification godoc
//
//	@Summary		Request verification code
//	@Description	Request verification code
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.RequestCodeDto	true	"Request verification code"
//	@Success		200
//	@Failure		400	{object}	utils.HttpError
//	@Router			/auth/request-verify [post]
func RequestVerification(c *gin.Context) {
	body, _ := c.MustGet("body").(models.RequestCodeDto)

	var user models.User

	if user.FindOne("email", body.Email) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	verificationCode := models.VerificationCode{
		UserID:    user.ID,
		Email:     user.Email,
		Code:      strconv.Itoa(services.GenerateVerificationCode()),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if verificationCode.Save() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save verification code"})
		return
	}

	if services.SendVerificationEmail(user.Email, verificationCode.Code) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
}

// RequestPasswordReset godoc
//
//	@Summary		Request password reset
//	@Description	Request password reset
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.RequestCodeDto	true	"Request password reset"
//	@Success		200
//	@Failure		400	{object}	utils.HttpError
//	@Router			/auth/request-password-reset [post]
func RequestPasswordReset(c *gin.Context) {
	body, _ := c.MustGet("body").(models.RequestCodeDto)

	var user models.User

	if user.FindOne("email", body.Email) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	resetCode := models.VerificationCode{
		UserID:    user.ID,
		Email:     user.Email,
		Code:      strconv.Itoa(services.GenerateVerificationCode()),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if resetCode.Save() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save password reset code"})
		return
	}

	if services.SendEmailWithCode(user.Email, services.PasswordResetEmail, resetCode.Code) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset code sent successfully"})
}

// ResetPassword godoc
//
//	@Summary		Reset password
//	@Description	Reset user password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			reset	body		models.ResetPasswordDto	true	"Reset password"
//	@Success		200
//	@Failure		400	{object}	utils.HttpError
//	@Router			/auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	var verificationCode models.VerificationCode
	var user models.User

	body, _ := c.MustGet("body").(models.ResetPasswordDto)

	if verificationCode.FindOneByCodeAndEmail(body.Code, body.Email) != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Verification code not found"})
		return
	}

	if verificationCode.IsExpired() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification code"})
		return
	}

	if user.FindOneById(verificationCode.UserID) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if user.HashPassword(body.NewPassword) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if user.Save() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	if verificationCode.Delete() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// Logout godoc
//
//	@Summary		Logout user
//	@Description	Logout user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	utils.HttpError
//	@Router			/auth/logout [post]
func Logout(c *gin.Context) {
	token, _ := c.MustGet("token").(models.Token)

	if token.Delete() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
		return
	}

	clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Helper functions

// setAuthCookies sets authentication cookies for the user.
func setAuthCookies(c *gin.Context, token *models.Token) {
	cookieSecure, _ := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	c.SetCookie("access_token", token.AccessToken, 3600, "/", "", cookieSecure, true)
	c.SetCookie("refresh_token", token.RefreshToken, 3600*24*30, "/", "", cookieSecure, true)
}

// clearAuthCookies clears authentication cookies for the user.
func clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}

// isUserExists checks if a user with the given email or username already exists.
func isUserExists(body models.CreateUserDto) bool {
	emailExists, _ := models.CountUsersByEmail(body.Email)
	usernameExists, _ := models.CountUsersByUsername(body.Username)
	return emailExists > 0 || usernameExists > 0
}

// sendWelcomeAndVerificationEmails sends welcome and verification emails to the user.
func sendWelcomeAndVerificationEmails(user models.User, c *gin.Context) {
	if services.SendEmail(user.Email, services.WelcomeEmail) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	verificationCode := models.VerificationCode{
		UserID:    user.ID,
		Email:     user.Email,
		Code:      strconv.Itoa(services.GenerateVerificationCode()),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if verificationCode.Save() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save verification code"})
		return
	}

	if services.SendVerificationEmail(user.Email, verificationCode.Code) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
