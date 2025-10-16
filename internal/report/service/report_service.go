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
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService interface {
	GetAll(ctx context.Context) ([]response.ReportResponse, error)
	UploadReport4App(ctx context.Context, req *request.UploadReport4AppRequest) error
	UploadReport4Web(ctx context.Context, req *request.UploadReport4AWebRequest) error
	GetReport4App(ctx context.Context, req *request.GetReportRequest4App) (response.ReportResponse, error)
	GetReport4Web(ctx context.Context, req *request.GetReportRequest4Web) (response.ReportResponse, error)
	GetTeacherReportTasks(ctx context.Context) ([]response.GetTeacherReportTasksResponse, error)
	UploadClassroomReport(ctx context.Context, req request.UploadClassroomReport4WebRequest) error
	GetClassroomReports4Web(ctx context.Context, req request.GetClassroomReportRequest4Web) (*response.GetClassroomReportResponse4Web, error)
	ApplyTopicPlanTemplateIsSchool2Report(ctx context.Context, req request.ApplyTemplateIsSchoolToReportRequest) error
	ApplyTopicPlanTemplateIsClassroom2Report(ctx context.Context, req request.ApplyTemplateIsClassroomToReportRequest) error
}

type reportService struct {
	userGateway            gateway.UserGateway
	termGateway            gateway.TermGateway
	mediaGateway           gateway.MediaGateway
	classroomGateway       gateway.ClassroomGateway
	repo                   repository.ReportRepository
	historyRepo            repository.ReportHistoryRepository
	reportPlanTemplateRepo repository.ReportPlanTemplateRepositopry
}

