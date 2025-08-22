package handler

import (
	"net/http"
	"report-service/helper"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
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
