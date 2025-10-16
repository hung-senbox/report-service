package request

type GetReportOverViewRequest struct {
	TermID        string `json:"term_id" binding:"required"`
	ClassroomID   string `json:"classroom_id" binding:"required"`
	UniqueLangKey string `json:"unique_lang_key" binding:"required"`
}
