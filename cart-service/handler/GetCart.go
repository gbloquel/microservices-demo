package handler

import (
	"cart-service/logger"
	"cart-service/model"
	"cart-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCart is a handler function for retrieving a cart.
// It retrieves the cart ID from the request parameters and attempts to retrieve the corresponding cart from the repository.
// If an error occurs during the retrieval process, it returns an Internal Server Error response.
// If the cart is not found in the repository, a new cart with the specified ID is created.
// The retrieved or newly created cart is then returned in a JSON response.
func GetCart(ctx *gin.Context) {
	// Retrieve the cart ID from the request parameters
	cartID := ctx.Param("cartId")

	logger.Logger.Infof("GetCart - Retrieve cart with ID %s from database", cartID)

	// Attempt to retrieve the cart using the GetCart function from the repository
	cart, err := repository.GetCart(ctx, cartID)
	if err != nil {
		logger.Logger.Errorf("GetCart - Cart with ID %s - error - %s", cartID, err.Error())
		// Handle repository error and return Internal Server Error response
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// If the cart is not found in the repository, create a new cart with the specified ID
	if cart == nil {
		logger.Logger.Debugf("GetCart - Cart with ID %s - not found", cartID)
		logger.Logger.Debugf("GetCart - Cart with ID %s - create an empty one", cartID)
		cart = &model.Cart{
			ID: cartID,
		}
	}

	logger.Logger.Infof("GetCart - Cart with ID %s - successful", cartID)
	// Return a JSON response with the retrieved or newly created cart
	ctx.JSON(http.StatusOK, gin.H{
		"cart": cart,
	})
}
