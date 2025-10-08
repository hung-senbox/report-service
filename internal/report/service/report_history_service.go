package service

import (
	"context"
	"report-service/helper"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/mapper"
	"report-service/internal/report/repository"
	"report-service/pkg/constants"
)

type ReportHistoryService interface {
	GetByEditor4App(ctx context.Context) ([]response.ReportHistoryResponse, error)
}

type reportHistoryService struct {
	repo repository.ReportHistoryRepository
}

func NewReportHistoryService(repo repository.ReportHistoryRepository) ReportHistoryService {
	return &reportHistoryService{repo: repo}
}

func (s reportHistoryService) GetByEditor4App(ctx context.Context) ([]response.ReportHistoryResponse, error) {
	editorID := helper.GetUserID(ctx)
	editorRole := string(constants.ReportHistoryRoleTeacher)
	histories, err := s.repo.GetByEditor(ctx, editorID, editorRole)
	if err != nil {
		return nil, err
	}
	return mapper.MapReportHistoryListToRes4App(histories), nil
}
