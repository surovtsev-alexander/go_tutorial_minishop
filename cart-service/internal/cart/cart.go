package cart

import "time"

// Item — товар в корзине
type Item struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Cart — корзина пользователя
type Cart struct {
	UserID    string    `json:"userId"`
	Items     []Item    `json:"items"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AddItem добавляет товар в корзину
func (c *Cart) AddItem(item Item) {
	for i := range c.Items {
		if c.Items[i].ProductID == item.ProductID {
			c.Items[i].Quantity += item.Quantity
			return
		}
	}
	c.Items = append(c.Items, item)
}

// GetTotal возвращает общую сумму
func (c *Cart) GetTotal() float64 {
	var total float64
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}