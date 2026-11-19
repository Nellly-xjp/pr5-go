package handlers

import (
	"errors"
	"net/http"
	"pr-5-go/database"
	"pr-5-go/dto"
	"pr-5-go/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func formatValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			errs[fe.Field()] = "Порушено правило валідації: " + fe.Tag() + " " + fe.Param()
		}
		return errs
	}
	errs["general"] = err.Error()
	return errs
}


func CreateIncident(c *gin.Context) {
	var input dto.CreateIncidentDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	incident := models.Incident{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      "open",
	}

	if err := database.DB.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка збереження в БД"})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func GetAllIncidents(c *gin.Context) {
	var incidents []models.Incident
	if err := database.DB.Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка читання з БД"})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func GetIncidentByID(c *gin.Context) {
	id := c.Param("id")
	var incident models.Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Інцидент не знайдено"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка бази даних"})
		return
	}

	c.JSON(http.StatusOK, incident)
}


func UpdateIncidentStatus(c *gin.Context) {
	id := c.Param("id")
	var incident models.Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Інцидент не знайдено"})
		return
	}

	var input dto.UpdateIncidentStatusDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	incident.Status = input.Status
	database.DB.Save(&incident)

	c.JSON(http.StatusOK, incident)
}

func DeleteIncident(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Incident{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка видалення"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Інцидент не знайдено"})
		return
	}

	c.Status(http.StatusNoContent)
}