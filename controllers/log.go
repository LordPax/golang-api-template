package controllers

import (
	"golang-api/models"
	"golang-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLogs(c *gin.Context) {
	query, _ := c.MustGet("query").(services.QueryFilter)

	logs, err := models.FindAllLogs(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func GetLog(c *gin.Context) {
	log, _ := c.MustGet("log").(*models.Log)
	c.JSON(http.StatusOK, log)
}
