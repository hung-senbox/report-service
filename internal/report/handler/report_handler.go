package handler

import (
	"net/http"
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

// UploadReport godoc
// @Summary Upload or update report
// @Description Upload report data by key
// @Tags Reports
// @Accept json
// @Produce json
// @Param request body request.UploadReportRequestDTO true "Report upload request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/admin/reports [post]
func (h *ReportHandler) UploadReport(c *gin.Context) {
	var req request.UploadReportRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.UploadReport(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "report uploaded successfully",
	})
}
