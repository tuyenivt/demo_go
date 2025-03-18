package handlers

import (
	"context"
	"fhir-gateway/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/sirupsen/logrus"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
)

type PatientHandler struct {
	DB    *gorm.DB
	CACHE valkey.Client
}

func (h *PatientHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")

	// Check cache first
	ctx := context.Background()
	cacheKey := "patient:" + id
	cached, err := h.CACHE.Do(ctx, h.CACHE.B().Get().Key(cacheKey).Build()).ToString()
	if err == nil {
		patient, err := fhir.UnmarshalPatient([]byte(cached))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process cached patient data"})
			return
		}
		logrus.Info("Retrieved patient data from cache")
		c.JSON(http.StatusOK, patient)
		return
	}

	var patientRecord database.Patient
	if err := h.DB.Where("id = ?", id).First(&patientRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Cache the result
	err = h.CACHE.Do(ctx, h.CACHE.B().Set().Key(cacheKey).Value(string(patientRecord.Data)).Build()).Error()
	if err != nil {
		logrus.Errorf("Failed to cache patient: %v", err)
	} else {
		logrus.Info("Cached patient data")
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

	// Cache the result
	ctx := context.Background()
	cacheKey := "patient:" + *patient.Id
	err = h.CACHE.Do(ctx, h.CACHE.B().Set().Key(cacheKey).Value(string(data)).Build()).Error()
	if err != nil {
		logrus.Errorf("Failed to cache patient: %v", err)
	} else {
		logrus.Info("Cached patient data")
	}

	c.JSON(http.StatusCreated, gin.H{"id": patient.Id})
}