func NewReportService(
	userGateway gateway.UserGateway,
	termGateway gateway.TermGateway,
	mediaGateway gateway.MediaGateway,
	classroomGateway gateway.ClassroomGateway,
	repo repository.ReportRepository,
	historyRepo repository.ReportHistoryRepository,
	reportPlanTemplateRepo repository.ReportPlanTemplateRepositopry,
) ReportService {
	return &reportService{
		userGateway:            userGateway,
		termGateway:            termGateway,
		mediaGateway:           mediaGateway,
		classroomGateway:       classroomGateway,
		repo:                   repo,
		historyRepo:            historyRepo,
		reportPlanTemplateRepo: reportPlanTemplateRepo,
	}
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

func (s *reportService) GetClassroomReports4Web(ctx context.Context, req request.GetClassroomReportRequest4Web) (*response.GetClassroomReportResponse4Web, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user")
	}

	if currentUser == nil {
		return nil, fmt.Errorf("current user not found")
	}
	res := &response.GetClassroomReportResponse4Web{}
	var reports = make([]response.ClassroomReportResponse4Web, 0)

	// Lấy template
	reportTemplateSchool, _ := s.reportPlanTemplateRepo.GetSchoolTemplate(ctx, req.TermID, req.TopicID, req.UniqueLangKey, currentUser.OrganizationAdmin.ID)
	reportTemplateClassroom, _ := s.reportPlanTemplateRepo.GetClassroomTemplate(ctx, req.TermID, req.TopicID, req.UniqueLangKey, req.ClassroomID, currentUser.OrganizationAdmin.ID)

	var schoolTemplate model.Template
	var classroomTemplate model.Template
	if reportTemplateSchool != nil {
		schoolTemplate = model.Template{
			Title:          reportTemplateSchool.Template.Title,
			Introduction:   reportTemplateSchool.Template.Introduction,
			CurriculumArea: reportTemplateSchool.Template.CurriculumArea,
		}
	}
	if reportTemplateClassroom != nil {
		classroomTemplate = model.Template{
			Title:          reportTemplateClassroom.Template.Title,
			Introduction:   reportTemplateClassroom.Template.Introduction,
			CurriculumArea: reportTemplateClassroom.Template.CurriculumArea,
		}
	}

	res.SchoolTemplate = schoolTemplate
	res.ClassroomTempate = classroomTemplate

	// Lấy danh sách học sinh trong lớp
	students, _ := s.classroomGateway.GetStudents4ClassroomReport(ctx, req.TermID, req.ClassroomID, req.TeacherID)
	if len(students) == 0 {
		return res, nil
	}

	// Lấy thông tin editor
	editor, err := s.userGateway.GetUserByTeacher(ctx, req.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get editor (teacher): %w", err)
	}

	// Lấy thông tin teacher
	teacher, _ := s.userGateway.GetTeacherInfo(ctx, editor.ID, students[0].OrganizationID)

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
				previousTermReport, _ := s.repo.GetByStudentTopicTermLanguageAndEditor(
					ctx,
					report.StudentID,
					report.TopicID,
					previousTerm.ID,
					report.Language,
					report.EditorID,
				)
				if previousTermReport != nil {
					managerCommentPreviousTerm.TermTitle = previousTerm.Title
					teacherReportPrevioiusTerm.TermTitle = previousTerm.Title

					if nowData, ok := previousTermReport.ReportData["now"].(primitive.M); ok {
						if comment, ok := nowData["manager_comment"].(string); ok {
							managerCommentPreviousTerm.Now = comment
						}
						if r, ok := nowData["teacher_report"].(string); ok {
							teacherReportPrevioiusTerm.Now = r
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
						if r, ok := conclusionData["teacher_report"].(string); ok {
							teacherReportPrevioiusTerm.Conclusion = r
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

			reportRes = mapper.MapReportToResDTO(
				report,
				teacher,
				managerCommentPreviousTerm,
				teacherReportPrevioiusTerm,
				"",
			)
		} else {
			reportRes = response.ReportResponse{}
		}

		reports = append(reports, response.ClassroomReportResponse4Web{
			StudentID:   std.StudentID,
			StudentName: std.StudentName,
			Report:      reportRes,
		})
	}

	res.Reports = reports

	return res, nil
}

func (s *reportService) ApplyTopicPlanTemplateIsSchool2Report(ctx context.Context, req request.ApplyTemplateIsSchoolToReportRequest) error {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current user")
	}

	if currentUser.IsSuperAdmin {
		return errors.New("super admin cannot apply template to report")
	}

	// tao report plan template
	rpt := &model.ReportPlanTemplate{
		OrganizationID: currentUser.OrganizationAdmin.ID,
		TopicID:        req.TopicID,
		TermID:         req.TermID,
		Language:       req.UniqueLangKey,
		IsSchool:       true,
		Template: model.Template{
			Title:          req.Title,
			Introduction:   req.Introduction,
			CurriculumArea: req.CurriculumArea,
		},
	}
	if err := s.reportPlanTemplateRepo.CreateOrUpdate(ctx, rpt); err != nil {
		return fmt.Errorf("failed to create report plan template: %w", err)
	}

	// Lấy danh sách reports theo term, topic, language
	reports, err := s.repo.GetTopicsByTermTopicLanguage(ctx, req.TermID, req.TopicID, req.UniqueLangKey)
	if err != nil {
		return fmt.Errorf("failed to get reports")
	}
	if len(reports) == 0 {
		return errors.New("no reports found to apply template")
	}

	// Áp dụng template cho từng report
	for _, report := range reports {
		report.ReportData = toBsonM(report.ReportData)

		// --- Chuẩn bị dữ liệu template ---
		title := toBsonM(report.ReportData["title"])
		title["content"] = req.Title
		report.ReportData["title"] = title

		intro := toBsonM(report.ReportData["introduction"])
		intro["content"] = req.Introduction
		report.ReportData["introduction"] = intro

		cur := toBsonM(report.ReportData["curriculum_area"])
		cur["content"] = req.CurriculumArea
		report.ReportData["curriculum_area"] = cur

		// --- Gọi repository update ---
		if err := s.repo.ApplyTopicPlanTemplate(ctx, report); err != nil {
			// nếu không tìm thấy report thì bỏ qua
			if strings.Contains(err.Error(), "report not found") {
				continue
			}
			return fmt.Errorf("failed to apply template to report")
		}
	}

	return nil
}

func (s *reportService) ApplyTopicPlanTemplateIsClassroom2Report(ctx context.Context, req request.ApplyTemplateIsClassroomToReportRequest) error {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current user")
	}

	if currentUser.IsSuperAdmin {
		return errors.New("super admin cannot apply template to report")
	}

	// tao report plan template
	rpt := &model.ReportPlanTemplate{
		OrganizationID: currentUser.OrganizationAdmin.ID,
		TopicID:        req.TopicID,
		TermID:         req.TermID,
		Language:       req.UniqueLangKey,
		ClassroomID:    req.ClassroomID,
		IsSchool:       false,
		Template: model.Template{
			Title:          req.Title,
			Introduction:   req.Introduction,
			CurriculumArea: req.CurriculumArea,
		},
	}
	if err := s.reportPlanTemplateRepo.CreateOrUpdate(ctx, rpt); err != nil {
		return fmt.Errorf("failed to create report plan template")
	}

	// get students by classroom id from gw
	students, err := s.classroomGateway.GetStudentsByClassroomID(ctx, req.ClassroomID, req.TermID)
	if err != nil {
		return fmt.Errorf("failed to get students")
	}
	if len(students) == 0 {
		return fmt.Errorf("students not found")
	}

	var reports []*model.Report

	for _, student := range students {
		report, _ := s.repo.GetByStudentTopicTermAndLanguage(
			ctx,
			student.StudentID,
			req.TopicID,
			req.TermID,
			req.UniqueLangKey,
		)
		if report != nil {
			reports = append(reports, report)
		}
	}

	// Áp dụng template cho từng report
	for _, report := range reports {
		report.ReportData = toBsonM(report.ReportData)

		// --- Chuẩn bị dữ liệu template ---
		title := toBsonM(report.ReportData["title"])
		title["content"] = req.Title
		report.ReportData["title"] = title

		intro := toBsonM(report.ReportData["introduction"])
		intro["content"] = req.Introduction
		report.ReportData["introduction"] = intro

		cur := toBsonM(report.ReportData["curriculum_area"])
		cur["content"] = req.CurriculumArea
		report.ReportData["curriculum_area"] = cur

		// --- Gọi repository update ---
		if err := s.repo.ApplyTopicPlanTemplate(ctx, report); err != nil {
			// nếu không tìm thấy report thì bỏ qua
			if strings.Contains(err.Error(), "report not found") {
				continue
			}
			return fmt.Errorf("failed to apply template to report %s: %w", report.ID.Hex(), err)
		}
	}

	return nil
}

func (s *reportService) GetReportOverView(ctx context.Context, req request.GetReportOverViewRequest) error {
	return nil
}
