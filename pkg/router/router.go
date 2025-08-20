package router

import (
	"report-service/internal/term/handler"
	"report-service/internal/term/repository"
	"report-service/internal/term/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(mongoCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// Setup dependency injection
	repo := repository.NewTermRepository(mongoCollection)
	svc := service.NewTermService(repo)
	h := handler.NewHandler(svc)

	// Register routes
	h.RegisterRoutes(r)

	return r
}
