package router

import (
	"report-service/internal/report/handler"
	"report-service/internal/report/repository"
	"report-service/internal/report/route"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(reportCollection *mongo.Collection, reportHistoryCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// Setup dependency injection
	reportRepo := repository.NewReportRepository(reportCollection)
	historyRepo := repository.NewReportHistoryRepository(reportHistoryCollection)

	// reoport
	reportService := service.NewReportService(reportRepo, historyRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// report history
	reportHistoryService := service.NewReportHistoryService(historyRepo)
	reportHistoryHandler := handler.NewReportHistoryHandler(reportHistoryService)

	// Register routes
	route.RegisterReportRoutes(r, reportHandler, reportHistoryHandler)
	return r
}
