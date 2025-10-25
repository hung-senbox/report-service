package service

import (
	"context"
	"errors"
	"report-service/internal/gateway"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"
	"report-service/internal/report/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportTranslateService interface {
	UploadReportTranslate4Web(ctx context.Context, req request.UploadReportTranslateRequest) error
	GetReportTranslate4WebByTopicAndLang(ctx context.Context, studentID, topicID, termID, lang string) (*model.ReportTranslationData, error)
	GetReportTranslate4WebByReport(ctx context.Context, studentID, termID, lang string) ([]*response.ReportTranslateResponse, error)
}

type reportTranslateService struct {
	ReportTranslateRepo repository.ReportTranslateRepo
	TopicGateWay        gateway.MediaGateway
}

func NewReportTranslateService(repo repository.ReportTranslateRepo, TopicGateWay gateway.MediaGateway) ReportTranslateService {
	return &reportTranslateService{
		ReportTranslateRepo: repo,
		TopicGateWay:        TopicGateWay,
	}
}

func (s *reportTranslateService) UploadReportTranslate4Web(ctx context.Context, req request.UploadReportTranslateRequest) error {

	if req.StudentID == "" {
		return errors.New("student id is required")
	}

	if req.TopicID == "" {
		return errors.New("topic id is required")
	}

	if req.TermID == "" {
		return errors.New("term id is required")
	}

	if req.Language == "" {
		return errors.New("language is required")
	}

	if req.ReportData == nil {
		return errors.New("report data is required")
	}

	var translateData model.ReportTranslationData

	if before, ok := req.ReportData["before"].(string); ok {
		translateData.Before = before
	} else {
		translateData.Before = ""
	}

	if conclusion, ok := req.ReportData["conclusion"].(string); ok {
		translateData.Conclusion = conclusion
	} else {
		translateData.Conclusion = ""
	}

	if now, ok := req.ReportData["now"].(string); ok {
		translateData.Now = now
	} else {
		translateData.Now = ""
	}

	existing, err := s.ReportTranslateRepo.FindByStudentTopicTerm(ctx, req.StudentID, req.TopicID, req.TermID)
	if err != nil {
		return err
	}

	if existing == nil {
		newData := &model.ReportTranslation{
			ID:           primitive.NewObjectID(),
			StudentID:    req.StudentID,
			TopicID:      req.TopicID,
			TermID:       req.TermID,
			Translations: map[string]model.ReportTranslationData{req.Language: translateData},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.ReportTranslateRepo.Create(ctx, newData); err != nil {
			return err
		}
	} else {

		existing.Translations[req.Language] = translateData
		existing.UpdatedAt = time.Now()

		if err := s.ReportTranslateRepo.Update(ctx, existing); err != nil {
			return err
		}
	}

	return nil

}

func (s *reportTranslateService) GetReportTranslate4WebByTopicAndLang(ctx context.Context, studentID, topicID, termID, lang string) (*model.ReportTranslationData, error) {

	if studentID == "" {
		return nil, errors.New("student id is required")
	}

	if topicID == "" {
		return nil, errors.New("topic id is required")
	}

	if termID == "" {
		return nil, errors.New("term id is required")
	}

	if lang == "" {
		return nil, errors.New("language is required")
	}

	reportTranslate, err := s.ReportTranslateRepo.FindByStudentTopicTerm(ctx, studentID, topicID, termID)
	if err != nil {
		return nil, err
	}

	var result model.ReportTranslationData
	if reportTranslate == nil {
		result = model.ReportTranslationData{
			Before:     "",
			Conclusion: "",
			Now:        "",
		}
	} else {
		if translate, ok := reportTranslate.Translations[lang]; ok {
			result = translate
		} else {
			result = model.ReportTranslationData{
				Before:     "",
				Conclusion: "",
				Now:        "",
			}
		}
	}

	return &result, nil
}

func (s *reportTranslateService) GetReportTranslate4WebByReport(ctx context.Context, studentID, termID, lang string) ([]*response.ReportTranslateResponse, error) {

	if studentID == "" {
		return nil, errors.New("student id is required")
	}

	topics, err := s.TopicGateWay.GetTopicByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	var result []*response.ReportTranslateResponse

	for _, t := range topics {
		reportTranslate, err := s.ReportTranslateRepo.FindByStudentTopicTerm(ctx, studentID, t.ID, termID)
		if err != nil {
			return nil, err
		}

		resp := &response.ReportTranslateResponse{
			Topic: response.Topic{
				ID:           t.ID,
				Title:        t.Title,
				MainImageUrl: t.MainImageUrl,
			},
			Translations: map[string]model.ReportTranslationData{},
		}

		if reportTranslate == nil {
			result = append(result, resp)
			continue
		}

		if lang != "" {
			if data, ok := reportTranslate.Translations[lang]; ok {
				resp.Translations[lang] = data
			}
		} else {
			resp.Translations = reportTranslate.Translations
		}

		result = append(result, resp)
	}

	return result, nil
}
