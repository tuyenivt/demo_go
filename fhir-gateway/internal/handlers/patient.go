package handlers

import (
	"fhir-gateway/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"gorm.io/gorm"
)

type PatientHandler struct {
	DB *gorm.DB
}

func (h *PatientHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	var patientRecord database.Patient
	if err := h.DB.Where("id = ?", id).First(&patientRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	patient, err := fhir.UnmarshalPatient(patientRecord.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process patient data"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var patient fhir.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient data"})
		return
	}

	data, err := patient.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process patient data"})
		return
	}

	if patient.Id == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	if err := h.DB.Create(&database.Patient{ID: *patient.Id, Data: data}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": patient.Id})
}
