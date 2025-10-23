package mapper

import (
	"encoding/json"
	"fmt"
	"report-service/helper"
	gw_response "report-service/internal/gateway/dto/response"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"

	"go.mongodb.org/mongo-driver/bson"
)

func MapReportToResDTO(
	report *model.Report,
	teacher *gw_response.TeacherResponse,
	managerCmPrevious response.ManagerCommentPreviousTerm,
	teacherRpPrevious response.TeacherReportPreviousTerm,
	latestDataTermId string,
) response.ReportResponse {

	// --- Đảm bảo report.ReportData không nil ---
	if report.ReportData == nil {
		report.ReportData = bson.M{}
	}

	// --- Thêm "title" nếu chưa có ---
	if _, ok := report.ReportData["title"]; !ok {
		report.ReportData["title"] = bson.M{
			"content":    "",
			"updated_at": "",
		}
	}

	// --- Thêm "goal" nếu chưa có ---
	if _, ok := report.ReportData["goal"]; !ok {
		report.ReportData["goal"] = bson.M{
			"content":    "",
			"updated_at": "",
		}
	}

	// --- Thêm "curriculum_area" nếu chưa có ---
	if _, ok := report.ReportData["curriculum_area"]; !ok {
		report.ReportData["curriculum_area"] = bson.M{
			"content":    "",
			"updated_at": "",
		}
	}

	// --- Đảm bảo phần "now" có các key manager_* ---
	if nowData, ok := report.ReportData["now"].(bson.M); ok {
		if _, ok := nowData["manager_note"]; !ok {
			nowData["manager_note"] = ""
		}
		if _, ok := nowData["manager_comment"]; !ok {
			nowData["manager_comment"] = ""
		}
		if _, ok := nowData["manager_updated_at"]; !ok {
			nowData["manager_updated_at"] = ""
		}
		report.ReportData["now"] = nowData
	}

	// --- Đảm bảo phần "before" có các key manager_* ---
	if beforeData, ok := report.ReportData["before"].(bson.M); ok {
		if _, ok := beforeData["manager_note"]; !ok {
			beforeData["manager_note"] = ""
		}
		if _, ok := beforeData["manager_comment"]; !ok {
			beforeData["manager_comment"] = ""
		}
		if _, ok := beforeData["manager_updated_at"]; !ok {
			beforeData["manager_updated_at"] = ""
		}
		report.ReportData["before"] = beforeData
	}

	// --- Đảm bảo phần "conclusion" có các key manager_* ---
	if conclusionData, ok := report.ReportData["conclusion"].(bson.M); ok {
		if _, ok := conclusionData["manager_note"]; !ok {
			conclusionData["manager_note"] = ""
		}
		if _, ok := conclusionData["manager_comment"]; !ok {
			conclusionData["manager_comment"] = ""
		}
		if _, ok := conclusionData["manager_updated_at"]; !ok {
			conclusionData["manager_updated_at"] = ""
		}
		report.ReportData["conclusion"] = conclusionData
	}

	// --- Đảm bảo phần "note" có các key manager_* ---
	if noteData, ok := report.ReportData["note"].(bson.M); ok {
		if _, ok := noteData["manager_note"]; !ok {
			noteData["manager_note"] = ""
		}
		if _, ok := noteData["manager_comment"]; !ok {
			noteData["manager_comment"] = ""
		}
		if _, ok := noteData["manager_updated_at"]; !ok {
			noteData["manager_updated_at"] = ""
		}
		report.ReportData["note"] = noteData
	}

	// --- Map editor ---
	var teacherEditor gw_response.TeacherResponse
	if teacher != nil {
		teacherEditor = *teacher
	}

	editing := true
	if report.Editing != nil {
		editing = *report.Editing
	}

	// --- Thêm latest_update_time cho từng section ---
	sections := []string{"now", "before", "conclusion", "note"}

	for _, section := range sections {
		if sectionData, ok := report.ReportData[section].(bson.M); ok {
			updatedAt, _ := sectionData["updated_at"].(string)
			managerUpdatedAt, _ := sectionData["manager_updated_at"].(string)

			latestTimeStr := helper.GetLatestTimeStr(updatedAt, managerUpdatedAt)
			sectionData["latest_update_time"] = latestTimeStr

			report.ReportData[section] = sectionData
		}
	}

	return response.ReportResponse{
		ID:                         report.ID.Hex(),
		StudentID:                  report.StudentID,
		TopicID:                    report.TopicID,
		TermID:                     report.TermID,
		Language:                   report.Language,
		Status:                     report.Status,
		Editing:                    editing,
		ReportData:                 report.ReportData,
		CreatedAt:                  report.CreatedAt,
		Editor:                     teacherEditor,
		ManagerCommentPreviousTerm: managerCmPrevious,
		TeacherReportPreviousTerm:  teacherRpPrevious,
		LatestDataTermID:           latestDataTermId,
	}
}

