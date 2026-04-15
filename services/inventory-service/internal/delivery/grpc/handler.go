package grpc

import (
	"context"

	// Import your generated protobuf code
	pb "github.com/durgaprasadraju/AeroStore/api/proto"
	"github.com/durgaprasadraju/AeroStore/services/inventory-service/internal/domain"
)

// inventoryHandler implements the generated gRPC interface
type inventoryHandler struct {
	pb.UnimplementedInventoryServiceServer
	usecase domain.InventoryUseCase
}

// NewInventoryHandler creates the gRPC handler
func NewInventoryHandler(uc domain.InventoryUseCase) *inventoryHandler {
	return &inventoryHandler{usecase: uc}
}

// ReserveStock is the actual network endpoint called by other microservices!
func (h *inventoryHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {

	// 1. Convert the Protobuf request into our internal Domain struct
	var items []domain.InventoryItem
	for _, item := range req.Items {
		items = append(items, domain.InventoryItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
		})
	}

	// 2. Call our Business Logic (UseCase)
	err := h.usecase.ReserveItems(ctx, req.OrderId, items)
	if err != nil {
		if err == domain.ErrNotEnoughStock {
			return &pb.ReserveStockResponse{Success: false, Message: "Not enough stock"}, nil
		}
		return nil, err // Return an actual gRPC error
	}

	// 3. Return success!
	return &pb.ReserveStockResponse{Success: true, Message: "Stock reserved successfully!"}, nil
}
