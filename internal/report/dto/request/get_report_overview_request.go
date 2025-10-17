package request

type GetReportOverViewAllClassroomRequest struct {
	TermID string `json:"term_id" binding:"required"`
}
