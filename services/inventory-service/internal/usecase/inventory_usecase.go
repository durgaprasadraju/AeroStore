package usecase

import (
	"context"
	"log"

	"github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/domain"
)

type inventoryUseCase struct {
	repo domain.InventoryRepository
}

func NewInventoryUseCase(r domain.InventoryRepository) domain.InventoryUseCase {
	return &inventoryUseCase{repo: r}
}

func (u *inventoryUseCase) ReserveItems(ctx context.Context, orderID string, items []domain.InventoryItem) error {
	log.Printf("Reserving items for Order: %s", orderID)

	// Business Rule: You could add logic here to check if the user is blacklisted,
	// or if the warehouse is closed.

	return u.repo.DeductStock(ctx, items)
}

func (u *inventoryUseCase) ReleaseItems(ctx context.Context, orderID string) error {
	log.Printf("Releasing items for Order: %s (Payment Failed)", orderID)
	// Implementation would involve fetching the items from the order and calling ReleaseStock
	return nil
}
