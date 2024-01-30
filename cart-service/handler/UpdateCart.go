package handler

import (
	"cart-service/logger"
	"cart-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateCart is a handler function for updating a cart with new items.
// It retrieves the cart ID from the request parameters and the items to update from the JSON request body.
// If the request body is invalid or does not match the expected format, it returns a Bad Request response.
// It then calls the UpdateCart function from the repository to update the cart with the specified items.
// If an error occurs during the update process, it returns an Internal Server Error response.
// Otherwise, it returns an Accepted response indicating the successful update of the cart.
func UpdateCart(ctx *gin.Context) {
	// Retrieve the cart ID from the request parameters
	cartID := ctx.Param("cartId")
	logger.Logger.Infof("UpdateCart - Update cart with ID %s in database", cartID)

	// Parse the JSON request body into an UpdateCartRequest struct
	var updateCartRequest UpdateCartRequest
	if err := ctx.ShouldBindJSON(&updateCartRequest); err != nil {
		logger.Logger.Warnf("UpdateCart with cartId %s - error - %s", cartID, err.Error())

		// Handle JSON binding error and return Bad Request response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the UpdateCart function from the repository to update the cart with the specified items
	err := repository.UpdateCart(ctx, cartID, updateCartRequest.Items)
	if err != nil {
		logger.Logger.Errorf("UpdateCart with cartId %s - error - %s", cartID, err.Error())

		// Handle repository error and return Internal Server Error response
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Infof("UpdateCart with cartId %s - successful", cartID)
	// Return Accepted response indicating the successful update of the cart
	ctx.Status(http.StatusAccepted)
}

// UpdateCartRequest represents the structure of a JSON request
// expected when updating a cart. It includes a list of items to be updated.
// The struct tags define the corresponding JSON key names for serialization and deserialization.
type UpdateCartRequest struct {
	// Items is a list of strings representing the items to be updated in the cart.
	// It is expected to be provided in the JSON request as "items".
	Items []string `json:"items"`
}
