package handler

import (
	"cart-service/logger"
	"cart-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteCart is a handler function for deleting a cart.
// It retrieves the cart ID from the request parameters and attempts to delete the corresponding cart.
// If an error occurs during the deletion process, it returns an Internal Server Error response.
// Otherwise, it returns an Accepted response indicating the successful deletion of the cart.
func DeleteCart(ctx *gin.Context) {
	// Retrieve the cart ID from the request parameters
	cartID := ctx.Param("cartId")
	logger.Logger.Infof("DeleteCart - Delete cart with ID %s from database", cartID)

	// Attempt to delete the cart using the DeleteCart function from the repository
	err := repository.DeleteCart(ctx, cartID)
	if err != nil {
		// Handle repository error and return Internal Server Error response
		logger.Logger.Errorf("DeleteCart with cartId %s - error - %s", cartID, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return Accepted response indicating the successful deletion of the cart
	logger.Logger.Infof("DeleteCart with cartId %s - successful", cartID)
	ctx.Status(http.StatusAccepted)
}
