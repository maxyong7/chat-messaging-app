package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type userProfileRoute struct {
	t usecase.UserProfile
	l logger.Interface
}

func newUserProfile(handler *gin.RouterGroup, t usecase.UserProfile, l logger.Interface) {
	route := &userProfileRoute{t, l}

	h := handler.Group("/profile")
	{
		h.GET("", route.getUserProfile)
		h.PATCH("", route.updateUserProfile)
	}
}

func (r *userProfileRoute) getUserProfile(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userProfile, err := r.t.GetUserInfo(c.Request.Context(), userId)
	if err != nil {
		r.l.Error(err, "http - v1 - getContacts - GetContacts")
		handleCustomErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

func (r *userProfileRoute) updateUserProfile(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request entity.UserInfoDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - updateUserProfile")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	request.UserUUID = userId

	err = r.t.UpdateUserProfile(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err, "http - v1 - updateUserProfile - updateUserProfiles")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
