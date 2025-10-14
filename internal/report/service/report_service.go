package service

import (
	"context"
	"errors"
	"fmt"
	"report-service/helper"
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
	UploadReport4Web(ctx context.Context, req *request.UploadReport4AWebRequest) error
	GetReport4App(ctx context.Context, req *request.GetReportRequest4App) (response.ReportResponse, error)
	GetReport4Web(ctx context.Context, req *request.GetReportRequest4Web) (response.ReportResponse, error)
	GetTeacherReportTasks(ctx context.Context) ([]response.GetTeacherReportTasksResponse, error)
	UploadClassroomReport(ctx context.Context, req request.UploadClassroomReport4WebRequest) error
	GetClassroomReports4Web(ctx context.Context, req request.GetClassroomReportRequest4Web) ([]*response.ClassroomReportResponse4Web, error)
}

type reportService struct {
	userGateway      gateway.UserGateway
	termGateway      gateway.TermGateway
	mediaGateway     gateway.MediaGateway
	classroomGateway gateway.ClassroomGateway
	repo             repository.ReportRepository
	historyRepo      repository.ReportHistoryRepository
}

func NewReportService(
	userGateway gateway.UserGateway,
	termGateway gateway.TermGateway,
	mediaGateway gateway.MediaGateway,
	classroomGateway gateway.ClassroomGateway,
	repo repository.ReportRepository,
	historyRepo repository.ReportHistoryRepository,
) ReportService {
	return &reportService{
		userGateway:      userGateway,
		termGateway:      termGateway,
		mediaGateway:     mediaGateway,
		classroomGateway: classroomGateway,
		repo:             repo,
		historyRepo:      historyRepo,
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

	// get student info
	student, _ := s.userGateway.GetStudentInfo(ctx, req.StudentID)
	if student == nil {
		return errors.New("student not found")
	}

	// get teacher by usser id and organization if of student
	editorID := helper.GetUserID(ctx)
	teacher, _ := s.userGateway.GetTeacherInfo(ctx, editorID, student.OrganizationID)
	if teacher == nil {
		return errors.New("teacher not found")
	}

	report := &model.Report{
		StudentID:  req.StudentID,
		TopicID:    req.TopicID,
		TermID:     req.TermID,
		Language:   req.Language,
		Status:     req.Status,
		ReportData: req.ReportData,
	}

	if editorID != "" {
		report.EditorID = editorID
	}

	// create or update report
	err := s.repo.CreateOrUpdateStudentView4App(ctx, report)
	if err != nil {
		return err
	}

	// save report history
	history := &model.ReportHistory{
		ID:         primitive.NewObjectID(),
		ReportID:   report.ID,
		EditorID:   report.EditorID,
		Type:       string(constants.ReportHistoryTypeAppStudentView),
		EditorRole: string(constants.ReportHistoryRoleTeacher),
		Report:     report,
		Timestamp:  time.Now(),
	}

	if err := s.historyRepo.Create(ctx, history); err != nil {
		return err
	}

	return nil
}

func (s *reportService) GetReport4App(ctx context.Context, req *request.GetReportRequest4App) (response.ReportResponse, error) {
	report, err := s.repo.GetByStudentTopicTermAndLanguage(ctx, req.StudentID, req.TopicID, req.TermID, req.Language)
	if err != nil {
		return response.ReportResponse{}, err
	}
	if report == nil {
		return response.ReportResponse{}, errors.New("report not found")
	}

	// get student info
	student, _ := s.userGateway.GetStudentInfo(ctx, report.StudentID)
	if student == nil {
		return response.ReportResponse{}, errors.New("student not found")
	}

	var managerCommentPreviousTerm response.ManagerCommentPreviousTerm
	var teacherReportPrevioiusTerm response.TeacherReportPreviousTerm
	previousTerm, _ := s.termGateway.GetPreviousTerm(ctx, report.TermID, student.OrganizationID)
	if previousTerm != nil {
		previousTermReport, _ := s.repo.GetByStudentTopicTermLanguageAndEditor(ctx, report.StudentID, report.TopicID, previousTerm.ID, report.Language, report.EditorID)
		if previousTermReport != nil {
			managerCommentPreviousTerm.TermTitle = previousTerm.Title

			//fmt.Printf("[DEBUG] type of ReportData[now]: %T\n", previousTermReport.ReportData["now"])
			//fmt.Printf("[DEBUG] value: %#v\n", previousTermReport.ReportData["now"])
			if nowData, ok := previousTermReport.ReportData["now"].(primitive.M); ok {
				if comment, ok := nowData["manager_comment"].(string); ok {
					managerCommentPreviousTerm.Now = comment
				}
				if report, ok := nowData["teacher_report"].(string); ok {
					teacherReportPrevioiusTerm.Now = report
				}
			}

			if conclusionData, ok := previousTermReport.ReportData["conclusion"].(primitive.M); ok {
				if comment, ok := conclusionData["manager_comment"].(string); ok {
					managerCommentPreviousTerm.Conclusion = comment
				}
				if report, ok := conclusionData["teacher_report"].(string); ok {
					teacherReportPrevioiusTerm.Conclusion = report
				}
			}
		}

	}

	return mapper.MapReportToResDTO(report, nil, managerCommentPreviousTerm, teacherReportPrevioiusTerm, ""), nil
}

func (s *reportService) GetReport4Web(ctx context.Context, req *request.GetReportRequest4Web) (response.ReportResponse, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return response.ReportResponse{}, errors.New("get current user failed")
	}

	if currentUser.IsSuperAdmin {
		return response.ReportResponse{}, errors.New("super admin can't get report")
	}

	// get edtior by teacher id
	editor, err := s.userGateway.GetUserByTeacher(ctx, req.TeacherID)
	if err != nil {
		return response.ReportResponse{}, err
	}

	report, err := s.repo.GetByStudentTopicTermLanguageAndEditor(ctx, req.StudentID, req.TopicID, req.TermID, req.UniqueLangKey, editor.ID)
	if err != nil {
		return response.ReportResponse{}, err
	}
	if report == nil {
		return response.ReportResponse{}, errors.New("report not found")
	}

	// get student info
	student, _ := s.userGateway.GetStudentInfo(ctx, report.StudentID)
	if student == nil {
		return response.ReportResponse{}, errors.New("student not found")
	}

	// get teacher
	teacher, _ := s.userGateway.GetTeacherInfo(ctx, report.EditorID, student.OrganizationID)

	var managerCommentPreviousTerm response.ManagerCommentPreviousTerm
	var teacherReportPrevioiusTerm response.TeacherReportPreviousTerm
	previousTerm, _ := s.termGateway.GetPreviousTerm(ctx, report.TermID, student.OrganizationID)
	if previousTerm != nil {
		previousTermReport, _ := s.repo.GetByStudentTopicTermLanguageAndEditor(ctx, report.StudentID, report.TopicID, previousTerm.ID, report.Language, report.EditorID)
		if previousTermReport != nil {
			managerCommentPreviousTerm.TermTitle = previousTerm.Title
			teacherReportPrevioiusTerm.TermTitle = previousTerm.Title

			if nowData, ok := previousTermReport.ReportData["now"].(primitive.M); ok {
				if comment, ok := nowData["manager_comment"].(string); ok {
					managerCommentPreviousTerm.Now = comment
				}
				if report, ok := nowData["teacher_report"].(string); ok {
					teacherReportPrevioiusTerm.Now = report
				}
				if managerUpdatedAt, ok := nowData["manager_updated_at"].(string); ok {
					managerCommentPreviousTerm.NowUpdatedAt = managerUpdatedAt
				}
				if updatedAt, ok := nowData["updated_at"].(string); ok {
					teacherReportPrevioiusTerm.NowUpdatedAt = updatedAt
				}
			}

			if conclusionData, ok := previousTermReport.ReportData["conclusion"].(primitive.M); ok {
				if comment, ok := conclusionData["manager_comment"].(string); ok {
					managerCommentPreviousTerm.Conclusion = comment
				}
				if report, ok := conclusionData["teacher_report"].(string); ok {
					teacherReportPrevioiusTerm.Conclusion = report
				}
				if managerUpdatedAt, ok := conclusionData["manager_updated_at"].(string); ok {
					managerCommentPreviousTerm.ConclusionUpdatedAt = managerUpdatedAt
				}
				if updatedAt, ok := conclusionData["updated_at"].(string); ok {
					teacherReportPrevioiusTerm.ConclusionUpdatedAt = updatedAt
				}
			}
		}

	}

	res := mapper.MapReportToResDTO(report, teacher, managerCommentPreviousTerm, teacherReportPrevioiusTerm, "")

	return res, nil
}

