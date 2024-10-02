package cart

import (
	"fmt"

	"github.com/ardhptr21/ecomgo/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	ids := make([]int, len(items))
	for i, v := range items {
		if v.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", v.ProductID)
		}

		ids[i] = v.ProductID
	}

	return ids, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	prodMap := make(map[int]types.Product)
	for _, product := range ps {
		prodMap[product.ID] = product
	}

	if err := checkIfCartIsInStock(items, prodMap); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(items, prodMap)
	for _, item := range items {
		product := prodMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	orderID, err := h.store.CreateOrder(userID, types.CreateOrderPayload{
		Total:   totalPrice,
		Status:  "pending",
		Address: "dummy address",
	})
	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.store.CreateOrderItem(orderID, types.CreateOrderItemPayload{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     prodMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d is available in the quantity requested", item.ProductID)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var totalPrice float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		totalPrice += product.Price * float64(item.Quantity)
	}

	return totalPrice
}
