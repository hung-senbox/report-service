package service

import (
	"context"
	"errors"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/mapper"
	"report-service/internal/report/model"
	"report-service/internal/report/repository"
)

type ReportService interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	GetByID(ctx context.Context, id string) (*model.Report, error)
	Update(ctx context.Context, id string, report *model.Report) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]response.ReportResponse, error)
	UploadReport(ctx context.Context, req *request.UploadReportRequestDTO) error
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) Create(ctx context.Context, report *model.Report) (*model.Report, error) {
	if report == nil {
		return nil, errors.New("report is nil")
	}
	if report.Key == "" {
		return nil, errors.New("report key is required")
	}
	return s.repo.Create(ctx, report)
}

func (s *reportService) GetByID(ctx context.Context, id string) (*model.Report, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *reportService) Update(ctx context.Context, id string, report *model.Report) error {
	if id == "" {
		return errors.New("id is required")
	}
	if report == nil {
		return errors.New("report is nil")
	}
	return s.repo.Update(ctx, id, report)
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
	return res, nil
}

func (s *reportService) UploadReport(ctx context.Context, req *request.UploadReportRequestDTO) error {

	report := &model.Report{
		Key:        req.Key,
		Note:       req.Note,
		ReportData: req.ReportData,
	}

	// create or update
	return s.repo.CreateOrUpdate(ctx, report)
}
