package repository

import (
	"context"
	"report-service/internal/report/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type ReportHistoryRepository interface {
	Create(ctx context.Context, history *model.ReportHistory) error
}

type reportHistoryRepository struct {
	collection *mongo.Collection
}

func NewReportHistoryRepository(collection *mongo.Collection) ReportHistoryRepository {
	return &reportHistoryRepository{collection}
}

func (r *reportHistoryRepository) Create(ctx context.Context, history *model.ReportHistory) error {
	_, err := r.collection.InsertOne(ctx, history)
	return err
}
