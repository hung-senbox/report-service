package repository

import (
	"context"
	"report-service/internal/report/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportHistoryRepository interface {
	Create(ctx context.Context, history *model.ReportHistory) error
	GetAll(ctx context.Context) ([]*model.ReportHistory, error)
	GetByEditor(ctx context.Context, editorID string, editorRole string) ([]*model.ReportHistory, error)
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

func (r *reportHistoryRepository) GetAll(ctx context.Context) ([]*model.ReportHistory, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []*model.ReportHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *reportHistoryRepository) GetByEditor(ctx context.Context, editorID string, editorRole string) ([]*model.ReportHistory, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"editor_id": editorID, "editor_role": editorRole})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []*model.ReportHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, err
	}

	return histories, nil
}
