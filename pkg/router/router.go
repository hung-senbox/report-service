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

func SetupRouter(consulClient *api.Client, reportCollection, reportHistoryCollection, reportPlanTemplateCollection, reportTranslateCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// gateway
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	termGateway := gateway.NewTermGateway("term-service", consulClient)
	mediaGateway := gateway.NewMediaGateway("media-service", consulClient)
	classroomGateway := gateway.NewClassroomGateway("classroom-service", consulClient)
	fileGateway := gateway.NewFileGateway("go-main-service", consulClient)

	// Setup dependency injection
	reportRepo := repository.NewReportRepository(reportCollection)
	historyRepo := repository.NewReportHistoryRepository(reportHistoryCollection)
	reportPlanTemplateRepo := repository.NewReportPlanTemplateRepository(reportPlanTemplateCollection)

	// report
	reportAppUseCase := usecase.NewReportAppUseCase(reportRepo, historyRepo, userGateway, classroomGateway, termGateway, mediaGateway)
	reportWebUseCase := usecase.NewReportWebUsecase(reportRepo, historyRepo, reportPlanTemplateRepo, userGateway, classroomGateway, termGateway, mediaGateway, fileGateway)
	reportService := service.NewReportService(reportAppUseCase, reportWebUseCase)
	reportHandler := handler.NewReportHandler(reportService)

	// report history
	reportHistoryService := service.NewReportHistoryService(historyRepo)
	reportHistoryHandler := handler.NewReportHistoryHandler(reportHistoryService)

	// report plan template
	reportPlanTemplateService := service.NewReportPlanTemplateService(reportPlanTemplateRepo, userGateway)
	reportPlanTemplateHandler := handler.NewReportPlanTemplateHandler(reportPlanTemplateService)

	// report translate
	reportTranslateRepo := repository.NewReportTranslateRepo(reportTranslateCollection)
	reportTranslateService := service.NewReportTranslateService(reportTranslateRepo, mediaGateway)
	reportTranslateHandler := handler.NewReportTranslateHandler(reportTranslateService)

	// Register routes
	route.RegisterReportRoutes(r, reportHandler, reportHistoryHandler, reportPlanTemplateHandler, reportTranslateHandler)
	return r
}
