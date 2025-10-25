package route

import (
	"report-service/internal/gateway"
	"report-service/internal/middleware"
	"report-service/internal/report/handler"

	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(r *gin.Engine, h *handler.ReportHandler, rh *handler.ReportHistoryHandler, rph *handler.ReportPlanTemplateHandler, userGw gateway.UserGateway) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured(userGw))
	{
		reportsAdmin := adminGroup.Group("/reports")
		{
			reportsAdmin.POST("", h.UploadReport4Web)
			reportsAdmin.POST("/get-report", h.GetReport4Web)
			reportsAdmin.GET("/overview", h.GetReportOverViewAllClassroom4Web)

			// report history
			reportsAdmin.GET("/histories", rh.GetByEditor4App)

			// plan template
			reportsClassroomAdmin := reportsAdmin.Group("/classrooms")
			{
				reportsClassroomAdmin.POST("/plan-templates", rph.UploadReportPlanTemplate)
				reportsClassroomAdmin.POST("", h.UploadClassroomReport4Web)
				reportsClassroomAdmin.POST("/get-report", h.GetClassroomReports4Web)
				reportsClassroomAdmin.POST("/templates/school/apply", h.ApplyTopicPlanTemplateIsSchool2Report)
				reportsClassroomAdmin.POST("/templates/classroom/apply", h.ApplyTopicPlanTemplateIsClassroom2Report)
				reportsClassroomAdmin.GET("/overview", h.GetReportOverViewByClassroom4Web)
			}
		}
	}

	// user routes
	userGroup := r.Group("/api/v1/user")
	userGroup.Use(middleware.Secured(userGw))
	{
		reportsUser := userGroup.Group("/reports")
		{
			reportsUser.POST("", h.UploadReport4App)
			reportsUser.POST("/get-report", h.GetReport4App)
			reportsUser.GET("/tasks", h.GetTeacherReportTasks4App)
			reportsUser.GET("/histories", rh.GetByEditor4App)
		}
	}
}
