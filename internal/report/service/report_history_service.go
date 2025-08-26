package service

import (
	"context"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/mapper"
	"report-service/internal/report/model"
	"report-service/internal/report/repository"
)

type ReportHistoryService interface {
	GetAll(ctx context.Context) ([]response.ReportHistoryResponse, error)
	Create(ctx context.Context, history *model.ReportHistory) (response.ReportHistoryResponse, error)
}

type reportHistoryService struct {
	repo repository.ReportHistoryRepository
}

func NewReportHistoryService(repo repository.ReportHistoryRepository) ReportHistoryService {
	return &reportHistoryService{repo: repo}
}

func (s *reportHistoryService) Create(ctx context.Context, history *model.ReportHistory) (response.ReportHistoryResponse, error) {
	err := s.repo.Create(ctx, history)
	if err != nil {
		return response.ReportHistoryResponse{}, err
	}

	return mapper.MapReportHistoryToResDTO(history), nil
}

func (s *reportHistoryService) GetAll(ctx context.Context) ([]response.ReportHistoryResponse, error) {
	histories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.MapReportHistoryListToResDTO(histories), nil
}
