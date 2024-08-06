package middlewares

import (
	"golang-api/models"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

func Validate[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body T

		if err := c.BindJSON(&body); err != nil {
			models.ErrorLogf([]string{"middleware", "Validate"}, "%s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		validate := validator.New()
		if err := validate.Struct(body); err != nil {
			models.ErrorLogf([]string{"middleware", "Validate"}, "%s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("body", body)
	}
}

func Get[T models.Model](name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var model T
		var err error

		id, err := strconv.Atoi(c.Param(name))
		if err != nil {
			models.ErrorLogf([]string{"middleware", "Get"}, "%s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			c.Abort()
			return
		}

		instance := reflect.New(reflect.TypeOf(model).Elem()).Interface().(models.Model)
		if instance.FindOneById(id) != nil {
			models.ErrorLogf([]string{"middleware", "Get"}, "Instance not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			c.Abort()
			return
		}

		c.Set(name, instance)
		c.Next()
	}
}
