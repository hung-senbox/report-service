package service

import (
	"context"
	"errors"
	"report-service/internal/gateway"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/model"
	"report-service/internal/report/repository"
)

type ReportPlanTemplateService interface {
	Upload(ctx context.Context, req request.UploadReportPlanTemplateRequest) error
}

type reportPlanTemplateService struct {
	repo        repository.ReportPlanTemplateRepositopry
	userGateway gateway.UserGateway
}

func NewReportPlanTemplateService(repo repository.ReportPlanTemplateRepositopry, userGateway gateway.UserGateway) ReportPlanTemplateService {
	return &reportPlanTemplateService{
		repo:        repo,
		userGateway: userGateway,
	}
}

func (s *reportPlanTemplateService) Upload(ctx context.Context, req request.UploadReportPlanTemplateRequest) error {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	if currentUser.IsSuperAdmin {
		return errors.New("super admin can't upload report plan template")
	}

	rpPlanTemp := &model.ReportPlanTemplate{
		OrganizationID: currentUser.OrganizationAdmin.ID,
		TopicID:        req.TopicID,
		TermID:         req.TermID,
		Language:       req.Language,
		IsSchool:       req.IsSchool,
		Template: model.Template{
			Title:          req.Title,
			Introduction:   req.Introduction,
			CurriculumArea: req.CurriculumArea,
		},
	}

	return s.repo.CreateOrUpdate(ctx, rpPlanTemp)
}
