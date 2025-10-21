package router

import (
	"report-service/internal/gateway"
	"report-service/internal/report/handler"
	"report-service/internal/report/repository"
	"report-service/internal/report/route"
	"report-service/internal/report/service"
	"report-service/internal/report/usecase"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(consulClient *api.Client, reportCollection *mongo.Collection, reportHistoryCollection *mongo.Collection, reportPlanTemplateCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// gateway
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	termGateway := gateway.NewTermGateway("term-service", consulClient)
	mediaGateway := gateway.NewMediaGateway("media-service", consulClient)
	classroomGateway := gateway.NewClassroomGateway("classroom-service", consulClient)

	// Setup dependency injection
	reportRepo := repository.NewReportRepository(reportCollection)
	historyRepo := repository.NewReportHistoryRepository(reportHistoryCollection)
	reportPlanTemplateRepo := repository.NewReportPlanTemplateRepository(reportPlanTemplateCollection)

	// report
	reportAppUseCase := usecase.NewReportAppUseCase(reportRepo, historyRepo, userGateway, classroomGateway, termGateway, mediaGateway)
	reportWebUseCase := usecase.NewReportWebUsecase(reportRepo, historyRepo, reportPlanTemplateRepo, userGateway, classroomGateway, termGateway, mediaGateway)
	reportService := service.NewReportService(reportAppUseCase, reportWebUseCase)
	reportHandler := handler.NewReportHandler(reportService)

	// report history
	reportHistoryService := service.NewReportHistoryService(historyRepo)
	reportHistoryHandler := handler.NewReportHistoryHandler(reportHistoryService)

	// report plan template
	reportPlanTemplateService := service.NewReportPlanTemplateService(reportPlanTemplateRepo, userGateway)
	reportPlanTemplateHandler := handler.NewReportPlanTemplateHandler(reportPlanTemplateService)

	// Register routes
	route.RegisterReportRoutes(r, reportHandler, reportHistoryHandler, reportPlanTemplateHandler)
	return r
}
