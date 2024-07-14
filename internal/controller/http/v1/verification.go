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
		h.POST("/register", r.registerUser)
	}
}

type webserverResponse struct {
	Token string `json:"token,omitempty"`
}

type userRequest struct {
	Username string `json:"username"       binding:"required"  example:"username"`
	Password string `json:"password"       binding:"required"  example:"password"`
	Email    string `json:"email"       binding:"required"  example:"email"`
}

// @Summary     Verification
// @Description Verify a user
// @ID          do-verification
// @Tags  	    verification
// @Accept      json
// @Produce     json
// @Param       request body userRequest true "Set up verification"
// @Success     200 {object} webserverResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /user/verify [post]
func (r *webserverRoutes) verifyCredentials(c *gin.Context) {
	var request entity.UserRegistration
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userUuid, isValid, err := r.t.VerifyCredentials(
		c.Request.Context(),
		entity.UserCredentials{
			Username: request.Username,
			Password: request.Password,
			Email:    request.Email,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - verifyCredentials")
		handleCustomErrors(c, err)
		return
	}

	if isValid && userUuid != "" {
		token, err := createToken(userUuid)
		if err != nil {
			r.l.Error(err, "http - v1 - createToken")
			handleCustomErrors(c, err)
			return
		}

		// c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusOK, webserverResponse{
			Token: token,
		})
		return
	}

	c.JSON(http.StatusNotFound, webserverResponse{})
}

// @Summary     RegisterUser
// @Description Register user's credentials
// @ID          registerUser
// @Tags  	    registration
// @Accept      json
// @Produce     json
// @Param       request body userRequest true "Set up verification"
// @Success     200 {object} webserverResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /user/verify [post]
func (r *webserverRoutes) registerUser(c *gin.Context) {
	var request entity.UserRegistration
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := r.t.RegisterUser(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err, "http - v1 - registerUser")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}
