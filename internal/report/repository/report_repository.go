package repository

import (
	"context"
	"errors"
	"fmt"
	"report-service/internal/report/model"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReportRepository interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	GetByID(ctx context.Context, id string) (*model.Report, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*model.Report, error)
	CreateOrUpdate(ctx context.Context, report *model.Report) error
	GetByStudentTopicTermAndLanguage(ctx context.Context, studentID, topicID, termID, language string) (*model.Report, error)
	GetByStudentTopicTermLanguageAndEditor(ctx context.Context, studentID, topicID, termID, language, editorID string) (*model.Report, error)
	GetAllByEditorID(ctx context.Context, editorID string) ([]*model.Report, error)
	CreateOrUpdate4App(ctx context.Context, report *model.Report) error
	CreateOrUpdate4Web(ctx context.Context, report *model.Report) error
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

	// filter theo bộ unique (student_id + topic_id + term_id)
	filter := bson.M{
		"student_id": report.StudentID,
		"topic_id":   report.TopicID,
		"term_id":    report.TermID,
		"language":   report.Language,
	}

	update := bson.M{
		"$set": bson.M{
			"student_id":  report.StudentID,
			"editor_id":   report.EditorID,
			"topic_id":    report.TopicID,
			"term_id":     report.TermID,
			"language":    report.Language,
			"status":      report.Status,
			"note":        report.Note,
			"report_data": report.ReportData,
			"updated_at":  report.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":        report.ID,
			"created_at": report.CreatedAt,
		},
	}

	// upsert = true → nếu chưa có thì insert, có rồi thì update
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *reportRepository) GetByStudentTopicTermAndLanguage(ctx context.Context, studentID, topicID, termID, language string) (*model.Report, error) {
	filter := bson.M{
		"student_id": studentID,
		"topic_id":   topicID,
		"term_id":    termID,
		"language":   language,
	}

	var report model.Report
	err := r.collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) GetByStudentTopicTermLanguageAndEditor(ctx context.Context, studentID, topicID, termID, language, editorID string) (*model.Report, error) {
	filter := bson.M{
		"student_id": studentID,
		"topic_id":   topicID,
		"term_id":    termID,
		"language":   language,
		"editor_id":  editorID,
	}

	var report model.Report
	err := r.collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) GetAllByEditorID(ctx context.Context, editorID string) ([]*model.Report, error) {
	filter := bson.M{
		"editor_id": editorID,
	}

	cur, err := r.collection.Find(ctx, filter)
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

func (r *reportRepository) CreateOrUpdate4App(ctx context.Context, report *model.Report) error {
	filter := bson.M{
		"student_id": report.StudentID,
		"topic_id":   report.TopicID,
		"term_id":    report.TermID,
		"language":   report.Language,
	}

	update := bson.M{
		"$set": bson.M{
			"status":     report.Status,
			"editor_id":  report.EditorID,
			"updated_at": time.Now(),
		},
	}

	// merge report_data
	for section, data := range report.ReportData {
		if section == "goal" || section == "title" || section == "sub_title" {
			continue
		}
		subData, ok := data.(map[string]interface{})
		if !ok {
			continue
		}

		for k, v := range subData {
			if strings.HasPrefix(k, "manager_") {
				continue
			}
			update["$set"].(bson.M)[fmt.Sprintf("report_data.%s.%s", section, k)] = v
		}
	}

	// merge note
	for k, v := range report.Note {
		if strings.HasPrefix(k, "manager_") {
			continue
		}
		update["$set"].(bson.M)[fmt.Sprintf("note.%s", k)] = v
	}

	update["$setOnInsert"] = bson.M{
		"created_at": time.Now(),
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("create or update report failed: %w", err)
	}

	return nil
}

func (r *reportRepository) CreateOrUpdate4Web(ctx context.Context, report *model.Report) error {
	filter := bson.M{
		"student_id": report.StudentID,
		"topic_id":   report.TopicID,
		"term_id":    report.TermID,
		"language":   report.Language,
	}

	update := bson.M{
		"$set": bson.M{
			"status":     report.Status,
			"editor_id":  report.EditorID,
			"updated_at": time.Now(),
		},
	}

	// --- merge report_data: chỉ field manager_ và status ---
	for section, data := range report.ReportData {
		subData, ok := data.(map[string]interface{})
		if !ok {
			continue
		}

		for k, v := range subData {
			// Cho phép update nếu là manager_* hoặc status
			if strings.HasPrefix(k, "manager_") || k == "status" {
				update["$set"].(bson.M)[fmt.Sprintf("report_data.%s.%s", section, k)] = v
			}
		}
	}

	// --- merge note: chỉ field manager_* và status ---
	for k, v := range report.Note {
		if strings.HasPrefix(k, "manager_") || k == "status" {
			update["$set"].(bson.M)[fmt.Sprintf("note.%s", k)] = v
		}
	}

	opts := options.Update().SetUpsert(false)
	res, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update report (web) failed: %w", err)
	}
	if res.MatchedCount == 0 {
		return errors.New("report not found")
	}

	return nil
}
