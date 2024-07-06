package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type webserverRoutes struct {
	t usecase.Verification
	l logger.Interface
}

func newUserVerificationRoute(handler *gin.RouterGroup, t usecase.Verification, l logger.Interface) {
	r := &webserverRoutes{t, l}

	h := handler.Group("/user")
	{
		h.POST("/verify", r.verifyCredentials)
	}
}

type webserverResponse struct {
	History []entity.Translation `json:"history"`
}

type verifyUserRequest struct {
	Username string `json:"username"       binding:"required"  example:"username"`
	Password string `json:"password"       binding:"required"  example:"password"`
	Email    string `json:"email"       binding:"required"  example:"email"`
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *webserverRoutes) verifyCredentials(c *gin.Context) {
	var request verifyUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	verification, err := r.t.VerifyCredentials(
		c.Request.Context(),
		entity.Verification{
			Username: request.Username,
			Password: request.Password,
			Email:    request.Email,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "verification service problems")

		return
	}

	c.JSON(http.StatusOK, verification)
}
