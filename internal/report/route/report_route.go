package route

import (
	"report-service/internal/middleware"
	"report-service/internal/report/handler"

	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(r *gin.Engine, h *handler.ReportHandler, rh *handler.ReportHistoryHandler) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured())
	{
		reportsAdmin := adminGroup.Group("/reports")
		{
			reportsAdmin.POST("", h.UploadReport4Web)
			reportsAdmin.POST("/get-report", h.GetReport4Web)

			// report history
			reportsAdmin.GET("/history", rh.GetAllReportHistory)
		}
	}

	// user routes
	userGroup := r.Group("/api/v1/user")
	userGroup.Use(middleware.Secured())
	{
		reportsUser := userGroup.Group("/reports")
		{
			reportsUser.POST("", h.UploadReport4App)
			reportsUser.POST("/get-report", h.GetReport4App)
			reportsUser.GET("/tasks", h.GetTeacherReportTasks)
		}
	}
}
