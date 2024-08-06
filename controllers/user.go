package controllers

import (
	"golang-api/models"
	"golang-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
//
//	@Summary		get all users
//	@Description	get all users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.SanitizedUser
//	@Failure		500	{object}	utils.HttpError
//	@Router			/users/ [get]
func GetUsers(c *gin.Context) {
	var sanitized []models.SanitizedUser

	query, _ := c.MustGet("query").(services.QueryFilter)

	users, err := models.FindAllUsers(query)
	if err != nil {
		models.ErrorLogf([]string{"user", "GetUsers"}, "%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, user := range users {
		sanitized = append(sanitized, user.Sanitize(true))
	}

	models.PrintLogf([]string{"user", "GetUsers"}, "Fetched %d users", len(sanitized))
	c.JSON(http.StatusOK, sanitized)
}

// GetUser godoc
//
//	@Summary		get users by id
//	@Description	get users by id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.SanitizedUser
//	@Failure		404	{object}	utils.HttpError
//	@Failure		500	{object}	utils.HttpError
//	@Router			/users/{id} [get]
func GetUser(c *gin.Context) {
	user, _ := c.MustGet("user").(*models.User)

	sanitized := user.Sanitize(true)

	models.PrintLogf([]string{"user", "GetUser"}, "Fetched user %d", user.ID)
	c.JSON(http.StatusOK, sanitized)
}

// GetUserMe godoc
//
//	@Summary		get connected user
//	@Description	get connected user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.SanitizedUser
//	@Router			/users/me [get]
func GetUserMe(c *gin.Context) {
	connectedUser, _ := c.MustGet("connectedUser").(models.User)
	c.JSON(http.StatusOK, connectedUser)
}

// UpdateUser godoc
//
//	@Summary		update user
//	@Description	update user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UpdateUserDto	true	"User"
//	@Param			id		path		int						true	"User ID"
//	@Success		200		{object}	models.SanitizedUser
//	@Failure		400		{object}	utils.HttpError
//	@Failure		401		{object}	utils.HttpError
//	@Failure		404		{object}	utils.HttpError
//	@Failure		500	{object}	utils.HttpError
//	@Router			/users/{id} [patch]
func UpdateUser(c *gin.Context) {
	user, _ := c.MustGet("user").(*models.User)
	body, _ := c.MustGet("body").(models.UpdateUserDto)

	if count, err := models.CountUsersByEmail(body.Email); err != nil || count > 0 {
		models.ErrorLogf([]string{"user", "UpdateUser"}, "Email %s already exists", body.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	if count, err := models.CountUsersByUsername(body.Username); err != nil || count > 0 {
		models.ErrorLogf([]string{"user", "UpdateUser"}, "Username %s already exists", body.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	if body.Firstname != "" {
		user.Firstname = body.Firstname
	}
	if body.Lastname != "" {
		user.Lastname = body.Lastname
	}
	if body.Username != "" {
		user.Username = body.Username
	}
	if body.Email != "" {
		user.Email = body.Email
	}

	if err := user.Save(); err != nil {
		models.ErrorLogf([]string{"user", "UpdateUser"}, "%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sanitized := user.Sanitize(true)
	models.PrintLogf([]string{"user", "UpdateUser"}, "User %d updated", user.ID)
	c.JSON(http.StatusOK, sanitized)
}

// UpdateUserImage godoc
//
//	@Summary		update user image
//	@Description	update user image
//	@Tags			user
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user		path	int		true	"User ID"
//	@Param			upload[]		formData	file	true	"Image"
//	@Success		200			{object}	models.SanitizedUser
//	@Failure		400			{object}	utils.HttpError
//	@Failure		401			{object}	utils.HttpError
//	@Failure		404			{object}	utils.HttpError
//	@Failure		500			{object}	utils.HttpError
//	@Router			/users/{user}/image [post]
func UploadUserImage(c *gin.Context) {
	user, _ := c.MustGet("user").(*models.User)
	files, _ := c.MustGet("files").([]string)

	user.ProfileImageUrl = files[0]

	if err := user.Save(); err != nil {
		models.ErrorLogf([]string{"user", "UploadUserImage"}, "%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	models.PrintLogf([]string{"user", "UploadUserImage"}, "User %d image updated", user.ID)
	c.JSON(http.StatusOK, user.Sanitize(false))
}

// DeleteUser godoc
//
//	@Summary		delete user
//	@Description	delete user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	utils.HttpError
//	@Failure		401	{object}	utils.HttpError
//	@Failure		404	{object}	utils.HttpError
//	@Failure		500	{object}	utils.HttpError
//	@Router			/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	user, _ := c.MustGet("user").(*models.User)

	if err := user.Delete(); err != nil {
		models.ErrorLogf([]string{"user", "DeleteUser"}, "%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	models.PrintLogf([]string{"user", "DeleteUser"}, "User %d deleted", user.ID)
	c.JSON(http.StatusNoContent, nil)
}

func GetUsersStats(c *gin.Context) {
	ws := services.GetWebsocket()
	stats := models.CountStatsUsers(ws)
	c.JSON(http.StatusOK, stats)
}
