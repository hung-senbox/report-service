package route

import (
	"report-service/internal/middleware"
	"report-service/internal/report/handler"

	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(r *gin.Engine, h *handler.ReportHandler) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured(), middleware.RequireAdmin())
	{
		reportsAdmin := adminGroup.Group("/reports")
		{
			reportsAdmin.POST("", h.UploadReport)
			reportsAdmin.GET("", h.GetAllReports)
		}
	}

	// user routes
	userGroup := r.Group("/api/v1/user")
	userGroup.Use(middleware.Secured())
	{
		reportsUser := userGroup.Group("/reports")
		{
			reportsUser.POST("", h.UploadReport)
			reportsUser.GET("", h.GetAllReports)
		}
	}
}