func (s *reportService) GetTeacherReportTasks(ctx context.Context) ([]response.GetTeacherReportTasksResponse, error) {
	userID := helper.GetUserID(ctx)
	// Lấy tất cả reports do editor này phụ trách
	reports, err := s.repo.GetAllByEditorID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get reports failed: %w", err)
	}

	var results []response.GetTeacherReportTasksResponse

	for _, r := range reports {

		if r.ReportData != nil {
			reportData := toBsonM(r.ReportData)
			for key, val := range reportData {
				section := toBsonM(val)
				status, _ := section["status"].(string)

				if status == "teacher" || status == "empty" {
					termTitle := ""
					topicTitle := ""
					stdName := ""
					term, _ := s.termGateway.GetTermByID(ctx, r.TermID)
					topic, _ := s.mediaGateway.GetTopicByID(ctx, r.TopicID)
					student, _ := s.userGateway.GetStudentInfo(ctx, r.StudentID)

					if term != nil {
						termTitle = term.Title
					}
					if topic != nil {
						topicTitle = topic.Title
					}
					if student != nil {
						stdName = student.Name
					}

					results = append(results, response.GetTeacherReportTasksResponse{
						Term:        termTitle,
						Topic:       topicTitle,
						StudentName: stdName,
						Deadline:    "empty",
						Task:        constants.TeacherReportTask(key),
						Status:      status,
					})
				}
			}
		}

	}

	return results, nil
}

