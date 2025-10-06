package service

import (
	"context"
	"errors"
	"fmt"
	"report-service/internal/gateway"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/mapper"
	"report-service/internal/report/model"
	"report-service/internal/report/repository"
	"report-service/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	GetByID(ctx context.Context, id string) (*model.Report, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]response.ReportResponse, error)
	UploadReport4App(ctx context.Context, req *request.UploadReport4AppRequest) error
	GetReport4App(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error)
	GetReport4Web(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error)
}

type reportService struct {
	userGateway  gateway.UserGateway
	termGateway  gateway.TermGateway
	mediaGateway gateway.MediaGateway
	repo         repository.ReportRepository
	historyRepo  repository.ReportHistoryRepository
}

func NewReportService(
	userGateway gateway.UserGateway,
	termGateway gateway.TermGateway,
	mediaGateway gateway.MediaGateway,
	repo repository.ReportRepository,
	historyRepo repository.ReportHistoryRepository,
) ReportService {
	return &reportService{
		userGateway:  userGateway,
		termGateway:  termGateway,
		mediaGateway: mediaGateway,
		repo:         repo,
		historyRepo:  historyRepo,
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
		// get student info
		student, _ := s.userGateway.GetStudentInfo(ctx, report.StudentID)
		editor, _ := s.userGateway.GetTeacherByUserAndOrganization(ctx, report.EditorID, student.OrganizationID)

		if editor != nil {
			res[i].Editor = *editor
		}
	}
	return res, nil
}

func (s *reportService) UploadReport4App(ctx context.Context, req *request.UploadReport4AppRequest) error {
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

func (s *reportService) GetReport4App(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error) {
	report, err := s.repo.GetByStudentTopicTermAndLanguage(ctx, req.StudentID, req.TopicID, req.TermID, req.Language)
	if err != nil {
		return response.ReportResponse{}, err
	}
	if report == nil {
		return response.ReportResponse{}, errors.New("report not found")
	}
	return mapper.MapReportToResDTO(report), nil
}

func (s *reportService) GetReport4Web(ctx context.Context, req *request.GetReportRequest) (response.ReportResponse, error) {
	// get edtior by teacher id
	editor, err := s.userGateway.GetUserByTeacher(ctx, req.TeacherID)
	if err != nil {
		return response.ReportResponse{}, err
	}

	report, err := s.repo.GetByStudentTopicTermLanguageAndEditor(ctx, req.StudentID, req.TopicID, req.TermID, req.Language, editor.ID)
	if err != nil {
		return response.ReportResponse{}, err
	}
	if report == nil {
		return response.ReportResponse{}, errors.New("report not found")
	}

	res := mapper.MapReportToResDTO(report)

	return res, nil
}

func (s *reportService) GetTeacherReportTasks(ctx context.Context, teacherID string) ([]response.GetTeacherReportTasksResponse, error) {
	// Lấy thông tin editor từ teacherID
	editor, err := s.userGateway.GetUserByTeacher(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("get teacher failed: %w", err)
	}

	// Lấy tất cả reports do editor này phụ trách
	reports, err := s.repo.GetAllByEditorID(ctx, editor.ID)
	if err != nil {
		return nil, fmt.Errorf("get reports failed: %w", err)
	}

	var results []response.GetTeacherReportTasksResponse

	for _, r := range reports {
		// Bỏ qua nếu main status không phù hợp
		if r.Status != "Teacher" && r.Status != "Empty" {
			continue
		}

		// Duyệt qua các field trong report_data để xem phần nào thuộc "Teacher" hoặc "Empty"
		if reportData, ok := r.ReportData["report_data"].(bson.M); ok {
			for key, val := range reportData {
				if section, ok := val.(bson.M); ok {
					status, _ := section["status"].(string)
					term, _ := s.termGateway.GetTermByID(ctx, r.TermID)
					topic, _ := s.mediaGateway.GetTopicByID(ctx, r.TopicID)
					student, _ := s.userGateway.GetStudentInfo(ctx, r.StudentID)
					if status == "Teacher" || status == "Empty" {
						results = append(results, response.GetTeacherReportTasksResponse{
							Term:        term.Title,
							Topic:       topic.Title,
							StudentName: student.Name,
							Deadline:    "123",
							Task:        constants.TeacherReportTask(key),
							Status:      status,
						})
					}
				}
			}
		}
	}

	return results, nil
}
