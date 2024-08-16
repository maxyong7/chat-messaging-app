package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type contactRoute struct {
	t usecase.Contact
	l logger.Interface
}

// Handles api routes for contacts functionality
func newContactRoute(handler *gin.RouterGroup, t usecase.Contact, l logger.Interface) {
	route := &contactRoute{t, l}

	// Group the routes under the "/contact" path.
	h := handler.Group("/contact")
	{
		// Define the endpoints for the contacts functionality.
		h.GET("", route.getContacts)
		h.POST("/:username/add", route.addContact)
		h.POST("/:username/remove", route.removeContact)
		h.PATCH("/:username", route.updateBlockContact)
	}
}

func (r *contactRoute) getContacts(c *gin.Context) {
	// Get user_uuid from context
	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Calls GetContacts method from contact entity object
	contacts, err := r.t.GetContacts(c.Request.Context(), userId)
	if err != nil {
		// Logs the error
		r.l.Error(err, "http - v1 - getContacts - GetContacts")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Writes the status code provided in the argument.
	// It also writes a JSON body returned by the contact entity object.
	c.JSON(http.StatusOK, contacts)
}

func (r *contactRoute) addContact(c *gin.Context) {
	// Get user_uuid from context
	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get username from URL parameter
	contactUserName := c.Param("username")
	if contactUserName == "" {
		errorResponse(c, http.StatusUnprocessableEntity, "missing username in parameter")
		return
	}

	// Calls AddContact method from contact entity object
	err = r.t.AddContact(c.Request.Context(), contactUserName, userId)
	if err != nil {
		// Logs the error
		r.l.Error(err, "http - v1 - addContact - AddContacts")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Sends an HTTP response header with the provided status code.
	c.Writer.WriteHeader(http.StatusCreated)
}

func (r *contactRoute) removeContact(c *gin.Context) {
	// Get user_uuid from context
	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get username from URL parameter
	contactUserName := c.Param("username")
	if contactUserName == "" {
		errorResponse(c, http.StatusUnprocessableEntity, "missing username in parameter")
		return
	}
	// Calls RemoveContact method from contact entity object
	err = r.t.RemoveContact(c.Request.Context(), contactUserName, userId)
	if err != nil {
		// Logs the error
		r.l.Error(err, "http - v1 - addContact - AddContacts")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Sends an HTTP response header with the provided status code.
	c.Writer.WriteHeader(http.StatusNoContent)
}

func (r *contactRoute) updateBlockContact(c *gin.Context) {
	// Get user_uuid from context
	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		// If the user UUID cannot be retrieved, return an unauthorized error response.
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get username from URL parameter
	contactUserName := c.Param("username")
	if contactUserName == "" {
		errorResponse(c, http.StatusUnprocessableEntity, "missing username in parameter")
		return
	}

	// Get 'block' value from URL query
	blockString, ok := c.GetQuery("block")
	if !ok {
		errorResponse(c, http.StatusUnprocessableEntity, "block query missing")
		return
	}
	// Convert 'blockString' into lower case. Then set 'block' value accordingly
	var block bool
	switch strings.ToLower(blockString) {
	case "true":
		block = true
	case "false":
		block = false
	default:
		errorResponse(c, http.StatusUnprocessableEntity, "block query missing")
		return
	}

	// Calls UpdateBlockContact method from contact entity object
	err = r.t.UpdateBlockContact(c.Request.Context(), contactUserName, userId, block)
	if err != nil {
		// Logs the error
		r.l.Error(err, "http - v1 - updateBlockContact - UpdateBlockContact")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Sends an HTTP response header with the provided status code.
	c.Writer.WriteHeader(http.StatusNoContent)
}
