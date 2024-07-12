// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/maxyong7/chat-messaging-app/docs"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type RouterUseCases struct {
	Translation  usecase.Translation
	Verification usecase.Verification
	Conversation usecase.Conversation
	Inbox        usecase.Inbox
}

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, uc RouterUseCases) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Home route
	handler.GET("/", serveHome)

	// Routers
	h := handler.Group("/v1")
	{
		newTranslationRoutes(h, uc.Translation, l)
		newUserVerificationRoute(h, uc.Verification, l)
		newConversationRoute(h, uc.Conversation, l)
		newInboxRoute(h, uc.Inbox, l)
	}

}

func serveHome(c *gin.Context) {
	if c.Request.URL.Path != "/" {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	c.File("home.html")
}
