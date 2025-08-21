package repository

import (
	"context"
	"report-service/internal/report/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReportRepository interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	GetByID(ctx context.Context, id string) (*model.Report, error)
	Update(ctx context.Context, id string, report *model.Report) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*model.Report, error)
	CreateOrUpdate(ctx context.Context, report *model.Report) error
}

type reportRepository struct {
	collection *mongo.Collection
}

func NewReportRepository(collection *mongo.Collection) ReportRepository {
	return &reportRepository{collection}
}

func (r *reportRepository) Create(ctx context.Context, report *model.Report) (*model.Report, error) {
	now := time.Now()
	report.ID = primitive.NewObjectID()
	report.CreatedAt = now
	report.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (r *reportRepository) GetByID(ctx context.Context, id string) (*model.Report, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var report model.Report
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) Update(ctx context.Context, id string, report *model.Report) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	report.UpdatedAt = time.Now()
	update := bson.M{
		"key":         report.Key,
		"note":        report.Note,
		"report_data": report.ReportData,
		"updated_at":  report.UpdatedAt,
	}

	_, err = r.collection.UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	return err
}

func (r *reportRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *reportRepository) GetAll(ctx context.Context) ([]*model.Report, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var reports []*model.Report
	if err = cur.All(ctx, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *reportRepository) CreateOrUpdate(ctx context.Context, report *model.Report) error {
	now := time.Now()
	report.UpdatedAt = now

	// nếu chưa có ID thì tạo mới
	if report.ID.IsZero() {
		report.ID = primitive.NewObjectID()
		report.CreatedAt = now
	}

	filter := bson.M{"key": report.Key}
	update := bson.M{
		"$set": bson.M{
			"note":        report.Note,
			"report_data": report.ReportData,
			"updated_at":  report.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":        report.ID,
			"created_at": report.CreatedAt,
			"key":        report.Key,
		},
	}

	// upsert = true → nếu chưa có thì insert, có rồi thì update
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
