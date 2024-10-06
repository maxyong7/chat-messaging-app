package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type groupChatRoute struct {
	t usecase.GroupChat
	l logger.Interface
}

// Handles api routes for groupchat functionality
func newGroupChatRoute(handler *gin.RouterGroup, t usecase.GroupChat, l logger.Interface) {
	route := &groupChatRoute{t, l}

	// Group the routes under the "/groupchat" path.
	h := handler.Group("/groupchat")
	{
		// Define the endpoints for the conversation functionality.
		h.POST("/create", route.createGroupChat)
		h.POST("/add", route.addParticipant)
		h.POST("/remove", route.removeParticipant)
		h.PATCH("/title", route.updateGroupTitle)
	}
}

// Handles the creation of a new group chat.
func (r *groupChatRoute) createGroupChat(c *gin.Context) {
	// Get user_uuid from context
	userUUID, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Bind the incoming JSON request body to the GroupChatCreationForm struct.
	var request boundary.GroupChatCreationForm
	if err := c.ShouldBindJSON(&request); err != nil {
		// If the request body is invalid, log the error and return a bad request response.
		r.l.Error(err, "http - v1 - createGroupChat")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Calls CreateGroupChat method from group chat entity object
	err = r.t.CreateGroupChat(c.Request.Context(), request.ToGroupChat(userUUID))
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - createGroupChat - CreateGroupChat")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// If the group chat is successfully created, return a "Created" status code.
	c.Writer.WriteHeader(http.StatusCreated)
}

// Handles adding a participant to a group chat.
func (r *groupChatRoute) addParticipant(c *gin.Context) {
	// Get user_uuid from context
	userUUID, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Bind the incoming JSON request body to the GroupChatCreationForm struct.
	var request boundary.GroupChatCreationForm
	if err := c.ShouldBindJSON(&request); err != nil {
		// If the request body is invalid, log the error and return a bad request response.
		r.l.Error(err, "http - v1 - addParticipant")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Calls AddParticipant method from group chat entity object
	err = r.t.AddParticipant(c.Request.Context(), request.ToGroupChat(userUUID))
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - addParticipant - AddParticipant")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// Handles removing a participant from a group chat.
func (r *groupChatRoute) removeParticipant(c *gin.Context) {
	// Get user_uuid from context
	userUUID, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Bind the incoming JSON request body to the GroupChatCreationForm struct.
	var request boundary.GroupChatCreationForm
	if err := c.ShouldBindJSON(&request); err != nil {
		// If the request body is invalid, log the error and return a bad request response.
		r.l.Error(err, "http - v1 - removeParticipant")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Calls RemoveParticipant method from group chat entity object
	err = r.t.RemoveParticipant(c.Request.Context(), request.ToGroupChat(userUUID))
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - removeParticipant - RemoveParticipant")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// If the participant is successfully removed, return an "OK" status code.
	c.Writer.WriteHeader(http.StatusOK)
}

// Handles updating the title of a group chat.
func (r *groupChatRoute) updateGroupTitle(c *gin.Context) {
	// Get user_uuid from context
	userUUID, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Bind the incoming JSON request body to the GroupChatCreationForm struct.
	var request boundary.GroupChatCreationForm
	if err := c.ShouldBindJSON(&request); err != nil {
		// If the request body is invalid, log the error and return a bad request response.
		r.l.Error(err, "http - v1 - updateGroupTitle")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Calls UpdateGroupTitle method from group chat entity object
	err = r.t.UpdateGroupTitle(c.Request.Context(), request.ToGroupChat(userUUID))
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - updateGroupTitle - UpdateGroupTitle")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// If the title is successfully updated, return an "OK" status code.
	c.Writer.WriteHeader(http.StatusOK)
}
