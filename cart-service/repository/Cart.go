package repository

import (
	"cart-service/logger"
	"cart-service/model"
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// GetCart retrieves a cart from the Redis cache using the provided cart ID.
// If the cart is not found in the cache, it returns nil and no error.
// If an error occurs during the retrieval process, it returns the error.
func GetCart(ctx context.Context, cartID string) (cart *model.Cart, err error) {
	// Retrieve the raw cart data from the Redis cache
	cartRaw, err := client.Get(ctx, cartID).Result()
	if err != nil {
		// Handle the case where the cart is not found in the cache
		if errors.Is(err, redis.Nil) {
			logger.Logger.Debugf("GetCart - Cart with ID %s not found in cache", cartID)
			return nil, nil
		}
		logger.Logger.Debugf("GetCart - Cart with ID %s - Error %s", cartID, err)
		return nil, err
	}

	logger.Logger.Debugf("GetCart - Cart with ID %s found in cache", cartID)
	// Decode the raw cart data into a Cart struct
	cart, err = model.DecodeCart(cartRaw)
	if err != nil {
		logger.Logger.Debugf("GetCart - Cart with ID %s - cannot decode cart - Error %s", cartID, err)
		return nil, err
	}

	logger.Logger.Debugf("GetCart - Cart with ID %s - cart %s", cartID, cart)
	return cart, nil
}

// UpdateCart updates a cart in the Redis cache with the specified items.
// It retrieves the current cart, updates the items, encodes the updated cart, and stores it back in the cache.
// If the cart does not exist, a new cart is created with the provided cart ID.
// If an error occurs during the update process, it returns the error.
func UpdateCart(ctx context.Context, cartID string, items []string) error {
	// Retrieve the current cart from the cache
	cart, err := GetCart(ctx, cartID)
	if err != nil {
		logger.Logger.Debugf("UpdateCart - Cart with ID %s - Error %s", cartID, err)
		return err
	}

	// If the cart does not exist, create a new one with the provided cart ID
	if cart == nil {
		logger.Logger.Debugf("UpdateCart - Cart with ID %s not found - create an empty one", cartID)
		cart = &model.Cart{
			ID:    cartID,
			Items: []string{},
		}
	}

	// Update the cart items
	cart.Items = items

	// Encode the updated cart
	cartRaw, err := model.EncodeCart(cart)
	if err != nil {
		logger.Logger.Debugf("UpdateCart - Cart with ID %s - cannot encode cart - Error %s", cartID, err)
		return err
	}

	// Store the updated cart back in the Redis cache
	err = client.Set(ctx, cartID, cartRaw, 0).Err()
	if err != nil {
		logger.Logger.Debugf("UpdateCart - Cart with ID %s - Error %s", cartID, err)
		return err
	}
	logger.Logger.Debugf("UpdateCart - Cart with ID %s - update successful", cartID)

	return nil
}

// DeleteCart removes a cart from the Redis cache using the provided cart ID.
// If an error occurs during the deletion process, it returns the error.
func DeleteCart(ctx context.Context, cartID string) error {
	// Delete the cart from the Redis cache
	err := client.Del(ctx, cartID).Err()
	if err != nil {
		logger.Logger.Debugf("DeleteCart - Cart with ID %s - Error %s", cartID, err)
		return err
	}

	logger.Logger.Debugf("DeleteCart - Cart with ID %s - delete successful", cartID)
	return nil
}

// AddToCart adds an item to the cart in the Redis cache using the provided cart ID.
// It retrieves the current cart, appends the new item, encodes the updated cart, and stores it back in the cache.
// If the cart does not exist, a new cart is created with the provided cart ID.
// If an error occurs during the process, it returns the error.
func AddToCart(ctx context.Context, cartID string, itemID string) error {
	// Retrieve the current cart from the cache
	cart, err := GetCart(ctx, cartID)
	if err != nil {
		logger.Logger.Debugf("AddToCart - Cart with ID %s, Article with ID %s - Error %s", cartID, itemID, err)
		return err
	}

	// If the cart does not exist, create a new one with the provided cart ID
	if cart == nil {
		logger.Logger.Debugf("AddToCart - Cart with ID %s not found - create an empty one", cartID)
		cart = &model.Cart{
			ID:    cartID,
			Items: []string{},
		}
	}

	// Append the new item to the cart
	cart.Items = append(cart.Items, itemID)

	// Encode the updated cart
	cartRaw, err := model.EncodeCart(cart)
	if err != nil {
		logger.Logger.Debugf(
			"AddToCart - Cart with ID %s, Article with ID %s - cannot encode cart - Error %s",
			cartID,
			err)
		return err
	}

	// Store the updated cart back in the Redis cache
	err = client.Set(ctx, cartID, cartRaw, 0).Err()
	if err != nil {
		logger.Logger.Debugf(
			"AddToCart - Cart with ID %s, Article with ID %s - cannot store cart - Error %s",
			cartID,
			err)
		return err
	}

	logger.Logger.Debugf("UpdateCart - Cart with ID %s - update successful", cartID)
	return nil
}
