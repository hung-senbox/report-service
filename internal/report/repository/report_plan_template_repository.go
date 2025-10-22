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

type ReportPlanTemplateRepository interface {
	Create(ctx context.Context, rpt *model.ReportPlanTemplate) error
	CreateOrUpdate(ctx context.Context, rpt *model.ReportPlanTemplate) error
	GetSchoolTemplate(ctx context.Context, termID, topicID, language, organizationID string) (*model.ReportPlanTemplate, error)
	GetClassroomTemplate(ctx context.Context, termID, topicID, language, classroomID, organizationID string) (*model.ReportPlanTemplate, error)
}

type reportPlanTemplateRepository struct {
	collection *mongo.Collection
}

func NewReportPlanTemplateRepository(collection *mongo.Collection) ReportPlanTemplateRepository {
	return &reportPlanTemplateRepository{collection}
}

func (r *reportPlanTemplateRepository) Create(ctx context.Context, rpt *model.ReportPlanTemplate) error {
	now := time.Now().Unix()
	rpt.CreatedAt = now
	rpt.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, rpt)
	return err
}

func (r *reportPlanTemplateRepository) CreateOrUpdate(ctx context.Context, rpt *model.ReportPlanTemplate) error {
	filter := bson.M{
		"organization_id": rpt.OrganizationID,
		"topic_id":        rpt.TopicID,
		"term_id":         rpt.TermID,
		"language":        rpt.Language,
		"is_school":       rpt.IsSchool,
	}

	now := time.Now().Unix()
	rpt.UpdatedAt = now
	if rpt.ID == primitive.NilObjectID {
		rpt.ID = primitive.NewObjectID()
	}

	update := bson.M{
		"$set": bson.M{
			"template":        rpt.Template,
			"organization_id": rpt.OrganizationID,
			"topic_id":        rpt.TopicID,
			"term_id":         rpt.TermID,
			"language":        rpt.Language,
			"is_school":       rpt.IsSchool,
			"classroom_id":    rpt.ClassroomID,
			"updated_at":      rpt.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":        rpt.ID,
			"created_at": now,
		},
	}

	// upsert = true: tạo mới nếu không tồn tại
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *reportPlanTemplateRepository) GetSchoolTemplate(ctx context.Context, termID, topicID, language, organizationID string) (*model.ReportPlanTemplate, error) {
	filter := bson.M{
		"term_id":         termID,
		"topic_id":        topicID,
		"language":        language,
		"organization_id": organizationID,
		"is_school":       true,
	}

	var reportTemplate *model.ReportPlanTemplate
	err := r.collection.FindOne(ctx, filter).Decode(&reportTemplate)
	return reportTemplate, err
}

func (r *reportPlanTemplateRepository) GetClassroomTemplate(ctx context.Context, termID, topicID, language, organizationID, classroomID string) (*model.ReportPlanTemplate, error) {
	filter := bson.M{
		"term_id":         termID,
		"topic_id":        topicID,
		"language":        language,
		"organization_id": organizationID,
		"is_school":       false,
		"classroom_id":    classroomID,
	}

	var reportTemplate *model.ReportPlanTemplate
	err := r.collection.FindOne(ctx, filter).Decode(&reportTemplate)
	return reportTemplate, err
}
