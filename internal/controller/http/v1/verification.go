package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type userRoutes struct {
	t usecase.User
	l logger.Interface
}

// Handles api routes for user functionality
func newUserVerificationRoute(handler *gin.RouterGroup, t usecase.User, l logger.Interface) {
	r := &userRoutes{t, l}

	// Group the routes under the "/user" path.
	h := handler.Group("/user")
	{
		// Define the endpoints for the user functionality.
		h.POST("/login", r.loginUser)
		h.POST("/register", r.registerUser)
		h.POST("/logout", r.logoutUser)
	}
}

// loginUser handles the login process for users.
func (r *userRoutes) loginUser(c *gin.Context) {
	// Bind the incoming JSON request body to the LoginForm struct.
	var request boundary.LoginForm
	if err := c.ShouldBindJSON(&request); err != nil {
		// If the request body is invalid, log the error and return a bad request response.
		r.l.Error(err, "http - v1 - loginUser")
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}

	// Verify user credentials by calling VerifyCredentials method from user entity object
	userUuid, isValid, err := r.t.VerifyCredentials(c.Request.Context(), request.ToUserCredentials())
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - loginUser")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// If the credentials are valid and a user UUID is returned, create a session token.
	if isValid && userUuid != "" {
		token, err := createToken(userUuid, 24)
		if err != nil {
			// If an error occurs while creating the token, logs error message
			r.l.Error(err, "http - v1 - createToken")

			// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
			// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
			handleCustomErrors(c, err)
			return
		}
		// Return the token as JSON with a status code of 200 (OK).
		c.JSON(http.StatusOK, boundary.LogoutScreen{
			Token: token,
		})
		return
	}

	// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
	// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
	handleCustomErrors(c, err)
}

// registerUser handles the registration process for new users.
func (r *userRoutes) registerUser(c *gin.Context) {
	// Bind the incoming JSON request body to the RegistrationForm struct.
	var request boundary.RegistrationForm
	if err := c.ShouldBindJSON(&request); err != nil {
		// If the request body is invalid, log the error and return a bad request response.
		r.l.Error(err, "http - v1 - registerUser")
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}

	// Call RegisterUser method from user entity object
	err := r.t.RegisterUser(c.Request.Context(), request.ToUserRegistration())
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - registerUser")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	// Return a "Created" status code to indicate successful registration.
	c.Writer.WriteHeader(http.StatusCreated)
}

// logoutUser handles the logout process for users.
func (r *userRoutes) logoutUser(c *gin.Context) {
	// Invalidate the current session token by creating an expired token.
	token, err := createToken("", 0)
	if err != nil {
		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		r.l.Error(err, "http - v1 - createToken")
		handleCustomErrors(c, err)
		return
	}

	// Return the invalidated token as JSON with a status code of 200 (OK).
	c.JSON(http.StatusOK, boundary.LogoutScreen{
		Token: token,
	})
}
