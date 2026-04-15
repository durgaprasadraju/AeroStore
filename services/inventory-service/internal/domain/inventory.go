package domain

import (
	"context"
	"time"
	"errors"
)

// Custom errors for inventory service
var (
	ErrNotEnoughStock = errors.New("not enough stock available")
	ErrItemNotFound   = errors.New("item not found in inventory")
)

// InventoryItem represents a product in our warehouse

type InventoryItem struct {
	ID string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity int `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InventoryRepository is the contract for the Database (Postgres/Redis)
// The domain doesn't care IF it's Postgres, it just knows these functions must exist.
type InventoryRepository interface {
	GetStock(ctx context.Context, productID string) (*InventoryItem, error)
	DeductStock(ctx context.Context, items []InventoryItem) error
	ReleaseStock(ctx context.Context, items []InventoryItem) error
}

// InventoryUseCase is the contract for our Business Logic
type InventoryUseCase interface {
	ReserveItems(ctx context.Context, orderID string, items []InventoryItem) error
	ReleaseItems(ctx context.Context, orderID string) error
}