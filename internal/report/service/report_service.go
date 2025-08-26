package service

import (
	"context"
	"errors"
	"report-service/internal/gateway"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/mapper"
	"report-service/internal/report/model"
	"report-service/internal/report/repository"
	"report-service/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	GetByID(ctx context.Context, id string) (*model.Report, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]response.ReportResponse, error)
	UploadReport(ctx context.Context, req *request.UploadReportRequestDTO) error
	Get4Report(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error)
	GetReport4Admin(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error)
}

type reportService struct {
	userGateway gateway.UserGateway
	repo        repository.ReportRepository
	historyRepo repository.ReportHistoryRepository
}

func NewReportService(
	userGateway gateway.UserGateway,
	repo repository.ReportRepository,
	historyRepo repository.ReportHistoryRepository,
) ReportService {
	return &reportService{
		userGateway: userGateway,
		repo:        repo,
		historyRepo: historyRepo,
	}
}

func (s *reportService) Create(ctx context.Context, report *model.Report) (*model.Report, error) {
	if report == nil {
		return nil, errors.New("report is nil")
	}
	return s.repo.Create(ctx, report)
}

func (s *reportService) GetByID(ctx context.Context, id string) (*model.Report, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *reportService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *reportService) GetAll(ctx context.Context) ([]response.ReportResponse, error) {
	// Lấy danh sách report từ repository
	reports, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi sang DTO response
	res := mapper.MapReportListToResDTO(reports)

	// get editor info for each report

	for i, report := range reports {
		editor, _ := s.userGateway.GetTeacherByUser(ctx, report.EditorID)

		if editor != nil {
			res[i].Editor = *editor
		}
	}
	return res, nil
}

func (s *reportService) UploadReport(ctx context.Context, req *request.UploadReportRequestDTO) error {
	report := &model.Report{
		StudentID:  req.StudentID,
		TopicID:    req.TopicID,
		TermID:     req.TermID,
		Language:   req.Language,
		Status:     req.Status,
		Note:       req.Note,
		ReportData: req.ReportData,
	}

	// get editor_id from context (from middleware)
	if editorID, ok := ctx.Value(constants.UserID).(string); ok {
		report.EditorID = editorID
	}

	// create or update report
	err := s.repo.CreateOrUpdate(ctx, report)
	if err != nil {
		return err
	}

	// save report history
	history := &model.ReportHistory{
		ID:        primitive.NewObjectID(),
		ReportID:  report.ID,
		EditorID:  report.EditorID,
		Report:    report,
		Timestamp: time.Now().Unix(),
	}

	if err := s.historyRepo.Create(ctx, history); err != nil {
		return err
	}

	return nil
}

func (s *reportService) Get4Report(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error) {
	report, err := s.repo.GetByStudentTopicTerm(ctx, req.StudentID, req.TopicID, req.TermID, req.Language)
	if err != nil {
		return response.ReportResponse{}, err
	}
	if report == nil {
		return response.ReportResponse{}, errors.New("report not found")
	}
	return mapper.MapReportToResDTO(report), nil
}

func (s *reportService) GetReport4Admin(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error) {
	report, err := s.repo.GetByStudentTopicTerm(ctx, req.StudentID, req.TopicID, req.TermID, req.Language)
	if err != nil {
		return response.ReportResponse{}, err
	}
	if report == nil {
		return response.ReportResponse{}, errors.New("report not found")
	}
	res := mapper.MapReportToResDTO(report)

	// get editor info
	editor, _ := s.userGateway.GetTeacherByUser(ctx, report.EditorID)
	if editor != nil {
		res.Editor = *editor
	}
	return res, nil
}
