package model

import "encoding/json"

// Cart represents the structure of a shopping cart.
// It includes an ID and a list of items in the cart.
// The struct tags define the corresponding JSON key names for serialization and deserialization.
type Cart struct {
	// ID is the unique identifier of the cart.
	// It is expected to be provided in the JSON representation as "id".
	ID string `json:"id"`

	// Items is a list of strings representing the items in the cart.
	// It is expected to be provided in the JSON representation as "items".
	Items []string `json:"items"`
}

// DecodeCart parses a JSON representation of a cart from a string.
// It returns a pointer to a Cart struct and an error, if any.
func DecodeCart(rawCart string) (*Cart, error) {
	cart := &Cart{}
	err := json.Unmarshal([]byte(rawCart), cart)
	return cart, err
}

// EncodeCart converts a Cart struct into its JSON representation as a string.
// It returns the JSON string and an error, if any.
func EncodeCart(cart *Cart) (string, error) {
	rawCart, err := json.Marshal(cart)
	return string(rawCart), err
}
