package cart

type CartEventType string

const (
	CartItemAdded   CartEventType = "cart_item_added"
	CartItemRemoved CartEventType = "cart_item_removed"
)
