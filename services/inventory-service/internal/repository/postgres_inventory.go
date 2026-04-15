package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) domain.InventoryRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetStock(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	query := `SELECT id, product_id, quantity, created_at, updated_at FROM inventory WHERE product_id = $1`

	var item domain.InventoryItem
	err := r.db.QueryRow(ctx, query, productID).Scan(
		&item.ID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *postgresRepository) DeductStock(ctx context.Context, items []domain.InventoryItem) error {
	// Start a transaction: Critical for stability
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Rollback if the function returns early

	for _, item := range items {
		// Advanced SQL: Check quantity >= requested amount in the same query to prevent race conditions
		query := `UPDATE inventory SET quantity = quantity - $1, updated_at = NOW() 
		          WHERE product_id = $2 AND quantity >= $1`

		cmd, err := tx.Exec(ctx, query, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return fmt.Errorf("product %s: %w", item.ProductID, domain.ErrNotEnoughStock)
		}
	}

	return tx.Commit(ctx)
}

func (r *postgresRepository) ReleaseStock(ctx context.Context, items []domain.InventoryItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		query := `UPDATE inventory SET quantity = quantity + $1, updated_at = NOW() WHERE product_id = $2`
		_, err := tx.Exec(ctx, query, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