func (s *reportService) UploadReport4Web(ctx context.Context, req *request.UploadReport4AWebRequest) error {
	report := &model.Report{
		StudentID:  req.StudentID,
		TopicID:    req.TopicID,
		TermID:     req.TermID,
		Language:   req.UniqueLangKey,
		Status:     req.Status,
		ReportData: req.ReportData,
	}

	// check report da duoc tao tu app chua ?
	reportExist, _ := s.repo.GetByStudentTopicTermAndLanguage(ctx, req.StudentID, req.TopicID, req.TermID, req.UniqueLangKey)
	if reportExist == nil {
		return errors.New("report not found, need to create report from teacher")
	}

	// create or update report
	err := s.repo.CreateOrUpdateStudentView4Web(ctx, report)
	if err != nil {
		return err
	}

	// save report history
	history := &model.ReportHistory{
		ID:         primitive.NewObjectID(),
		ReportID:   report.ID,
		EditorID:   helper.GetUserID(ctx),
		Type:       string(constants.ReportHistoryTypeWebStudentView),
		EditorRole: string(constants.ReportHistoryRoleManager),
		Report:     report,
		Timestamp:  time.Now(),
	}

	if err := s.historyRepo.Create(ctx, history); err != nil {
		return err
	}

	return nil
}

func toBsonM(v interface{}) bson.M {
	if m, ok := v.(bson.M); ok {
		return m
	}
	if m, ok := v.(map[string]interface{}); ok {
		return bson.M(m)
	}
	return bson.M{}
}

// func (s *reportService) getLatestDataTermID(
// 	ctx context.Context,
// 	termID string,
// 	organizationID string,
// 	uploadReq request.UploadReport4AWebRequest,
// ) (string, error) {

// 	// Lấy danh sách các term trước đó, sắp theo thứ tự gần nhất -> xa nhất
// 	previousTerms, err := s.termGateway.GetPreviousTerms(ctx, termID, organizationID)
// 	if err != nil {
// 		return "", fmt.Errorf("get previous terms failed: %w", err)
// 	}

// 	// Duyệt từng term để tìm term có teacher_report != ""
// 	for _, term := range previousTerms {
// 		report, err := s.repo.GetByStudentTopicTermAndLanguage(ctx,
// 			uploadReq.StudentID, uploadReq.TopicID, term.ID, uploadReq.UniqueLangKey)
// 		if err != nil {
// 			continue // bỏ qua lỗi từng report, không dừng toàn bộ
// 		}

// 		if report == nil || report.ReportData == nil {
// 			continue
// 		}

// 		reportData := toBsonM(report.ReportData)

// 		// danh sách các section có thể chứa teacher_report
// 		sections := []string{"before", "now", "note", "conclusion"}

// 		for _, section := range sections {
// 			if data, ok := reportData[section].(bson.M); ok {
// 				if teacherReport, ok := data["teacher_report"].(string); ok && teacherReport != "" {
// 					return term.ID, nil
// 				}
// 			}
// 		}
// 	}

// 	return "", nil
// }

