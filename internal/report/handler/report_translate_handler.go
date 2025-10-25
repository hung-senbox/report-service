package handler

import (
	"net/http"
	"report-service/helper"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
)

type ReportTranslateHandler struct {
	ReportTranslateService service.ReportTranslateService
}

func NewReportTranslateHandler(service service.ReportTranslateService) *ReportTranslateHandler {
	return &ReportTranslateHandler{
		ReportTranslateService: service,
	}
}

func (h *ReportTranslateHandler) UploadReportTranslate4Web(c *gin.Context) {

	var req request.UploadReportTranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.ReportTranslateService.UploadReportTranslate4Web(c.Request.Context(), req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report uploaded successfully", nil)

}

func (h *ReportTranslateHandler) GetReportTranslate4WebByTopicAndLang(c *gin.Context) {

	studentID := c.Query("student_id")
	topicID := c.Query("topic_id")
	termID := c.Query("term_id")
	lang := c.Query("lang_key")

	report, err := h.ReportTranslateService.GetReportTranslate4WebByTopicAndLang(c.Request.Context(), studentID, topicID, termID, lang)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report retrieved successfully", report)

}

func (h *ReportTranslateHandler) GetReportTranslate4WebByReport(c *gin.Context) {

	studentID := c.Query("student_id")
	lang := c.Query("lang_key")
	termID := c.Query("term_id")
	
	report, err := h.ReportTranslateService.GetReportTranslate4WebByReport(c.Request.Context(), studentID, termID, lang)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report retrieved successfully", report)

}
