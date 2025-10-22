package service

import (
	"context"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/usecase"
)

type ReportService interface {
	UploadReport4App(ctx context.Context, req *request.UploadReport4AppRequest) error
	UploadReport4Web(ctx context.Context, req *request.UploadReport4AWebRequest) error
	GetReport4App(ctx context.Context, req *request.GetReportRequest4App) (response.ReportResponse, error)
	GetReport4Web(ctx context.Context, req *request.GetReportRequest4Web) (response.ReportResponse, error)
	GetTeacherReportTasks4App(ctx context.Context) ([]response.GetTeacherReportTasksResponse4App, error)
	UploadClassroomReport4Web(ctx context.Context, req request.UploadClassroomReport4WebRequest) error
	GetClassroomReports4Web(ctx context.Context, req request.GetClassroomReportRequest4Web) (*response.GetClassroomReportResponse4Web, error)
	ApplyTopicPlanTemplateIsSchool2Report(ctx context.Context, req request.ApplyTemplateIsSchoolToReportRequest) error
	ApplyTopicPlanTemplateIsClassroom2Report(ctx context.Context, req request.ApplyTemplateIsClassroomToReportRequest) error
	GetReportOverViewAllClassroom4Web(ctx context.Context, req request.GetReportOverViewAllClassroomRequest) (*response.GetReportOverviewAllClassroomResponse4Web, error)
	GetReportOverViewByClassroom4Web(ctx context.Context, req request.GetReportOverViewByClassroomRequest) (*response.GetReportOverviewByClassroomResponse4Web, error)
}

type reportService struct {
	appUsecase usecase.ReportAppUseCase
	webUsecase usecase.ReportWebUseCase
}

func NewReportService(
	appUsecase usecase.ReportAppUseCase,
	webUsecase usecase.ReportWebUseCase,
) ReportService {
	return &reportService{
		appUsecase: appUsecase,
		webUsecase: webUsecase,
	}
}

func (s *reportService) UploadReport4App(ctx context.Context, req *request.UploadReport4AppRequest) error {
	return s.appUsecase.UploadReport4App(ctx, req)
}

func (s *reportService) GetReport4App(ctx context.Context, req *request.GetReportRequest4App) (response.ReportResponse, error) {
	return s.appUsecase.GetReport4App(ctx, req)
}

func (s *reportService) GetTeacherReportTasks4App(ctx context.Context) ([]response.GetTeacherReportTasksResponse4App, error) {
	return s.appUsecase.GetTeacherReportTasks4App(ctx)
}

func (s *reportService) GetReport4Web(ctx context.Context, req *request.GetReportRequest4Web) (response.ReportResponse, error) {
	return s.webUsecase.GetReport4Web(ctx, req)
}

func (s *reportService) UploadReport4Web(ctx context.Context, req *request.UploadReport4AWebRequest) error {
	return s.webUsecase.UploadReport4Web(ctx, req)
}

func (s *reportService) UploadClassroomReport4Web(ctx context.Context, req request.UploadClassroomReport4WebRequest) error {
	return s.webUsecase.UploadClassroomReport4Web(ctx, req)
}

func (s *reportService) ApplyTopicPlanTemplateIsSchool2Report(ctx context.Context, req request.ApplyTemplateIsSchoolToReportRequest) error {
	return s.webUsecase.ApplyTopicPlanTemplateIsSchool2Report(ctx, req)
}

func (s *reportService) ApplyTopicPlanTemplateIsClassroom2Report(ctx context.Context, req request.ApplyTemplateIsClassroomToReportRequest) error {
	return s.webUsecase.ApplyTopicPlanTemplateIsClassroom2Report(ctx, req)
}

func (s *reportService) GetReportOverViewAllClassroom4Web(ctx context.Context, req request.GetReportOverViewAllClassroomRequest) (*response.GetReportOverviewAllClassroomResponse4Web, error) {
	return s.webUsecase.GetReportOverViewAllClassroom4Web(ctx, req)
}

func (s *reportService) GetReportOverViewByClassroom4Web(ctx context.Context, req request.GetReportOverViewByClassroomRequest) (*response.GetReportOverviewByClassroomResponse4Web, error) {
	return s.webUsecase.GetReportOverViewByClassroom4Web(ctx, req)
}

func (s *reportService) GetClassroomReports4Web(ctx context.Context, req request.GetClassroomReportRequest4Web) (*response.GetClassroomReportResponse4Web, error) {
	return s.webUsecase.GetClassroomReports4Web(ctx, req)
}
