package mapper

import (
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

	// --- Thêm "sub_title" nếu chưa có ---
	if _, ok := report.ReportData["sub_title"]; !ok {
		report.ReportData["sub_title"] = bson.M{
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

	// --- Đảm bảo phần "introduction" có các key manager_* ---
	if introductionData, ok := report.ReportData["introduction"].(bson.M); ok {
		if _, ok := introductionData["manager_note"]; !ok {
			introductionData["manager_note"] = ""
		}
		if _, ok := introductionData["manager_comment"]; !ok {
			introductionData["manager_comment"] = ""
		}
		if _, ok := introductionData["manager_updated_at"]; !ok {
			introductionData["manager_updated_at"] = ""
		}
		report.ReportData["introduction"] = introductionData
	} else {
		report.ReportData["introduction"] = bson.M{
			"manager_note":       "",
			"manager_comment":    "",
			"manager_updated_at": "",
		}
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

	return response.ReportResponse{
		ID:                         report.ID.Hex(),
		StudentID:                  report.StudentID,
		TopicID:                    report.TopicID,
		TermID:                     report.TermID,
		Language:                   report.Language,
		Status:                     report.Status,
		ReportData:                 report.ReportData,
		CreatedAt:                  report.CreatedAt,
		Editor:                     teacherEditor,
		ManagerCommentPreviousTerm: managerCmPrevious,
		TeacherReportPreviousTerm:  teacherRpPrevious,
	}
}

// MapReportListToResDTO maps slice of model.Report to slice of ReportResponse
func MapReportListToResDTO(reports []*model.Report) []response.ReportResponse {
	result := make([]response.ReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, MapReportToResDTO(r, nil, response.ManagerCommentPreviousTerm{}, response.TeacherReportPreviousTerm{}))
	}
	return result
}
