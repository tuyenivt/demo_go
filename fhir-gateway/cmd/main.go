package main

import (
	"fhir-gateway/internal/cache"
	"fhir-gateway/internal/config"
	"fhir-gateway/internal/database"
	"fhir-gateway/internal/handlers"
	"fhir-gateway/internal/middleware"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Optimize for all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load config
	config.Load()

	// Connect to database
	db := database.Connect()

	// Connect to cache server
	valkey := cache.Connect()
	defer valkey.Close()

	// gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.ReleaseMode)

	// Setup routes
	router := gin.Default()

	// Setup monitoring
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Setup middleware
	router.Use(middleware.AuthMiddleware(db))
	router.Use(middleware.LoggingMiddleware())

	// Setup handlers
	patientHandler := handlers.PatientHandler{DB: db, CACHE: valkey}
	fhirGroup := router.Group("/fhir/r4")
	{
		fhirGroup.GET("/Patient/:id", patientHandler.GetPatient)
		fhirGroup.POST("/Patient", patientHandler.CreatePatient)
	}

	// Start server
	router.Run(":8080")
}
