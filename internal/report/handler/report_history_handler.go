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

func (h *ReportHistoryHandler) GetByEditor4App(c *gin.Context) {
	histories, err := h.service.GetByEditor4App(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "History retrieved successfully", histories)
}
