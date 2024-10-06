package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type userProfileRoute struct {
	t usecase.UserProfile
	l logger.Interface
}

// Handles api routes for user profile functionality
func newUserProfile(handler *gin.RouterGroup, t usecase.UserProfile, l logger.Interface) {
	route := &userProfileRoute{t, l}

	// Group the routes under the "/profile" path.
	h := handler.Group("/profile")
	{
		// Define the endpoints for the user profile functionality.
		h.GET("", route.getUserProfile)
		h.PATCH("", route.updateUserProfile)
	}
}

// getUserProfile handles the retrieval of the user's profile information.
func (r *userProfileRoute) getUserProfile(c *gin.Context) {
	// Get user_uuid from context
	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Call GetUserProfile method from user profile entity object
	userProfile, err := r.t.GetUserProfile(c.Request.Context(), userId)
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - getUserProfile - GetUserInfo")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Return the user's profile as JSON with a status code of 200 (OK).
	c.JSON(http.StatusOK, userProfile)
}

// updateUserProfile handles updating the user's profile information.
func (r *userProfileRoute) updateUserProfile(c *gin.Context) {
	// Get user_uuid from context
	userUUID, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Bind the incoming JSON request body to the ProfileSettingScreen struct.
	var request boundary.ProfileSettingScreen
	if err := c.ShouldBindJSON(&request); err != nil {
		// If the request body is invalid, log the error and return a bad request response.
		r.l.Error(err, "http - v1 - updateUserProfile")
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}

	// Call UpdateUserProfile method from user profile entity object
	err = r.t.UpdateUserProfile(c.Request.Context(), request.ToUserInfo(userUUID))
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - updateUserProfile - updateUserProfiles")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Return an "OK" status code to indicate the profile was successfully updated.
	c.Writer.WriteHeader(http.StatusOK)
}
