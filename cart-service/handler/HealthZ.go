package handler

import (
	"cart-service/logger"
	"cart-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthZ is a handler function for checking the health status of the application.
// It calls the PingDatabase function to check the connection to the database.
// If the database is reachable, it returns a JSON response with the status "mongodb: true" and HTTP status OK (200).
// If the database is not reachable, it returns a JSON response with the status "mongodb: false" and HTTP status Gateway Timeout (504).
func HealthZ(ctx *gin.Context) {
	logger.Logger.Debugf("Healthz - checking database connection")
	// Check the health status of the database by calling PingDatabase
	if err := repository.PingDatabase(); err != nil {
		// Handle database unavailability, return JSON response with "mongodb: false" and Gateway Timeout status
		logger.Logger.Errorf("HealthZ - %s", err.Error())

		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"redis": false,
		})
		return
	}

	logger.Logger.Debugf("Healthz - OK")

	// Database is reachable, return JSON response with "mongodb: true" and OK status
	ctx.JSON(http.StatusOK, gin.H{
		"redis": true,
	})
}