// MapReportListToResDTO maps slice of model.Report to slice of ReportResponse
func MapReportListToResDTO(reports []*model.Report) []response.ReportResponse {
	result := make([]response.ReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, MapReportToResDTO(r, nil, response.ManagerCommentPreviousTerm{}, response.TeacherReportPreviousTerm{}, ""))
	}
	return result
}

func MapReportToStruct(report *model.Report) (*model.Reportstruct, error) {
	if report == nil {
		return nil, fmt.Errorf("report is nil")
	}

	// Parse report_data (bson.M → model.ReportData)
	var rd model.ReportData
	dataBytes, err := json.Marshal(report.ReportData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report.ReportData: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &rd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to model.ReportData: %w", err)
	}

	// Map sang Reportstruct
	res := &model.Reportstruct{
		ID:         report.ID.Hex(),
		StudentID:  report.StudentID,
		TopicID:    report.TopicID,
		TermID:     report.TermID,
		EditorID:   report.EditorID,
		Language:   report.Language,
		Status:     report.Status,
		ReportData: rd,
		CreatedAt:  report.CreatedAt,
	}

	return res, nil
}

func MapReportsToStruct(reports []*model.Report) ([]*model.Reportstruct, error) {
	if len(reports) == 0 {
		return nil, nil
	}

	statusMap := map[string]int{
		"Empty":    0,
		"Teacher":  10,
		"Manager":  15,
		"Done":     20,
		"Approved": 25,
	}

	var result []*model.Reportstruct

	for _, report := range reports {
		if report == nil {
			continue
		}

		dataBytes, err := json.Marshal(report.ReportData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal report.ReportData: %w", err)
		}

		var rd model.ReportData
		if err := json.Unmarshal(dataBytes, &rd); err != nil {
			return nil, fmt.Errorf("failed to unmarshal to model.ReportData: %w", err)
		}

		// Tính progress
		progress := 0
		progress += statusMap[rd.Before.Status]
		progress += statusMap[rd.Now.Status]
		progress += statusMap[rd.Conclusion.Status]

		res := &model.Reportstruct{
			ID:         report.ID.Hex(),
			StudentID:  report.StudentID,
			TopicID:    report.TopicID,
			TermID:     report.TermID,
			EditorID:   report.EditorID,
			Language:   report.Language,
			Status:     report.Status,
			ReportData: rd,
			CreatedAt:  report.CreatedAt,
			Progress:   progress,
		}

		result = append(result, res)
	}

	return result, nil
}

func MapReport2Print(report *model.Report) *response.GetReport2Print {
	reportData := helper.ToBsonM(report.ReportData)

	getContent := func(section string) string {
		if sec, ok := reportData[section].(bson.M); ok {
			if content, ok := sec["teacher_report"].(string); ok {
				return content
			}
		}
		return ""
	}

	return &response.GetReport2Print{
		Before:     getContent("before"),
		Now:        getContent("now"),
		Conclusion: getContent("conclusion"),
	}
}
