package service

import (
	"context"
	"fmt"
)

type PointServer struct {
}

func (s *PointServer) Point(ctx context.Context, point *PointRequest) (*PointResponse, error) {
	fmt.Println(point)
	return &PointResponse{}, nil
}
