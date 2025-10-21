package service

import (
	"context"
	"errors"
	"fmt"
	"report-service/helper"
	"report-service/internal/gateway"
	dto "report-service/internal/gateway/dto/response"
	mockdata "report-service/internal/gateway/mock_data"
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
	GetReportOverViewAllClassroom(ctx context.Context, req request.GetReportOverViewAllClassroomRequest) (*response.GetReportOverviewAllClassroomResponse, error)
}

type reportService struct {
	userGateway            gateway.UserGateway
	termGateway            gateway.TermGateway
	mediaGateway           gateway.MediaGateway
	classroomGateway       gateway.ClassroomGateway
	repo                   repository.ReportRepository
	historyRepo            repository.ReportHistoryRepository
	reportPlanTemplateRepo repository.ReportPlanTemplateRepository
}

func NewReportService(
	userGateway gateway.UserGateway,
	termGateway gateway.TermGateway,
	mediaGateway gateway.MediaGateway,
	classroomGateway gateway.ClassroomGateway,
	repo repository.ReportRepository,
	historyRepo repository.ReportHistoryRepository,
	reportPlanTemplateRepo repository.ReportPlanTemplateRepository,
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
		Editing:    &req.Editing,
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
	// get editor from teacher id
	user, _ := s.userGateway.GetUserByTeacher(ctx, req.TeacherID)
	if user == nil {
		return errors.New("upload report failed, teacher not found")
	}
	reportExist, _ := s.repo.GetByStudentTopicTermLanguageAndEditor(ctx, req.StudentID, req.TopicID, req.TermID, req.UniqueLangKey, user.ID)
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

	// Get students assigned to classroom
	assigned, _ := s.classroomGateway.GetClassroomAssignedTemplate(ctx, req.TermID, req.ClassroomID)
	if assigned == nil {
		return errors.New("assignment template not found")
	}
	if len(assigned.Students) == 0 {
		return errors.New("students assignment template not found")
	}

	var reports []*model.Report

	for _, student := range assigned.Students {
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

type topicAgg struct {
	Status response.AllClassroomTopicStatus
	Count  int
}

func (s *reportService) GetReportOverViewAllClassroom(ctx context.Context, req request.GetReportOverViewAllClassroomRequest) (*response.GetReportOverviewAllClassroomResponse, error) {

	var res response.GetReportOverviewAllClassroomResponse
	res.Reports = make([]response.AllClassroomReport, 0)

	allClassroomMockData := mockdata.FakeAllClassroomAssignTemplate()

	for _, class := range allClassroomMockData {
		topicsByClass := make(map[string]topicAgg)

		for _, assign := range class.AssignTemplates {
			editor, err := s.userGateway.GetUserByTeacher(ctx, assign.TeacherID)
			if err != nil || editor == nil {
				continue
			}

			reports, err := s.repo.GetByEditorIDAndStudentIDAndTermID(ctx, editor.ID, assign.StudentID, req.TermID)
			if err != nil || len(reports) == 0 {
				continue
			}

			// Gọi hàm phụ để xử lý gom dữ liệu
			classTopics, err := s.aggregateTopicsByClassroom(ctx, reports)
			if err != nil {
				continue
			}

			// Gộp kết quả vào topicsByClass tổng
			for topicID, agg := range classTopics {
				if existing, ok := topicsByClass[topicID]; ok {
					// Gộp dữ liệu trung bình giữa các nhóm
					newCount := existing.Count + agg.Count
					existing.Status.Before = (existing.Status.Before*float32(existing.Count) + agg.Status.Before*float32(agg.Count)) / float32(newCount)
					existing.Status.Now = (existing.Status.Now*float32(existing.Count) + agg.Status.Now*float32(agg.Count)) / float32(newCount)
					existing.Status.Conclusion = (existing.Status.Conclusion*float32(existing.Count) + agg.Status.Conclusion*float32(agg.Count)) / float32(newCount)
					existing.Status.MainStatus = (existing.Status.MainStatus*float32(existing.Count) + agg.Status.MainStatus*float32(agg.Count)) / float32(newCount)
					existing.Status.MainPercentage = (existing.Status.MainPercentage*float32(existing.Count) + agg.Status.MainPercentage*float32(agg.Count)) / float32(newCount)
					existing.Count = newCount
					topicsByClass[topicID] = existing
				} else {
					topicsByClass[topicID] = agg
				}
			}
		}

		// map → slice
		topicsSlice := make([]response.AllClassroomTopicStatus, 0, len(topicsByClass))
		for _, topic := range topicsByClass {
			topicsSlice = append(topicsSlice, topic.Status)
		}

		res.Reports = append(res.Reports, response.AllClassroomReport{
			ClassName: class.ClassroomName,
			DOB:       "EMPTY",
			Age:       0,
			Class:     0.0,
			Topics:    topicsSlice,
		})
	}

	return &res, nil
}

func (s *reportService) aggregateTopicsByClassroom(ctx context.Context, reports []*model.Report) (map[string]topicAgg, error) {

	topicsByClass := make(map[string]topicAgg)

	for _, r := range reports {
		reportStruct, err := mapper.MapReportToStruct(r)
		if err != nil || reportStruct == nil {
			continue
		}

		before := constants.MapStatusValue(reportStruct.ReportData.Before.Status)
		now := constants.MapStatusValue(reportStruct.ReportData.Now.Status)
		conclusion := constants.MapStatusValue(reportStruct.ReportData.Conclusion.Status)
		mainStatus := constants.MapStatusValue(reportStruct.Status)
		mainPercentage := before + now + conclusion + mainStatus

		if existing, ok := topicsByClass[r.TopicID]; ok {
			newCount := existing.Count + 1

			existing.Status.Before = (existing.Status.Before*float32(existing.Count) + before) / float32(newCount)
			existing.Status.Now = (existing.Status.Now*float32(existing.Count) + now) / float32(newCount)
			existing.Status.Conclusion = (existing.Status.Conclusion*float32(existing.Count) + conclusion) / float32(newCount)
			existing.Status.MainStatus = (existing.Status.MainStatus*float32(existing.Count) + mainStatus) / float32(newCount)
			existing.Status.MainPercentage = (existing.Status.MainPercentage*float32(existing.Count) + mainPercentage) / float32(newCount)

			existing.Count = newCount
			topicsByClass[r.TopicID] = existing

		} else {
			topic, _ := s.mediaGateway.GetTopicByID(ctx, r.TopicID)
			topicTitle := ""
			topicMainImageUrl := ""

			if topic != nil {
				topicTitle = topic.Title
				topicMainImageUrl = topic.MainImageUrl
			}

			topicsByClass[r.TopicID] = topicAgg{
				Status: response.AllClassroomTopicStatus{
					TopicID:           r.TopicID,
					TopicTitle:        topicTitle,
					TopicMainImageUrl: topicMainImageUrl,
					Before:            before,
					Now:               now,
					Conclusion:        conclusion,
					MainPercentage:    mainPercentage,
					MainStatus:        mainStatus,
				},
				Count: 1,
			}
		}
	}

	return topicsByClass, nil
}

func aggregateReportsSummary(reports []response.ReportResponse) response.ReportSummary {
	var total float32
	var beforeSum, nowSum, conclusionSum float32

	for _, r := range reports {
		rd := r.ReportData
		if rd == nil {
			continue
		}

		beforeData := toBsonM(rd["before"])
		nowData := toBsonM(rd["now"])
		conclusionData := toBsonM(rd["conclusion"])

		if status, ok := beforeData["status"].(string); ok {
			beforeSum += constants.MapStatusValue(status)
		}
		if status, ok := nowData["status"].(string); ok {
			nowSum += constants.MapStatusValue(status)
		}
		if status, ok := conclusionData["status"].(string); ok {
			conclusionSum += constants.MapStatusValue(status)
		}

		total++
	}

	if total == 0 {
		return response.ReportSummary{}
	}

	beforeAvg := beforeSum / total
	nowAvg := nowSum / total
	conclusionAvg := conclusionSum / total

	mainPercentage := beforeAvg + nowAvg + conclusionAvg

	return response.ReportSummary{
		MainPercentage: mainPercentage,
		Status:         mainPercentage,
		Before:         beforeAvg,
		Now:            nowAvg,
		Conclusion:     conclusionAvg,
	}
}

func (s *reportService) GetClassroomReports4Web(ctx context.Context, req request.GetClassroomReportRequest4Web) (*response.GetClassroomReportResponse4Web, error) {

	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil || currentUser == nil {
		return nil, fmt.Errorf("failed to get current user")
	}

	res := &response.GetClassroomReportResponse4Web{}

	// Load school & classroom templates.
	res.SchoolTemplate = s.getTemplateIfExists(
		func() (*model.ReportPlanTemplate, error) {
			return s.reportPlanTemplateRepo.GetSchoolTemplate(ctx,
				req.TermID, req.TopicID, req.UniqueLangKey, currentUser.OrganizationAdmin.ID)
		})

	res.ClassroomTempate = s.getTemplateIfExists(
		func() (*model.ReportPlanTemplate, error) {
			return s.reportPlanTemplateRepo.GetClassroomTemplate(ctx,
				req.TermID, req.TopicID, req.UniqueLangKey, req.ClassroomID, currentUser.OrganizationAdmin.ID)
		})

	// Get students assigned to classroom
	assigned, _ := s.classroomGateway.GetClassroomAssignedTemplate(ctx, req.TermID, req.ClassroomID)
	if assigned == nil {
		return res, nil
	}
	if len(assigned.Students) == 0 {
		return res, nil
	}

	// Build student reports
	for _, std := range assigned.Students {
		report := s.getStudentReport(ctx, req, std)
		if report.ID != "" {
			reports := response.ClassroomReportResponse4Web{
				Student: response.StudentReportClassroom{
					StudentID:     std.StudentID,
					StudentName:   std.StudentName,
					AvatarMainUrl: std.Avatar.ImageUrl,
				},
				Teacher: response.TeacherReportClassroom{
					TeacherID:     report.Editor.ID,
					TeacherName:   report.Editor.Name,
					AvatarMainUrl: report.Editor.Avatar.ImageUrl,
				},
				Report: report,
			}
			res.Reports = append(res.Reports, reports)
		}
	}

	// tinh main percentage
	var reportList []response.ReportResponse
	for _, r := range res.Reports {
		reportList = append(reportList, r.Report)
	}

	summary := aggregateReportsSummary(reportList)
	res.MainPercentage = summary.MainPercentage

	return res, nil
}

func (s *reportService) getTemplateIfExists(getter func() (*model.ReportPlanTemplate, error)) model.Template {
	t, err := getter()
	if err != nil || t == nil {
		return model.Template{}
	}
	return model.Template{
		Title:          t.Template.Title,
		Introduction:   t.Template.Introduction,
		CurriculumArea: t.Template.CurriculumArea,
	}
}

func (s *reportService) getStudentReport(ctx context.Context, req request.GetClassroomReportRequest4Web, std dto.StudentTemplate) response.ReportResponse {

	report, _ := s.repo.GetByStudentTopicTermAndLanguage(
		ctx,
		std.StudentID,
		req.TopicID,
		req.TermID,
		req.UniqueLangKey,
	)

	if report == nil {
		return response.ReportResponse{}
	}

	managerPrev, teacherPrev := s.getPreviousTermReports(ctx, report, std.OrganizationID)

	// get teacher info
	techerInfo, _ := s.userGateway.GetTeacherInfo(ctx, report.EditorID, std.OrganizationID)

	return mapper.MapReportToResDTO(
		report,
		techerInfo,
		managerPrev,
		teacherPrev,
		"",
	)
}

func (s *reportService) getPreviousTermReports(ctx context.Context, currentReport *model.Report, orgID string) (response.ManagerCommentPreviousTerm, response.TeacherReportPreviousTerm) {

	var mgr response.ManagerCommentPreviousTerm
	var tch response.TeacherReportPreviousTerm

	prevTerm, _ := s.termGateway.GetPreviousTerm(ctx, currentReport.TermID, orgID)
	if prevTerm == nil {
		return mgr, tch
	}

	prevReport, _ := s.repo.GetByStudentTopicTermLanguageAndEditor(
		ctx,
		currentReport.StudentID,
		currentReport.TopicID,
		prevTerm.ID,
		currentReport.Language,
		currentReport.EditorID,
	)
	if prevReport == nil {
		return mgr, tch
	}

	mgr.TermTitle = prevTerm.Title
	tch.TermTitle = prevTerm.Title

	parseTermData := func(section string, dstMgr *response.ManagerCommentPreviousTerm, dstTch *response.TeacherReportPreviousTerm) {
		if data, ok := prevReport.ReportData[section].(primitive.M); ok {
			if v, ok := data["manager_comment"].(string); ok {
				if section == "now" {
					dstMgr.Now = v
				} else {
					dstMgr.Conclusion = v
				}
			}
			if v, ok := data["teacher_report"].(string); ok {
				if section == "now" {
					dstTch.Now = v
				} else {
					dstTch.Conclusion = v
				}
			}
			if v, ok := data["manager_updated_at"].(string); ok {
				if section == "now" {
					dstMgr.NowUpdatedAt = v
				} else {
					dstMgr.ConclusionUpdatedAt = v
				}
			}
			if v, ok := data["updated_at"].(string); ok {
				if section == "now" {
					dstTch.NowUpdatedAt = v
				} else {
					dstTch.ConclusionUpdatedAt = v
				}
			}
		}
	}

	parseTermData("now", &mgr, &tch)
	parseTermData("conclusion", &mgr, &tch)

	return mgr, tch
}
