package handler

import (
	"errors"
	"net/http"
	"report-service/helper"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

func (h *ReportHandler) UploadReport(c *gin.Context) {
	var req request.UploadReportRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.UploadReport(c.Request.Context(), &req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report uploaded successfully", nil)
}

func (h *ReportHandler) GetAllReports(c *gin.Context) {
	reports, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Reports retrieved successfully", reports)
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	// get report by student_id, topic_id, term_id, language from search params
	studentID := c.Query("student_id")
	topicID := c.Query("topic_id")
	termID := c.Query("term_id")
	language := c.Query("language")

	if studentID == "" || topicID == "" || termID == "" || language == "" {
		helper.SendError(c, http.StatusBadRequest, nil, "student_id, topic_id, term_id and language are required")
		return
	}

	var req request.GetReportRequest
	req.StudentID = studentID
	req.TopicID = topicID
	req.TermID = termID
	req.Language = language

	report, err := h.service.Get4Report(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			helper.SendSuccess(c, http.StatusNotFound, "Report not found", nil)
			return
		}
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report retrieved successfully", report)
}
