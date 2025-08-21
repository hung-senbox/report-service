package router

import (
	"report-service/internal/report/handler"
	"report-service/internal/report/repository"
	"report-service/internal/report/route"
	"report-service/internal/report/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(mongoCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// Setup dependency injection
	repo := repository.NewReportRepository(mongoCollection)
	svc := service.NewReportService(repo)
	h := handler.NewReportHandler(svc)

	// Register routes
	route.RegisterReportRoutes(r, h)
	return r
}