func (s *reportService) UploadClassroomReport(ctx context.Context, req request.UploadClassroomReport4WebRequest) error {
	report := &model.Report{
		StudentID:  req.StudentID,
		TopicID:    req.TopicID,
		TermID:     req.TermID,
		Language:   req.UniqueLangKey,
		Status:     req.Status,
		ReportData: req.ReportData,
	}

	// check report da duoc tao tu app chua ?
	reportExist, _ := s.repo.GetByStudentTopicTermAndLanguage(ctx, req.StudentID, req.TopicID, req.TermID, req.UniqueLangKey)
	if reportExist == nil {
		return errors.New("report not found, need to create report from teacher")
	}

	// create or update report
	err := s.repo.CreateOrUpdateClassroomView4Web(ctx, report)
	if err != nil {
		return err
	}

	// save report history
	history := &model.ReportHistory{
		ID:          primitive.NewObjectID(),
		ReportID:    report.ID,
		EditorID:    helper.GetUserID(ctx),
		ClassroomID: req.ClassroomID,
		Type:        string(constants.ReportHistoryTypeWebClassroomView),
		EditorRole:  string(constants.ReportHistoryRoleManager),
		Report:      report,
		Timestamp:   time.Now(),
	}

	if err := s.historyRepo.Create(ctx, history); err != nil {
		return err
	}

	return nil
}

func (s *reportService) GetClassroomReports4Web(ctx context.Context, req request.GetClassroomReportRequest4Web) ([]*response.ClassroomReportResponse4Web, error) {
	var res = make([]*response.ClassroomReportResponse4Web, 0)

	// Lấy danh sách học sinh trong lớp
	students, err := s.classroomGateway.GetStudents4ClassroomReport(ctx, req.TermID, req.ClassroomID, req.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get students in classroom: %w", err)
	}
	if len(students) == 0 {
		return res, errors.New("no students found in classroom")
	}

	// Lấy thông tin editor
	editor, err := s.userGateway.GetUserByTeacher(ctx, req.TeacherID)
	// get teacher
	teacher, _ := s.userGateway.GetTeacherInfo(ctx, editor.ID, students[0].OrganizationID)

	if err != nil {
		return nil, fmt.Errorf("failed to get editor (teacher): %w", err)
	}

	for _, std := range students {
		report, _ := s.repo.GetByStudentTopicTermLanguageAndEditor(
			ctx,
			std.StudentID,
			req.TopicID,
			req.TermID,
			req.UniqueLangKey,
			editor.ID,
		)

		var reportRes response.ReportResponse
		if report != nil {
			var managerCommentPreviousTerm response.ManagerCommentPreviousTerm
			var teacherReportPrevioiusTerm response.TeacherReportPreviousTerm
			previousTerm, _ := s.termGateway.GetPreviousTerm(ctx, report.TermID, std.OrganizationID)
			if previousTerm != nil {
				previousTermReport, _ := s.repo.GetByStudentTopicTermLanguageAndEditor(ctx, report.StudentID, report.TopicID, previousTerm.ID, report.Language, report.EditorID)
				if previousTermReport != nil {
					managerCommentPreviousTerm.TermTitle = previousTerm.Title
					teacherReportPrevioiusTerm.TermTitle = previousTerm.Title

					if nowData, ok := previousTermReport.ReportData["now"].(primitive.M); ok {
						if comment, ok := nowData["manager_comment"].(string); ok {
							managerCommentPreviousTerm.Now = comment
						}
						if report, ok := nowData["teacher_report"].(string); ok {
							teacherReportPrevioiusTerm.Now = report
						}
						if managerUpdatedAt, ok := nowData["manager_updated_at"].(string); ok {
							managerCommentPreviousTerm.NowUpdatedAt = managerUpdatedAt
						}
						if updatedAt, ok := nowData["updated_at"].(string); ok {
							teacherReportPrevioiusTerm.NowUpdatedAt = updatedAt
						}
					}

					if conclusionData, ok := previousTermReport.ReportData["conclusion"].(primitive.M); ok {
						if comment, ok := conclusionData["manager_comment"].(string); ok {
							managerCommentPreviousTerm.Conclusion = comment
						}
						if report, ok := conclusionData["teacher_report"].(string); ok {
							teacherReportPrevioiusTerm.Conclusion = report
						}
						if managerUpdatedAt, ok := conclusionData["manager_updated_at"].(string); ok {
							managerCommentPreviousTerm.ConclusionUpdatedAt = managerUpdatedAt
						}
						if updatedAt, ok := conclusionData["updated_at"].(string); ok {
							teacherReportPrevioiusTerm.ConclusionUpdatedAt = updatedAt
						}
					}
				}

			}
			reportRes = mapper.MapReportToResDTO(report, teacher, managerCommentPreviousTerm, teacherReportPrevioiusTerm, "")
		} else {
			reportRes = response.ReportResponse{}
		}

		res = append(res, &response.ClassroomReportResponse4Web{
			StudentID:   std.StudentID,
			StudentName: std.StudentName,
			Report:      reportRes,
		})
	}

	return res, nil
}
