package order

var (
	orderStatuses = []string{
		"Order placed",
		"Processing",
		"Baking",
		"Quality Check",
		"Ready",
	}
	pizzaTypes = []string{
		"Margherita",
		"Pepperoni",
		"Vegetarian",
		"Hawaiian",
		"Bbq chicken",
		"Meat lover",
		"Buffalo chicken",
		"Supreme",
		"Truffle Mushroom",
		"Four Cheese",
	}
	pizzaSizes = []string{
		"Small",
		"Medium",
		"Large",
		"X-Large",
	}
)

func GetOrderStatuses() []string {
	result := make([]string, len(orderStatuses))
	copy(result, orderStatuses)
	return result
}

func GetPizzaTypes() []string {
	result := make([]string, len(pizzaTypes))
	copy(result, pizzaTypes)
	return result
}

func GetPizzaSizes() []string {
	result := make([]string, len(pizzaSizes))
	copy(result, pizzaSizes)
	return result
}
