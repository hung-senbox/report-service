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

func (h *ReportHandler) UploadReport4App(c *gin.Context) {
	var req request.UploadReport4AppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.UploadReport4App(c.Request.Context(), &req); err != nil {
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

func (h *ReportHandler) GetReport4App(c *gin.Context) {
	var req request.GetReportRequest4App
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	report, err := h.service.GetReport4App(c.Request.Context(), &req)
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

func (h *ReportHandler) GetReport4Web(c *gin.Context) {
	var req request.GetReportRequest4Web
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	report, err := h.service.GetReport4Web(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			helper.SendSuccess(c, http.StatusOK, "Report not found", nil)
			return
		}
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report retrieved successfully", report)
}

func (h *ReportHandler) GetTeacherReportTasks(c *gin.Context) {

	reports, err := h.service.GetTeacherReportTasks(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report tasks retrieved successfully", reports)
}

func (h *ReportHandler) UploadReport4Web(c *gin.Context) {
	var req request.UploadReport4AWebRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.UploadReport4Web(c.Request.Context(), &req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report uploaded successfully", nil)
}

func (h *ReportHandler) UploadClassroomReport(c *gin.Context) {
	var req request.UploadClassroomReport4WebRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.UploadClassroomReport(c.Request.Context(), req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report uploaded successfully", nil)
}

func (h *ReportHandler) GetClassroomReports4Web(c *gin.Context) {
	var req request.GetClassroomReportRequest4Web
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	reports, err := h.service.GetClassroomReports4Web(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			helper.SendSuccess(c, http.StatusOK, "Report not found", nil)
			return
		}
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report retrieved successfully", reports)
}

func (h *ReportHandler) ApplyTopicPlanTemplateIsSchool2Report(c *gin.Context) {
	var req request.ApplyTemplateIsSchoolToReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.ApplyTopicPlanTemplateIsSchool2Report(c.Request.Context(), req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report template applied successfully", nil)
}

func (h *ReportHandler) ApplyTopicPlanTemplateIsClassroom2Report(c *gin.Context) {
	var req request.ApplyTemplateIsClassroomToReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.ApplyTopicPlanTemplateIsClassroom2Report(c.Request.Context(), req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report template applied successfully", nil)
}

func (h *ReportHandler) GetReportOverViewAllClassroom(c *gin.Context) {
	termID := c.Query("term_id")
	if termID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("termID is required"), helper.ErrInvalidRequest)
		return
	}

	var req request.GetReportOverViewAllClassroomRequest
	req.TermID = termID

	reports, err := h.service.GetReportOverViewAllClassroom(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			helper.SendSuccess(c, http.StatusOK, "Report not found", nil)
			return
		}
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Report retrieved successfully", reports)
}
