package request

type GetReportOverViewAllClassroomRequest struct {
	TermID string `json:"term_id" binding:"required"`
}

type GetReportOverViewByClassroomRequest struct {
	TermID      string `json:"term_id" binding:"required"`
	ClassroomID string `json:"classroom_id" binding:"required"`
}
