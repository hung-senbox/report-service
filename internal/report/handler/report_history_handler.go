package handler

import (
	"net/http"
	"report-service/helper"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
)

type ReportHistoryHandler struct {
	service service.ReportHistoryService
}

func NewReportHistoryHandler(s service.ReportHistoryService) *ReportHistoryHandler {
	return &ReportHistoryHandler{service: s}
}

func (h *ReportHistoryHandler) GetAllReportHistory(c *gin.Context) {
	reportHis, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report history retrieved successfully", reportHis)
}
