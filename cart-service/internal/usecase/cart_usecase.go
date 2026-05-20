package usecase

import (
	"time"

	"cart-service/internal/cart"
)

type CartRepository interface {
	Save(cart *cart.Cart) error
	FindByID(userID string) (*cart.Cart, error)
}

type CartUsecase struct {
	repo CartRepository
}

func NewCartUsecase(repo CartRepository) *CartUsecase {
	return &CartUsecase{repo: repo}
}

func (uc *CartUsecase) AddItem(userID, productID string, quantity int, price float64) error {
	cart, err := uc.repo.FindByID(userID)
	if err != nil {
		cart = &cart.Cart{UserID: userID, UpdatedAt: time.Now()}
	}

	cart.AddItem(Item{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	})
	cart.UpdatedAt = time.Now()

	return uc.repo.Save(cart)
}

func (uc *CartUsecase) GetCart(userID string) (*cart.Cart, error) {
	return uc.repo.FindByID(userID)
}