package router

import (
	"report-service/internal/gateway"
	"report-service/internal/report/handler"
	"report-service/internal/report/repository"
	"report-service/internal/report/route"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(consulClient *api.Client, reportCollection *mongo.Collection, reportHistoryCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// gateway
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	termGateway := gateway.NewTermGateway("term-service", consulClient)
	mediaGateway := gateway.NewMediaGateway("media-service", consulClient)

	// Setup dependency injection
	reportRepo := repository.NewReportRepository(reportCollection)
	historyRepo := repository.NewReportHistoryRepository(reportHistoryCollection)

	// reoport
	reportService := service.NewReportService(userGateway, termGateway, mediaGateway, reportRepo, historyRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// report history
	reportHistoryService := service.NewReportHistoryService(historyRepo)
	reportHistoryHandler := handler.NewReportHistoryHandler(reportHistoryService)

	// Register routes
	route.RegisterReportRoutes(r, reportHandler, reportHistoryHandler)
	return r
}
