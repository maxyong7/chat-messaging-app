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
	Verification usecase.Verification
	Conversation usecase.Conversation
	Contact      usecase.Contact
	Message      usecase.Message
	GroupChat    usecase.GroupChat
	UserProfile  usecase.UserProfile
	Reaction     usecase.Reaction
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

	handler.Use(CORSMiddleware())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	publicHandler := handler.Group("")
	{
		// Home route
		handler.GET("/", serveHome)
		newUserVerificationRoute(publicHandler, uc.Verification, l)
	}

	// Routers
	protectedHandler := handler.Group("/v1")
	protectedHandler.Use(authMiddleware)
	{
		newConversationRoute(protectedHandler, uc.Conversation, uc.UserProfile, uc.Message, uc.Reaction, l)
		newContactRoute(protectedHandler, uc.Contact, l)
		newMessageRoute(protectedHandler, uc.Message, l)
		newGroupChatRoute(protectedHandler, uc.GroupChat, l)
		newUserProfile(protectedHandler, uc.UserProfile, l)
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
