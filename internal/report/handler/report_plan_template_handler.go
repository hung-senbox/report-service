package handler

import (
	"net/http"
	"report-service/helper"
	"report-service/internal/report/dto/request"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
)

type ReportPlanTemplateHandler struct {
	service service.ReportPlanTemplateService
}

func NewReportPlanTemplateHandler(s service.ReportPlanTemplateService) *ReportPlanTemplateHandler {
	return &ReportPlanTemplateHandler{service: s}
}

func (h *ReportPlanTemplateHandler) UploadReportPlanTemplate(c *gin.Context) {
	var req request.UploadReportPlanTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.Upload(c.Request.Context(), req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report plan template uploaded successfully", nil)
}
