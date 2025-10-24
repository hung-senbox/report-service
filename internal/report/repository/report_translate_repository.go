package repository

import (
	"context"
	"report-service/internal/report/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportTranslateRepo interface {
	Create(ctx context.Context, data *model.ReportTranslation) error
	Update(ctx context.Context, data *model.ReportTranslation) error
	FindByStudentTopicTerm(ctx context.Context, studentID, topicID, termID string) (*model.ReportTranslation, error)
}

type reportTranslateRepo struct {
	ReportTranslateCollection *mongo.Collection
}

func NewReportTranslateRepo(collection *mongo.Collection) ReportTranslateRepo {
	return &reportTranslateRepo{
		ReportTranslateCollection: collection,
	}
}

func (r *reportTranslateRepo) Create(ctx context.Context, data *model.ReportTranslation) error {
	_, err := r.ReportTranslateCollection.InsertOne(ctx, data)
	return err
}

func (r *reportTranslateRepo) Update(ctx context.Context, data *model.ReportTranslation) error {
	_, err := r.ReportTranslateCollection.UpdateOne(ctx, bson.M{"_id": data.ID}, bson.M{"$set": data})
	return err
}

func (r *reportTranslateRepo) FindByStudentTopicTerm(ctx context.Context, studentID, topicID, termID string) (*model.ReportTranslation, error) {

	var report model.ReportTranslation

	filter := bson.M{
		"student_id": studentID,
		"topic_id":   topicID,
		"term_id":    termID,
	}

	err := r.ReportTranslateCollection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &report, nil

}