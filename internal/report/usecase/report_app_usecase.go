package usecase

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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportAppUseCase interface {
	GetReport4App(ctx context.Context, req *request.GetReportRequest4App) (response.ReportResponse, error)
	UploadReport4App(ctx context.Context, req *request.UploadReport4AppRequest) error
	GetTeacherReportTasks4App(ctx context.Context) ([]response.GetTeacherReportTasksResponse4App, error)
}

type reportAppUseCase struct {
	reportRepo  repository.ReportRepository
	historyRepo repository.ReportHistoryRepository
	userGw      gateway.UserGateway
	classroomGw gateway.ClassroomGateway
	termGw      gateway.TermGateway
	mediaGw     gateway.MediaGateway
}

func NewReportAppUseCase(
	reportRepo repository.ReportRepository,
	historyRepo repository.ReportHistoryRepository,
	userGw gateway.UserGateway,
	classroomGw gateway.ClassroomGateway,
	termGw gateway.TermGateway,
	mediaGw gateway.MediaGateway,
) ReportAppUseCase {
	return &reportAppUseCase{
		reportRepo:  reportRepo,
		historyRepo: historyRepo,
		userGw:      userGw,
		classroomGw: classroomGw,
		termGw:      termGw,
		mediaGw:     mediaGw,
	}
}

func (u *reportAppUseCase) GetReport4App(ctx context.Context, req *request.GetReportRequest4App) (response.ReportResponse, error) {
	report, err := u.reportRepo.GetByStudentTopicTermAndLanguage(ctx, req.StudentID, req.TopicID, req.TermID, req.Language)
	if err != nil {
		return response.ReportResponse{}, err
	}
	if report == nil {
		return response.ReportResponse{}, errors.New("report not found")
	}

	// get student info
	student, _ := u.userGw.GetStudentInfo(ctx, report.StudentID)
	if student == nil {
		return response.ReportResponse{}, errors.New("student not found")
	}

	var managerCommentPreviousTerm response.ManagerCommentPreviousTerm
	var teacherReportPrevioiusTerm response.TeacherReportPreviousTerm
	previousTerm, _ := u.termGw.GetPreviousTerm(ctx, report.TermID, student.OrganizationID)
	if previousTerm != nil {
		previousTermReport, _ := u.reportRepo.GetByStudentTopicTermLanguageAndEditor(ctx, report.StudentID, report.TopicID, previousTerm.ID, report.Language, report.EditorID)
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

func (u *reportAppUseCase) UploadReport4App(ctx context.Context, req *request.UploadReport4AppRequest) error {
	// get student info
	student, _ := u.userGw.GetStudentInfo(ctx, req.StudentID)
	if student == nil {
		return errors.New("student not found")
	}

	// get teacher by usser id and organization if of student
	editorID := helper.GetUserID(ctx)
	teacher, _ := u.userGw.GetTeacherInfo(ctx, editorID, student.OrganizationID)
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
	err := u.reportRepo.CreateOrUpdateStudentView4App(ctx, report)
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

	if err := u.historyRepo.Create(ctx, history); err != nil {
		return err
	}

	return nil
}

func (u *reportAppUseCase) GetTeacherReportTasks4App(ctx context.Context) ([]response.GetTeacherReportTasksResponse4App, error) {
	userID := helper.GetUserID(ctx)
	// Lấy tất cả reports do editor này phụ trách
	reports, err := u.reportRepo.GetAllByEditorID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get reports failed: %w", err)
	}

	var results []response.GetTeacherReportTasksResponse4App

	for _, r := range reports {

		if r.ReportData != nil {
			reportData := helper.ToBsonM(r.ReportData)
			for key, val := range reportData {
				section := helper.ToBsonM(val)
				status, _ := section["status"].(string)

				if status == "teacher" || status == "empty" {
					termTitle := ""
					topicTitle := ""
					stdName := ""
					term, _ := u.termGw.GetTermByID(ctx, r.TermID)
					topic, _ := u.mediaGw.GetTopicByID(ctx, r.TopicID)
					student, _ := u.userGw.GetStudentInfo(ctx, r.StudentID)

					if term != nil {
						termTitle = term.Title
					}
					if topic != nil {
						topicTitle = topic.Title
					}
					if student != nil {
						stdName = student.Name
					}

					results = append(results, response.GetTeacherReportTasksResponse4App{
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
