package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type messageRoute struct {
	t usecase.Message
	l logger.Interface
}

// Handles api routes for message functionality
func newMessageRoute(handler *gin.RouterGroup, t usecase.Message, l logger.Interface) {
	route := &messageRoute{t, l}

	// Group the routes under the "/message" path.
	h := handler.Group("/message")
	{
		// Define the endpoints for the message functionality.
		h.GET("/:conversation_uuid", route.getMessagesFromConversation)
		h.GET("/:conversation_uuid/search", route.searchMessage)
		h.GET("/status/:message_uuid", route.getMessageStatus)
	}
}

// getMessagesFromConversation handles fetching messages from a specific conversation.
func (r *messageRoute) getMessagesFromConversation(c *gin.Context) {
	// Get conversation_uuid from URL parameter
	convUUID := c.Param("conversation_uuid")
	if convUUID == "" {
		// If the conversation UUID is missing, return an error response.
		errorResponse(c, http.StatusUnprocessableEntity, "missing conversation_uuid")
		return
	}

	// Get decoded 'cursor' value from URL query
	cursor, err := queryParamCursor(c)
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - getMessagesFromConversation - cursor validation error")
		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
	}
	// Get 'limit' value from URL query and convert into integer type
	// If not value was provided, default to 20
	limit := queryParamInt(c, "limit", 20)

	// Get user_uuid from context
	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Build request params entity object
	requestParams := entity.RequestParams{
		Cursor: cursor,
		Limit:  limit,
		UserID: userId,
	}

	// Call GetMessagesFromConversation method from message entity object
	messages, err := r.t.GetMessagesFromConversation(c.Request.Context(), requestParams, convUUID)
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - getMessagesFromConversation - GetMessagesFromConversation")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Build seen status entity object
	seenStatusEntity := entity.SeenStatus{
		UserUUID:         userId,
		ConversationUUID: convUUID,
	}
	// Update the seen status of the conversation for the current user
	// by calling UpdateSeenStatus method from message entity object
	err = r.t.UpdateSeenStatus(c.Request.Context(), seenStatusEntity)
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - getMessagesFromConversation - UpdateSeenStatus")
		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Prepare the cursor for pagination, if there are more messages to load.
	var encodedCursor string
	if len(messages) == limit {
		encodedCursor = encodeCursor(&messages[len(messages)-1].CreatedAt)
	}

	// Build the response structure containing the messages and pagination information.
	msgResp := boundary.ConversationScreen{
		Data: boundary.ConversationData{
			Messages: messages,
		},
		Pagination: boundary.Pagination{
			Cursor: encodedCursor,
			Limit:  limit,
		},
	}

	// Return the response as JSON with a status code of 200 (OK).
	c.JSON(http.StatusOK, msgResp)
}

// searchMessage handles searching for messages within a specific conversation.
func (r *messageRoute) searchMessage(c *gin.Context) {
	// Get conversation_uuid from URL parameter
	convUUID := c.Param("conversation_uuid")
	if convUUID == "" {
		// If the conversation UUID is missing, return an error response.
		errorResponse(c, http.StatusUnprocessableEntity, "missing conversation_uuid")
		return
	}

	// Get 'keyword' value from URL query
	keyword, ok := c.GetQuery("keyword")
	if keyword == "" || !ok {
		// If the keyword is missing or invalid, return an error response.
		errorResponse(c, http.StatusUnprocessableEntity, "keyword query missing")
		return
	}

	// Calls SearchMessage method from message entity object
	messages, err := r.t.SearchMessage(c.Request.Context(), keyword, convUUID)
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - getMessagesFromConversation - GetMessagesFromConversation")
		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Construct the response structure containing the search results.
	searchMsgResp := boundary.MessageSearchScreen{
		Data: boundary.MessageSearchData{
			Messages: messages,
		},
	}

	// Return the search results as JSON with a status code of 200 (OK).
	c.JSON(http.StatusOK, searchMsgResp)
}

func (r *messageRoute) getMessageStatus(c *gin.Context) {
	// Get message_uuid from URL parameter
	msgUUID := c.Param("message_uuid")
	if msgUUID == "" {
		// If the message UUID is missing, return an error response.
		errorResponse(c, http.StatusUnprocessableEntity, "missing message_uuid")
		return
	}

	// Call GetSeenStatus method from message entity object
	seenStatus, err := r.t.GetSeenStatus(c.Request.Context(), msgUUID)
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - getMessagesFromConversation - GetMessagesFromConversation")
		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Construct the response structure containing the seen status.
	seenStatusResp := boundary.MessageStatusIndicator{
		SeenStatus: seenStatus,
	}

	// Return the seen status as JSON with a status code of 200 (OK).
	c.JSON(http.StatusOK, seenStatusResp)
}
